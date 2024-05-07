package depp

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"sync"
)

var __indentPattern *regexp.Regexp
var __depPattern *regexp.Regexp
var __projectPattern *regexp.Regexp
var __taskPrefixPattern *regexp.Regexp
var __prepareProjectPattern *regexp.Regexp
var __projectDepPattern *regexp.Regexp
var __initPatterns sync.Once

func initPatterns() {
	__initPatterns.Do(func() {
		__prepareProjectPattern = regexp.MustCompile(`^-+$`)
		__taskPrefixPattern = regexp.MustCompile(`^[A-Za-z0-9.:_-]+\s+-`)
		__indentPattern = regexp.MustCompile(`^[ |/\\+-]+`)
		__projectPattern = regexp.MustCompile(`(?:[Rr]oot project|[Pp]roject)\s*([A-Za-z0-9:._-]+|'[A-Za-z0-9:._-]+')?`)
		__depPattern = regexp.MustCompile(`[A-Za-z0-9.-]+:[A-Za-z0-9.-]+:[A-Za-z0-9.-]+(?:\s+->\s+[A-Za-z0-9.-]+)?`)
		__projectDepPattern = regexp.MustCompile(`project [A-Za-z0-9:.-]+`)
	})
}

type TreeNode struct {
	G string
	A string
	V string
	C []TreeNode
}

func Parse(reader io.Reader, commit func(project string, task string, data []TreeNode)) (e error) {
	initPatterns()
	var scanner = bufio.NewScanner(reader)
	scanner.Buffer(nil, 4096)
	scanner.Split(bufio.ScanLines)
	var projectName string
	var projectNameSet bool
	var prepareProjectName bool
	var taskName string
	_ = taskName
	var record depRecord
	var doCollect = func() {
		if taskName == "" {
			return
		}
		commit(projectName, taskName, record.roots)
		taskName = ""
		record = depRecord{}
	}
	for scanner.Scan() {
		if e != nil {
			return
		}
		if scanner.Err() != nil {
			e = scanner.Err()
			return
		}
		var line = strings.TrimSpace(scanner.Text())

		// Match dependencies
		if taskName != "" {
			var cIndent = len(__indentPattern.FindString(line))
			if cIndent == 0 {
				doCollect()
			} else {
				var dep = __depPattern.FindString(line)
				var projDep = __projectDepPattern.FindString(line)
				if dep != "" {
					g, a, v := parseDepMatch(dep)
					if g != "" && a != "" {
						record.add(cIndent, g, a, v)
					}
				} else if projDep != "" {
					projDep = parseProjectDepMatch(projDep)
					if projDep != "" {
						record.add(cIndent, "", projDep, "")
					}
				}
			}
		}

		// Match project name
		{
			if __prepareProjectPattern.MatchString(line) {
				prepareProjectName = true
				continue
			}
			var projectMatch = __projectPattern.FindStringSubmatch(line)
			if prepareProjectName && len(projectMatch) > 0 {
				projectName = strings.Trim(projectMatch[1], "'")
				projectNameSet = true
				prepareProjectName = false
				continue
			}
			prepareProjectName = false
			if !projectNameSet {
				continue
			}
		}

		// Match task name
		{
			var taskPrefix = __taskPrefixPattern.FindString(line)
			if taskPrefix != "" {
				doCollect()
				taskName = strings.TrimSuffix(taskPrefix, "-")
				continue
			}
		}
	}
	return
}

func parseProjectDepMatch(s string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "project"))
}

func parseDepMatch(s string) (group, artifact, version string) {
	var split = strings.Split(s, "->")
	if len(split) > 2 {
		return
	}
	var alterVer string
	if len(split) == 2 {
		alterVer = strings.TrimSpace(split[1])
	}
	split = strings.Split(strings.TrimSpace(split[0]), ":")
	if len(split) >= 2 {
		group = split[0]
		artifact = split[1]
	}
	if len(split) >= 3 {
		version = split[2]
	}
	if alterVer != "" {
		version = alterVer
	}
	return
}

type depRecord struct {
	stack         []int
	deepListStack []*[]TreeNode
	roots         []TreeNode
}

func (d *depRecord) add(indent int, g, a, v string) {
	if len(d.stack) == 0 {
		//whereToAppend = &d.roots
		d.stack = []int{indent}
		d.deepListStack = []*[]TreeNode{&d.roots}
	} else if indent > d.stack[len(d.stack)-1] {
		// push
		d.stack = append(d.stack, indent)
		var currentList = *d.deepListStack[len(d.deepListStack)-1]
		var lastChildrenList = &currentList[len(currentList)-1].C
		d.deepListStack = append(d.deepListStack, lastChildrenList)
	} else if indent < d.stack[len(d.stack)-1] {
		// pop
		for indent < d.stack[len(d.stack)-1] && len(d.stack) > 0 {
			d.stack = d.stack[:len(d.stack)-1]
			d.deepListStack = d.deepListStack[:len(d.deepListStack)-1]
		}
	}
	if len(d.stack) == 0 {
		// re-init
		d.stack = []int{indent}
		d.deepListStack = []*[]TreeNode{&d.roots}
	}
	*d.deepListStack[len(d.deepListStack)-1] = append(*d.deepListStack[len(d.deepListStack)-1], TreeNode{
		G: g,
		A: a,
		V: v,
	})
}
