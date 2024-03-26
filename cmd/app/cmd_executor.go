package main

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/deploy"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/registry"
	"github.com/ichenhe/cert-deployer/trigger/filetrigger"
)

// commandExecutor executes intent operation with given arguments. Normally it should be called from
// a commandDispatcher.
type commandExecutor interface {
	// executeDeployments executes deployments who have specific id defined in profile.
	executeDeployments(providers map[string]domain.CloudProvider, deployments map[string]domain.Deployment, deploymentIds []string)

	// customDeploy reads configurations from command line and processes a one-time deployment.
	customDeploy(providers map[string]domain.CloudProvider, rawTypes []string, cert []byte, key []byte)

	// registerTriggers finds the deployment and starts the trigger if it has at least one valid deployment.
	// Returns started triggers.
	registerTriggers(providers map[string]domain.CloudProvider, deployments map[string]domain.Deployment, triggerDefs map[string]domain.TriggerDefiner) (registeredTriggers []domain.Trigger)
}

type defaultCommandExecutor struct {
}

func newCommandExecutor() commandExecutor {
	return &defaultCommandExecutor{}
}

func (d *defaultCommandExecutor) customDeploy(providers map[string]domain.CloudProvider, rawTypes []string, cert []byte, key []byte) {
	// convert asset types
	types := make([]domain.AssetType, 0, len(rawTypes))
	for _, str := range rawTypes {
		if t, err := domain.AssetTypeFromString(str); err == nil {
			types = append(types, t)
		} else {
			logger.Warnf("ignore invalid asset type '%s'", str)
		}
	}

	uniDeployer := registry.NewUnionDeployer(logger, providers)

	// search assets
	logger.Infof("looking for %d types of applicable assets: %v", len(types), types)
	allAssets := make([]domain.Asseter, 0, 64)
	for _, t := range types {
		if assets, err := uniDeployer.ListApplicableAssets(t, cert); err != nil {
			logger.Warnf("failed to search assets for type '%s': %v", t, err)
		} else {
			logger.Debugf("%d assets for type '%s' were acquired", len(assets), t)
			allAssets = append(allAssets, assets...)
		}
	}
	logger.Infof("a total of %d assets were acquired, deploying...", len(allAssets))

	deployed, _ := uniDeployer.Deploy(allAssets, cert, key)
	logger.Infof("%d/%d assets deployed successfully: %v", len(deployed), len(allAssets),
		domain.MapSlice(deployed, func(s domain.Asseter) string {
			i := s.GetBaseInfo()
			return fmt.Sprintf("%s-%s@%s", i.Type, i.Name, i.Provider)
		}))
}

func (d *defaultCommandExecutor) executeDeployments(providers map[string]domain.CloudProvider, deployments map[string]domain.Deployment, deploymentIds []string) {
	for _, deploymentId := range deploymentIds {
		if d, ex := deployments[deploymentId]; !ex {
			logger.Warnf("deployment '%s' does not exist, ignroe", deploymentId)
			continue
		} else {
			logger.Debugf("execute deployment '%s'...", deploymentId)
			err := deploy.NewDeploymentExecutor(logger, providers).ExecuteDeployment(d)
			if err != nil {
				logger.Warnf("failed to deploy '%s': %v", deploymentId, err)
			} else {
				logger.Infof("deployment ‘%s’ completed", deploymentId)
			}
		}
	}
	logger.Infof("all %d deployment completed", len(deploymentIds))
}

func (d *defaultCommandExecutor) registerTriggers(providers map[string]domain.CloudProvider,
	deployments map[string]domain.Deployment, triggerDefs map[string]domain.TriggerDefiner) (registeredTriggers []domain.Trigger) {

	registeredTriggers = make([]domain.Trigger, 0)
	for name, triggerDef := range triggerDefs {
		// retrieve triggered deployments
		triggeredDeployments := make([]domain.Deployment, 0, len(triggerDef.GetDeploymentIds()))
		for _, deploymentId := range triggerDef.GetDeploymentIds() {
			if d, ok := deployments[deploymentId]; ok {
				triggeredDeployments = append(triggeredDeployments, d)
			}
		}
		if len(triggeredDeployments) == 0 {
			logger.Infof("no valid deployments in trigger '%s', ignore register", name)
			continue
		}

		l := logger.With("trigger", name)
		executor := deploy.NewDeploymentExecutor(l, providers)

		var trigger domain.Trigger
		switch triggerDef.GetType() {
		case "file_monitoring":
			fileTriggerDef := triggerDef.(*domain.FileMonitoringTriggerDef)
			trigger = filetrigger.NewFileTrigger(l, name, executor, fileTriggerDef.Options, triggeredDeployments)
		default:
			logger.Warnf("invalid type '%s' of trigger '%s', ignore register", triggerDef.GetType(), name)
			continue
		}
		if err := trigger.StartMonitoring(); err != nil {
			logger.Warnf("failed to start trigger '%s': %v", name, err)
			continue
		}
		logger.Debugf("trigger '%s' started successfully", name)
		registeredTriggers = append(registeredTriggers, trigger)
	}
	logger.Infof("%d/%d triggers started successfully", len(registeredTriggers), len(triggerDefs))
	return
}
