package tencent

import (
	"github.com/ichenhe/cert-deployer/domain"
)

const Provider = "TencentCloud"

type CdnAsset struct {
	domain.Asset
	Domain string
}
