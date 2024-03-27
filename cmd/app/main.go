package main

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/config"
	"github.com/ichenhe/cert-deployer/domain"
	_ "github.com/ichenhe/cert-deployer/plugins"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
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
	core := zapcore.NewCore(createLoggerEncoder("fluent"), zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	logger = zap.New(core, zap.AddCaller()).Sugar()
}

func main() {
	defer func() { _ = logger.Sync() }()

	defaultInitializer := initializerFunc(func(c *cli.Context, modifier func(k *koanf.Koanf)) (*domain.AppConfig, error) {
		if profile, err := config.CreateWithModifier(c, modifier); err == nil {
			if err := setupLogger(profile.LogDrivers); err != nil {
				return nil, fmt.Errorf("failed to setup logger: %w", err)
			}
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
					Usage:    "asset type, e.g. cdn",
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
