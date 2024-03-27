package tencent

import (
	"github.com/ichenhe/cert-deployer/domain"
)

const Provider = "TencentCloud"

var supportedAssetTypes = map[string]struct{}{
	CDN: {},
}

const (
	CDN = "cdn"
)

func (d *deployer) IsAssetTypeSupported(assetType string) bool {
	_, ok := supportedAssetTypes[assetType]
	return ok
}

type CdnAsset struct {
	domain.Asset
	Domain string
}
