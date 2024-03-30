package main

import (
	"context"
	"github.com/ichenhe/cert-deployer/deploy"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/trigger/filetrigger"
)

// commandExecutor executes intent operation with given arguments. Normally it should be called from
// a commandDispatcher.
type commandExecutor interface {
	// executeDeployments executes deployments who have specific id defined in profile.
	executeDeployments(providers map[string]domain.CloudProvider, deployments map[string]domain.Deployment, deploymentIds []string)

	// registerTriggers finds the deployment and starts the trigger if it has at least one valid deployment.
	// Returns started triggers.
	registerTriggers(providers map[string]domain.CloudProvider, deployments map[string]domain.Deployment, triggerDefs map[string]domain.TriggerDefiner) (registeredTriggers []domain.Trigger)
}

type defaultCommandExecutor struct {
}

func newCommandExecutor() commandExecutor {
	return &defaultCommandExecutor{}
}

func (d *defaultCommandExecutor) executeDeployments(providers map[string]domain.CloudProvider, deployments map[string]domain.Deployment, deploymentIds []string) {
	for _, deploymentId := range deploymentIds {
		if d, ex := deployments[deploymentId]; !ex {
			logger.Warnf("deployment '%s' does not exist, ignroe", deploymentId)
			continue
		} else {
			logger.Debugf("execute deployment '%s'...", deploymentId)
			err := deploy.NewDeploymentExecutor(logger.With("deployment", deploymentId), providers).ExecuteDeployment(context.TODO(), d)
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
