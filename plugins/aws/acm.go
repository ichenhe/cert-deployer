package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmTypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/ichenhe/cert-deployer/domain"
)

var _ acmApi = &acm.Client{}

// acmApi abstracts the functions used in acm.Client for testing purposes.
//
// All functions in this interface must be the same as in acm.Client to make sure acm.Client is a
// valid implementation.
type acmApi interface {
	DescribeCertificate(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error)

	ListCertificates(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error)

	ImportCertificate(ctx context.Context, params *acm.ImportCertificateInput, optFns ...func(*acm.Options)) (*acm.ImportCertificateOutput, error)

	ListTagsForCertificate(ctx context.Context, params *acm.ListTagsForCertificateInput, optFns ...func(*acm.Options)) (*acm.ListTagsForCertificateOutput, error)

	DeleteCertificate(ctx context.Context, params *acm.DeleteCertificateInput, optFns ...func(*acm.Options)) (*acm.DeleteCertificateOutput, error)
}

// acmManager wraps the acm.Client to manage certificates in ACM (AWS Certificate Manager).
type acmManager interface {
	// FindCertInACM finds given cert from ACM (AWS Certificate Manager).
	// Returns certification's ARN (Amazon Resource Name) and the source (whether from cache) if
	// found, otherwise "".
	// Returning "",false,nil means no errors but not found.
	FindCertInACM(ctx context.Context, certBundle domain.CertificateBundle) (arn string, fromCache bool, err error)

	// DeleteManagedCertFromAcmIfUnused deletes a cert from ACM.
	//
	// The deleted certificate must meet ALL conditions:
	//   - Unused.
	//   - A imported cert. i.e. not issued by amazon.
	//   - Managed by the cert-deployer (has acmManagedTagKey tag).
	DeleteManagedCertFromAcmIfUnused(ctx context.Context, certArn *string) (deleted bool, err error)

	ImportCertificate(ctx context.Context, certBundle domain.CertificateBundle, key []byte) (arn string, err error)

	RemoveCertFromCache(arn string)
}

var _ acmManager = &cachedAcmManager{}

func newCachedAcmCertFinder(acmApi acmApi) *cachedAcmManager {
	return &cachedAcmManager{
		api: acmApi,

		cachedCertSummary: make(map[string]acmTypes.CertificateSummary),
		cachedSnToArn:     make(map[string]string),
	}
}

// cachedAcmManager caches the certificate list to speed up the query.
// Usage point:
//
//   - Cannot be used in parallel.
//   - Must be notified immediately after adding/removing a certificate.
//   - Recommended for continuous deployment only - do not cache for a long time.
type cachedAcmManager struct {
	api acmApi

	// runtime
	cachedCertSummary map[string]acmTypes.CertificateSummary // arn -> summary
	cachedSnToArn     map[string]string                      // cert serial number (hex string without colons or prefix) -> ARN
	cachedArnToSn     map[string]string
}

func (f *cachedAcmManager) ImportCertificate(ctx context.Context, certBundle domain.CertificateBundle, key []byte) (arn string, err error) {
	if result, err := f.api.ImportCertificate(ctx, &acm.ImportCertificateInput{
		Certificate:      certBundle.GetRawCert(),
		PrivateKey:       key,
		CertificateChain: certBundle.GetRawChain(),
		Tags:             []acmTypes.Tag{{Key: aws.String(acmManagedTagKey), Value: aws.String(acmManagedTagValue)}},
	}); err != nil {
		return "", err
	} else {
		// update cache
		f.cachedArnToSn[*result.CertificateArn] = certBundle.GetSerialNumberHexString()
		f.cachedSnToArn[certBundle.GetSerialNumberHexString()] = *result.CertificateArn

		return *result.CertificateArn, nil
	}
}

func (f *cachedAcmManager) RemoveCertFromCache(arn string) {
	delete(f.cachedCertSummary, arn)
	if sn, ex := f.cachedArnToSn[arn]; ex {
		delete(f.cachedSnToArn, sn)
		delete(f.cachedArnToSn, arn)
	}
}

