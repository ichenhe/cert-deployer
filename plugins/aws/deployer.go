package aws

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/registry"
	"go.uber.org/zap"
)

func init() {
	registry.MustRegister(Provider, func(options domain.Options) (s domain.Deployer, err error) {
		defer domain.RecoverFromInvalidOptionError(func(e *domain.InvalidOptionError) {
			err = e
			s = nil
		})

		secretId := options.MustReadString("secretId")
		secretKey := options.MustReadString("secretKey")
		logger := options.MustReadLogger()
		if logger == nil {
			logger = zap.NewNop().Sugar()
		}
		return newAwsDeployer(secretId, secretKey, logger)
	})
}

var _ domain.Deployer = &deployer{}

func newAwsDeployer(secretId string, secretKey string, logger *zap.SugaredLogger) (*deployer, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(secretId, secretKey, "")))
	if err != nil {
		return nil, err
	}
	return &deployer{
		cfg:    cfg,
		logger: logger,
	}, nil
}

type deployer struct {
	cfg    aws.Config
	logger *zap.SugaredLogger
}

func (d *deployer) ListAssets(ctx context.Context, assetType string) ([]domain.Asseter, error) {
	switch assetType {
	case CloudFront:
		return d.listCloudFrontAssets(ctx, nil)
	}
	return nil, nil
}

func (d *deployer) ListApplicableAssets(ctx context.Context, assetType string, cert []byte) ([]domain.Asseter, error) {
	certBundle, err := domain.NewCertificateBundle(cert)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert: %w", err)
	}

	switch assetType {
	case CloudFront:
		return d.listCloudFrontAssets(ctx, certBundle)
	}
	return nil, nil
}

func (d *deployer) Deploy(ctx context.Context, assets []domain.Asseter, cert []byte, key []byte, callback *domain.DeployCallback) error {
	certBundle, err := domain.NewCertificateBundle(cert)
	if err != nil {
		return err
	}

	onDeployResult := func(asset domain.Asseter, err error) {
		if callback != nil && callback.ResultCallback != nil {
			callback.ResultCallback(asset, err)
		}
	}

	acmClient := acm.NewFromConfig(d.cfg)
	acmCertFinder := newCachedAcmCertFinder(acmClient)
	cloudfrontClient := cloudfront.NewFromConfig(d.cfg)

	for _, asset := range assets {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if callback != nil && callback.PreExecuteCallback != nil {
			callback.PreExecuteCallback(asset)
		}

		info := asset.GetBaseInfo()
		if info.Provider != Provider {
			onDeployResult(asset, errors.New("not a AWS asset"))
			continue
		}
		if !info.Available {
			onDeployResult(asset, errors.New("asset not available"))
			continue
		}

		switch info.Type {
		case CloudFront:
			if cfAsset, ok := asset.(*cloudFrontDistribution); !ok {
				onDeployResult(asset, errors.New("can not convert asset to CloudFrontDistribution"))
			} else {
				err := d.deployCloudFrontCert(ctx, cloudfrontClient, acmCertFinder, cfAsset, certBundle, key)
				onDeployResult(asset, err)
			}
		}
	}

	return nil
}
