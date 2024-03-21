package main

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/registry"
	"github.com/ichenhe/cert-deployer/utils"
)

// commandExecutor executes intent operation with given arguments. Normally it should be called from
// a commandDispatcher.
type commandExecutor interface {
	// executeDeployments executes deployments who have specific id defined in profile.
	executeDeployments(appConfig *domain.AppConfig, deploymentIds []string)

	// customDeploy reads configurations from command line and processes a one-time deployment.
	customDeploy(providers map[string]domain.CloudProvider, rawTypes []string, cert []byte, key []byte)
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
		utils.MapSlice(deployed, func(s domain.Asseter) string {
			i := s.GetBaseInfo()
			return fmt.Sprintf("%s-%s@%s", i.Type, i.Name, i.Provider)
		}))
}

func (d *defaultCommandExecutor) executeDeployments(appConfig *domain.AppConfig, deploymentIds []string) {
	for _, deploymentId := range deploymentIds {
		if d, ex := appConfig.Deployments[deploymentId]; !ex {
			logger.Warnf("deployment '%s' does not exist, ignroe", deploymentId)
			continue
		} else {
			logger.Debugf("execute deployment '%s'...", deploymentId)
			err := newDeploymentExecutor().executeDeployment(appConfig.CloudProviders, d)
			if err != nil {
				logger.Warnf("failed to deploy '%s': %v", deploymentId, err)
			} else {
				logger.Infof("deployment ‘%s’ completed", deploymentId)
			}
		}
	}
	logger.Infof("all %d deployment completed", len(deploymentIds))
}
