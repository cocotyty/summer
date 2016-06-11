package gen

import (
	"gopkg.in/fsnotify.v1"
	"time"
	"regexp"
	"strings"
)

type Watcher struct {
	Ignores []*regexp.Regexp
	Root    string
	Callback func (string)
}

func (watcher *Watcher) Init(Root string,callback func(string)) {
	watcher.Ignores = []*regexp.Regexp{}
	watcher.Root = Root
	watcher.Callback=callback
}
func (Watcher *Watcher)Ignore(path string) {
	Watcher.Ignores = append(Watcher.Ignores, regexp.MustCompile(path))
}
func (watcher *Watcher)Loop() {
	for {
		watcher.work(watcher.Root)
	}
}
func (watcher *Watcher)work(name string) {
	w, _ := fsnotify.NewWatcher()
	defer w.Close()
	set := map[string]struct{}{}
	FindAllDir(name, set)
	for v, _ := range set {
		ignore:=false
		for _, i := range watcher.Ignores {
			if i.MatchString(strings.Replace(v[len(name):], L, "/", -1)) {
				ignore=true
				break
			}
		}
		if !ignore{
			w.Add(v)
		}
	}
	Listen:
	e := <-w.Events
	for _, v := range watcher.Ignores {
		if v.MatchString(strings.Replace(e.Name[len(name):], L, "/", -1)) {
			goto Listen
		}
	}
	time.Sleep(500 * time.Millisecond)
	go watcher.Callback(e.Name[:strings.LastIndex(e.Name,L)])
}
