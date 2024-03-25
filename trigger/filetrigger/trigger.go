package filetrigger

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/ichenhe/cert-deployer/domain"
	"go.uber.org/zap"
	"hash"
	"io"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const mb = 1024 * 1024

var _ domain.Trigger = &FileTrigger{}

func NewFileTrigger(logger *zap.SugaredLogger, name string, deploymentExecutor domain.DeploymentExecutor,
	options domain.FileMonitoringTriggerOptions, deployments []domain.Deployment) *FileTrigger {
	return &FileTrigger{
		logger:             logger,
		deploymentExecutor: deploymentExecutor,
		options:            options,
		deployments:        deployments,
		name:               name,

		monitorFactory: func() (fileMonitor, error) { return newFileMonitor() },
		filer:          defaultFiler{},
		hasher:         calculateHash,
	}
}

type FileTrigger struct {
	logger             *zap.SugaredLogger
	deploymentExecutor domain.DeploymentExecutor
	options            domain.FileMonitoringTriggerOptions
	deployments        []domain.Deployment
	name               string
	mu                 sync.Mutex

	monitor        fileMonitor
	fileHash       string // hash value of watched file
	closeWaitGroup sync.WaitGroup

	//  ---------- for testing only

	monitorFactory monitorFactory
	filer          filer
	hasher         func(file string) (string, error)
}

func (t *FileTrigger) GetName() string {
	return t.name
}

func (t *FileTrigger) StartMonitoring() (err error) {
	if t.filer.IsDir(t.options.File) {
		return fmt.Errorf("can not watch a dir: %s", t.options.File)
	}

	if fileHash, err := t.hasher(t.options.File); err != nil {
		t.fileHash = ""
	} else {
		t.fileHash = fileHash
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.monitor != nil {
		// already watching
		return nil
	}

	if w, err := t.monitorFactory(); err != nil {
		return fmt.Errorf("failed to create a watcher: %w", err)
	} else {
		t.monitor = w
	}

	t.closeWaitGroup.Add(1)
	go t.fileLoop(t.monitor.ErrorCh(), t.monitor.EventCh(), t.options.File)

	if err = t.monitor.Add(filepath.Dir(t.options.File)); err != nil {
		return fmt.Errorf("faliled to watch '%s': %w", t.options.File, err)
	}

	return nil
}

// Close closes the monitor. Triggered deployment will continue to run but any deployments in waiting
// list will be canceled.
func (t *FileTrigger) Close() {
	t.mu.Lock()
	if t.monitor != nil {
		if err := t.monitor.Close(); err != nil {
			t.logger.Warn("failed to close watcher: ", err)
		}
		t.monitor = nil
	}
	t.mu.Unlock()
}

// closeAndWait closes the monitor and waits for the event handler to exit. Be aware that deployments
// in waiting list will still be canceled.
// This method is for testing.
func (t *FileTrigger) closeAndWait() {
	t.Close()
	t.closeWaitGroup.Wait()
}

func (t *FileTrigger) fileLoop(errors chan error, events chan fsnotify.Event, file string) {
	defer func() {
		t.closeWaitGroup.Done()
	}()

	waitFor := time.Duration(t.options.Wait) * time.Millisecond
	var tm *time.Timer

	eventHandler := func(e fsnotify.Event) {
		if t.monitor == nil {
			return // already closed
		}
		if currentHash, err := t.hasher(e.Name); err != nil {
			t.logger.Debugf("ignore the event: cannot hash file '%s': %v", e.Name, err)
			return // file does not exist, or content not changed
		} else {
			if currentHash == t.fileHash {
				return // file content not changed
			}
			t.fileHash = currentHash
		}
		for _, deployment := range t.deployments {
			if err := t.deploymentExecutor.ExecuteDeployment(deployment); err != nil {
				t.logger.Warnf("failed to execute deployment '%s': %v", deployment.Name, err)
			} else {
				t.logger.Infof("deployment '%s' completed", deployment.Name)
			}
		}
	}

	for {
		select {
		case _, ok := <-errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
		case e, ok := <-events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			if e.Name != file {
				continue
			}
			if !e.Has(fsnotify.Create) && !e.Has(fsnotify.Write) {
				continue
			}

			if waitFor <= 0 {
				eventHandler(e)
				continue
			}

			if tm == nil {
				tm = time.AfterFunc(math.MaxInt64, func() { eventHandler(e) })
				tm.Stop()
			}
			tm.Reset(waitFor)
		}
	}
}

// calculateHash calculates the hash value of given file.
func calculateHash(path string) (string, error) {
	var (
		err  error
		info os.FileInfo
		f    *os.File
	)
	if info, err = os.Stat(path); err != nil {
		return "", err
	}
	if f, err = os.Open(path); err != nil {
		return "", err
	}
	defer func(f *os.File) { _ = f.Close() }(f)

	var hasher hash.Hash
	if info.Size() <= mb*50 {
		hasher = sha256.New()
	} else {
		hasher = md5.New()
	}
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x\n", hasher.Sum(nil)), nil
}
