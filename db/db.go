package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// DB 数据库连接（全局）
var DB *gorm.DB

// Config 数据库连接配置，同时支持 MySQL 与 PostgreSQL
type Config struct {
	Driver   string `toml:"driver" json:"driver"` // "mysql" | "postgres"
	Host     string `toml:"host" json:"host"`
	Port     int    `toml:"port" json:"port"`
	Username string `toml:"username" json:"username"`
	Password string `toml:"password" json:"password"`
	Database string `toml:"database" json:"database"`

	// Postgres only，默认 "disable"
	SSLMode string `toml:"sslmode" json:"sslmode"`

	// 连接池（0 表示用 database/sql 默认值）
	MaxOpenConns int           `toml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns int           `toml:"max_idle_conns" json:"max_idle_conns"`
	ConnMaxLife  time.Duration `toml:"conn_max_life" json:"conn_max_life"`
}

// Init 根据 cfg.Driver 初始化全局 DB
func Init(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("db config is nil")
	}
	dialector, err := buildDialector(cfg)
	if err != nil {
		return err
	}
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLife > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLife)
	}
	DB = gormDB
	return nil
}

func buildDialector(cfg *Config) (gorm.Dialector, error) {
	switch cfg.Driver {
	case "mysql":
		return buildMysqlDialector(cfg), nil
	case "postgres", "postgresql":
		return buildPostgresDialector(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported db driver: %q", cfg.Driver)
	}
}
