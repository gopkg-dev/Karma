package gormx

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	sdmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type ResolverConfig struct {
	DBType   string   // mysql/postgres/sqlite3
	Sources  []string //
	Replicas []string //
	Tables   []string //
}

// Config 配置参数
type Config struct {
	Debug                                    bool             // 是否开启调试模式
	PrepareStmt                              bool             //
	DBType                                   string           // 数据库类型,mysql/postgres/sqlite3
	DSN                                      string           // 数据库链接字符串
	MaxLifetime                              int              // 连接最长存活期,超过这个时间连接将不再被复用
	MaxIdleTime                              int              // 设置连接空闲的最大时间
	MaxOpenConns                             int              // 数据库最大连接数
	MaxIdleConns                             int              // 最大空闲连接数
	TablePrefix                              string           // 表名前缀
	DisableForeignKeyConstraintWhenMigrating bool             // 迁移时禁用外键约束
	Resolver                                 []ResolverConfig //
}

// New 创建DB实例
func New(cfg Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(cfg.DBType) {
	case "mysql":
		if err := createDatabaseWithMySQL(cfg.DSN); err != nil {
			return nil, err
		}
		dialector = mysql.Open(cfg.DSN)
	case "postgres":
		dialector = postgres.Open(cfg.DSN)
	case "sqlite3":
		_ = os.MkdirAll(filepath.Dir(cfg.DSN), os.ModePerm)
		dialector = sqlite.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}

	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix,
			SingularTable: true,
		},
		PrepareStmt:                              cfg.PrepareStmt,
		Logger:                                   logger.Discard,
		DisableForeignKeyConstraintWhenMigrating: cfg.DisableForeignKeyConstraintWhenMigrating,
	}

	if cfg.Debug {
		config.Logger = logger.Default
	}

	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}

	if len(cfg.Resolver) > 0 {
		resolver := &dbresolver.DBResolver{}
		for _, r := range cfg.Resolver {
			resolverCfg := dbresolver.Config{}
			var open func(dsn string) gorm.Dialector
			dbType := strings.ToLower(r.DBType)
			switch dbType {
			case "mysql":
				open = mysql.Open
			case "postgres":
				open = postgres.Open
			case "sqlite3":
				open = sqlite.Open
			default:
				continue
			}

			for _, replica := range r.Replicas {
				if dbType == "sqlite3" {
					_ = os.MkdirAll(filepath.Dir(cfg.DSN), os.ModePerm)
				}
				resolverCfg.Replicas = append(resolverCfg.Replicas, open(replica))
			}
			for _, source := range r.Sources {
				if dbType == "sqlite3" {
					_ = os.MkdirAll(filepath.Dir(cfg.DSN), os.ModePerm)
				}
				resolverCfg.Sources = append(resolverCfg.Sources, open(source))
			}
			tables := stringSliceToInterfaceSlice(r.Tables)
			resolver.Register(resolverCfg, tables...)
			fmt.Printf("Use resolver, #tables: %v, #replicas: %v, #sources: %v \n",
				tables, r.Replicas, r.Sources)
		}

		resolver.SetMaxIdleConns(cfg.MaxIdleConns).
			SetMaxOpenConns(cfg.MaxOpenConns).
			SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second).
			SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Second)
		if err = db.Use(resolver); err != nil {
			return nil, err
		}
	}

	if cfg.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Second)

	return db, nil
}

func stringSliceToInterfaceSlice(s []string) []interface{} {
	r := make([]interface{}, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

// AutoMigrate 自动映射数据表
func AutoMigrate(db *gorm.DB, dst ...interface{}) error {
	if db.Dialector.Name() == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	return db.AutoMigrate(dst...)
}

func createDatabaseWithMySQL(dsn string) error {
	cfg, err := sdmysql.ParseDSN(dsn)
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", cfg.User, cfg.Passwd, cfg.Addr))
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET = `utf8mb4`;", cfg.DBName)
	_, err = db.Exec(query)
	return err
}
