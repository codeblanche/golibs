package fsutils

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/howeyc/fsnotify"
)

// Watch for changes in the specified directory recursively
func Watch(dir string, handler func(name string)) error {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return err
	}

	go watch(watcher, dir, handler)

	return nil
}

func watch(watcher *fsnotify.Watcher, dir string, handler func(name string)) {
	defer watcher.Close()

	err := recursive(watcher, dir)

	if err != nil {
		return
	}

	for {
		select {
		case ev := <-watcher.Event:
			if ev.IsCreate() {
				recursive(watcher, ev.Name)
			}

			handler(ev.Name)
		case <-watcher.Error:
			// nothing to do for now
		}
	}
}

// recursive watches the given dir and recurses into all subdirectories as well
func recursive(watcher *fsnotify.Watcher, dir string) error {
	info, err := os.Stat(dir)

	if err == nil && !info.Mode().IsDir() {
		err = errors.New("Watching a file is not supported. Expected a directory")
	}

	// Watch the specified dir
	if err == nil {
		err = watcher.Watch(dir)
	}

	var list []os.FileInfo

	// Grab list of subdirs
	if err == nil {
		list, err = ioutil.ReadDir(dir)
	}

	// Call recursive for each dir in list
	if err == nil {
		for _, file := range list {
			recursive(watcher, filepath.Join(dir, file.Name()))
		}
	}

	return err
}
