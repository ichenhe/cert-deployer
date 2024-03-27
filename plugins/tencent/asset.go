package tencent

import (
	"github.com/ichenhe/cert-deployer/domain"
)

const Provider = "TencentCloud"

const (
	CDN = "cdn"
)

type CdnAsset struct {
	domain.Asset
	Domain string
}
