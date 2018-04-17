package runner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func watchFolder(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal(err)
	}

	// defer watcher.Close()

	// done := make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Events:

				if isWatchedFile(ev.Name) {
					watcherLog("event: %s", ev)

					if ev.Op&fsnotify.Write == fsnotify.Write {
						watcherLog("modified file: %s", ev)
						//log.Println("modified file:", event.Name)
						//startChannel <- ev.String()
					}

				}
			case err := <-watcher.Errors:
				watcherLog("error: %s", err)
			}
		}
	}()

	watcherLog("Watching %s", path)
	err = watcher.Add(path)

	if err != nil {
		fatal(err)
	}
	// <-done
}

func watch() {
	root := root()
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() && !isTmpDir(path) {
			if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}

			if isIgnoredFolder(path) {
				watcherLog("Ignoring %s", path)
				return filepath.SkipDir
			}

			watchFolder(path)
		}

		return err
	})
}
