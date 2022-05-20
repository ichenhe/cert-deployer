package search

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/asset"
)

type AssetSearcherConstructor = func(options map[string]interface{}) (AssetSearcher, error)

var assetSearcherConstructors = make(map[string]AssetSearcherConstructor)

// AssetSearcher is used to list all assets belong to your account.
type AssetSearcher interface {
	// List fetches all assets that match the given type.
	// The assetType should be one of constants 'asset.Type*', e.g. asset.TypeCdn.
	List(assetType string) ([]asset.Asseter, error)

	// ListApplicable fetch all assets that match the given type and cert.
	// The assetType should be one of constants 'asset.Type*', e.g. asset.TypeCdn.
	ListApplicable(assetType string, cert []byte) ([]asset.Asseter, error)
}

// MustRegister will mustRegister a searcher constructor corresponding to the provider to the list.
// The provider should be a globally unique human-readable identifier which will be used as
// Asset.Provider and the value of 'cloud-providers[i].provider' in profile. e.g. TencentCloud
//
// All AssetSearcher must call this function in init function to mustRegister itself.
func MustRegister(provider string, searcher AssetSearcherConstructor) {
	if _, ex := assetSearcherConstructors[provider]; ex {
		panic(fmt.Errorf("[AsseSearcher] provider '%s' is already registered", provider))
	} else {
		assetSearcherConstructors[provider] = searcher
	}
}
