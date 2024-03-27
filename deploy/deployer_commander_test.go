package deploy

import (
	"errors"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/mocker"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_cachedDeployerCommander_DeployToAsset(t *testing.T) {
	type args struct {
		assetId               []string
		fetchedAssetsProvider func(ty string) ([]domain.Asseter, error) // mock the result of deployer.listAssets
		deployResultProvider  func() []*domain.DeployError
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
				fetchedAssetsProvider: func(ty string) ([]domain.Asseter, error) {
					return []domain.Asseter{
						&domain.Asset{Type: ty, Id: "id1", Available: true},
					}, nil
				},
				deployResultProvider: func() []*domain.DeployError {
					return nil
				},
			},
			wantErr: []bool{false},
		},
		{
			name: "use assets cache",
			args: args{
				assetId: []string{"id1", "id1", "id1"},
				fetchedAssetsProvider: func(ty string) ([]domain.Asseter, error) {
					return []domain.Asseter{
						&domain.Asset{Type: ty, Id: "id1", Available: true},
					}, nil
				},
				deployResultProvider: func() []*domain.DeployError {
					return nil
				},
			},
			wantErr: []bool{false, false, false},
		},
		{
			name: "failed to list assets",
			args: args{
				assetId: []string{"id1"},
				fetchedAssetsProvider: func(ty string) ([]domain.Asseter, error) {
					return nil, errors.New("failed to list asserts")
				},
				deployResultProvider: nil,
			},
			wantErr: []bool{true},
		},
		{
			name: "target asset does not exist",
			args: args{
				assetId: []string{"id2"},
				fetchedAssetsProvider: func(ty string) ([]domain.Asseter, error) {
					return []domain.Asseter{
						&domain.Asset{Type: ty, Id: "id1", Available: true},
					}, nil
				},
				deployResultProvider: nil,
			},
			wantErr: []bool{true},
		},
		{
			name: "deployment failure",
			args: args{
				assetId: []string{"id1"},
				fetchedAssetsProvider: func(ty string) ([]domain.Asseter, error) {
					return []domain.Asseter{
						&domain.Asset{Type: ty, Id: "id1", Available: true},
					}, nil
				},
				deployResultProvider: func() []*domain.DeployError {
					return []*domain.DeployError{domain.NewDeployError(nil, errors.New("err"))}
				},
			},
			wantErr: []bool{true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deployer := mocker.NewMockDeployer(t)
			deployer.EXPECT().ListAssets("test").RunAndReturn(tt.args.fetchedAssetsProvider).Once()
			if tt.args.deployResultProvider != nil {
				deployer.EXPECT().Deploy(mock.Anything, mock.Anything, mock.Anything).Return(nil, tt.args.deployResultProvider())
			}

			cmder := newCachedDeployerCommander(deployer)
			for i, targetId := range tt.args.assetId {
				if err := cmder.DeployToAsset("test", targetId, nil, nil); (err != nil) != tt.wantErr[i] {
					t.Errorf("DeployToAsset() error = %v, wantErr = %v", err, tt.wantErr)
				}
			}
		})
	}
}
