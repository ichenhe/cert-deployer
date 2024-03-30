package deploy

import (
	"context"
	"errors"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/mocker"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_cachedDeployerCommander_DeployToAsset(t *testing.T) {
	type args struct {
		assetId               []string
		fetchedAssetsProvider func(ctx context.Context, ty string) ([]domain.Asseter, error) // mock the result of deployer.listAssets
		deployError           error
	}

	tests := []struct {
		name    string
		args    args
		wantErr []bool
	}{
		{
			name: "success",
			args: args{
				assetId: []string{"id1"},
				fetchedAssetsProvider: func(ctx context.Context, ty string) ([]domain.Asseter, error) {
					return []domain.Asseter{
						&domain.Asset{Type: ty, Id: "id1", Available: true},
					}, nil
				},
			},
			wantErr: []bool{false},
		},
		{
			name: "use assets cache",
			args: args{
				assetId: []string{"id1", "id1", "id1"},
				fetchedAssetsProvider: func(ctx context.Context, ty string) ([]domain.Asseter, error) {
					return []domain.Asseter{
						&domain.Asset{Type: ty, Id: "id1", Available: true},
					}, nil
				},
			},
			wantErr: []bool{false, false, false},
		},
		{
			name: "failed to list assets",
			args: args{
				assetId: []string{"id1"},
				fetchedAssetsProvider: func(ctx context.Context, ty string) ([]domain.Asseter, error) {
					return nil, errors.New("failed to list asserts")
				},
			},
			wantErr: []bool{true},
		},
		{
			name: "target asset does not exist",
			args: args{
				assetId: []string{"id2"},
				fetchedAssetsProvider: func(ctx context.Context, ty string) ([]domain.Asseter, error) {
					return []domain.Asseter{
						&domain.Asset{Type: ty, Id: "id1", Available: true},
					}, nil
				},
			},
			wantErr: []bool{true},
		},
		{
			name: "deployment failure",
			args: args{
				assetId: []string{"id1"},
				fetchedAssetsProvider: func(ctx context.Context, ty string) ([]domain.Asseter, error) {
					return []domain.Asseter{
						&domain.Asset{Type: ty, Id: "id1", Available: true},
					}, nil
				},
				deployError: errors.New("err"),
			},
			wantErr: []bool{true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deployer := mocker.NewMockDeployer(t)
			deployer.EXPECT().ListAssets(mock.Anything, "test").RunAndReturn(tt.args.fetchedAssetsProvider).Once()
			deployer.EXPECT().Deploy(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, _ []domain.Asseter, _ []byte, _ []byte, callback *domain.DeployCallback) error {
				callback.ResultCallback(nil, tt.args.deployError)
				return nil
			}).Maybe()

			cmder := newCachedDeployerCommander(deployer)
			for i, targetId := range tt.args.assetId {
				if err := cmder.DeployToAsset(context.Background(), "test", targetId, nil, nil); (err != nil) != tt.wantErr[i] {
					t.Errorf("DeployToAsset() error = %v, wantErr = %v", err, tt.wantErr)
				}
			}
		})
	}
}
