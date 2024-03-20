package tencent

import (
	"encoding/base64"
	"github.com/ichenhe/cert-deployer/domain"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func (d *deployer) ListAssets(assetType domain.AssetType) ([]domain.Asseter, error) {
	switch assetType {
	case domain.TypeCdn:
		return d.listCDNAssets()
	}
	return nil, nil
}

func (d *deployer) ListApplicableAssets(assetType domain.AssetType, cert []byte) ([]domain.Asseter, error) {
	switch assetType {
	case domain.TypeCdn:
		return d.listApplicableCDNAssets(cert)
	}
	return nil, nil
}

func (d *deployer) listCDNAssets() ([]domain.Asseter, error) {
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
		assets = append(assets, &CdnAsset{
			Asset: domain.Asset{
				Id:       *domainName.ResourceId,
				Name:     *domainName.Domain,
				Type:     domain.TypeCdn,
				Provider: Provider,
				Available: *domainName.Disable == "normal" &&
					(*domainName.Status == "online" || *domainName.Status == "processing"),
			},
			Domain: *domainName.Domain,
		})
	}
	return assets, nil
}

func (d *deployer) listApplicableCDNAssets(cert []byte) ([]domain.Asseter, error) {
	client, err := cdn.NewClient(d.newCredential(), "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	req := cdn.NewDescribeCertDomainsRequest()
	req.Cert = common.StringPtr(base64.StdEncoding.EncodeToString(cert))
	req.Product = common.StringPtr("cdn")
	resp, err := client.DescribeCertDomains(req)
	if err != nil {
		return nil, err
	}

	allCDNs, err := d.listCDNAssets()
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
