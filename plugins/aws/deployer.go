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
	certBundle, err := newCertificateBundle(cert)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert: %w", err)
	}

	switch assetType {
	case CloudFront:
		return d.listCloudFrontAssets(ctx, certBundle)
	}
	return nil, nil
}

func (d *deployer) Deploy(ctx context.Context, assets []domain.Asseter, cert []byte, key []byte) (deployedAssets []domain.Asseter, deployErrs []*domain.DeployError) {
	certBundle, err := newCertificateBundle(cert)
	if err != nil {
		return nil, []*domain.DeployError{{Err: err}}
	}

	acmClient := acm.NewFromConfig(d.cfg)
	acmCertFinder := newCachedAcmCertFinder(acmClient)
	cloudfrontClient := cloudfront.NewFromConfig(d.cfg)

	for _, asset := range assets {
		info := asset.GetBaseInfo()
		if info.Provider != Provider {
			d.logger.Debugf("not a AWS asset, ignore: %v", asset)
			continue
		}
		if !info.Available {
			d.logger.Debugf("asset not available, ignore: %v", assets)
			continue
		}

		switch info.Type {
		case CloudFront:
			if cfAsset, ok := asset.(*cloudFrontDistribution); !ok {
				deployErrs = append(deployErrs, domain.NewDeployError(asset,
					errors.New("can not convert asset to CloudFrontDistribution")))
			} else if err := d.deployCloudFrontCert(ctx, cloudfrontClient, acmCertFinder, cfAsset, certBundle, key); err != nil {
				deployErrs = append(deployErrs, domain.NewDeployError(asset, err))
			} else {
				deployedAssets = append(deployedAssets, asset)
			}
		}
	}

	return deployedAssets, nil
}
