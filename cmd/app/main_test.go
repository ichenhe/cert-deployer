package main

import (
	"github.com/ichenhe/cert-deployer/config"
	"github.com/ichenhe/cert-deployer/domain"
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
			args: []string{"--profile", "/a.yml", "deploy", "--type", "cdn", "--cert", "cert.pem", "--key", "key.pem"},
			executor: func(t *testing.T) commandExecutor {
				e := NewMockcommandExecutor(t)
				e.EXPECT().customDeploy(mock.Anything, []string{"cdn"}, mock.Anything, mock.Anything).Return().Once()
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
