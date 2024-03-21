package registry

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"go.uber.org/zap"
)

type AssetDeployerConstructor = func(options domain.Options) (domain.Deployer, error)

var assetDeployerConstructors = make(map[string]AssetDeployerConstructor)

// MustRegister registers a deployer constructor corresponding to the provider to the list.
// The provider should be a globally unique human-readable identifier which will be used as
// Asset.Provider and the value of 'cloud-providers[i].provider' in profile. e.g. TencentCloud
//
// All Deployer must call this function in init function to mustRegister itself.
func MustRegister(provider string, deployerConstructor AssetDeployerConstructor) {
	if _, ex := assetDeployerConstructors[provider]; ex {
		panic(fmt.Errorf("[AssetDeployer] provider '%s' is already registered", provider))
	} else {
		assetDeployerConstructors[provider] = deployerConstructor
	}
}

type defaultDeployerFactory struct {
}

func NewDeployerFactory() domain.DeployerFactory {
	return &defaultDeployerFactory{}
}

// NewDeployer creates a deployer corresponding to the given cloudProvider.
func (f *defaultDeployerFactory) NewDeployer(logger *zap.SugaredLogger, cloudProvider domain.CloudProvider) (domain.Deployer, error) {
	var deployerConstructor AssetDeployerConstructor
	if c, ok := assetDeployerConstructors[cloudProvider.Provider]; !ok {
		return nil, fmt.Errorf("provider '%s' not supported", cloudProvider.Provider)
	} else {
		deployerConstructor = c
	}
	options := map[string]interface{}{
		"secretId":              cloudProvider.SecretId,
		"secretKey":             cloudProvider.SecretKey,
		domain.OptionsKeyLogger: logger,
	}
	if deployer, err := deployerConstructor(options); err != nil {
		return nil, err
	} else {
		return deployer, nil
	}
}
