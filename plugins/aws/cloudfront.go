package aws

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmTypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/ichenhe/cert-deployer/domain"
	"strings"
)

// listCloudFrontAssets lists all cloud front distributions that match the given certificate.
// Hostname checking is ignored if the fullChain is nil.
//
// Distributions without aliases always match no certificates.
func (d *deployer) listCloudFrontAssets(ctx context.Context, certBundle *certificateBundle) ([]domain.Asseter, error) {
	client := cloudfront.NewFromConfig(d.cfg)

	// extractName generates a name for the DistributionSummary
	extractName := func(item *types.DistributionSummary) string {
		if *item.Aliases.Quantity == 1 {
			return item.Aliases.Items[0]
		}
		if *item.Aliases.Quantity > 1 {
			return fmt.Sprintf("[%s]", strings.Join(item.Aliases.Items, ", "))
		}
		return *item.DomainName
	}

	// paging query
	var marker *string
	assets := make([]domain.Asseter, 0)
	for {
		result, err := client.ListDistributions(ctx, &cloudfront.ListDistributionsInput{
			Marker: marker,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range result.DistributionList.Items {
			if certBundle != nil {
				if *item.Aliases.Quantity == 0 {
					continue // only default domain name provided by aws
				}
				if !certBundle.VerifyHostnames(item.Aliases.Items) {
					continue // cert not match
				}
			}
			asset := cloudFrontDistribution{
				Asset: domain.Asset{
					Id:        *item.Id,
					Name:      extractName(&item),
					Type:      CloudFront,
					Provider:  Provider,
					Available: true,
				},
			}
			assets = append(assets, &asset)
		}

		if result.DistributionList.IsTruncated != nil && *result.DistributionList.IsTruncated {
			marker = result.DistributionList.NextMarker
		} else {
			break
		}
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
func (d *deployer) deployCloudFrontCert(ctx context.Context, certFinder acmCertFinder, asset *cloudFrontDistribution, cert *certificateBundle, key []byte) error {
	// get current cloud front config
	cfClient := cloudfront.NewFromConfig(d.cfg)
	result, err := cfClient.GetDistributionConfig(ctx, &cloudfront.GetDistributionConfigInput{
		Id: &asset.Id,
	})
	if err != nil {
		return fmt.Errorf("failed to get distribution config: %w", err)
	}
	// verify domain name matching
	// WARN: AWS does not verify the matching of certificate!
	// Deploying the wrong certificate can cause client access to fail.
	if result.DistributionConfig.Aliases == nil {
		return errors.New("cert not match")
	} else {
		for _, item := range result.DistributionConfig.Aliases.Items {
			if !cert.VerifyHostname(item) {
				return fmt.Errorf("cert dose not match %s", item)
			}
		}
	}

	// try to find cert in ACM
	acmClient := certFinder.GetAcmClient()
	var certARN *string
	var newImportedCert = false // whether the cert deployed this time is the new one
	if arn, err := certFinder.FindCertInACM(ctx, cert); err != nil {
		return fmt.Errorf("failed to find cert from ACM: %w", err)
	} else if arn != "" {
		certARN = &arn
	}

	if certARN == nil {
		// cert does not exist in ACM, import it
		tagKey := aws.String(acmManagedTagKey)
		tagValue := aws.String(acmManagedTagValue)
		if result, err := acmClient.ImportCertificate(ctx, &acm.ImportCertificateInput{
			Certificate:      cert.ClientCertRaw(),
			PrivateKey:       key,
			CertificateChain: cert.ChainRaw,
			Tags:             []acmTypes.Tag{{Key: tagKey, Value: tagValue}},
		}); err != nil {
			return fmt.Errorf("failed to import cert to ACM: %w", err)
		} else {
			newImportedCert = true
			certARN = result.CertificateArn
			certFinder.NotifyCertAdded(cert, *result.CertificateArn)
		}
	}

	if certARN == nil {
		return fmt.Errorf("unknown error: cert ARN should not be nil")
	}

	// create new cert config
	var oldCertArn *string
	oldCertConfig := result.DistributionConfig.ViewerCertificate
	newCertConfig := &types.ViewerCertificate{}
	result.DistributionConfig.ViewerCertificate = newCertConfig
	if oldCertConfig != nil {
		oldCertArn = oldCertConfig.ACMCertificateArn
		newCertConfig.MinimumProtocolVersion = oldCertConfig.MinimumProtocolVersion
		newCertConfig.SSLSupportMethod = oldCertConfig.SSLSupportMethod
	} else {
		newCertConfig.MinimumProtocolVersion = types.MinimumProtocolVersionTLSv1
		newCertConfig.SSLSupportMethod = types.SSLSupportMethodSniOnly
	}
	newCertConfig.ACMCertificateArn = certARN
	newCertConfig.CloudFrontDefaultCertificate = aws.Bool(false)

	// submit
	if _, err := cfClient.UpdateDistribution(ctx, &cloudfront.UpdateDistributionInput{
		DistributionConfig: result.DistributionConfig,
		Id:                 &asset.Id,
		IfMatch:            result.ETag,
	}); err != nil {
		return fmt.Errorf("failed to update distribution: %w", err)
	}

	// delete old cert from ACM
	if newImportedCert && oldCertArn != nil {
		if deleted, err := d.deleteManagedCertFromAcmIfUnused(ctx, acmClient, oldCertArn); err != nil {
			d.logger.Warnf("failed to delete unused cert '%s' from ACM: %v", *oldCertArn, err)
		} else if deleted {
			certFinder.NotifyCertDeleted(*oldCertArn)
			d.logger.Debugf("deleted unused cert '%s' from ACM", *oldCertArn)
		}
	}
	return nil
}
