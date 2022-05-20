package deploy

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/asset"
)

type DeployError struct {
	Asset asset.Asseter
	Err   error
}

func (d *DeployError) Error() string {
	return fmt.Sprintf("failed to deploy %v: %v", d.Asset, d.Err)
}

func NewDeployError(asset asset.Asseter, err error) *DeployError {
	return &DeployError{
		Asset: asset,
		Err:   err,
	}
}
