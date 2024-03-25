package filetrigger

import "github.com/fsnotify/fsnotify"

type monitorFactory = func() (fileMonitor, error)

// fileMonitor is the abstraction of fsnotify watcher for testing.
type fileMonitor interface {
	Add(name string) error
	Close() error
	ErrorCh() chan error
	EventCh() chan fsnotify.Event
}

func newFileMonitor() (fileMonitor, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &fsnotifyFileMonitor{
		watcher: w,
	}, nil
}

type fsnotifyFileMonitor struct {
	watcher *fsnotify.Watcher
}

func (f fsnotifyFileMonitor) Add(name string) error {
	return f.watcher.Add(name)
}

func (f fsnotifyFileMonitor) Close() error {
	return f.watcher.Close()
}

func (f fsnotifyFileMonitor) ErrorCh() chan error {
	return f.watcher.Errors
}

func (f fsnotifyFileMonitor) EventCh() chan fsnotify.Event {
	return f.watcher.Events
}
