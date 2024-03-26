package main

import (
	"github.com/ichenhe/cert-deployer/config"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v2"
	"os"
	"syscall"
	"testing"
	"time"
)

func Test_run(t *testing.T) {
	profileLoader := func(c *cli.Context) (*domain.AppConfig, error) {
		return config.DefaultConfig(), nil
	}
	tests := []struct {
		name       string
		args       []string
		executor   func(t *testing.T) commandExecutor
		fileReader domain.FileReader
		waitFor    int
		wantErr    bool
	}{
		{
			name: "deploy with 3 deployments",
			args: []string{"--profile", "/a.yml", "deploy", "--deployment", "a,b,c"},
			executor: func(t *testing.T) commandExecutor {
				e := NewMockcommandExecutor(t)
				e.EXPECT().executeDeployments(mock.Anything, mock.Anything, []string{"a", "b", "c"}).Return().Once()
				return e
			},
			wantErr: false,
		},
		{
			name: "custom deploy",
			args: []string{"deploy", "--type", "cdn", "--cert", "cert.pem", "--key", "key.pem", "--provider", "TencentCloud", "--secret-id", "x", "--secret-key", "y"},
			executor: func(t *testing.T) commandExecutor {
				e := NewMockcommandExecutor(t)
				verifier := func(providers map[string]domain.CloudProvider, deployments map[string]domain.Deployment, deploymentIds []string) {
					assert.Len(t, deploymentIds, 1)
					assert.Len(t, providers, 1)
					for _, provider := range providers {
						assert.Equal(t, "TencentCloud", provider.Provider)
						assert.Equal(t, "x", provider.SecretId)
						assert.Equal(t, "y", provider.SecretKey)
					}
					assert.Len(t, deployments, 1)
					for _, deployment := range deployments {
						assert.Equal(t, "cert.pem", deployment.Cert)
						assert.Equal(t, "key.pem", deployment.Key)
						assert.Equal(t, deploymentIds[0], deployment.Name)
						assert.Len(t, deployment.Assets, 1)
						assert.Equal(t, "cdn", deployment.Assets[0].Type)
					}
				}
				e.EXPECT().executeDeployments(mock.Anything, mock.Anything, mock.Anything).Run(verifier).Return().Once()
				return e
			},
			fileReader: domain.FileReaderFunc(func(name string) ([]byte, error) {
				return nil, nil
			}),
			wantErr: false,
		},
		{
			name: "run trigger",
			args: []string{"--profile", "/a.yml", "run"},
			executor: func(t *testing.T) commandExecutor {
				e := NewMockcommandExecutor(t)
				e.EXPECT().registerTriggers(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				return e
			},
			wantErr: false,
			waitFor: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdDispatcher := newCommandDispatcher(profileLoader, tt.fileReader, tt.executor(t))
			args := os.Args[0:1]
			args = append(args, tt.args...)
			if tt.waitFor > 0 {
				tm := time.AfterFunc(time.Millisecond*time.Duration(tt.waitFor), func() {
					_ = syscall.Kill(os.Getpid(), syscall.SIGINT)
				})
				defer tm.Stop()
			}
			if err := run(args, cmdDispatcher); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
