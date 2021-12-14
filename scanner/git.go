package scanner

import (
	"fmt"
	"io/ioutil"
	"murphysec-cli-simple/util"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

// file: todo: rewrite

func getProjectNameFromGitConfigFile(gitConfigPath string) (projectName string, gitRemoteUrl string, err error) {
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
func isGitProject(projectDir string) bool {
	p := filepath.Join(projectDir, ".git", "config")
	return util.IsPathExist(p) && !util.IsDir(p)
}
