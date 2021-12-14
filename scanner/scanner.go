package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

type ProjectInfo struct {
	ProjectName  string `json:"name"`
	GitRemoteUrl string `json:"git_remote_url"`
	GitBranch    string `json:"git_branch"`
	ProjectType  string `json:"project_type"`
}

func GetProjectInfo(projectDir string) *ProjectInfo {
	gitConfigPath := filepath.Join(projectDir, ".git", "config")
	if isGitProject(projectDir) {
		var projectName, gitRemoteUrl, err = getProjectNameFromGitConfigFile(gitConfigPath)
		if err == nil {
			return &ProjectInfo{
				ProjectName:  projectName,
				GitRemoteUrl: gitRemoteUrl,
				GitBranch:    getGitBranch(projectDir),
				ProjectType:  "git-remote",
			}
		}
	}
	return &ProjectInfo{
		ProjectName: getProjectNameFromPath(projectDir),
		ProjectType: "local",
	}
}

func getProjectNameFromPath(Path string) string {
	parts := strings.Split(Path, string(os.PathSeparator))
	dirname := parts[len(parts)-1]
	return dirname
}
