package deploy

import (
	"context"
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"sync"
)

// deployerCommander manipulates domain.Deployer to deploy given cert to the target assets.
// It retrieves assets and submit them to the domain.Deployer.
//
// Typically, the implementation should contain field to save the domain.Deployer.
type deployerCommander interface {
	IsAssetTypeSupported(assetType string) bool

	// DeployToAsset deploys the certificate to a specific asset using the inner deployer.
	DeployToAsset(ctx context.Context, assetType string, assetId string, cert []byte, key []byte) error

	// DeployToAssetType deploys the certificate to all assets with given type.
	DeployToAssetType(ctx context.Context, assetType string, cert, key []byte, onAssetsAcquired func(assets []domain.Asseter), onDeployResult func(asset domain.Asseter, err error)) error
}

var _ deployerCommander = &cachedDeployerCommander{}

// cachedDeployerCommander caches assets retrieved by the domain.Deployer and submit them directly
// during the asset-id-based deployment.
type cachedDeployerCommander struct {
	deployer domain.Deployer

	cachedAssets map[string]domain.Asseter
	cachedTypes  map[string][]string // assetType -> assetIds
	mu           sync.Mutex
}

func newCachedDeployerCommander(deployer domain.Deployer) *cachedDeployerCommander {
	return &cachedDeployerCommander{
		deployer:     deployer,
		cachedAssets: make(map[string]domain.Asseter),
		cachedTypes:  make(map[string][]string),
	}
}

func (c *cachedDeployerCommander) IsAssetTypeSupported(assetType string) bool {
	return c.deployer.IsAssetTypeSupported(assetType)
}

func (c *cachedDeployerCommander) refreshCache(ctx context.Context, assetType string) error {
	assets, err := c.deployer.ListAssets(ctx, assetType)
	if err != nil {
		return fmt.Errorf("failed to list assests: %w", err)
	}

	for _, id := range c.cachedTypes[assetType] {
		delete(c.cachedAssets, id)
	}
	delete(c.cachedTypes, assetType)
	c.cachedTypes[assetType] = make([]string, 0, len(assets))
	for _, asset := range assets {
		id := asset.GetBaseInfo().Id
		c.cachedAssets[id] = asset
		c.cachedTypes[assetType] = append(c.cachedTypes[assetType], id)
	}
	return nil
}

func (c *cachedDeployerCommander) addToCache(assetType string, assets []domain.Asseter) {
	for _, asset := range assets {
		id := asset.GetBaseInfo().Id
		c.cachedAssets[id] = asset
		c.cachedTypes[assetType] = append(c.cachedTypes[assetType], id)
	}
}

// DeployToAssetType deploys the cert to all assets with given type.
// This function does not use the cache but updates the cache with assets it acquired.
func (c *cachedDeployerCommander) DeployToAssetType(ctx context.Context, assetType string, cert, key []byte,
	onAssetsAcquired func(assets []domain.Asseter),
	onDeployResult func(asset domain.Asseter, err error)) error {

	assets, err := c.deployer.ListApplicableAssets(ctx, assetType, cert)
	if err != nil {
		return fmt.Errorf("failed to list assests: %w", err)
	}
	c.mu.Lock()
	c.addToCache(assetType, assets)
	c.mu.Unlock()

	if onAssetsAcquired != nil {
		onAssetsAcquired(assets)
	}
	for _, asset := range assets {
		_, errors := c.deployer.Deploy(ctx, []domain.Asseter{asset}, cert, key)
		var err error
		if errors != nil && len(errors) > 0 {
			err = errors[0]
		}
		if onDeployResult != nil {
			onDeployResult(asset, err)
		}
	}
	return nil
}

func (c *cachedDeployerCommander) retrieveAsset(ctx context.Context, assetType string, assetId string) (domain.Asseter, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if asset, ex := c.cachedAssets[assetId]; ex {
		return asset, nil
	}
	// asset not found in cache, refresh
	if err := c.refreshCache(ctx, assetType); err != nil {
		return nil, err
	}
	if asset, ex := c.cachedAssets[assetId]; ex {
		return asset, nil
	} else {
		return nil, fmt.Errorf("asset does not exist")
	}
}

// DeployToAsset deploys the cert to asset with given id.
func (c *cachedDeployerCommander) DeployToAsset(ctx context.Context, assetType string, assetId string, cert []byte, key []byte) error {
	asset, err := c.retrieveAsset(ctx, assetType, assetId)
	if err != nil {
		return err
	}

	_, errors := c.deployer.Deploy(ctx, []domain.Asseter{asset}, cert, key)
	if errors != nil && len(errors) > 0 {
		return errors[0]
	}
	return nil
}
