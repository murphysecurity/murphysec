package findfile

import (
	"container/list"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type Predication func(name string, dir string) bool

type Option struct {
	MaxDepth    int
	ExcludeFile bool
	ExcludeDir  bool
	Predication Predication
}

func FileNameRegexp(re *regexp.Regexp) Predication {
	return func(name string, dir string) bool {
		return re.MatchString(name)
	}
}

func Find(baseDir string, option Option) *FileIterator {
	type d struct {
		path  string
		depth int
	}
	iter := &FileIterator{
		cond:       *sync.NewCond(&sync.Mutex{}),
		list:       list.List{},
		terminated: false,
	}
	go func() {
		visited := map[string]struct{}{}
		q := list.New()
		q.PushBack(d{baseDir, 0})
		for q.Len() != 0 {
			curr := q.Front().Value.(d)
			q.Remove(q.Front())
			entries, e := os.ReadDir(curr.path)
			if e != nil {
				iter.append(e)
				iter.cond.Signal()
				continue
			}
			var rs []interface{}
			for _, entry := range entries {
				path := filepath.Join(curr.path, entry.Name())
				if _, ok := visited[path]; ok {
					continue
				}
				visited[path] = struct{}{}
				if entry.IsDir() && (option.MaxDepth == 0 || curr.depth < option.MaxDepth) {
					q.PushBack(d{path, curr.depth + 1})
				}
				if option.ExcludeFile && !entry.IsDir() {
					continue
				}
				if option.ExcludeDir && entry.IsDir() {
					continue
				}
				if option.Predication == nil || option.Predication(entry.Name(), curr.path) {
					rs = append(rs, path)
				}
			}
			iter.append(rs...)
			iter.cond.Signal()
		}
		iter.cond.L.Lock()
		defer iter.cond.L.Unlock()
		iter.terminated = true
		iter.cond.Broadcast()
	}()
	return iter
}

type FileIterator struct {
	cond       sync.Cond
	list       list.List
	terminated bool
	curr       interface{}
}

func (f *FileIterator) append(v ...interface{}) {
	f.cond.L.Lock()
	defer f.cond.L.Unlock()
	for _, it := range v {
		f.list.PushBack(it)
	}
}

func (f *FileIterator) Err() error {
	if e, ok := f.curr.(error); ok {
		return e
	}
	return nil
}

func (f *FileIterator) Path() string {
	if s, ok := f.curr.(string); ok {
		return s
	}
	return ""
}

func (f *FileIterator) Next() bool {
	f.cond.L.Lock()
	defer f.cond.L.Unlock()
	for f.list.Len() < 1 && !f.terminated {
		f.cond.Wait()
	}
	if f.terminated && f.list.Len() == 0 {
		return false
	}
	f.curr = f.list.Front().Value
	f.list.Remove(f.list.Front())
	return true
}
