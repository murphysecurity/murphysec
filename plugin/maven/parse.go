package maven

import (
	"fmt"
	"regexp"
	"strings"
)

func parseMvnCommandResult(cmdResult string) (interface{}, error) {
	errCharRegex := regexp.MustCompile(`(?m)^\[ERROR\]`)
	logLabelRegex := regexp.MustCompile(`(?m)\[\w+\]\s*`)
	errContentRegex := regexp.MustCompile(`(?m)\[ERROR\] Failed to execute.*`)

	if matched := errCharRegex.MatchString(cmdResult); matched == true {
		err := fmt.Errorf("%s", "The maven plugin execution failed ")
		errContents := errContentRegex.FindAllString(cmdResult, -1)
		if len(errContents) > 0 {
			errContent := errContents[0]
			err = fmt.Errorf("mvn scan failed:\n\n%s", errContent)
		}
		return nil, err
	}
	text := logLabelRegex.ReplaceAllString(cmdResult, "")
	return getRootProjectInfo(text)
}

func getRootProjectInfo(text string) (interface{}, error) {
	digraphRegex := regexp.MustCompile(`(?m)digraph([\s\S]*?)}`)
	projects := digraphRegex.FindAllString(text, -1)
	if len(projects) == 0 {
		err := fmt.Errorf("没有找到项目使用的开源组件信息")
		return nil, err
	}
	rootProject := getProjectInfo(projects[0])
	for i := 1; i < len(projects); i++ {
		subProjectInfo := getProjectInfo(projects[i])
		if (subProjectInfo != nil) && (subProjectInfo["name"] != nil) {
			subProjectName := subProjectInfo["name"].(string)
			rootProject["dependencies"].(map[string]interface{})[subProjectName] = subProjectInfo
		}
	}
	return rootProject, nil
}

func getPackageInfo(packageStr string) map[string]interface{} {
	parts := strings.Split(packageStr, ":")
	result := map[string]interface{}{
		"name":         parts[0] + ":" + parts[1],
		"version":      parts[3],
		"dependencies": make(map[string]interface{}),
	}
	if len(parts) > 4 {
		result["scope"] = parts[len(parts)-1]
		result["version"] = parts[len(parts)-2]
	}
	return result
}

func getProjectInfo(projectText string) map[string]interface{} {
	lines := strings.Split(projectText, "\n")
	projectTarget := getQuoteValue(lines[0])
	deps := make(map[string][]string)
	for i := 1; i < len(lines)-1; i++ {
		line := strings.Trim(lines[i], " ")
		if strings.Contains(line, "->") == false {
			continue
		}
		lineParts := strings.Split(line, "->")
		source := getQuoteValue(lineParts[0])
		target := getQuoteValue(lineParts[1])
		if deps[source] == nil {
			deps[source] = []string{target}
		} else {
			deps[source] = append(deps[source], target)
		}
	}
	return assembleProjectDeps(projectTarget, deps)
}

func assembleProjectDeps(source string, pkgDeps map[string][]string) map[string]interface{} {
	sourcePackage := getPackageInfo(source)
	if sourcePackage["scope"] == "test" {
		return nil
	}
	sourceDeps := pkgDeps[source]
	if sourceDeps == nil {
		return sourcePackage
	}
	for _, dep := range sourceDeps {
		depPackage := assembleProjectDeps(dep, pkgDeps)
		if depPackage["name"] != nil {
			depPackageName := depPackage["name"].(string)
			sourcePackage["dependencies"].(map[string]interface{})[depPackageName] = depPackage
		}
	}
	return sourcePackage
}

func getQuoteValue(str string) string {
	return str[strings.Index(str, "\"")+1 : strings.LastIndex(str, "\"")]
}
