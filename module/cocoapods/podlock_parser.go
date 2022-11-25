package cocoapods

import (
	"container/list"
	"encoding/json"
	"github.com/murphysecurity/murphysec/model"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type symbol int

const (
	_SEnter symbol = iota + 1
	_SExit
)

type pdTokens []any

type intStack []int

func (t *intStack) push(i int) {
	*t = append(*t, i)
}

func (t *intStack) empty() bool {
	return len(*t) == 0
}

func (t *intStack) peek() int {
	return (*t)[len(*t)-1]
}

func (t *intStack) pop() int {
	r := t.peek()
	*t = (*t)[:len(*t)-1]
	return r
}

func tokenizePodLocks(input string) (pdTokens, error) {
	var rs pdTokens
	sk := &intStack{}
	sk.push(0)
	lines := strings.Split(input, "\n")
	pos := 0
	for pos < len(lines) {
		if sk.empty() {
			return nil, errors.New("Bad indent")
		}
		if l := strings.TrimSpace(lines[pos]); l == "" || strings.HasPrefix(l, "<<<<") || strings.HasPrefix(l, ">>>>") || strings.HasPrefix(l, "====") {
			// empty line or git conflict
			pos++
			continue
		}
		indent := len(lines[pos]) - len(strings.TrimLeft(lines[pos], " "))
		if indent > sk.peek() {
			// enter
			rs = append(rs, _SEnter)
			sk.push(indent)
			continue
		}
		if indent < sk.peek() {
			// exit
			rs = append(rs, _SExit)
			sk.pop()
			continue
		}
		rs = append(rs, strings.TrimSpace(lines[pos]))
		pos++
	}
	return rs, nil
}

func getDepFromLock(input string) ([]model.DependencyItem, error) {
	tree, e := parse(input)
	if e != nil {
		return nil, e
	}
	namePattern := regexp.MustCompile(`[\w\\\/.-]+`)
	// ([\w.\\\/-]+)\s*\(([\w.\\\/-]+)
	// catch group: name, version
	nvPattern := regexp.MustCompile(`([\w.\\\/-]+)\s*\(([\w.\\\/-]+)`)

	// finding all direct dependencies
	var directDep []string
	if n := tree.get("DEPENDENCIES:"); n != nil {
		for _, it := range n.children {
			m := namePattern.FindString(strings.Trim(it.text, " \t-"))
			if m == "" {
				continue
			}
			directDep = append(directDep, m)
		}
	}

	// find all edges & resolved version
	versionMap := map[string]string{}
	graph := map[string][]string{}
	if n := tree.get("PODS:"); n != nil {
		for _, left := range n.children {
			m := nvPattern.FindStringSubmatch(strings.TrimLeft(left.text, " -\""))
			if m == nil {
				continue
			}
			leftName := m[1]
			leftVersion := m[2]
			versionMap[leftName] = leftVersion
			for _, right := range left.children {
				rightName := namePattern.FindString(strings.TrimLeft(right.text, " -\""))
				if rightName == "" {
					continue
				}
				graph[leftName] = append(graph[leftName], rightName)
			}
		}
	}
	var rs []model.DependencyItem
	for _, it := range directDep {
		t := _buildTree(graph, versionMap, map[string]struct{}{}, it)
		if t == nil {
			continue
		}
		rs = append(rs, *t)
	}
	return rs, nil
}

func _buildTree(graph map[string][]string, versions map[string]string, visited map[string]struct{}, target string) *model.DependencyItem {
	if _, ok := visited[target]; ok {
		return nil
	}
	visited[target] = struct{}{}
	defer delete(visited, target)

	r := &model.DependencyItem{
		Component: model.Component{
			CompName:    target,
			CompVersion: versions[target],
			EcoRepo:     EcoRepo,
		},
		Dependencies: nil,
	}
	for _, it := range graph[target] {
		t := _buildTree(graph, versions, visited, it)
		if t == nil {
			continue
		}
		r.Dependencies = append(r.Dependencies, *t)
	}
	return r
}

func parse(input string) (*node, error) {
	tokens, e := tokenizePodLocks(input)
	if e != nil {
		return nil, e
	}
	root := &node{}
	q := list.New()
	q.PushBack(root)
	for _, token := range tokens {
		if q.Len() == 0 {
			return nil, errors.New("stack is empty")
		}
		top := q.Back().Value.(*node)
		if s, ok := token.(string); ok {
			top.children = append(top.children, node{text: s})
		}
		if token == _SEnter {
			if len(top.children) == 0 {
				return nil, errors.New("unexpected _SEnter")
			}
			q.PushBack(&top.children[len(top.children)-1])
		} else if token == _SExit {
			q.Remove(q.Back())
		}
	}
	return root, nil
}

type node struct {
	text     string
	children []node
}

func (t node) get(p ...string) *node {
	cur := t
o:
	for _, n := range p {
		for _, it := range cur.children {
			if it.text == n {
				cur = it
				continue o
			}
		}
		return nil
	}
	return &cur
}

func (n node) MarshalJSON() ([]byte, error) {
	m := map[string]node{}
	for _, it := range n.children {
		m[it.text] = it
	}
	return json.Marshal(m)
}
