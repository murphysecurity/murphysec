package npm

import (
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"strings"
)

type v3Lockfile struct {
	Name            string               `json:"name"`
	Version         string               `json:"version"`
	LockfileVersion int                  `json:"lockfileVersion"`
	Packages        map[string]v3Package `json:"packages"`
}

type v3Package struct {
	Name                 string            `json:"name"`
	Version              string            `json:"version"`
	Dependencies         map[string]string `json:"dependencies"`
	DevDependencies      map[string]string `json:"devDependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
	Dev                  bool              `json:"dev"`
}

func processLockfileV3(data []byte) (r *model.DependencyItem, e error) {
	var lockfile v3Lockfile
	if e := json.Unmarshal(data, &lockfile); e != nil {
		return nil, fmt.Errorf("parse lockfile failed: %w", e)
	}
	if lockfile.LockfileVersion != 3 {
		return nil, fmt.Errorf("unsupported lockfile version: %d", lockfile.LockfileVersion)
	}
	if lockfile.Packages == nil {
		lockfile.Packages = make(map[string]v3Package)
	}
	var handler visitV3Handler[*model.DependencyItem] = func(theValue *model.DependencyItem, pred, succ [2]string, isDev bool, doNext func(nextValue *model.DependencyItem)) {
		var dep = model.DependencyItem{
			Component: model.Component{CompName: succ[0], CompVersion: succ[1], EcoRepo: EcoRepo},
		}
		dep.IsOnline.SetOnline(!isDev)
		doNext(&dep)
		theValue.Dependencies = append(theValue.Dependencies, dep)
	}
	var rootNode model.DependencyItem
	_visitV3[*model.DependencyItem](&lockfile, nil, nil, make(map[string]struct{}), &rootNode, handler, make(map[string]struct{}))
	rootNode.CompName = lockfile.Name
	rootNode.CompVersion = lockfile.Version
	for i := range rootNode.Dependencies {
		rootNode.Dependencies[i].IsDirectDependency = true
	}
	return &rootNode, nil
}

type visitV3Handler[T any] func(theValue T, pred, succ [2]string, isDev bool, doNext func(nextValue T))

func _visitV3[T any](lockfile *v3Lockfile, pred *v3Package, rPath []string, pathVisited map[string]struct{}, theValue T, handler visitV3Handler[T], pruneSet map[string]struct{}) {
	if pred == nil {
		var root, ok = lockfile.Packages[""]
		if ok {
			pred = &root
		} else {
			return
		}
	}
	var predV = [2]string{pred.Name, pred.Version}
	var traversalDependencies = func(isDev bool, m map[string]string) {
		for succName := range m {
			var succPath, succSegments, succ, ok = npmFindCorrectPair(lockfile.Packages, rPath, succName)
			if !ok {
				continue
			}
			if _, ok := pathVisited[succPath]; ok {
				continue
			}
			pathVisited[succPath] = struct{}{}
			var succV = [2]string{succ.Name, succ.Version}
			if succV[0] == "" {
				succV[0] = succName
			}
			var doNext = func(v T) { _visitV3(lockfile, &succ, succSegments, pathVisited, v, handler, pruneSet) }
			if pruneSet != nil {
				if _, ok := pruneSet[succPath]; ok {
					doNext = func(v T) {}
				}
				pruneSet[succPath] = struct{}{}
			}
			handler(theValue, predV, succV, isDev, doNext)
			delete(pathVisited, succPath)
		}
	}
	traversalDependencies(false, pred.Dependencies)
	traversalDependencies(false, pred.OptionalDependencies)
	traversalDependencies(true, pred.DevDependencies)
}

func npmFindCorrectPair[U any](m map[string]U, rp []string, name string) (key string, segments []string, value U, ok bool) {
	rp = utils.CopySlice(rp)
	for {
		segments = append(rp, "node_modules", name)
		key = strings.Join(segments, "/")
		value, ok = m[key]
		if ok {
			return
		}
		segments = append(rp, name)
		key = strings.Join(segments, "/")
		value, ok = m[key]
		if ok {
			return
		}
		if len(rp) == 0 {
			break
		}
		rp = rp[:len(rp)-1]
	}
	key = ""
	segments = nil
	return
}
