package tencent

import (
	"context"
	"errors"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/registry"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"go.uber.org/zap"
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

type deployer struct {
	secretId  string
	secretKey string
	logger    *zap.SugaredLogger
}

func init() {
	registry.MustRegister(Provider, func(options domain.Options) (s domain.Deployer, err error) {
		defer domain.RecoverFromInvalidOptionError(func(e *domain.InvalidOptionError) {
			err = e
			s = nil
		})

		secretId := options.MustReadString("secretId")
		secretKey := options.MustReadString("secretKey")
		logger := options.MustReadLogger()
		s = newTencentDeployer(secretId, secretKey, logger)
		return
	})
}

func newTencentDeployer(secretId string, secretKey string, logger *zap.SugaredLogger) domain.Deployer {
	return &deployer{secretId: secretId, secretKey: secretKey, logger: logger}
}

func (d *deployer) newCredential() *common.Credential {
	return common.NewCredential(d.secretId, d.secretKey)
}

func (d *deployer) Deploy(ctx context.Context, assets []domain.Asseter, cert []byte, key []byte, callback *domain.DeployCallback) error {
	onDeployResult := func(asset domain.Asseter, err error) {
		if callback != nil && callback.ResultCallback != nil {
			callback.ResultCallback(asset, err)
		}
	}

	for _, item := range assets {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if callback != nil && callback.PreExecuteCallback != nil {
			callback.PreExecuteCallback(item)
		}

		info := item.GetBaseInfo()
		if info.Provider != Provider {
			onDeployResult(item, errors.New("not a tencent asset"))
			continue
		}
		if !info.Available {
			onDeployResult(item, errors.New("asset not available"))
			continue
		}

		switch info.Type {
		case CDN:
			if cdnAsset, ok := item.(*CdnAsset); !ok {
				onDeployResult(item, errors.New("can not convert asset to TencentCdnAsset"))
			} else {
				err := d.deployCdnCert(cdnAsset, cert, key)
				onDeployResult(item, err)
			}
		}
	}
	return nil
}

func (d *deployer) ListAssets(ctx context.Context, assetType string) ([]domain.Asseter, error) {
	switch assetType {
	case CDN:
		return d.listCDNAssets(ctx)
	}
	return nil, nil
}

func (d *deployer) ListApplicableAssets(ctx context.Context, assetType string, cert []byte) ([]domain.Asseter, error) {
	switch assetType {
	case CDN:
		return d.listApplicableCDNAssets(ctx, cert)
	}
	return nil, nil
}
