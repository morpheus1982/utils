package debug

import (
	"fmt"
	"net/http"

	"github.com/morpheus1982/utils/logger"

	_ "net/http/pprof"

	"go.uber.org/zap"
)

type Config struct {
	Enabled bool `toml:"enabled" json:"enabled"`
	Port    int  `toml:"port" json:"port"`
}

func Serve(port int) {
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	if err != nil {
		logger.Warn("Start debug server failed!", zap.Error(err))
	}
}
