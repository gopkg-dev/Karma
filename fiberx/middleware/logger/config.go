package logger

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

// Config defines the config for middleware.
type Config struct {
	MaxRequestBodyLen  int
	MaxResponseBodyLen int
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool
	// Done is a function that is called after the log string for a request is written to Output,
	// and pass the log string as parameter.
	//
	// Optional. Default: nil
	Done func(c *fiber.Ctx, args fiber.Map)
	// CustomAttr defines the custom tag action
	//
	// Optional. Default: []CustomFunc
	CustomAttr   []CustomFunc
	BuiltinAttrs []string
	Logger       *slog.Logger
}

type CustomFunc func(*fiber.Ctx, error) []slog.Attr

var HiddenRequestHeaders = map[string]struct{}{
	"authorization": {},
	"cookie":        {},
	"set-cookie":    {},
	"x-auth-token":  {},
	"x-csrf-token":  {},
	"x-xsrf-token":  {},
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	MaxRequestBodyLen:  1024 * 1024,
	MaxResponseBodyLen: 1024 * 1024,
	Next:               nil,
	Done:               nil,
	CustomAttr:         nil,
	BuiltinAttrs: []string{
		TagPid,
		TagReferer,
		TagProtocol,
		TagPort,
		TagIP,
		TagIPs,
		TagHost,
		TagMethod,
		TagPath,
		TagUrl,
		TagUA,
		TagLatency,
		TagStatus,
		TagResBody,
		TagReqHeaders,
		TagQueryStringParams,
		TagBody,
		TagBytesSent,
		TagBytesReceived,
		TagRoute,
		TagError,
	},
	Logger: slog.Default(),
}

func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}

	if cfg.Done == nil {
		cfg.Done = ConfigDefault.Done
	}

	if cfg.BuiltinAttrs == nil {
		cfg.BuiltinAttrs = ConfigDefault.BuiltinAttrs
	}

	if cfg.Logger == nil {
		cfg.Logger = ConfigDefault.Logger
	}

	return cfg
}
