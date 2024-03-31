package alibaba

import (
	"fmt"
	cdn "github.com/alibabacloud-go/cdn-20180510/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/ichenhe/cert-deployer/domain"
	"strings"
)

// The id indicates domain name.
type cdnAsset struct {
	domain.Asset
}

var _ cdnApi = &cdn.Client{}

// cdnApi abstracts the functions used in client that cdn.NewClient() returned for testing purposes.
//
// All functions in this interface must be the same as in cdn.Client to make sure
// the original client is a valid implementation.
type cdnApi interface {
	DescribeUserDomains(request *cdn.DescribeUserDomainsRequest) (*cdn.DescribeUserDomainsResponse, error)

	SetCdnDomainSSLCertificate(request *cdn.SetCdnDomainSSLCertificateRequest) (*cdn.SetCdnDomainSSLCertificateResponse, error)
}

// listCdnAssets lists all cdn assets that match the given certificate, regardless of their status.
// Hostname checking is ignored if the certBundle is nil.
func (d *deployer) listCdnAssets(cdnApi cdnApi, certBundle domain.CertificateBundle) ([]domain.Asseter, error) {
	var wildcard = false
	if certBundle != nil {
		// check if the certificate is a wildcard certificate
		for _, name := range certBundle.GetDomains() {
			if strings.Index(name, "*") != -1 {
				wildcard = true
				break
			}
		}
	}

	// create request parameter
	var request *cdn.DescribeUserDomainsRequest
	if certBundle != nil && !wildcard {
		if len(certBundle.GetDomains()) > 50 {
			return nil, fmt.Errorf("too many domains in the certificate")
		}
		request = &cdn.DescribeUserDomainsRequest{
			DomainName:       tea.String(strings.Join(certBundle.GetDomains(), ",")),
			DomainSearchType: tea.String("full_match"),
		}
	} else {
		// no specific cert or contains wildcard, let's list all
		request = &cdn.DescribeUserDomainsRequest{}
	}

	domains, err := cdnApi.DescribeUserDomains(request)
	if err != nil {
		return nil, err
	}

	assets := make([]domain.Asseter, 0, *domains.Body.TotalCount/2)
	for _, item := range domains.Body.Domains.PageData {
		if certBundle != nil && wildcard {
			// perform manual filtering
			if !certBundle.VerifyHostname(*item.DomainName) {
				continue
			}
		}
		available := *item.DomainStatus != "offline" && *item.DomainStatus != "deleting"
		asset := &cdnAsset{
			Asset: domain.Asset{
				Name:      *item.DomainName,
				Id:        *item.DomainName,
				Type:      CDN,
				Provider:  Provider,
				Available: available,
			},
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

// deployCloudFrontCert imports the cert to ACM (AWS Certificate Manager) and update the target
// cloud front distribution to use it.
//
// If a certificate with the same serial number is found in ACM, it will be reused rather than
// imported again, even if the certificate is not managed by cert-deployer.
//
// The previous cert will be deleted from ACM if it is unused anymore and managed by cert-deployer.
func (d *deployer) deployToCdn(cdnApi cdnApi, asset *cdnAsset, cert domain.CertificateBundle, key []byte) error {
	if !cert.VerifyHostname(asset.Id) {
		return fmt.Errorf("certificate does not match the asset")
	}

	config := &cdn.SetCdnDomainSSLCertificateRequest{
		DomainName:  tea.String(asset.Id),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(string(cert.GetRaw())),
		SSLPri:      tea.String(string(key)),
	}

	_, err := cdnApi.SetCdnDomainSSLCertificate(config)
	if err != nil {
		return err
	}
	return nil
}
