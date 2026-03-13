package nuget

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/murphysecurity/murphysec/infra/pathignore"
	"github.com/murphysecurity/murphysec/utils"
)

// 判断项目中是否有packages.config文件
func checkPackagesIsExistence(fileName string) bool {
	return utils.IsFile(fileName)
}

func isDotnetBuildTarget(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".sln", ".csproj", ".fsproj", ".vbproj":
		return true
	default:
		return false
	}
}

func normalizeBuildTargets(candidates []string) []string {
	seen := map[string]struct{}{}
	var targets []string
	for _, p := range candidates {
		if !isDotnetBuildTarget(p) {
			continue
		}
		if !utils.IsFile(p) {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		targets = append(targets, p)
	}
	sort.Strings(targets)
	return targets
}

func findCLNList(dir string) (filePath []string, err error) {
	var slnPaths []string
	var projectPaths []string

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if pathignore.DirName(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".sln" {
			slnPaths = append(slnPaths, path)
			return nil
		}
		if isDotnetBuildTarget(path) {
			projectPaths = append(projectPaths, path)
		}
		return nil // 继续搜索
	})

	if err != nil {
		return nil, err
	}
	slnPaths = normalizeBuildTargets(slnPaths)
	if len(slnPaths) > 0 {
		return slnPaths, nil
	}
	projectPaths = normalizeBuildTargets(projectPaths)
	return projectPaths, nil
}

func slnHasDockerComposeProject(slnPath string) bool {
	data, err := os.ReadFile(slnPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "docker-compose.dcproj")
}

func splitSlnPathsByDockerCompose(slnPaths []string) (preferred []string, fallback []string) {
	for _, p := range slnPaths {
		if slnHasDockerComposeProject(p) {
			fallback = append(fallback, p)
		} else {
			preferred = append(preferred, p)
		}
	}
	return preferred, fallback
}

func validateBuildTargets(targets []string) error {
	for _, p := range targets {
		if !isDotnetBuildTarget(p) {
			return fmt.Errorf("invalid nuget build target: %s", p)
		}
		if !utils.IsFile(p) {
			return fmt.Errorf("nuget build target not found: %s", p)
		}
	}
	return nil
}

var projectReferencePattern = regexp.MustCompile(`(?is)<\s*ProjectReference\b[^>]*\bInclude\s*=\s*["']([^"']+)["']`)

func findMissingProjectReferences(projectPath string) ([]string, error) {
	ext := strings.ToLower(filepath.Ext(projectPath))
	switch ext {
	case ".csproj", ".fsproj", ".vbproj":
	default:
		return nil, nil
	}
	data, err := os.ReadFile(projectPath)
	if err != nil {
		return nil, err
	}
	matches := projectReferencePattern.FindAllStringSubmatch(string(data), -1)
	if len(matches) == 0 {
		return nil, nil
	}
	baseDir := filepath.Dir(projectPath)
	seen := map[string]struct{}{}
	var missing []string
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		ref := filepath.Clean(m[1])
		refPath := ref
		if !filepath.IsAbs(refPath) {
			refPath = filepath.Join(baseDir, refPath)
		}
		refPath = filepath.Clean(refPath)
		if utils.IsFile(refPath) {
			continue
		}
		if _, ok := seen[refPath]; ok {
			continue
		}
		seen[refPath] = struct{}{}
		missing = append(missing, refPath)
	}
	sort.Strings(missing)
	return missing, nil
}
