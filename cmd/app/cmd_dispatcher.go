package main

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
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
	loadProfile profileLoader
	fileReader  domain.FileReader
	cmdExecutor commandExecutor
}

func newCommandDispatcher(loadProfile profileLoader, fileReader domain.FileReader, cmdExecutor commandExecutor) commandDispatcher {
	return &defaultCommandDispatcher{
		loadProfile: loadProfile,
		fileReader:  fileReader,
		cmdExecutor: cmdExecutor,
	}
}

func (d *defaultCommandDispatcher) deploy(c *cli.Context) error {
	appConfig, err := d.loadProfile(c)
	if err != nil {
		return err
	}
	setLogger(appConfig)

	if deploymentIds := c.StringSlice("deployment"); deploymentIds != nil && len(deploymentIds) > 0 {
		d.cmdExecutor.executeDeployments(appConfig, deploymentIds)
		return nil
	}

	// check arguments
	requiredFlags := []string{"cert", "key", "type"}
	for _, flag := range requiredFlags {
		if c.Generic(flag) == nil {
			return fmt.Errorf("flags %v must be provided without --deployment", requiredFlags)
		}
	}

	certData, err := d.fileReader.ReadFile(c.Path("cert"))
	if err != nil {
		return fmt.Errorf("invalid public cert: %w", err)
	}
	keyData, err := d.fileReader.ReadFile(c.Path("key"))
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}

	d.cmdExecutor.customDeploy(appConfig.CloudProviders, c.StringSlice("type"), certData, keyData)
	return nil
}

func (d *defaultCommandDispatcher) run(c *cli.Context) error {
	appConfig, err := d.loadProfile(c)
	if err != nil {
		return err
	}
	setLogger(appConfig)

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
