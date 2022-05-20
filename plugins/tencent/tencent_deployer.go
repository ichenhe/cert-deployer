package tencent

import (
	"errors"
	"fmt"
	"github.com/ichenhe/cert-deployer/asset"
	"github.com/ichenhe/cert-deployer/deploy"
	"github.com/ichenhe/cert-deployer/utils"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"go.uber.org/zap"
)

type Deploy struct {
	secretId  string
	secretKey string
	logger    *zap.SugaredLogger
}

func init() {
	deploy.MustRegister(Provider, func(options map[string]interface{}) (s deploy.Deployer, err error) {
		defer func() {
			if e, ok := recover().(error); ok {
				err = e
				s = nil
			}
		}()
		secretId := utils.MustReadStringOption(options, "secretId")
		secretKey := utils.MustReadStringOption(options, "secretKey")
		logger := utils.MustReadOption[*zap.SugaredLogger](options, "logger")
		s = NewTencentDeploy(secretId, secretKey, logger)
		return
	})
}

func NewTencentDeploy(secretId string, secretKey string, logger *zap.SugaredLogger) *Deploy {
	return &Deploy{secretId: secretId, secretKey: secretKey, logger: logger}
}

func (d *Deploy) newCredential() *common.Credential {
	return common.NewCredential(d.secretId, d.secretKey)
}

func (d *Deploy) Deploy(assets []asset.Asseter, cert []byte, key []byte) (
	deployedAsseters []asset.Asseter, deployErrs []*deploy.DeployError) {
	for _, item := range assets {
		info := item.GetBaseInfo()
		if info.Provider != Provider {
			d.logger.Warnf("not a tencent asset, ignore: %v", item)
			continue
		}
		if !info.Available {
			d.logger.Warnf("asset not available, ignore: %v", item)
			continue
		}

		switch info.Type {
		case asset.TypeCdn:
			if cdnAsset, ok := item.(*CdnAsset); !ok {
				deployErrs = append(deployErrs, deploy.NewDeployError(item,
					errors.New("can not convert asset to TencentCdnAsset")))
			} else if err := d.deployCdnCert(cdnAsset, cert, key); err != nil {
				deployErrs = append(deployErrs, deploy.NewDeployError(item, err))
			} else {
				deployedAsseters = append(deployedAsseters, item)
			}
		}
	}
	if len(deployErrs) == 0 {
		deployErrs = nil
	}
	return
}

func (d *Deploy) deployCdnCert(asset *CdnAsset, cert []byte, key []byte) error {
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
