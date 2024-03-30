package tencent

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type CdnAsset struct {
	domain.Asset
	Domain string
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

func (d *deployer) listCDNAssets(ctx context.Context) ([]domain.Asseter, error) {
	client, err := cdn.NewClient(d.newCredential(), "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	resp, err := client.DescribeDomains(cdn.NewDescribeDomainsRequest())
	if err != nil {
		return nil, err
	}
	assets := make([]domain.Asseter, 0, *resp.Response.TotalNumber)
	for _, domainName := range resp.Response.Domains {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		assets = append(assets, &CdnAsset{
			Asset: domain.Asset{
				Id:       *domainName.ResourceId,
				Name:     *domainName.Domain,
				Type:     CDN,
				Provider: Provider,
				Available: *domainName.Disable == "normal" &&
					(*domainName.Status == "online" || *domainName.Status == "processing"),
			},
			Domain: *domainName.Domain,
		})
	}
	return assets, nil
}

func (d *deployer) listApplicableCDNAssets(ctx context.Context, cert []byte) ([]domain.Asseter, error) {
	client, err := cdn.NewClient(d.newCredential(), "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	req := cdn.NewDescribeCertDomainsRequest()
	req.Cert = common.StringPtr(base64.StdEncoding.EncodeToString(cert))
	req.Product = common.StringPtr(CDN)
	resp, err := client.DescribeCertDomains(req)
	if err != nil {
		return nil, err
	}

	allCDNs, err := d.listCDNAssets(ctx)
	if err != nil {
		return nil, err
	}

	domainSets := make(map[string]struct{})
	for _, domainNames := range resp.Response.Domains {
		domainSets[*domainNames] = struct{}{}
	}
	result := make([]domain.Asseter, 0)
	for _, cdnItem := range allCDNs {
		if _, ex := domainSets[(cdnItem.(*CdnAsset)).Domain]; ex {
			result = append(result, cdnItem)
		}
	}
	return result, nil
}
