package deploy

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/asset"
)

type AssetDeployerConstructor = func(options Options) (Deployer, error)

var assetDeployerConstructors = make(map[string]AssetDeployerConstructor)

type Deployer interface {
	// ListAssets fetches all assets that match the given type.
	// The assetType should be one of constants 'asset.Type*', e.g. asset.TypeCdn.
	ListAssets(assetType string) ([]asset.Asseter, error)

	// ListApplicableAssets fetch all assets that match the given type and cert.
	// The assetType should be one of constants 'asset.Type*', e.g. asset.TypeCdn.
	ListApplicableAssets(assetType string, cert []byte) ([]asset.Asseter, error)

	// Deploy the given pem cert to the all assets.
	//
	// Returns assets that were successfully deployed and errors. Please note that there is no
	// guarantee that len(deployedAsseters)+len(deployErrs)=len(assets), because some minor
	// problems do not count as errors, such as provider mismatch.
	Deploy(assets []asset.Asseter, cert []byte, key []byte) (deployedAssets []asset.Asseter,
		deployErrs []*DeployError)
}

// MustRegister will mustRegister a deployer constructor corresponding to the provider to the list.
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
