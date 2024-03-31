package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

func Test_deployer_deployCloudFrontCert(t *testing.T) {

	type args struct {
		cfApi      func(t *testing.T) cloudfrontApi
		certFinder func(t *testing.T) acmManager
	}

	// fixed args
	certBundle, _ := domain.NewCertificateBundle(testCert)
	deployAsset := &cloudFrontDistribution{Asset: domain.Asset{Id: "id"}}

	// helpers
	createGetDistributionConfigOutput := func(aliases ...string) *cloudfront.GetDistributionConfigOutput {
		return &cloudfront.GetDistributionConfigOutput{
			DistributionConfig: &types.DistributionConfig{
				Aliases: &types.Aliases{
					Quantity: aws.Int32(int32(len(aliases))),
					Items:    aliases,
				},
			},
		}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "deploy to a cloudfront with only amazon domain",
			args: args{
				cfApi: func(t *testing.T) cloudfrontApi {
					f := NewMockcloudfrontApi(t)
					f.EXPECT().GetDistributionConfig(mock.Anything, mock.Anything).Return(createGetDistributionConfigOutput(), nil)
					return f
				},
				certFinder: func(t *testing.T) acmManager {
					return NewMockacmManager(t)
				},
			},
			wantErr: true,
		},
		{
			name: "deploy to a cloudfront with domain does not match the cert",
			args: args{
				cfApi: func(t *testing.T) cloudfrontApi {
					f := NewMockcloudfrontApi(t)
					f.EXPECT().GetDistributionConfig(mock.Anything, mock.Anything).
						Return(createGetDistributionConfigOutput("www.chenhe.me", "not-in-the-cert.xyz"), nil)
					return f
				},
				certFinder: func(t *testing.T) acmManager {
					return NewMockacmManager(t)
				},
			},
			wantErr: true,
		},
		{
			name: "reuse the cert in ACM",
			args: args{
				cfApi: func(t *testing.T) cloudfrontApi {
					f := NewMockcloudfrontApi(t)
					f.EXPECT().GetDistributionConfig(mock.Anything, mock.Anything).
						Return(createGetDistributionConfigOutput("chenhe.me"), nil)
					f.EXPECT().UpdateDistribution(mock.Anything, mock.Anything).Return(nil, nil)
					return f
				},
				certFinder: func(t *testing.T) acmManager {
					f := NewMockacmManager(t)
					f.EXPECT().FindCertInACM(mock.Anything, mock.Anything).Return("arn", false, nil)
					return f
				},
			},
			wantErr: false,
		},
		{
			name: "import a new cert to ACM",
			args: args{
				cfApi: func(t *testing.T) cloudfrontApi {
					f := NewMockcloudfrontApi(t)
					f.EXPECT().GetDistributionConfig(mock.Anything, mock.Anything).
						Return(createGetDistributionConfigOutput("chenhe.me"), nil)
					f.EXPECT().UpdateDistribution(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, input *cloudfront.UpdateDistributionInput, f ...func(*cloudfront.Options)) (*cloudfront.UpdateDistributionOutput, error) {
						assert.Equal(t, "new-arn", *input.DistributionConfig.ViewerCertificate.ACMCertificateArn)
						return nil, nil
					})
					return f
				},
				certFinder: func(t *testing.T) acmManager {
					m := NewMockacmManager(t)
					m.EXPECT().ImportCertificate(mock.Anything, mock.Anything, mock.Anything).Return("new-arn", nil)
					m.EXPECT().FindCertInACM(mock.Anything, mock.Anything).Return("", false, nil)
					return m
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, _ := newAwsDeployer("", "", zap.NewNop().Sugar())
			err := d.deployCloudFrontCert(context.Background(), tt.args.cfApi(t), tt.args.certFinder(t), deployAsset, certBundle, make([]byte, 0))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
