package tencent

import (
	"context"
	"errors"
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/registry"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"go.uber.org/zap"
)

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

func (d *deployer) deployCdnCert(asset *CdnAsset, cert []byte, key []byte) error {
	client, err := cdn.NewClient(d.newCredential(), "", profile.NewClientProfile())
	if err != nil {
		return err
	}
	// query original config
	queryReq := cdn.NewDescribeDomainsConfigRequest()
	queryReq.Filters = []*cdn.DomainFilter{{
		Name:  common.StringPtr("domain"),
		Value: []*string{common.StringPtr(asset.Domain)},
	}}
	var queryResp *cdn.DescribeDomainsConfigResponse
	queryResp, err = client.DescribeDomainsConfig(queryReq)
	if err != nil {
		return fmt.Errorf("failed to query domain config: %w", err)
	} else if len(queryResp.Response.Domains) != 1 {
		return fmt.Errorf("failed to query domain config: expect 1 result, actual is %d",
			len(queryResp.Response.Domains))
	}
	// update https config
	httpsConfig := queryResp.Response.Domains[0].Https
	httpsConfig.CertInfo = &cdn.ServerCert{
		Certificate: common.StringPtr(string(cert)),
		PrivateKey:  common.StringPtr(string(key)),
		Message:     common.StringPtr("deployed by cert-deployer"),
	}
	req := cdn.NewUpdateDomainConfigRequest()
	req.Domain = common.StringPtr(asset.Domain)
	req.Https = httpsConfig
	_, err = client.UpdateDomainConfig(req)
	return err
}
