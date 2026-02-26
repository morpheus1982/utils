package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/morpheus1982/utils/db"
	"github.com/morpheus1982/utils/debug"
	"github.com/morpheus1982/utils/logger"

	"github.com/pelletier/go-toml"
	"golang.org/x/net/proxy"
)

// Socks5Config SOCKS5代理配置
type Socks5Config struct {
	Enabled  bool   `toml:"enabled" json:"enabled"`   // 是否启用代理
	Host     string `toml:"host" json:"host"`         // 代理服务器地址
	Port     int    `toml:"port" json:"port"`         // 代理服务器端口
	Username string `toml:"username" json:"username"` // 用户名（可选）
	Password string `toml:"password" json:"password"` // 密码（可选）
}

// Address 返回代理地址 host:port
func (s *Socks5Config) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Dialer 返回 SOCKS5 代理拨号器
func (s *Socks5Config) Dialer() (proxy.Dialer, error) {
	if !s.Enabled {
		return proxy.Direct, nil
	}

	var auth *proxy.Auth
	if s.Username != "" {
		auth = &proxy.Auth{
			User:     s.Username,
			Password: s.Password,
		}
	}

	return proxy.SOCKS5("tcp", s.Address(), auth, proxy.Direct)
}

// ReaderConfig reader 服务配置
type ReaderConfig struct {
	Port              int `toml:"port"`
	TaskTimeout       int `toml:"task_timeout"`
	MaxConcurrentTasks int `toml:"max_concurrent_tasks"`
	LogRetentionDays  int `toml:"log_retention_days"`
}

// GetTimeout 获取超时时间
func (c *ReaderConfig) GetTimeout() time.Duration {
	if c.TaskTimeout <= 0 {
		return 10 * time.Minute
	}
	return time.Duration(c.TaskTimeout) * time.Second
}

type Config struct {
	Logger logger.Config `toml:"logger" json:"logger"`
	Debug  debug.Config  `toml:"debug" json:"debug"`
	Mysql  db.Mysql      `toml:"mysql" json:"mysql"`
	Socks5 Socks5Config  `toml:"socks5" json:"socks5"`
	Reader ReaderConfig  `toml:"reader" json:"reader"`
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