func (f *cachedAcmManager) DeleteManagedCertFromAcmIfUnused(ctx context.Context, certArn *string) (deleted bool, err error) {
	if certArn == nil {
		return false, fmt.Errorf("certArn is nil")
	}
	if result, err := f.api.DescribeCertificate(ctx, &acm.DescribeCertificateInput{CertificateArn: certArn}); err != nil {
		return false, fmt.Errorf("failed to describe cert: %w", err)
	} else if len(result.Certificate.InUseBy) > 0 {
		return false, nil // still using
	} else if result.Certificate.Type != acmTypes.CertificateTypeImported {
		return false, nil // not imported by the user
	}

	// verify the cert is managed by this program
	if result, err := f.api.ListTagsForCertificate(ctx, &acm.ListTagsForCertificateInput{CertificateArn: certArn}); err != nil {
		return false, fmt.Errorf("failed to list tags: %w", err)
	} else {
		find := false
		for _, tag := range result.Tags {
			if *tag.Key == acmManagedTagKey {
				find = true
				break
			}
		}
		if !find {
			// not managed by the tool
			return false, nil
		}
	}

	// delete
	_, err = f.api.DeleteCertificate(ctx, &acm.DeleteCertificateInput{CertificateArn: certArn})
	deleted = err == nil
	if deleted {
		// delete from cache
		delete(f.cachedCertSummary, *certArn)
		if sn, ex := f.cachedArnToSn[*certArn]; ex {
			delete(f.cachedSnToArn, sn)
			delete(f.cachedArnToSn, *certArn)
		}
	}
	return
}

func (f *cachedAcmManager) FindCertInACM(ctx context.Context, certBundle domain.CertificateBundle) (arn string, fromCache bool, err error) {
	// check cached result first
	if arn, ex := f.cachedSnToArn[certBundle.GetSerialNumberHexString()]; ex {
		return arn, true, nil
	}

	verifyIsTheSameCert := func(summary *acmTypes.CertificateSummary) (bool, error) {
		// quick check before request for details
		if !summary.NotBefore.Equal(*certBundle.NotBefore()) || !summary.NotAfter.Equal(*certBundle.NotAfter()) {
			return false, nil
		}
		if !certBundle.ContainsAllDomains(summary.SubjectAlternativeNameSummaries) {
			return false, nil
		}

		// need serial number for final check
		if certDetails, err := f.api.DescribeCertificate(ctx, &acm.DescribeCertificateInput{
			CertificateArn: summary.CertificateArn,
		}); err != nil {
			return false, err
		} else if certDetails.Certificate.Serial != nil && certBundle.VerifySerialNumber(*certDetails.Certificate.Serial) {
			return true, nil
		}
		return false, nil
	}

	// find from cached list
	for currentArn, summary := range f.cachedCertSummary {
		if ok, err := verifyIsTheSameCert(&summary); err == nil && ok {
			// update cache
			f.cachedSnToArn[certBundle.GetSerialNumberHexString()] = currentArn

			// it is updated because it has just been verified via DescribeCertificate()
			return currentArn, false, nil
		} else if err != nil {
			// failed to fetch certification details from ACM, the cert could have been deleted
			// remove it from cache
			delete(f.cachedCertSummary, currentArn)
			if sn, ex := f.cachedArnToSn[currentArn]; ex {
				delete(f.cachedSnToArn, sn)
				delete(f.cachedArnToSn, currentArn)
			}
		}
	}

	// not found in cache
	certs, err := f.listCertificates(ctx)
	if err != nil {
		return "", false, err
	}
	// update cached cert list
	f.cachedCertSummary = make(map[string]acmTypes.CertificateSummary, len(certs))
	for _, cert := range certs {
		f.cachedCertSummary[*cert.CertificateArn] = cert
	}
	// find target
	for _, cert := range certs {
		if ok, err := verifyIsTheSameCert(&cert); err == nil && ok {
			// update cache
			f.cachedSnToArn[certBundle.GetSerialNumberHexString()] = *cert.CertificateArn
			return *cert.CertificateArn, false, nil
		}
	}

	return "", false, err
}

// listCertificates lists all certificates with 'issued' status, includes amazon issued.
func (f *cachedAcmManager) listCertificates(ctx context.Context) ([]acmTypes.CertificateSummary, error) {
	certs := make([]acmTypes.CertificateSummary, 0)

	var nextToken *string = nil
	for {
		result, err := f.api.ListCertificates(ctx, &acm.ListCertificatesInput{
			NextToken:           nextToken,
			CertificateStatuses: []acmTypes.CertificateStatus{acmTypes.CertificateStatusIssued},
			Includes: &acmTypes.Filters{
				KeyTypes: acmTypes.KeyAlgorithm.Values(""),
			},
		})
		if err != nil {
			return nil, err
		}

		certs = append(certs, result.CertificateSummaryList...)

		nextToken = result.NextToken
		if nextToken == nil || *nextToken == "" {
			break
		}
	}

	return certs, nil
}
