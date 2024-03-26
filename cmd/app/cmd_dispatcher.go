package main

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
)

// commandDispatcher dispatchers commands registered in cli to executor based on flags and arguments.
//
// If the command needs background running, corresponding function in this interface will not return,
// instead it catches the system sig to exit.
type commandDispatcher interface {
	deploy(c *cli.Context) error
	run(c *cli.Context) error
}

type defaultCommandDispatcher struct {
	initializer initializer
	fileReader  domain.FileReader
	cmdExecutor commandExecutor
}

func newCommandDispatcher(initializer initializer, fileReader domain.FileReader, cmdExecutor commandExecutor) commandDispatcher {
	return &defaultCommandDispatcher{
		initializer: initializer,
		fileReader:  fileReader,
		cmdExecutor: cmdExecutor,
	}
}

func (d *defaultCommandDispatcher) deploy(c *cli.Context) error {
	if deploymentIds := c.StringSlice("deployment"); deploymentIds != nil && len(deploymentIds) > 0 {
		appConfig, err := d.initializer.LoadProfileAndSetupLogger(c, nil)
		if err != nil {
			return err
		}
		d.cmdExecutor.executeDeployments(appConfig.CloudProviders, appConfig.Deployments, deploymentIds)
		return nil
	}

	// check arguments
	requiredFlags := []string{"provider", "secret-id", "secret-key", "cert", "key", "type"}
	for _, flag := range requiredFlags {
		if c.Generic(flag) == nil {
			return fmt.Errorf("flags %v must be provided without --deployment", requiredFlags)
		}
	}

	// load & generate profile
	appConfig, err := d.initializer.LoadProfileAndSetupLogger(c, func(k *koanf.Koanf) {
		_ = k.Set("cloud-providers.from-cli-1", domain.CloudProvider{
			Provider:  c.String("provider"),
			SecretId:  c.String("secret-id"),
			SecretKey: c.String("secret-key"),
		})

		_ = k.Set("deployments.from-cli-1", domain.Deployment{
			ProviderId: "from-cli-1",
			Cert:       c.Path("cert"),
			Key:        c.Path("key"),
			Assets:     []domain.DeploymentAsset{{Type: c.String("type")}},
		})

	})
	if err != nil {
		return err
	}

	d.cmdExecutor.executeDeployments(appConfig.CloudProviders, appConfig.Deployments, []string{"from-cli-1"})
	return nil
}

func (d *defaultCommandDispatcher) run(c *cli.Context) error {
	appConfig, err := d.initializer.LoadProfileAndSetupLogger(c, nil)
	if err != nil {
		return err
	}

	triggers := d.cmdExecutor.registerTriggers(appConfig.CloudProviders, appConfig.Deployments, appConfig.Triggers)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		logger.Infof("Shutting down...")
		for _, trigger := range triggers {
			trigger.Close()
		}
		done <- true
	}()
	<-done
	return nil
}
