package demo

import (
	"flag"
	"fmt"
	"log"

	"github.com/morpheus1982/utils/config"
	"github.com/morpheus1982/utils/logger"
)

var (
	help       bool
	configFile string
	printLog   string
	version    bool
)

func init() {
	flag.StringVar(&configFile, "c", "config.toml", "set configuration `file`")
	flag.StringVar(&printLog, "p", "", "print log")
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "print version")
}

var _BUILD_ = ""

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	if version {
		fmt.Println(_BUILD_)
		return
	}
	// 加载配置文件
	err := config.LoadConfig(configFile)
	if err != nil {
		log.Println("Load config failed! error: ", err)
		return
	}
	// 初始化日志记录器
	if err := logger.InitLogger(&config.Cfg.Logger, printLog); err != nil {
		log.Println("Init logger failed! error: ", err)
		return
	}
}
