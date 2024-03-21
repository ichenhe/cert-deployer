package registry

import (
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustRegister(t *testing.T) {
	deployerConstructor := func(options domain.Options) (domain.Deployer, error) {
		return mocker.NewMockDeployer(t), nil
	}

	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{
			name: "success",
			f: func(t *testing.T) {
				assert.NotPanics(t, func() {
					MustRegister("provider", deployerConstructor)
				})
				assert.Len(t, assetDeployerConstructors, 1)
			},
		}, {
			name: "providerDuplicate_panic",
			f: func(t *testing.T) {
				MustRegister("d", deployerConstructor)
				assert.Panics(t, func() {
					MustRegister("d", deployerConstructor)
				})
				assert.Len(t, assetDeployerConstructors, 1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clear(assetDeployerConstructors)
			tt.f(t)
		})
	}
}

func Test_defaultDeployerFactory_NewDeployer(t *testing.T) {
	tests := []struct {
		name                       string
		registeredDeployerProvider []string
		requestProvider            string
		wantErr                    bool
	}{
		{
			name:                       "ok",
			registeredDeployerProvider: []string{"a", "b", "c"},
			requestProvider:            "b",
			wantErr:                    false,
		},
		{
			name:                       "providerNotExists",
			registeredDeployerProvider: []string{"a"},
			requestProvider:            "b",
			wantErr:                    true,
		},
	}

	for _, tt := range tests {
		clear(assetDeployerConstructors)
		deployers := make(map[string]domain.Deployer)
		t.Run(tt.name, func(t *testing.T) {
			for _, p := range tt.registeredDeployerProvider {
				x := mocker.NewMockDeployer(t)
				deployers[p] = x
				deployerConstructor := func(options domain.Options) (domain.Deployer, error) {
					return x, nil
				}
				MustRegister(p, deployerConstructor)
			}

			d, e := NewDeployerFactory().NewDeployer(nil, domain.CloudProvider{
				Provider: tt.requestProvider,
			})

			if !tt.wantErr {
				assert.NoError(t, e)
				assert.Same(t, deployers[tt.requestProvider], d, "wrong deployer")
			} else {
				assert.NotNil(t, e, "expect an error")
			}
		})
	}
}
