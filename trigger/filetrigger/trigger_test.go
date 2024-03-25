package filetrigger

import (
	"github.com/fsnotify/fsnotify"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/ichenhe/cert-deployer/mocker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestFileTrigger_StartMonitoring(t *testing.T) {
	type args struct {
		file               string
		startRepeatability int
		waitFor            int
		eventAfter         []int // [50,20] means: 50ms-E1-20ms-E2
		triggeredTimes     int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "start once",
			args: args{
				file:               "/path/a.txt",
				startRepeatability: 1,
			},
		},
		{
			name: "multiple concurrent start",
			args: args{
				file:               "/path/a.txt",
				startRepeatability: 100,
			},
		},
		{
			name: "zero wait time",
			args: args{
				file:               "/path/a.txt",
				startRepeatability: 1,
				eventAfter:         []int{0, 0, 2},
				triggeredTimes:     3,
			},
		},
		{
			name: "wait for 20ms",
			args: args{
				file:               "/path/a.txt",
				startRepeatability: 1,
				waitFor:            20,
				eventAfter:         []int{0, 5, 30}, // canceled, triggered, triggered
				triggeredTimes:     2,
			},
		},
	}
	for _, tt := range tests {
		var hashCount int
		opt := domain.FileMonitoringTriggerOptions{File: tt.args.file, Wait: tt.args.waitFor}
		deployments := []domain.Deployment{
			{Name: "d1"},
			{Name: "d2"},
		}
		t.Run(tt.name, func(t *testing.T) {
			eventCh := make(chan fsnotify.Event)
			errorCh := make(chan error)
			defer func() {
				if eventCh != nil {
					close(eventCh)
				}
				if errorCh != nil {
					close(errorCh)
				}
			}()
			executor := mocker.NewMockDeploymentExecutor(t)
			monitor := NewMockfileMonitor(t)
			trigger := NewFileTrigger(zap.NewNop().Sugar(), tt.name, executor, opt, deployments)
			trigger.filer = filerFunc(func(path string) bool { return false }) // not dir
			trigger.monitorFactory = func() (fileMonitor, error) { return monitor, nil }
			// give different value for each hashing try
			trigger.hasher = func(file string) (string, error) {
				hashCount++
				return strconv.Itoa(hashCount), nil
			}
			monitor.EXPECT().ErrorCh().Return(errorCh).Maybe()
			monitor.EXPECT().EventCh().Return(eventCh).Maybe()

			// monitor itself only be closed once, no matter how many times the trigger is closed
			monitor.EXPECT().Close().RunAndReturn(func() error {
				close(eventCh)
				eventCh = nil
				close(errorCh)
				errorCh = nil
				return nil
			}).Once()

			// always watch once no matter how many times start
			monitor.EXPECT().Add(filepath.Dir(tt.args.file)).Return(nil).Once()

			if tt.args.triggeredTimes > 0 {
				executor.EXPECT().ExecuteDeployment(mock.Anything).Return(nil).Times(tt.args.triggeredTimes * len(deployments))
			}

			// start
			var wg sync.WaitGroup
			wg.Add(tt.args.startRepeatability)
			for range tt.args.startRepeatability {
				go func() {
					assert.NoError(t, trigger.StartMonitoring())
					wg.Done()
				}()
			}
			wg.Wait()

			// events
			for _, delay := range tt.args.eventAfter {
				if delay > 0 {
					time.Sleep(time.Duration(delay) * time.Millisecond)
				}
				eventCh <- fsnotify.Event{Name: tt.args.file, Op: fsnotify.Write}
			}
			// allow the event to propagate and wait for any deployments in waiting list
			time.Sleep(time.Duration(tt.args.waitFor+10) * time.Millisecond)

			// close
			wg.Add(tt.args.startRepeatability)
			for range tt.args.startRepeatability {
				go func() {
					trigger.closeAndWait()
					wg.Done()
				}()
			}
			wg.Wait()
		})
	}
}
