package inspector

import (
	"container/list"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var fileHashScanBlackList = map[string]bool{
	"node_modules": true,
}

func FileScan(baseDir string) *FilePathIterator {
	iter := NewFilePathIterator()
	go func() {
		defer iter.Terminate()
		m := map[string]struct{}{}
		q := list.New()
		q.PushBack(baseDir)
		for q.Len() > 0 {
			cur := q.Front().Value.(string)
			q.Remove(q.Front())

			// 去重
			if _, ok := m[cur]; ok {
				continue
			} else {
				m[cur] = struct{}{}
			}

			// 遍历
			entries, e := os.ReadDir(cur)
			if e != nil {
				iter.PushErr(e)
				continue
			}
			for _, entry := range entries {
				fileName := entry.Name()
				filePath := filepath.Join(cur, entry.Name())
				if strings.HasPrefix(fileName, ".") || fileHashScanBlackList[fileName] {
					continue
				}
				if entry.IsDir() {
					q.PushBack(filePath)
				} else {
					iter.PushPath(filePath)
				}
			}
		}
	}()
	return iter
}

type FilePathIterator struct {
	s    sync.Cond
	q    *list.List
	done bool
	curr filePathIteratorEntry
}

type filePathIteratorEntry struct {
	e error
	p string
}

func NewFilePathIterator() *FilePathIterator {
	return &FilePathIterator{
		s:    *sync.NewCond(&sync.Mutex{}),
		q:    list.New(),
		done: false,
		curr: filePathIteratorEntry{},
	}
}

func (f *FilePathIterator) PushErr(e error) {
	if e == nil {
		panic("e is nil")
	}
	f.push(filePathIteratorEntry{e: e})
}

func (f *FilePathIterator) PushPath(s string) {
	f.push(filePathIteratorEntry{p: s})
}

func (f *FilePathIterator) push(entry filePathIteratorEntry) {
	f.s.L.Lock()
	defer f.s.L.Unlock()
	if f.done {
		panic("done")
	}
	f.q.PushBack(entry)
	f.s.Signal()
}

func (f *FilePathIterator) Terminate() {
	f.s.L.Lock()
	defer f.s.L.Unlock()
	if f.done {
		panic("done")
	}
	f.done = true
	f.s.Broadcast()
}

func (f *FilePathIterator) Next() bool {
	f.s.L.Lock()
	defer f.s.L.Unlock()
	for f.q.Len() == 0 && !f.done {
		f.s.Wait()
	}
	if f.q.Len() > 0 {
		f.curr = f.q.Front().Value.(filePathIteratorEntry)
		f.q.Remove(f.q.Front())
		return true
	}
	return false
}

func (f *FilePathIterator) Err() error {
	f.s.L.Lock()
	defer f.s.L.Unlock()
	return f.curr.e
}

func (f *FilePathIterator) Path() string {
	f.s.L.Lock()
	defer f.s.L.Unlock()
	return f.curr.p
}
