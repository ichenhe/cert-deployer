package main

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
)

func createLoggerEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	switch format {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "fluent":
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func setupLogger(drivers []domain.LogDriver) error {
	parseLevel := func(level string) zapcore.Level {
		switch strings.ToLower(level) {
		case "debug":
			return zapcore.DebugLevel
		case "warn":
			return zapcore.WarnLevel
		case "error":
			return zap.ErrorLevel
		default:
			return zap.InfoLevel
		}
	}

	zapCores := make([]zapcore.Core, 0, len(drivers))
	for i, driver := range drivers {
		switch driver.Driver {
		case "stdout":
			zapCores = append(zapCores, zapcore.NewCore(createLoggerEncoder(driver.Format), zapcore.AddSync(os.Stdout), parseLevel(driver.Level)))
		case "file":
			if file, ex := driver.Options["file"]; !ex {
				return fmt.Errorf("missing log-drivers[%d].options.file", i)
			} else if fileStr, ok := file.(string); !ok {
				return fmt.Errorf("log-drivers[%d].options.file shoue be a string", i)
			} else if strings.HasSuffix(fileStr, "/") {
				return fmt.Errorf("log-drivers[%d].options.file can not be a directory", i)
			} else {
				dir := filepath.Dir(fileStr)
				if !domain.IsDir(dir) {
					if err := os.MkdirAll(dir, 0755); err != nil {
						return fmt.Errorf("failed to create log dir: %s", dir)
					}
				}
				if f, err := os.OpenFile(fileStr, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755); err != nil {
					return fmt.Errorf("failed to create log file: %s", fileStr)
				} else {
					zapCores = append(zapCores, zapcore.NewCore(createLoggerEncoder(driver.Format), zapcore.AddSync(f), parseLevel(driver.Level)))
				}
			}
		}
	}

	// create logger
	if len(zapCores) > 0 {
		tee := zapcore.NewTee(zapCores...)
		logger = zap.New(tee).Sugar()
	} else {
		logger = zap.NewNop().Sugar()
	}

	return nil
}
