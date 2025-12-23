package config

import (
	"path/filepath"

	"github.com/morpheus1982/utils/db"
	"github.com/morpheus1982/utils/debug"
	"github.com/morpheus1982/utils/logger"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Logger logger.Config `toml:"logger" json:"logger"`
	Debug  debug.Config  `toml:"debug" json:"debug"`
	Mysql  db.Mysql      `toml:"mysql" json:"mysql"`
}

var Cfg *Config

func LoadConfig(configFile string) error {
	filePath, err := filepath.Abs(configFile)
	if err != nil {
		return err
	}

	cfg := &Config{}

	tree, err := toml.LoadFile(filePath)
	if err != nil {
		return err
	}

	if err := tree.Unmarshal(cfg); err != nil {
		return err
	}

	Cfg = cfg

	return nil
}
