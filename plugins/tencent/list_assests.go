package tencent

import (
	"encoding/base64"
	"github.com/ichenhe/cert-deployer/asset"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func (d *deployer) ListAssets(assetType string) ([]asset.Asseter, error) {
	switch assetType {
	case asset.TypeCdn:
		return d.listCDNAssets()
	}
	return nil, nil
}

func (d *deployer) ListApplicableAssets(assetType string, cert []byte) ([]asset.Asseter, error) {
	switch assetType {
	case asset.TypeCdn:
		return d.listApplicableCDNAssets(cert)
	}
	return nil, nil
}

func (d *deployer) listCDNAssets() ([]asset.Asseter, error) {
	client, err := cdn.NewClient(d.newCredential(), "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	resp, err := client.DescribeDomains(cdn.NewDescribeDomainsRequest())
	if err != nil {
		return nil, err
	}
	assets := make([]asset.Asseter, 0, *resp.Response.TotalNumber)
	for _, domain := range resp.Response.Domains {
		assets = append(assets, &CdnAsset{
			Asset: asset.Asset{
				Id:       *domain.ResourceId,
				Name:     *domain.Domain,
				Type:     asset.TypeCdn,
				Provider: Provider,
				Available: *domain.Disable == "normal" &&
					(*domain.Status == "online" || *domain.Status == "processing"),
			},
			Domain: *domain.Domain,
		})
	}
	return assets, nil
}

func (d *deployer) listApplicableCDNAssets(cert []byte) ([]asset.Asseter, error) {
	client, err := cdn.NewClient(d.newCredential(), "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	req := cdn.NewDescribeCertDomainsRequest()
	req.Cert = common.StringPtr(base64.StdEncoding.EncodeToString(cert))
	req.Product = common.StringPtr(asset.TypeCdn)
	resp, err := client.DescribeCertDomains(req)
	if err != nil {
		return nil, err
	}

	allCDNs, err := d.listCDNAssets()
	if err != nil {
		return nil, err
	}

	domainSets := make(map[string]struct{})
	for _, domain := range resp.Response.Domains {
		domainSets[*domain] = struct{}{}
	}
	result := make([]asset.Asseter, 0)
	for _, cdnItem := range allCDNs {
		if _, ex := domainSets[(cdnItem.(*CdnAsset)).Domain]; ex {
			result = append(result, cdnItem)
		}
	}
	return result, nil
}
