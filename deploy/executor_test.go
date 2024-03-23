package deploy

import (
	"errors"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/mocker"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

func Test_defaultAssetDeployer_deployToAsset(t *testing.T) {
	type args struct {
		assetId               string
		fetchedAssetsProvider func(ty domain.AssetType) ([]domain.Asseter, error) // mock the result of deployer.listAssets
		deployResultProvider  func() []*domain.DeployError
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				assetId: "id1",
				fetchedAssetsProvider: func(ty domain.AssetType) ([]domain.Asseter, error) {
					return []domain.Asseter{
						&domain.Asset{Type: ty, Id: "id1", Available: true},
					}, nil
				},
				deployResultProvider: func() []*domain.DeployError {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "failed to list assets",
			args: args{
				assetId: "id1",
				fetchedAssetsProvider: func(ty domain.AssetType) ([]domain.Asseter, error) {
					return nil, errors.New("failed to list asserts")
				},
				deployResultProvider: nil,
			},
			wantErr: true,
		},
		{
			name: "target asset does not exist",
			args: args{
				assetId: "id2",
				fetchedAssetsProvider: func(ty domain.AssetType) ([]domain.Asseter, error) {
					return []domain.Asseter{
						&domain.Asset{Type: ty, Id: "id1", Available: true},
					}, nil
				},
				deployResultProvider: nil,
			},
			wantErr: true,
		},
		{
			name: "target asset unavailable",
			args: args{
				assetId: "id2",
				fetchedAssetsProvider: func(ty domain.AssetType) ([]domain.Asseter, error) {
					return []domain.Asseter{
						&domain.Asset{Type: ty, Id: "id2", Available: false},
					}, nil
				},
				deployResultProvider: nil,
			},
			wantErr: true,
		},
		{
			name: "deployment failure",
			args: args{
				assetId: "id2",
				fetchedAssetsProvider: func(ty domain.AssetType) ([]domain.Asseter, error) {
					return []domain.Asseter{
						&domain.Asset{Type: ty, Id: "id2", Available: true},
					}, nil
				},
				deployResultProvider: func() []*domain.DeployError {
					return []*domain.DeployError{domain.NewDeployError(nil, errors.New("err"))}
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deployer := mocker.NewMockDeployer(t)
			deployer.EXPECT().ListAssets(domain.AssetType("test")).RunAndReturn(tt.args.fetchedAssetsProvider).Once()
			if tt.args.deployResultProvider != nil {
				deployer.EXPECT().Deploy(mock.Anything, mock.Anything, mock.Anything).Return(nil, tt.args.deployResultProvider())
			}

			if err := newAssetDeployer().deployToAsset(deployer, "test", tt.args.assetId, nil, nil); (err != nil) != tt.wantErr {
				t.Errorf("deployToAsset() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultDeploymentExecutor_executeDeployment(t *testing.T) {
	type fields struct {
		fileReader      domain.FileReader
		deployerFactory domain.DeployerFactory
		assetDeployer   assetDeployer
	}
	type args struct {
		providers  map[string]domain.CloudProvider
		deployment domain.Deployment
	}
	successFileReader := domain.FileReaderFunc(func(name string) ([]byte, error) {
		return nil, nil
	})
	successDeployFactory := func() domain.DeployerFactory {
		f := mocker.NewMockDeployerFactory(t)
		f.EXPECT().NewDeployer(mock.Anything, mock.Anything).Return(mocker.NewMockDeployer(t), nil).Once()
		return f
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				fileReader:      successFileReader,
				deployerFactory: successDeployFactory(),
				assetDeployer: func() assetDeployer {
					d := NewMockassetDeployer(t)
					d.EXPECT().deployToAsset(mock.Anything, mock.Anything, mock.Anything, mock.Anything,
						mock.Anything).Return(nil).Times(2)
					return d
				}(),
			},

			args: args{
				providers: map[string]domain.CloudProvider{"demo": {}},
				deployment: domain.Deployment{
					ProviderId: "demo",
					Assets:     []domain.DeploymentAsset{{Type: "cdn"}, {Type: "cdn"}},
				},
			},
			wantErr: false,
		},
		{
			name: "all assets failed",
			fields: fields{
				fileReader:      successFileReader,
				deployerFactory: successDeployFactory(),
				assetDeployer: func() assetDeployer {
					d := NewMockassetDeployer(t)
					d.EXPECT().deployToAsset(mock.Anything, mock.Anything, mock.Anything, mock.Anything,
						mock.Anything).Return(errors.New("err")).Times(2)
					return d
				}(),
			},
			args: args{
				providers: map[string]domain.CloudProvider{"demo": {}},

				deployment: domain.Deployment{
					ProviderId: "demo",
					Assets:     []domain.DeploymentAsset{{Type: "cdn"}, {Type: "cdn"}},
				},
			},
			wantErr: false, // the overall deployment itself is success once it tried to deploy.
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewCustomDeploymentExecutor(zap.NewNop().Sugar(), tt.fields.fileReader, tt.fields.deployerFactory, tt.fields.assetDeployer)
			if err := n.ExecuteDeployment(tt.args.providers, tt.args.deployment); (err != nil) != tt.wantErr {
				t.Errorf("ExecuteDeployment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
