package tencent

import (
	"encoding/base64"
	"github.com/ichenhe/cert-deployer/asset"
	"github.com/ichenhe/cert-deployer/search"
	"github.com/ichenhe/cert-deployer/utils"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type Search struct {
	secretId  string
	secretKey string
}

func init() {
	search.MustRegister(Provider, func(options map[string]interface{}) (s search.AssetSearcher, err error) {
		defer func() {
			if e, ok := recover().(error); ok {
				err = e
				s = nil
			}
		}()
		secretId := utils.MustReadStringOption(options, "secretId")
		secretKey := utils.MustReadStringOption(options, "secretKey")
		s = NewTencentSearch(secretId, secretKey)
		return
	})
}

// NewTencentSearch is a constructor that provide basic information for TencentSearch.
func NewTencentSearch(secretId string, secretKey string) *Search {
	return &Search{secretId: secretId, secretKey: secretKey}
}

func (s *Search) newCredential() *common.Credential {
	return common.NewCredential(s.secretId, s.secretKey)
}

func (s *Search) List(assetType string) ([]asset.Asseter, error) {
	switch assetType {
	case asset.TypeCdn:
		return s.listCDNAssets()
	}
	return nil, nil
}

func (s *Search) ListApplicable(assetType string, cert []byte) ([]asset.Asseter, error) {
	switch assetType {
	case asset.TypeCdn:
		return s.listApplicableCDNAssets(cert)
	}
	return nil, nil
}

func (s *Search) listCDNAssets() ([]asset.Asseter, error) {
	client, err := cdn.NewClient(s.newCredential(), "", profile.NewClientProfile())
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

func (s *Search) listApplicableCDNAssets(cert []byte) ([]asset.Asseter, error) {
	client, err := cdn.NewClient(s.newCredential(), "", profile.NewClientProfile())
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

	allCDNs, err := s.listCDNAssets()
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
