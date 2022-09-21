package sbt

import (
	"bufio"
	"context"
	"fmt"
	list "github.com/bahlo/generic-list-go"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"io"
	"os/exec"
	"regexp"
)

func sbtDependencyTree(ctx context.Context, dir string) ([]model.Dependency, error) {
	var logger = utils.UseLogger(ctx)
	c := exec.CommandContext(ctx, "sbt", "-Dsbt.ci=true", "-Dsbt.color=never", "-Dsbt.progress=never", "-Dsbt.log.noformat=true", "-Dsbt.supershell=false")
	logger.Sugar().Infof("Execute command: %s", c)
	c.Dir = dir
	parser := newSbtDependencyTreeOutputParser()
	defer parser.Close()
	logAppender := utils.NewLogPipe(logger, "sbt")
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
	return root.Dependencies, nil
}

type sbtDependencyTreeOutputParser struct {
	io.WriteCloser
	err     error
	closeCh chan struct{}
	root    *model.Dependency
}

func (s *sbtDependencyTreeOutputParser) Result() (*model.Dependency, error) {
	<-s.closeCh
	return s.root, s.err
}

var treePattern = regexp.MustCompile("^\\[info] ([ +-]*)([\\w+.-]+:[\\w+.-]+):([\\w+.-]+)(?: \\[\\w+])?")

func newSbtDependencyTreeOutputParser() (r *sbtDependencyTreeOutputParser) {
	reader, writer := io.Pipe()

	r = &sbtDependencyTreeOutputParser{
		WriteCloser: writer,
		closeCh:     make(chan struct{}),
		root:        &model.Dependency{},
	}

	go func() {
		defer close(r.closeCh)
		type item struct {
			*model.Dependency
			indent int
		}
		q := list.New[item]()
		q.PushBack(item{
			Dependency: &model.Dependency{},
			indent:     0,
		})
		r.root = q.Back().Value.Dependency

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
				q.Back().Value.Dependencies = append(q.Back().Value.Dependencies, model.Dependency{
					Name:    name,
					Version: version,
				})
				continue
			}
			if indent > q.Back().Value.indent {
				// use last children as new stack top
				if len(q.Back().Value.Dependencies) == 0 {
					r.err = fmt.Errorf("bad indent")
					break
				}
				lastChildren := &q.Back().Value.Dependencies[len(q.Back().Value.Dependencies)-1]
				lastChildren.Dependencies = append(lastChildren.Dependencies, model.Dependency{
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
