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
	*rfsnotify.RWatcher
	c *config
}

func newWatcher(c *config) *watcher {
	return &watcher{c: c}
}

func (w *watcher) watch() (err error) {
	if w.RWatcher != nil {
		w.Close()
		w.RWatcher = nil
	}
	w.RWatcher, err = rfsnotify.NewWatcher()
	if err != nil {
		return
	}
	for _, input := range w.c.inputs {
		if input.recursive {
			w.AddRecursive(input.path)
		} else {
			w.Add(input.path)
		}
	}
	err = collect(w.c)
	if err != nil {
		return
	}
	for {
		select {
		case e := <-w.Events:
			if e.Op != fsnotify.Chmod {
				writeLog("%s changed; rebuilding %s", e.Name, w.c.output)
				err = collect(w.c)
			}
		case err = <-w.Errors:
			return
		}
	}
	return
}
