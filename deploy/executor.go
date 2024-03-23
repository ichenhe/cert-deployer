package deploy

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/registry"
	"go.uber.org/zap"
	"os"
)

// assetDeployer is responsible for deploying a certificate to an asset using specific deployer.
// assetDeployer itself is not a domain.Deployer.
type assetDeployer interface {
	// deployToAsset deploys the certificate to a specific asset using the given deployer.
	deployToAsset(deployer domain.Deployer, assetType domain.AssetType, assetId string, cert []byte, key []byte) error
}

type defaultAssetDeployer struct {
}

func newAssetDeployer() assetDeployer {
	return &defaultAssetDeployer{}
}

func (n *defaultAssetDeployer) deployToAsset(deployer domain.Deployer, assetType domain.AssetType, assetId string, cert []byte, key []byte) error {
	assets, err := deployer.ListAssets(assetType)
	if err != nil {
		return fmt.Errorf("failed to list assests: %w", err)
	}

	// found the one need to be deployed
	var target domain.Asseter
	for _, asset := range assets {
		if asset.GetBaseInfo().Type == assetType && asset.GetBaseInfo().Id == assetId {
			target = asset
			break
		}
	}

	if target == nil {
		return fmt.Errorf("asset does not exist")
	}

	if !target.GetBaseInfo().Available {
		return fmt.Errorf("asset unavailable")
	}
	_, errors := deployer.Deploy([]domain.Asseter{target}, cert, key)
	if errors != nil && len(errors) > 0 {
		return errors[0]
	}
	return nil
}

type defaultDeploymentExecutor struct {
	logger          *zap.SugaredLogger
	fileReader      domain.FileReader
	deployerFactory domain.DeployerFactory
	assetDeployer   assetDeployer
}

func NewDeploymentExecutor(logger *zap.SugaredLogger) domain.DeploymentExecutor {
	return NewCustomDeploymentExecutor(logger, domain.FileReaderFunc(os.ReadFile), registry.NewDeployerFactory(), newAssetDeployer())
}

func NewCustomDeploymentExecutor(logger *zap.SugaredLogger, fileReader domain.FileReader, deployerFactory domain.DeployerFactory, assetDeployer assetDeployer) domain.DeploymentExecutor {
	return &defaultDeploymentExecutor{logger: logger, fileReader: fileReader, deployerFactory: deployerFactory, assetDeployer: assetDeployer}
}

func (n *defaultDeploymentExecutor) ExecuteDeployment(providers map[string]domain.CloudProvider, deployment domain.Deployment) error {
	certData, err := n.fileReader.ReadFile(deployment.Cert)
	if err != nil {
		return fmt.Errorf("failed to read public cert: %w", err)
	}
	keyData, err := n.fileReader.ReadFile(deployment.Key)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}

	// find provider configuration
	provider, ex := providers[deployment.ProviderId]
	if !ex {
		return fmt.Errorf("provider '%s' does not exist", deployment.ProviderId)
	}

	// create deployer
	deployer, err := n.deployerFactory.NewDeployer(n.logger, provider)
	if err != nil {
		return fmt.Errorf(" failed to create deployer: %w", err)
	}

	for _, asset := range deployment.Assets {
		assetType, err := domain.AssetTypeFromString(asset.Type)
		if err != nil {
			n.logger.Warnf("invalid asset type '%s'", asset.Type)
			continue
		}

		n.logger.Debugf("deploying to %s asset '%s'...", asset.Type, asset.Id)
		err = n.assetDeployer.deployToAsset(deployer, assetType, asset.Id, certData, keyData)
		if err != nil {
			n.logger.Warnf("failed to deploy to %s asset '%s': %v", asset.Type, asset.Id, err)
			// continue to deploy other assets, won't be considered as a fail.
		}
	}
	return nil
}
