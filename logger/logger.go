package logger

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var elg *zap.Logger // 错误日志
var alg *zap.Logger // 访问日志

func InitLogger(cfg *Config, printLog string) error {
	l, err := createLogger(cfg, printLog)
	if err != nil {
		return err
	}

	elg = l

	if cfg.Access.Filename != "" {
		l, err = createAccessLogger(&cfg.Access)
		if err != nil {
			return err
		}

		alg = l
	}

	return nil
}

func createLogger(cfg *Config, printLog string) (*zap.Logger, error) {
	var cores []zapcore.Core

	if printLog != "" {
		level := zap.NewAtomicLevel()
		if err := level.UnmarshalText([]byte(printLog)); err != nil {
			return nil, err
		}
		consoleCore, err := initConsoleLogger(level)
		if err != nil {
			return nil, err
		}

		cores = append(cores, consoleCore)
	}

	if len(cfg.Filename) > 0 {
		level := zap.NewAtomicLevel()
		if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
			return nil, err
		}

		if cfg.Port > 0 {
			mux := http.NewServeMux()
			mux.HandleFunc("/logger/level", level.ServeHTTP)
			srv := &http.Server{
				Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
				Handler: mux,
			}
			go func() {
				defer func() {
					if err := recover(); err != nil {
						log.Println(err)
					}
				}()
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Println(err)
					return
				}
			}()
		}

		fileCore, err := initFileLogger(cfg, level)
		if err == nil {
			cores = append(cores, fileCore)
		}
	}

	var opts []zap.Option
	opts = append(opts, zap.AddCaller())
	core := zapcore.NewTee(cores...)
	return zap.New(core, opts...), nil
}

func createAccessLogger(cfg *AccessConfig) (*zap.Logger, error) {
	var cores []zapcore.Core

	if len(cfg.Filename) > 0 {
		level := zap.NewAtomicLevelAt(zap.DebugLevel)
		fileCore, err := initAccessFileLogger(cfg, level)
		if err == nil {
			cores = append(cores, fileCore)
		}
	}

	var opts []zap.Option
	opts = append(opts, zap.WithCaller(false))
	core := zapcore.NewTee(cores...)
	return zap.New(core, opts...), nil
}

func initConsoleLogger(level zap.AtomicLevel) (core zapcore.Core, err error) {
	output := zapcore.Lock(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), output, level), nil
}

func initFileLogger(cfg *Config, level zap.AtomicLevel) (core zapcore.Core, err error) {
	if st, err := os.Stat(cfg.Filename); err == nil {
		if st.IsDir() {
			return nil, errors.New("can't use directory as log file name")
		}
	}

	if cfg.MaxSize == 0 {
		cfg.MaxSize = defaultLogMaxSize
	}

	output := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxDays,
		LocalTime:  true,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), output, level), nil
}

func initAccessFileLogger(cfg *AccessConfig, level zap.AtomicLevel) (core zapcore.Core, err error) {
	if st, err := os.Stat(cfg.Filename); err == nil {
		if st.IsDir() {
			return nil, errors.New("can't use directory as log file name")
		}
	}

	if cfg.MaxSize == 0 {
		cfg.MaxSize = defaultLogMaxSize
	}

	output := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxDays,
		LocalTime:  true,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeLevel = nil
	encoderConfig.MessageKey = ""

	return zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), output, level), nil
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	s := t.Format("2006/01/02 15:04:05.000")
	enc.AppendString(s)
}
