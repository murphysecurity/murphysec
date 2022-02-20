package maven

import (
	"bufio"
	"io"
	"murphysec-cli-simple/logger"
	"regexp"
	"sort"
	"strings"
)

const _MaxLineSize = 128 * 1024

var modulePattern = regexp.MustCompile("--+< +(.+?) +>--+")
var nodePattern = regexp.MustCompile("^(\\d+) +([^ ]+)")
var edgePattern = regexp.MustCompile("^(\\d+) (\\d+) (\\w+)")

func parseOutput(reader io.Reader) map[Coordinate][]Dependency {
	input := bufio.NewScanner(reader)
	input.Split(bufio.ScanLines)
	input.Buffer(make([]byte, _MaxLineSize), _MaxLineSize)

	deps := map[Coordinate][]Dependency{}

	var nodes map[string]string              // node id -> node
	var graph map[string]map[string]struct{} // adjust list
	// parse TGF format graph
	for input.Scan() {
		logger.Debug.Println("mvn line:", input.Text())
		line := strings.TrimSpace(strings.TrimPrefix(input.Text(), "[INFO]"))
		if moduleLine := modulePattern.FindStringSubmatch(line); moduleLine != nil {
			// convert adjust list to dependency tree map
			if len(nodes) != 0 {
				if tree := convertAdjustListToDepTree(graph, nodes); tree != nil {
					deps[tree.Coordinate] = tree.Children
					logger.Debug.Println("Collect module:", tree.Coordinate)
				}
			}
			nodes = map[string]string{}
			graph = map[string]map[string]struct{}{}
			continue
		}
		if nodes != nil {
			if edgeLine := edgePattern.FindStringSubmatch(line); edgeLine != nil {
				leftId := edgeLine[1]
				rightId := edgeLine[2]
				if graph[leftId] == nil {
					graph[leftId] = map[string]struct{}{}
				}
				graph[leftId][rightId] = struct{}{}
				continue
			}
			if depLine := nodePattern.FindStringSubmatch(line); depLine != nil {
				nodeId := depLine[1]
				nodeT := depLine[2]
				if old, ok := nodes[nodeId]; ok {
					logger.Info.Println("Found repeat node in maven TGF output, replace:", nodeId, nodeT, old)
				}
				nodes[nodeId] = nodeT
				continue
			}
		}
	}
	// tail module
	if len(nodes) != 0 {
		if tree := convertAdjustListToDepTree(graph, nodes); tree != nil {
			deps[tree.Coordinate] = tree.Children
			logger.Debug.Println("Collect module:", tree.Coordinate)
		}
	}
	return deps
}

func convertAdjustListToDepTree(graph map[string]map[string]struct{}, nodeMapping map[string]string) *Dependency {
	// find head node
	var headNode string
	{
		nodeSet := map[string]struct{}{}
		// copy map
		for s := range nodeMapping {
			nodeSet[s] = struct{}{}
		}
		// iterate the graph, find head node
	outer:
		for _, m := range graph {
			for s := range m {
				if len(nodeSet) < 2 {
					break outer
				}
				delete(nodeSet, s)
			}
		}
		if len(nodeSet) == 1 {
			for s := range nodeSet {
				headNode = s
			}
		}
	}
	if headNode == "" {
		logger.Info.Println("head node not found in the graph")
		return nil
	}
	depObjMap := map[string]Dependency{}
	{
		// mapping node to depType
		for id, t := range nodeMapping {
			a := strings.Split(t, ":")
			if len(a) < 4 {
				continue
			}
			depObjMap[id] = Dependency{
				Coordinate: Coordinate{
					GroupId:    a[0],
					ArtifactId: a[1],
					Version:    a[3],
				},
				Children: nil,
			}
		}
	}

	d := _recursiveBuildMap(depObjMap, graph, headNode, nil, 0)
	return &d
}

func _recursiveBuildMap(nodes map[string]Dependency, graph map[string]map[string]struct{}, curNode string, circleDetect map[string]int, depth int) Dependency {
	if circleDetect == nil {
		circleDetect = map[string]int{}
	}
	cur := nodes[curNode]
	if _, ok := circleDetect[curNode]; ok {
		var idPath []string
		for s := range circleDetect {
			idPath = append(idPath, s)
		}
		sort.Slice(idPath, func(i, j int) bool {
			return circleDetect[idPath[i]] < circleDetect[idPath[j]]
		})
		logger.Info.Println("Circle detected in mvn dep tree:", strings.Join(idPath, "->"))
		return cur
	}
	circleDetect[curNode] = depth
	defer delete(circleDetect, curNode)
	for it := range graph[curNode] {
		cur.Children = append(cur.Children, _recursiveBuildMap(nodes, graph, it, circleDetect, depth+1))
	}
	return cur
}
