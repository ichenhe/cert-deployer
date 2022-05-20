package tencent

import "github.com/ichenhe/cert-deployer/asset"

const Provider = "TencentCloud"

type CdnAsset struct {
	asset.Asset
	Domain string
}
