package sbt

import (
	"bufio"
	"context"
	"fmt"
	list "github.com/bahlo/generic-list-go"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/logpipe"
	"io"
	"os/exec"
	"regexp"
)

func sbtDependencyTree(ctx context.Context, dir string) ([]Dep, error) {
	var logger = logctx.Use(ctx)
	c := exec.CommandContext(ctx, "sbt", "-Dsbt.ci=true", "-Dsbt.color=never", "-Dsbt.progress=never", "-Dsbt.log.noformat=true", "-Dsbt.supershell=false", "dependencyTree")
	logger.Sugar().Infof("Execute command: %s", c)
	c.Dir = dir
	parser := newSbtDependencyTreeOutputParser()
	defer parser.Close()
	logAppender := logpipe.New(logger, "sbt")
	defer logAppender.Close()
	c.Stdout = io.MultiWriter(logAppender, parser)
	c.Stderr = logAppender
	if e := c.Run(); e != nil {
		return nil, fmt.Errorf("runSbt: %w", e)
	}
	parser.Close()
	root, e := parser.Result()
	if e != nil {
		return nil, e
	}
	return root.Children, nil
}

type sbtDependencyTreeOutputParser struct {
	io.WriteCloser
	err     error
	closeCh chan struct{}
	root    *Dep
}

func (s *sbtDependencyTreeOutputParser) Result() (*Dep, error) {
	<-s.closeCh
	return s.root, s.err
}

var treePattern = regexp.MustCompile(`^\[info] ([ +-]*)([\w+.-]+:[\w+.-]+):([\w+.-]+)(?: \[\w+])?`)

func newSbtDependencyTreeOutputParser() (r *sbtDependencyTreeOutputParser) {
	reader, writer := io.Pipe()

	r = &sbtDependencyTreeOutputParser{
		WriteCloser: writer,
		closeCh:     make(chan struct{}),
		root:        &Dep{},
	}

	go func() {
		defer close(r.closeCh)
		type item struct {
			*Dep
			indent int
		}
		q := list.New[item]()
		q.PushBack(item{
			Dep:    &Dep{},
			indent: 0,
		})
		r.root = q.Back().Value.Dep

		var scanner = bufio.NewScanner(reader)
		scanner.Buffer(make([]byte, 4096), 4096)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			if e := scanner.Err(); e != nil {
				r.err = e
				break
			}
			if q.Len() == 0 {
				r.err = fmt.Errorf("bad indent")
				break
			}
			var line = scanner.Text()
			var m = treePattern.FindStringSubmatch(line)
			if m == nil {
				// invalid line, unwind to root
				for q.Len() > 1 {
					q.Remove(q.Back())
				}
				continue
			}
			var indent = len(m[1])
			var name = m[2]
			var version = m[3]
			if indent == q.Back().Value.indent {
				// append as child of current node
				q.Back().Value.Children = append(q.Back().Value.Children, Dep{
					Name:    name,
					Version: version,
				})
				continue
			}
			if indent > q.Back().Value.indent {
				// use last children as new stack top
				if len(q.Back().Value.Children) == 0 {
					r.err = fmt.Errorf("bad indent")
					break
				}
				lastChildren := &q.Back().Value.Children[len(q.Back().Value.Children)-1]
				lastChildren.Children = append(lastChildren.Children, Dep{
					Name:    name,
					Version: version,
				})
				q.PushBack(item{lastChildren, indent})
				continue
			}
			if indent < q.Back().Value.indent {
				// unwind until indent match
				for q.Back().Value.indent > indent {
					if q.Len() == 0 {
						r.err = fmt.Errorf("bad indent")
						break
					}
					q.Remove(q.Back())
				}
			}
		}
	}()

	return r
}
