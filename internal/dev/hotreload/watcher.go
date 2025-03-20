package hotreload

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	watcher  *fsnotify.Watcher
	debounce time.Duration
	events   chan Event
	errors   chan error
	done     chan struct{}
}

type Event struct {
	Path string
	Op   fsnotify.Op
}

func NewFileWatcher(debounce time.Duration) (*FileWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &FileWatcher{
		watcher:  w,
		debounce: debounce,
		events:   make(chan Event),
		errors:   make(chan error),
		done:     make(chan struct{}),
	}, nil
}

func (fw *FileWatcher) Watch(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return fw.watcher.Add(path)
		}
		return nil
	})
}

func (fw *FileWatcher) Start() {
	go fw.run()
}

func (fw *FileWatcher) Stop() {
	close(fw.done)
	fw.watcher.Close()
}

func (fw *FileWatcher) Events() <-chan Event {
	return fw.events
}

func (fw *FileWatcher) Errors() <-chan error {
	return fw.errors
}

func (fw *FileWatcher) run() {
	timer := time.NewTimer(0)
	<-timer.C

	for {
		select {
		case event := <-fw.watcher.Events:
			timer.Reset(fw.debounce)
			select {
			case <-timer.C:
				fw.events <- Event{
					Path: event.Name,
					Op:   event.Op,
				}
			case <-fw.done:
				return
			}
		case err := <-fw.watcher.Errors:
			fw.errors <- err
		case <-fw.done:
			return
		}
	}
}
