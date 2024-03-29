package aws

import "github.com/ichenhe/cert-deployer/domain"

const Provider = "AWS"

var supportedAssetTypes = map[string]struct{}{
	CloudFront: {},
}

const (
	CloudFront = "cloud_front"
)

// The key of tag indicates that this cert in ACM is managed by cert-deployer.
//
// Changing the content will break downward compatibility.
const acmManagedTagKey = "_cert-deployer.flag"

// Do NOT use this value for any verification: it could change at any time.
// Use the existence of key instead.
const acmManagedTagValue = "DO NOT REMOVE . managed by cert-deployer"

func (d *deployer) IsAssetTypeSupported(assetType string) bool {
	_, ok := supportedAssetTypes[assetType]
	return ok
}

type cloudFrontDistribution struct {
	domain.Asset
}
