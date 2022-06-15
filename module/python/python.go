package python

import (
	"bufio"
	"context"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var pyRequirementsPattern = regexp.MustCompile("^([A-Za-z0-9_-]+) *.?= *([^= \\n\\r]+)$")
var pyImportPattern1 = regexp.MustCompile("import\\s+(?:[A-Za-z_-][A-Za-z_0-9.-]*)(?:\\s*,\\s*(?:[A-Za-z_-][A-Za-z_0-9.-]*))")
var pyImportPattern2 = regexp.MustCompile("from\\s+([A-Za-z_-][A-Za-z_0-9-]*)")

type Inspector struct{}

func (i Inspector) String() string {
	return "PythonInspector"
}

func (i Inspector) CheckDir(dir string) bool {
	r, e := os.ReadDir(dir)
	if e == nil {
		for _, it := range r {
			if filepath.Ext(it.Name()) == ".py" || strings.HasPrefix(it.Name(), "requirements") || it.Name() == "pyproject.toml" {
				return true
			}
		}
	}
	return false
}

// returns package name list
func parsePyImport(input string) []string {
	var rs []string
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "import ") {
		// import aa, bb.cc
		for _, it := range strings.Split(strings.TrimPrefix(pyImportPattern1.FindString(input), "import"), ",") {
			it = strings.TrimSpace(it)
			s := strings.Split(it, ".")[0]
			if s != "" {
				rs = append(rs, s)
			}
		}
	}
	if strings.HasPrefix(input, "from ") {
		if m := pyImportPattern2.FindStringSubmatch(input); m != nil {
			rs = append(rs, m[1])
		}
	}
	return rs
}

func (i Inspector) InspectProject(ctx context.Context) error {
	dir := model.UseInspectorTask(ctx).ScanDir
	componentMap := map[string]string{}
	requirementsFiles := map[string]struct{}{}
	ignoreSet := map[string]struct{}{}
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			ignoreSet[d.Name()] = struct{}{}
			return nil
		}
		if (filepath.Ext(path) == ".txt" || filepath.Ext(path) == "") && strings.HasPrefix(d.Name(), "requirements") {
			requirementsFiles[path] = struct{}{}
			return nil
		}
		if filepath.Ext(path) != ".py" {
			return nil
		}
		f, e := os.Open(path)
		if e != nil {
			return e
		}
		defer f.Close()
		scanner := bufio.NewScanner(io.LimitReader(f, 4*1024*1024))
		scanner.Split(bufio.ScanLines)
		scanner.Buffer(make([]byte, 16*1024), 16*1024)
		for scanner.Scan() {
			if scanner.Err() != nil {
				return nil
			}
			t := strings.TrimSpace(scanner.Text())
			for _, pkg := range parsePyImport(t) {
				if pyPkgBlackList[pkg] {
					continue
				}
				componentMap[pkg] = ""
			}
		}
		return nil
	})
	for fp := range requirementsFiles {
		for k, v := range parsePythonRequirements(fp) {
			componentMap[k] = v
		}
	}

	tomlPath := filepath.Join(dir, "pyproject.toml")
	if utils.IsFile(tomlPath) {
		if list, e := tomlBuildSysFile(tomlPath); e != nil {
			logger.Warn.Println("Analyze pyproject.toml failed", e.Error())
		} else {
			mergeComponentInto(componentMap, list)
		}
	}

	for s := range ignoreSet {
		delete(componentMap, s)
	}
	if len(componentMap) == 0 {
		return nil
	}
	{
		m := model.Module{
			Name:           "Python",
			PackageManager: model.PMPip,
			Language:       model.Python,
			Dependencies:   []model.Dependency{},
			FilePath:       filepath.Join(dir),
		}
		for k, v := range componentMap {
			m.Dependencies = append(m.Dependencies, model.Dependency{
				Name:    k,
				Version: v,
			})
		}
		model.UseInspectorTask(ctx).AddModule(m)
		return nil
	}
}

func mergeComponentInto(source map[string]string, append []model.Dependency) {
	for _, it := range append {
		name, version := it.Name, it.Version
		if version == "" && source[name] != "" {
			continue
		}
		source[name] = version
	}
}

func parsePythonRequirements(p string) map[string]string {
	rs := map[string]string{}
	f, e := os.Open(p)
	if e != nil {
		logger.Warn.Println("Open file failed.", e.Error(), p)
		return nil
	}
	defer f.Close()
	scanner := bufio.NewScanner(io.LimitReader(f, 4*1024*1024))
	for scanner.Scan() {
		if scanner.Err() != nil {
			logger.Warn.Println("read file failed.", e.Error(), p)
			return nil
		}
		t := strings.TrimSpace(scanner.Text())
		m := pyRequirementsPattern.FindStringSubmatch(t)
		if m == nil {
			continue
		}
		rs[m[1]] = m[2]
	}
	return rs
}

func New() base.Inspector {
	return &Inspector{}
}
