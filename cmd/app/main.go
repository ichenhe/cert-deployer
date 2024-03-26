package main

import (
	"github.com/ichenhe/cert-deployer/config"
	"github.com/ichenhe/cert-deployer/domain"
	_ "github.com/ichenhe/cert-deployer/plugins"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"strings"
)

// global logger
var logger *zap.SugaredLogger

type initializer interface {
	// LoadProfileAndSetupLogger loads profile and setups logger if no errors.
	LoadProfileAndSetupLogger(c *cli.Context, modifier func(k *koanf.Koanf)) (*domain.AppConfig, error)
}

type initializerFunc func(c *cli.Context, modifier func(k *koanf.Koanf)) (*domain.AppConfig, error)

func (f initializerFunc) LoadProfileAndSetupLogger(c *cli.Context, modifier func(k *koanf.Koanf)) (*domain.AppConfig, error) {
	return f(c, modifier)
}

func init() {
	encoder := newZapLogEncoder()
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	l := zap.New(core, zap.AddCaller())
	logger = l.Sugar()
}

func main() {
	defer func() { _ = logger.Sync() }()

	defaultInitializer := initializerFunc(func(c *cli.Context, modifier func(k *koanf.Koanf)) (*domain.AppConfig, error) {
		if profile, err := config.CreateWithModifier(c, modifier); err == nil {
			setupLogger(&profile.Log)
			return profile, nil
		} else {
			return nil, err
		}
	})

	err := run(os.Args, newCommandDispatcher(defaultInitializer, domain.FileReaderFunc(os.ReadFile), newCommandExecutor()))
	if err != nil {
		logger.Fatal(err)
	}
}

func run(args []string, cmdDispatcher commandDispatcher) error {
	app := &cli.App{
		Name:  "cert-deployer",
		Usage: "Deployer your https cert to various cloud services assets",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "profile",
				Usage:    "specify the config file manually",
				Required: false,
			},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name: "run", Usage: "Register all triggers and start listening",
			Action: cmdDispatcher.run,
		},
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
				&cli.StringFlag{
					Name:     "type",
					Usage:    "asset types, e.g. cdn",
					Required: false,
					Category: "Custom:",
				},
				&cli.StringFlag{
					Name:     "provider",
					Usage:    "cloud service provider, must in support list, e.g. TencentCloud",
					Required: false,
					Category: "Custom:",
				},
				&cli.StringFlag{
					Name:     "secret-id",
					Usage:    "api secret id of the provider",
					Required: false,
					Category: "Custom:",
				},
				&cli.StringFlag{
					Name:     "secret-key",
					Usage:    "api secret key of the provider",
					Required: false,
					Category: "Custom:",
				},
			},
			Action: cmdDispatcher.deploy,
		},
	}
	return app.Run(args)
}

func setupLogger(logConfig *domain.LogConfig) {
	logLevel := zapcore.InfoLevel

	switch strings.ToLower(logConfig.Level) {
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
		if !domain.IsDir(logConfig.FileDir) {
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

func newZapLogEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	return encoder
}
