package magpie

import (
	"github.com/dietsche/rfsnotify"
	"gopkg.in/fsnotify.v1"
)

func (m *Magpie) Watch() {
	for _, c := range m.configs {
		w := newWatcher(c)
		go w.watch()
	}
	select {}
}

type watcher struct {
	c *config
	w *rfsnotify.RWatcher
}

func newWatcher(c *config) *watcher {
	return &watcher{c: c}
}

func (w *watcher) watch() (err error) {
	if w.w != nil {
		w.w.Close()
		w.w = nil
	}
	w.w, err = rfsnotify.NewWatcher()
	if err != nil {
		return
	}
	for _, input := range w.c.Inputs {
		if input.Recursive {
			w.w.AddRecursive(input.Path)
		} else {
			w.w.Add(input.Path)
		}
	}
	err = collect(w.c)
	if err != nil {
		return
	}
	for {
		select {
		case e := <-w.w.Events:
			if e.Op != fsnotify.Chmod {
				writeLog("%s changed; rebuilding %s", e.Name, w.c.Output)
				err = collect(w.c)
			}
		case err = <-w.w.Errors:
			return
		}
	}
	return
}
