package bundler

import (
	"bytes"
	"container/list"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type node struct {
	line     string
	children []node
}

func (n node) get(path ...string) *node {
	cur := n
o:
	for _, p := range path {
		for _, it := range cur.children {
			if it.line == p {
				cur = it
				continue o
			}
		}
		return nil
	}
	return &cur
}

type symbol int

const (
	_SEnter symbol = iota + 1
	_SExit
)

type _GemLockTokens []any

func (t _GemLockTokens) String() string {
	b := new(bytes.Buffer)
	c := 0
	b.WriteString("<<<\n")
	for _, it := range t {
		if s, ok := it.(string); ok {
			for i := 0; i < c; i++ {
				b.WriteString("  ")
			}
			b.WriteString(s)
			b.WriteRune('\n')
		}
		if s, ok := it.(symbol); ok {
			if s == _SEnter {
				c++
			}
			if s == _SExit {
				c--
			}
		}
	}
	return b.String()
}

func lexGemLock(input string) (_GemLockTokens, error) {
	var rs []any
	q := list.New()
	q.PushBack(0)
	for lineNum, line := range strings.Split(input, "\n") {
		indent := _calcIndent(line)
		lastIndent := q.Back().Value.(int)
		if indent == lastIndent {
			rs = append(rs, strings.TrimSpace(line))
			continue
		}
		if indent > lastIndent {
			// enter
			q.PushBack(indent)
			rs = append(rs, _SEnter)
			rs = append(rs, strings.TrimSpace(line))
			continue
		}
		if indent < lastIndent {
			// back
			for {
				if q.Len() == 0 {
					return nil, errors.WithMessage(ErrBadIndent, fmt.Sprintf("Line %d", lineNum))
				}
				t := q.Back().Value.(int)
				if indent == t {
					break
				}
				if indent > t {
					return nil, errors.WithMessage(ErrBadIndent, fmt.Sprintf("Line %d", lineNum))
				}
				rs = append(rs, _SExit)
				q.Remove(q.Back())
			}
			rs = append(rs, strings.TrimSpace(line))
		}
	}
	return rs, nil
}

func parseGemLock(input string) (*node, error) {
	root := &node{}
	q := list.New()
	q.PushBack(root)
	tokens, e := lexGemLock(input)
	if e != nil {
		return nil, e
	}
	for _, token := range tokens {
		if q.Len() == 0 {
			return nil, errors.Wrap(ErrParseFail, "stack is empty")
		}
		top := q.Back().Value.(*node)
		if s, ok := token.(string); ok {
			top.children = append(top.children, node{line: s})
			continue
		}
		if s, ok := token.(symbol); ok {
			if s == _SEnter {
				if len(top.children) == 0 {
					return nil, errors.Wrap(ErrParseFail, "unexpected _SEnter")
				}
				q.PushBack(&top.children[len(top.children)-1])

			} else if s == _SExit {
				q.Remove(q.Back())
			}
		}
	}
	return root, nil
}

func _calcIndent(s string) int {
	return len(s) - len(strings.TrimLeft(s, " \t"))
}

func getDepGraph(input string) ([]model.DependencyItem, error) {
	// ([\w.-]+)\s*\(([\w.-]+)\)
	// catch group: name, version
	var pattern = regexp.MustCompile(`([\w.-]+)\s*\(([\w.-]+)\)`)
	var p2 = regexp.MustCompile(`^[\w.-]+`)
	tree, e := parseGemLock(input)
	if e != nil {
		return nil, e
	}
	specs := tree.get("GEM", "specs:")
	if specs == nil {
		return nil, errors.New("No graph")
	}
	// build graph
	graph := map[string][]string{}
	versionMap := map[string]string{}
	for _, left := range specs.children {
		m := pattern.FindStringSubmatch(left.line)
		if m == nil {
			continue
		}
		name := m[1]
		version := m[2]
		versionMap[name] = version
		for _, right := range left.children {
			rightName := p2.FindString(right.line)
			if rightName == "" {
				continue
			}
			graph[name] = append(graph[name], rightName)
		}
	}
	// find root node
	rm := map[string]struct{}{}
	for left := range graph {
		rm[left] = struct{}{}
	}
	for _, rights := range graph {
		for _, right := range rights {
			delete(rm, right)
		}
	}
	var roots []string
	for left := range rm {
		roots = append(roots, left)
	}
	// build tree
	var rs []model.DependencyItem
	for _, it := range roots {
		if t := _buildCompTree(graph, versionMap, it, map[string]struct{}{}); t != nil {
			rs = append(rs, *t)
		}
	}
	return rs, nil
}

func _buildCompTree(graph map[string][]string, versionMap map[string]string, target string, visited map[string]struct{}) *model.DependencyItem {
	// avoid cycling
	if _, ok := visited[target]; ok {
		return nil
	}
	visited[target] = struct{}{}
	defer delete(visited, target)

	d := &model.DependencyItem{
		Component: model.Component{
			CompName:    target,
			CompVersion: versionMap[target],
			EcoRepo:     EcoRepo,
		},
		Dependencies: nil,
	}
	for _, tt := range graph[target] {
		r := _buildCompTree(graph, versionMap, tt, visited)
		if r == nil {
			continue
		}
		d.Dependencies = append(d.Dependencies, *r)
	}
	return d
}
