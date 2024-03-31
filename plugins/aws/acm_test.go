package aws

import (
	"context"
	_ "embed"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmTypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/mocker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// full_chain.pem for testing.
//   - Name: *.chenhe.me
//   - Alternative name: [*.chenhe.me, chenhe.me]
//   - Serial number: 9f7aa7f3f62a992d9364d7f5f47b52b1
//   - Algorithm: SHA384withECDSA
//
//go:embed test_fullchain.pem
var testCert []byte

func Test_cachedAcmManager_FindCertInACM(t *testing.T) {
	mocker.NewMockCertificateBundle(t)
	targetCertBundle, _ := domain.NewCertificateBundle(testCert)
	tests := []struct {
		name          string
		queryTimes    int
		api           func(t *testing.T) acmApi
		wantFromCache []bool
		wantArn       string
		wantErr       bool
	}{
		{
			name:       "found the cert from ACM",
			queryTimes: 1,
			api: func(t *testing.T) acmApi {
				api := NewMockacmApi(t)
				output := &acm.ListCertificatesOutput{
					CertificateSummaryList: []acmTypes.CertificateSummary{
						{
							CertificateArn: aws.String("arn1"),
							NotBefore:      targetCertBundle.NotBefore(),
							NotAfter:       targetCertBundle.NotAfter(),
						},
					},
				}
				api.EXPECT().ListCertificates(mock.Anything, mock.Anything, mock.Anything).Return(output, nil).Once()
				api.EXPECT().DescribeCertificate(mock.Anything, mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, input *acm.DescribeCertificateInput, f ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error) {
					assert.Equal(t, "arn1", *input.CertificateArn)
					return &acm.DescribeCertificateOutput{
						Certificate: &acmTypes.CertificateDetail{
							CertificateArn: input.CertificateArn,
							Serial:         aws.String(targetCertBundle.GetSerialNumberHexString()),
						},
					}, nil
				}).Once()
				return api
			},
			wantFromCache: []bool{false},
			wantArn:       "arn1",
			wantErr:       false,
		},
		{
			name:       "found the cert from cache",
			queryTimes: 2,
			api: func(t *testing.T) acmApi {
				api := NewMockacmApi(t)
				output := &acm.ListCertificatesOutput{
					CertificateSummaryList: []acmTypes.CertificateSummary{
						{
							CertificateArn: aws.String("arn1"),
							NotBefore:      targetCertBundle.NotBefore(),
							NotAfter:       targetCertBundle.NotAfter(),
						},
					},
				}
				api.EXPECT().ListCertificates(mock.Anything, mock.Anything, mock.Anything).Return(output, nil).Once()
				api.EXPECT().DescribeCertificate(mock.Anything, mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, input *acm.DescribeCertificateInput, f ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error) {
					assert.Equal(t, "arn1", *input.CertificateArn)
					return &acm.DescribeCertificateOutput{
						Certificate: &acmTypes.CertificateDetail{
							CertificateArn: input.CertificateArn,
							Serial:         aws.String(targetCertBundle.GetSerialNumberHexString()),
						},
					}, nil
				}).Once()
				return api
			},
			wantFromCache: []bool{false, true},
			wantArn:       "arn1",
			wantErr:       false,
		},
		{
			name:       "not found",
			queryTimes: 1,
			api: func(t *testing.T) acmApi {
				api := NewMockacmApi(t)
				output := &acm.ListCertificatesOutput{}
				api.EXPECT().ListCertificates(mock.Anything, mock.Anything, mock.Anything).Return(output, nil).Once()
				return api
			},
			wantFromCache: []bool{false},
			wantArn:       "",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := newCachedAcmCertFinder(tt.api(t))

			for i := range tt.queryTimes {
				gotArn, fromCachefalse, err := finder.FindCertInACM(context.Background(), targetCertBundle)
				if (err != nil) != tt.wantErr {
					t.Errorf("FindCertInACM() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				assert.Equal(t, tt.wantFromCache[i], fromCachefalse)
				if gotArn != tt.wantArn {
					t.Errorf("FindCertInACM() gotArn = %v, want %v", gotArn, tt.wantArn)
				}
			}
		})
	}
}
