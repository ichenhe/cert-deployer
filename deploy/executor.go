package deploy

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/registry"
	"go.uber.org/zap"
	"os"
	"sync"
)

type defaultDeploymentExecutor struct {
	logger    *zap.SugaredLogger
	providers map[string]domain.CloudProvider

	// --- for testing only

	fileReader               domain.FileReader
	deployerFactory          domain.DeployerFactory
	deployerCommanderFactory func(deployer domain.Deployer) deployerCommander

	// --- runtime

	mu        sync.Mutex                   // protect commander
	commander map[string]deployerCommander // providerId -> deployerCommander
}

func NewDeploymentExecutor(logger *zap.SugaredLogger, providers map[string]domain.CloudProvider) domain.DeploymentExecutor {
	return &defaultDeploymentExecutor{
		logger:    logger,
		providers: providers,

		fileReader:               domain.FileReaderFunc(os.ReadFile),
		deployerFactory:          registry.NewDeployerFactory(),
		deployerCommanderFactory: func(deployer domain.Deployer) deployerCommander { return newCachedDeployerCommander(deployer) },

		commander: make(map[string]deployerCommander),
	}
}

// getCommander finds and returns deployerCommander for the given provider id. Creates a new one if
// it does not exist.
func (n *defaultDeploymentExecutor) getCommander(providerId string) (deployerCommander, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if commander, ok := n.commander[providerId]; ok {
		return commander, nil
	}

	// create new commander
	if provider, ex := n.providers[providerId]; !ex {
		return nil, fmt.Errorf("provider '%s' does not exist", providerId)
	} else if deployer, err := n.deployerFactory.NewDeployer(n.logger, provider); err != nil {
		return nil, fmt.Errorf("failed to create deployer: %w", err)
	} else {
		commander := n.deployerCommanderFactory(deployer)
		n.commander[providerId] = commander
		return commander, nil
	}
}

func (n *defaultDeploymentExecutor) ExecuteDeployment(deployment domain.Deployment) error {
	certData, err := n.fileReader.ReadFile(deployment.Cert)
	if err != nil {
		return fmt.Errorf("failed to read public cert: %w", err)
	}
	keyData, err := n.fileReader.ReadFile(deployment.Key)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}

	commander, err := n.getCommander(deployment.ProviderId)
	if err != nil {
		return err
	}
	for _, asset := range deployment.Assets {
		assetType, err := domain.AssetTypeFromString(asset.Type)
		if err != nil {
			n.logger.Warnf("invalid asset type '%s', ignore", asset.Type)
			continue
		}

		if asset.Id != "" {
			n.logger.Debugf("deploying to %s asset '%s'...", asset.Type, asset.Id)
			if err = commander.DeployToAsset(assetType, asset.Id, certData, keyData); err != nil {
				n.logger.Warnf("failed to deploy to %s asset '%s': %v", asset.Type, asset.Id, err)
				// continue to deploy other assets, won't be considered as a fail.
			} else {
				n.logger.Infof("%s asset '%s' deployed successfully", asset.Type, asset.Id)
			}
		} else {
			n.logger.Debugf("deploying to all %s assets...", asset.Type)

			onAssetsAcquired := func(assets []domain.Asseter) {
				n.logger.Infof("a total of %d assets were acquired, deploying...", len(assets))
			}
			onDeployResult := func(asset domain.Asseter, err error) {
				info := asset.GetBaseInfo()
				if err == nil {
					n.logger.Infof("%s asset '%s' deployed successfully", info.Type, info.Id)
				} else {
					n.logger.Warnf("failed to deploy to %s asset '%s': %v", info.Type, info.Id, err)
				}
			}
			if err = commander.DeployToAssetType(assetType, certData, keyData, onAssetsAcquired, onDeployResult); err != nil {
				n.logger.Warnf("failed to deploy to all %s assets: %v", assetType, err)
				// continue to deploy other assets, won't be considered as a fail.
			}
		}
	}
	return nil
}
