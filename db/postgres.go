package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func buildPostgresDialector(cfg *Config) gorm.Dialector {
	sslmode := cfg.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, sslmode,
	)
	return postgres.Open(dsn)
}
