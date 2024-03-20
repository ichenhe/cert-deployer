package main

import (
	"errors"
	"fmt"
	"github.com/ichenhe/cert-deployer/config"
	"github.com/ichenhe/cert-deployer/domain"
	_ "github.com/ichenhe/cert-deployer/plugins"
	"github.com/ichenhe/cert-deployer/registry"
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
				&cli.GenericFlag{
					Name:     "cert",
					Aliases:  []string{"c"},
					Usage:    "full chain public key pem file",
					Value:    &fileType{},
					Required: true,
				},
				&cli.GenericFlag{
					Name:     "key",
					Aliases:  []string{"k"},
					Usage:    "private key pem file",
					Value:    &fileType{},
					Required: true,
				},
				&cli.StringSliceFlag{
					Name:     "type",
					Aliases:  []string{"t"},
					Usage:    "asset types, e.g. cdn",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				appConfig, err := readProfile(c)
				if err != nil {
					return err
				}
				setLogger(appConfig)

				certData := c.Generic("cert").(*fileType)
				keyData := c.Generic("key").(*fileType)

				types := make([]domain.AssetType, 0)
				for _, str := range c.StringSlice("type") {
					if t, err := domain.AssetTypeFromString(str); err == nil {
						types = append(types, t)
					} else {
						logger.Infof("ignore invalid asset type '%s'", str)
					}
				}

				hasError := false
				uniDeployer := registry.NewUnionDeployer(logger, appConfig.CloudProviders)

				logger.Infof("look for %d types of applicable assets: %v", len(types), types)
				allAssets := make([]domain.Asseter, 0, 64)
				for _, t := range types {
					if assets, err := uniDeployer.ListApplicableAssets(t, certData.data); err != nil {
						hasError = true
						logger.Errorf("failed to search assets for type '%s': %v", t, err)
					} else {
						for _, item := range assets {
							logger.Debugf(item.GetBaseInfo().String())
						}
						allAssets = append(allAssets, assets...)
					}
				}
				logger.Infof("a total of %d assets were acquired, deploying...", len(allAssets))

				deployed, hasDeployErr := uniDeployer.Deploy(allAssets, certData.data, keyData.data)
				if hasDeployErr {
					hasError = true
				} else {
					logger.Infof("%d assets deployed successfully: %v", len(deployed),
						utils.MapSlice(deployed, func(s domain.Asseter) string {
							i := s.GetBaseInfo()
							return fmt.Sprintf("%s@%s-%s", i.Type, i.Provider, i.Name)
						}))
				}

				if hasError {
					return errors.New("some errors have occurred")
				} else {
					return nil
				}
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}

func readProfile(c *cli.Context) (*domain.AppConfig, error) {
	file := c.Path("profile")
	appConfig, err := config.ReadConfig(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile: %w", err)
	}
	return appConfig, nil
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
