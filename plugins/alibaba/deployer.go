package alibaba

import (
	"context"
	"errors"
	"fmt"
	cdn "github.com/alibabacloud-go/cdn-20180510/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/registry"
	"go.uber.org/zap"
)

const Provider = "AlibabaCloud"

var supportedAssetTypes = map[string]struct{}{
	CDN: {},
}

const (
	CDN = "cdn"
)

func init() {
	registry.MustRegister(Provider, func(options domain.Options) (d domain.Deployer, err error) {
		defer domain.RecoverFromInvalidOptionError(func(e *domain.InvalidOptionError) {
			err = e
			d = nil
		})

		secretId := options.MustReadString("secretId")
		secretKey := options.MustReadString("secretKey")
		logger := options.MustReadLogger()
		if logger == nil {
			logger = zap.NewNop().Sugar()
		}
		return newAliDeployer(secretId, secretKey, logger), nil
	})
}

var _ domain.Deployer = &deployer{}

func newAliDeployer(secretId string, secretKey string, logger *zap.SugaredLogger) *deployer {
	return &deployer{secretId: secretId, secretKey: secretKey, logger: logger}
}

type deployer struct {
	secretId  string
	secretKey string
	logger    *zap.SugaredLogger

	domain.DeployerHelper
}

func (d *deployer) IsAssetTypeSupported(assetType string) bool {
	_, ok := supportedAssetTypes[assetType]
	return ok
}

type listAssetsResult struct {
	assets []domain.Asseter
	err    error
}

func (d *deployer) ListAssets(ctx context.Context, assetType string) ([]domain.Asseter, error) {
	result := make(chan listAssetsResult)

	go func() {
		assets, err := d.doListAssets(assetType, nil)
		result <- listAssetsResult{assets, err}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("coroutine cancelled: %w", ctx.Err())
		case result := <-result:
			return result.assets, result.err
		}
	}
}

func (d *deployer) ListApplicableAssets(ctx context.Context, assetType string, cert []byte) ([]domain.Asseter, error) {
	result := make(chan listAssetsResult)

	go func() {
		assets, err := d.doListAssets(assetType, cert)
		result <- listAssetsResult{assets, err}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("coroutine cancelled: %w", ctx.Err())
		case result := <-result:
			return result.assets, result.err
		}
	}
}

func (d *deployer) doListAssets(assetType string, cert []byte) ([]domain.Asseter, error) {
	var certBundle domain.CertificateBundle
	if cert != nil {
		if bundle, err := domain.NewCertificateBundle(cert); err != nil {
			return nil, fmt.Errorf("failed to parse cert: %w", err)
		} else {
			certBundle = bundle
		}
	}

	config := &openapi.Config{
		AccessKeyId:     tea.String(d.secretId),
		AccessKeySecret: tea.String(d.secretKey),
	}
	switch assetType {
	case CDN:
		if api, err := cdn.NewClient(config); err != nil {
			return nil, fmt.Errorf("failed to create cdn client: %w", err)
		} else {
			return d.listCdnAssets(api, certBundle)
		}
	}
	return nil, fmt.Errorf("unsupported asset type: %s", assetType)
}

func (d *deployer) Deploy(ctx context.Context, assets []domain.Asseter, cert []byte, key []byte, callback *domain.DeployCallback) error {
	certBundle, err := domain.NewCertificateBundle(cert)
	if err != nil {
		return fmt.Errorf("failed to parse cert: %w", err)
	}
	cdnApi, err := cdn.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(d.secretId),
		AccessKeySecret: tea.String(d.secretKey),
	})
	if err != nil {
		return fmt.Errorf("failed to create cdn client: %w", err)
	}

	for _, asset := range assets {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		d.OnPreDeploy(callback, asset)

		info := asset.GetBaseInfo()
		if info.Provider != Provider {
			d.OnDeployResult(callback, asset, errors.New("not a AlibabaCloud asset"))
			continue
		}
		if !info.Available {
			d.OnDeployResult(callback, asset, errors.New("asset not available"))
			continue
		}

		switch info.Type {
		case CDN:
			if cdnAsset, ok := asset.(*cdnAsset); !ok {
				d.OnDeployResult(callback, asset, errors.New("can not convert asset to CdnAsset"))
			} else {
				err := d.deployToCdn(cdnApi, cdnAsset, certBundle, key)
				d.OnDeployResult(callback, asset, err)
			}
		}
	}
	return nil
}
