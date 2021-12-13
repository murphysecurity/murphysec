package scanner

import (
	"fmt"
	"io/ioutil"
	"murphysec-cli-simple/util"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ProjectInfo struct {
	Name         string
	GitRemoteUrl string
	GitBranch    string
	ProjectType  string
}

func GetProjectInfo(targetAbsPath string) *ProjectInfo {
	gitConfigPath := filepath.Join(targetAbsPath, ".git", "config")
	if isGitProject(gitConfigPath) {
		var projectName, gitRemoteUrl, err = getProjectNameFromGit(gitConfigPath)
		if err == nil {
			return &ProjectInfo{
				Name:         projectName,
				GitRemoteUrl: gitRemoteUrl,
				GitBranch:    getGitBranch(targetAbsPath),
				ProjectType:  "git-remote",
			}
		}
	}
	return &ProjectInfo{
		Name:        getProjectNameFromPath(targetAbsPath),
		ProjectType: "local",
	}
}

func getProjectNameFromPath(Path string) string {
	parts := strings.Split(Path, string(os.PathSeparator))
	dirname := parts[len(parts)-1]
	return dirname
}

func getProjectNameFromGit(gitConfigPath string) (projectName string, gitRemoteUrl string, err error) {
	file, err := ioutil.ReadFile(gitConfigPath)
	if err != nil {
		return "", "", err
	}
	configFileContent := string(file)
	gitRemoteRegex := regexp.MustCompile(`(?m)url = ([\S]*)`)
	gitRemoteInfo := gitRemoteRegex.FindString(configFileContent)
	gitRemoteUrl = strings.Split(gitRemoteInfo, " ")[2]
	gitRemoteUrl = removeUserFromGitRemoteUrl(gitRemoteUrl)
	if strings.HasSuffix(gitRemoteUrl, ".git") {
		parts := strings.Split(gitRemoteUrl, "/")
		projectInfo := parts[len(parts)-1]
		projectName = strings.Replace(projectInfo, ".git", "", -1)
		return
	} else {
		err = fmt.Errorf("Failed to get git remote url ")
		return "", "", err
	}

}

// 针对http开头的git地址，去掉里面的用户名和密码[脱敏]
func removeUserFromGitRemoteUrl(gitRemoteUrl string) string {
	parseUrl, err := url.Parse(gitRemoteUrl)
	if err != nil {
		return gitRemoteUrl
	}
	if strings.HasPrefix("http", parseUrl.Scheme) {
		parseUrl.User = nil
		return parseUrl.String()
	}
	return gitRemoteUrl
}

// 通过git HEAD文件获取当前git分支信息
func getGitBranch(targetAbsPath string) string {
	gitHEADPath := filepath.Join(targetAbsPath, ".git", "HEAD")
	if gitHEADFileExist := util.IsPathExist(gitHEADPath); gitHEADFileExist == true {
		file, err := ioutil.ReadFile(gitHEADPath)
		if err != nil {
			return ""
		}
		gitHEADFileContent := string(file)
		if strings.HasPrefix(gitHEADFileContent, "ref:") {
			parts := strings.Split(gitHEADFileContent, "/")
			branch := parts[len(parts)-1]
			branch = strings.Replace(branch, "\n", "", -1)
			return branch
		}
	}
	return ""
}
func isGitProject(gitConfigPath string) bool {
	if gitConfigFileExist := util.IsPathExist(gitConfigPath); gitConfigFileExist == true {
		return true
	}
	return false
}
