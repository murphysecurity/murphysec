package inspector

import (
	"io/ioutil"
	"murphysec-cli-simple/logger"
	"path/filepath"
)

type ScanTreeNode struct {
	Dir         string         `json:"dir"`
	GitInfo     *GitInfo       `json:"git_info"`
	IsProject   bool           `json:"is_project"`
	Children    []ScanTreeNode `json:"children"`
	IsContainer bool           `json:"is_container"`
}

var dirBlackList = map[string]bool{
	"node_modules": true,
	".gradle":      true,
	"gradle":       true,
	".m2":          true,
	".idea":        true,
	".git":         true,
}

type DirScanOption struct {
	ExcludeDirs map[string]bool
}

func _scan(dir string, visited map[string]bool, depthLimit int, option DirScanOption) *ScanTreeNode {
	if depthLimit < 0 {
		return nil
	}
	logger.Debug.Println("Visit dir:", dir)
	if option.ExcludeDirs[dir] {
		logger.Info.Println("Dir in exclude set, skip")
		return nil
	}
	node := &ScanTreeNode{Dir: dir}
	if dirBlackList[filepath.Base(dir)] {
		logger.Debug.Println("Dir name in blacklist, skip")
		return nil
	}
	if visited[dir] {
		logger.Debug.Println("Dir visited, skip.")
		return nil
	}
	visited[dir] = true

	// try git
	gitInfo, e := getGitInfo(dir)
	if e != nil {
		logger.Warn.Println("Dir is git repo, but get git info failed.", e.Error())
		node.IsProject = true
		return node
	}
	if gitInfo != nil {
		logger.Debug.Println("Get git info succeed.")
		node.IsProject = true
		node.GitInfo = gitInfo
		return node
	}
	logger.Debug.Println("Not a git repo, enumerating sub dirs...")
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		logger.Warn.Println("Enumerating sub dirs failed.", e.Error())
		return nil
	}
	for _, it := range files {
		if !it.IsDir() {
			continue
		}
		if n := _scan(filepath.Join(dir, it.Name()), visited, depthLimit-1, option); n != nil {
			node.Children = append(node.Children, *n)
		}
	}
	// if any children is a project, current node is a container;
	// if any children is a container, current node is a container also.
	for _, it := range node.Children {
		node.IsContainer = node.IsContainer || it.IsContainer || it.IsProject
	}
	// if current node is a container, any non-container children should be a project.
	if node.IsContainer {
		for i := range node.Children {
			if !node.Children[i].IsContainer {
				node.Children[i].IsProject = true
			}
		}
	}
	return node
}

func ScanDir(dir string, depthLimit int, option DirScanOption) *ScanTreeNode {
	return _scan(dir, map[string]bool{}, depthLimit, option)
}
