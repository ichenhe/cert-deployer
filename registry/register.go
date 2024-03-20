package registry

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
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
