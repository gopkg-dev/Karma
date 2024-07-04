package logging

import (
	"context"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config 定义了 Logger 的配置
type Config struct {
	Debug      bool
	Level      string // debug/info/warn/error/dpanic/panic/fatal
	CallerSkip int
	File       struct {
		Enable     bool
		Path       string
		MaxSize    int
		MaxBackups int
	}
	Hooks []*HookConfig
}

// HookConfig 定义了 Hook 的配置
type HookConfig struct {
	Enable    bool
	Level     string
	Type      string // gorm
	MaxBuffer int
	MaxThread int
	Options   map[string]string
	Extra     map[string]string
}

// HookHandlerFunc 是一个用于处理 Hook 的函数类型
type HookHandlerFunc func(ctx context.Context, hookCfg *HookConfig) (*Hook, error)

// InitWithConfig 根据配置初始化 Logger
func InitWithConfig(ctx context.Context, cfg *Config, hookHandle ...HookHandlerFunc) (*zap.Logger, func(), error) {
	var config zap.Config
	if cfg.Debug {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	// 设置日志级别
	if err := setLogLevel(cfg, &config); err != nil {
		return nil, nil, err
	}

	var (
		logger   *zap.Logger
		cleanFns []func()
	)

	// 初始化文件日志
	if cfg.File.Enable {
		fileWriter, cleanFn, err := initFileLogger(cfg, &config)
		if err != nil {
			return nil, nil, err
		}
		cleanFns = append(cleanFns, cleanFn)
		logger = zap.New(fileWriter)
	} else {
		iLogger, err := config.Build()
		if err != nil {
			return nil, nil, err
		}
		logger = iLogger
	}

	// 设置调用者跳过级数
	//skip := cfg.CallerSkip
	//if skip <= 0 {
	//	skip = 1
	//}

	logger = logger.WithOptions(
		zap.WithCaller(true),
		zap.AddStacktrace(zap.ErrorLevel),
		//zap.AddCallerSkip(skip),
	)

	// 初始化 Hooks
	for _, h := range cfg.Hooks {
		if !h.Enable || len(hookHandle) == 0 {
			continue
		}

		writer, err := hookHandle[0](ctx, h)
		if err != nil {
			return logger, nil, err
		} else if writer == nil {
			continue
		}

		cleanFns = append(cleanFns, func() {
			writer.Flush()
		})

		hookCore := initHookCore(h, writer)
		logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, hookCore)
		}))
	}

	// 替换全局 logger
	zap.ReplaceGlobals(logger)

	// 返回 logger 和清理函数
	return logger, func() {
		for _, fn := range cleanFns {
			fn()
		}
	}, nil
}

// setLogLevel 设置日志级别
func setLogLevel(cfg *Config, config *zap.Config) error {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	config.Level.SetLevel(level)
	return nil
}

// initFileLogger 初始化文件日志
func initFileLogger(cfg *Config, config *zap.Config) (zapcore.Core, func(), error) {
	filename := cfg.File.Path
	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return nil, nil, err
	}

	fileWriter := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.File.MaxSize,
		MaxBackups: cfg.File.MaxBackups,
		Compress:   false,
		LocalTime:  true,
	}

	cleanFn := func() {
		_ = fileWriter.Close()
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		zapcore.AddSync(fileWriter),
		config.Level,
	)

	return core, cleanFn, nil
}

// initHookCore 初始化 Hook Core
func initHookCore(h *HookConfig, writer *Hook) zapcore.Core {
	hookLevel := zap.NewAtomicLevel()
	if level, err := zapcore.ParseLevel(h.Level); err == nil {
		hookLevel.SetLevel(level)
	} else {
		hookLevel.SetLevel(zap.InfoLevel)
	}
	hookEncoder := zap.NewProductionEncoderConfig()
	hookEncoder.EncodeTime = zapcore.EpochMillisTimeEncoder
	hookEncoder.EncodeDuration = zapcore.MillisDurationEncoder
	hookCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(hookEncoder),
		zapcore.AddSync(writer),
		hookLevel,
	)
	return hookCore
}
