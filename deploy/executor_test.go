package deploy

import (
	"context"
	"errors"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/mocker"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

func Test_defaultDeploymentExecutor_executeDeployment(t *testing.T) {
	type fields struct {
		fileReader        domain.FileReader
		deployerFactory   domain.DeployerFactory
		deployerCommander deployerCommander
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
			name: "deploy to 2 assets",
			fields: fields{
				fileReader:      successFileReader,
				deployerFactory: successDeployFactory(),
				deployerCommander: func() deployerCommander {
					c := NewMockdeployerCommander(t)
					c.EXPECT().IsAssetTypeSupported("cdn").Return(true)
					c.EXPECT().DeployToAsset(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(2)
					return c
				}(),
			},
			args: args{
				providers: map[string]domain.CloudProvider{"demo": {}},
				deployment: domain.Deployment{
					ProviderId: "demo",
					Assets:     []domain.DeploymentAsset{{Type: "cdn", Id: "x"}, {Type: "cdn", Id: "x"}},
				},
			},
			wantErr: false,
		},
		{
			name: "deploy to a type",
			fields: fields{
				fileReader:      successFileReader,
				deployerFactory: successDeployFactory(),
				deployerCommander: func() deployerCommander {
					c := NewMockdeployerCommander(t)
					c.EXPECT().IsAssetTypeSupported("cdn").Return(true)
					c.EXPECT().DeployToAssetType(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
					return c
				}(),
			},
			args: args{
				providers: map[string]domain.CloudProvider{"demo": {}},
				deployment: domain.Deployment{
					ProviderId: "demo",
					Assets:     []domain.DeploymentAsset{{Type: "cdn"}},
				},
			},
			wantErr: false,
		},
		{
			name: "all assets failed",
			fields: fields{
				fileReader:      successFileReader,
				deployerFactory: successDeployFactory(),
				deployerCommander: func() deployerCommander {
					c := NewMockdeployerCommander(t)
					c.EXPECT().IsAssetTypeSupported("cdn").Return(true)
					c.EXPECT().DeployToAsset(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("err")).Times(2)
					return c
				}(),
			},
			args: args{
				providers: map[string]domain.CloudProvider{"demo": {}},

				deployment: domain.Deployment{
					ProviderId: "demo",
					Assets:     []domain.DeploymentAsset{{Type: "cdn", Id: "x"}, {Type: "cdn", Id: "x"}},
				},
			},
			wantErr: false, // the overall deployment itself is success once it tried to deploy.
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewDeploymentExecutor(zap.NewNop().Sugar(), tt.args.providers).(*defaultDeploymentExecutor)
			e.fileReader = tt.fields.fileReader
			e.deployerFactory = tt.fields.deployerFactory
			e.deployerCommanderFactory = func(deployer domain.Deployer) deployerCommander {
				return tt.fields.deployerCommander
			}
			if err := e.ExecuteDeployment(context.Background(), tt.args.deployment); (err != nil) != tt.wantErr {
				t.Errorf("ExecuteDeployment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
