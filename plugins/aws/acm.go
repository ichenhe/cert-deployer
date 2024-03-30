package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmTypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
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

type acmCertFinder interface {
	GetAcmApi() acmApi

	// FindCertInACM finds cert saved in the bundle from ACM (AWS Certificate Manager).
	// Returns ARN (Amazon Resource Name) of the cert if found, otherwise "".
	// Returning "",nil means no errors but not found.
	FindCertInACM(ctx context.Context, certBundle *certificateBundle) (string, error)

	// NotifyCertAdded notifies the finder that a cert has been imported to ACM. Useful for caching.
	NotifyCertAdded(certBundle *certificateBundle, arn string)

	// NotifyCertDeleted notifies the finder that a cert has been deleted from ACM. Useful for caching.
	NotifyCertDeleted(arn string)
}

var _ acmCertFinder = &cachedAcmCertFinder{}

func newCachedAcmCertFinder(acmApi acmApi) *cachedAcmCertFinder {
	return &cachedAcmCertFinder{
		api: acmApi,

		cachedCertSummary: make(map[string]acmTypes.CertificateSummary),
		cachedCerts:       make(map[string]string),
	}
}

// Usage point:
//
//   - Cannot be used in parallel.
//   - Must be notified immediately after adding/removing a certificate.
//   - Recommended for continuous deployment only - do not cache for a long time.
type cachedAcmCertFinder struct {
	api acmApi

	// runtime
	cachedCertSummary map[string]acmTypes.CertificateSummary // arn -> summary
	cachedCerts       map[string]string                      // cert serial number (hex string without colons or prefix) -> ARN
}

func (f *cachedAcmCertFinder) NotifyCertAdded(certBundle *certificateBundle, arn string) {
	f.cachedCerts[fmt.Sprintf("%x", certBundle.Cert.SerialNumber)] = arn
}

func (f *cachedAcmCertFinder) NotifyCertDeleted(arn string) {
	delete(f.cachedCertSummary, arn)
	delete(f.cachedCerts, arn)
}

func (f *cachedAcmCertFinder) GetAcmApi() acmApi {
	return f.api
}

func (f *cachedAcmCertFinder) FindCertInACM(ctx context.Context, certBundle *certificateBundle) (arn string, err error) {
	// check cached result first
	if arn, ex := f.cachedCerts[fmt.Sprintf("%x", certBundle.Cert.SerialNumber)]; ex {
		return arn, nil
	}

	// fetch cert list if needed
	if len(f.cachedCertSummary) == 0 {
		if certs, err := f.listCertificates(ctx); err != nil {
			return "", err
		} else {
			for _, cert := range certs {
				f.cachedCertSummary[*cert.CertificateArn] = cert
			}
		}
	}

	// find in candidates
	for _, summary := range f.cachedCertSummary {
		// quick check before request for details
		if !summary.NotBefore.Equal(certBundle.Cert.NotBefore) || !summary.NotAfter.Equal(certBundle.Cert.NotAfter) {
			continue
		}
		if !certBundle.ContainsAllDomains(summary.SubjectAlternativeNameSummaries) {
			continue
		}

		// need serial number for final check
		if certDetails, err := f.api.DescribeCertificate(ctx, &acm.DescribeCertificateInput{
			CertificateArn: summary.CertificateArn,
		}); err != nil {
			return "", err
		} else if certDetails.Certificate.Serial != nil && certBundle.VerifySerialNumber(*certDetails.Certificate.Serial) {
			// update cache
			f.cachedCerts[fmt.Sprintf("%x", certBundle.Cert.SerialNumber)] = *summary.CertificateArn
			return *summary.CertificateArn, nil
		}
	}

	return "", err
}

// listCertificates lists all certificates with 'issued' status, includes amazon issued.
func (f *cachedAcmCertFinder) listCertificates(ctx context.Context) ([]acmTypes.CertificateSummary, error) {
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

// deleteManagedCertFromAcmIfUnused deletes a cert from ACM.
//
// The deleted certificate must meet the ALL conditions:
//   - Unused.
//   - A imported cert. i.e. not issued by amazon.
//   - Managed by the cert-deployer (has acmManagedTagKey tag).
func (d *deployer) deleteManagedCertFromAcmIfUnused(ctx context.Context, acmApi acmApi, certArn *string) (deleted bool, err error) {
	if certArn == nil {
		return false, nil
	}
	if result, err := acmApi.DescribeCertificate(ctx, &acm.DescribeCertificateInput{
		CertificateArn: certArn,
	}); err != nil {
		return false, err
	} else if len(result.Certificate.InUseBy) > 0 {
		return false, nil // still using
	} else if result.Certificate.Type != acmTypes.CertificateTypeImported {
		return false, nil // not imported by the user
	}

	// verify the cert is managed by this program
	if result, err := acmApi.ListTagsForCertificate(ctx, &acm.ListTagsForCertificateInput{
		CertificateArn: certArn,
	}); err != nil {
		return false, err
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
	_, err = acmApi.DeleteCertificate(ctx, &acm.DeleteCertificateInput{
		CertificateArn: certArn,
	})
	return err == nil, err
}
