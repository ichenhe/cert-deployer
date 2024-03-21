package main

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/config"
	"github.com/ichenhe/cert-deployer/domain"
	_ "github.com/ichenhe/cert-deployer/plugins"
	"github.com/ichenhe/cert-deployer/utils"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
)

var logger *zap.SugaredLogger

func newZapLogEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	return encoder
}

func init() {
	encoder := newZapLogEncoder()
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	l := zap.New(core, zap.AddCaller())
	logger = l.Sugar()
}

func main() {
	defer func() { _ = logger.Sync() }()

	defaultProfileLoader := func(c *cli.Context) (*domain.AppConfig, error) {
		file := c.Path("profile")
		appConfig, err := config.ReadConfig(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read profile: %w", err)
		}
		return appConfig, nil
	}

	err := run(os.Args, newCommandDispatcher(defaultProfileLoader, domain.FileReaderFunc(os.ReadFile), newCommandExecutor()))
	if err != nil {
		logger.Fatal(err)
	}
}

type profileLoader = func(c *cli.Context) (*domain.AppConfig, error)

func run(args []string, cmdDispatcher commandDispatcher) error {
	app := &cli.App{
		Name:  "cert-deployer",
		Usage: "Deployer your https cert to various cloud services assets",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:  "profile",
				Usage: "specify the config file manually",
				Value: "cert-deployer.yaml",
			},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name: "deploy", Usage: "Deploy certs to cloud services",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:     "deployment",
					Aliases:  []string{"d"},
					Usage:    "name(id) of deployments defined in the profile",
					Required: false,
				},
				&cli.PathFlag{
					Name:     "cert",
					Usage:    "/path/to/fullchain.pem",
					Required: false,
					Category: "Custom:",
				},
				&cli.PathFlag{
					Name:     "key",
					Usage:    "/path/to/private.pem",
					Required: false,
					Category: "Custom:",
				},
				&cli.StringSliceFlag{
					Name:     "type",
					Usage:    "asset types, e.g. cdn",
					Required: false,
					Category: "Custom:",
				},
			},
			Action: cmdDispatcher.deploy,
		},
	}
	return app.Run(args)
}

func setLogger(profile *domain.AppConfig) {
	logConfig := profile.Log
	logLevel := zapcore.InfoLevel
	switch logConfig.Level {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	default:
		logLevel = zap.InfoLevel
	}

	if logConfig.EnableFile {
		var f *os.File
		var err error
		if !utils.IsDir(logConfig.FileDir) {
			_ = os.MkdirAll(logConfig.FileDir, 0755)
		}
		if f, err = os.OpenFile(path.Join(logConfig.FileDir, "cert-deployer.log"),
			os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755); err != nil {
			logger.Panicf("failed to create log file: %v", err)
		}

		multiSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(f), zapcore.AddSync(os.Stdout))
		core := zapcore.NewCore(newZapLogEncoder(), multiSyncer, logLevel)
		l := zap.New(core)
		logger = l.Sugar()
	}
}
