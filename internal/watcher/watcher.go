package watcher

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	events  chan struct{}
	ignored map[string]bool
	done    chan struct{}
}

func New(root string, userIgnored []string) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	ignoreMap := make(map[string]bool)
	for _, i := range userIgnored {
		ignoreMap[i] = true
	}

	// System ignores
	systemIgnored := []string{".git", ".kata", ".task", "katas", "node_modules", "vendor", ".venv", "bin"}
	for _, i := range systemIgnored {
		ignoreMap[i] = true
	}

	events := make(chan struct{})
	done := make(chan struct{})

	// Walk and add
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// Permission denied or other errors should not break the whole watch?
			// But for WalkDir, returning err stops walk.
			// We log and continue?
			slog.Debug("watcher walk error", "path", path, "error", err)
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if ignoreMap[name] {
				return filepath.SkipDir
			}
			return w.Add(path)
		}
		return nil
	})
	if err != nil {
		w.Close()
		return nil, err
	}

	gw := &Watcher{
		watcher: w,
		events:  events,
		ignored: ignoreMap,
		done:    done,
	}

	go gw.loop()

	return gw, nil
}

func (w *Watcher) Events() <-chan struct{} {
	return w.events
}

func (w *Watcher) Close() {
	close(w.done)
	w.watcher.Close()
}

func (w *Watcher) loop() {
	var debounceTimer *time.Timer
	debounceDuration := 300 * time.Millisecond

	trigger := func() {
		select {
		case w.events <- struct{}{}:
		default:
			// Consumer busy, skip
		}
	}

	for {
		select {
		case <-w.done:
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			return
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// Handle new directories
			if event.Has(fsnotify.Create) {
				stat, err := os.Stat(event.Name)
				if err == nil && stat.IsDir() {
					base := filepath.Base(event.Name)
					if !w.ignored[base] {
						_ = w.watcher.Add(event.Name)
					}
				}
			}

			// Debounce
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceTimer = time.AfterFunc(debounceDuration, trigger)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			slog.Debug("watcher error", "error", err)
		}
	}
}
