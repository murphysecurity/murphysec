package inspector

import (
	"murphysec-cli-simple/api"
	java_import_scanner "murphysec-cli-simple/java-import-scanner"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/must"
	"path/filepath"
	"regexp"
	"strings"
)

func javaImportClauseScan(v *api.VoDetectResponse, projectDir string) {
	taskId := v.TaskId
	type id struct {
		compId   int
		moduleId int
	}
	idMap := map[string][]id{}
	for _, module := range v.Modules {
		for _, comp := range module.Comps {
			idMap[comp.CompName] = append(idMap[comp.CompName], id{comp.CompId, module.ModuleId})
		}
	}
	var compNameList []string
	for s := range idMap {
		compNameList = append(compNameList, s)
	}
	// map compName -> listOf path
	scanResult := scanJavaImport(projectDir, compNameList)
	{
		// api request body
		// moduleId -> listOf compObj{comp_id: int, import_path: []string}
		m := map[int][]map[string]interface{}{}
		for compName, fileList := range scanResult {
			for _, it := range idMap[compName] {
				if m[it.moduleId] == nil {
					m[it.moduleId] = []map[string]interface{}{}
				}
				m[it.moduleId] = append(m[it.moduleId], map[string]interface{}{
					"comp_id":     it.compId,
					"import_path": fileList,
				})
			}
		}
		modules := make([]map[string]interface{}, 0)
		for i := range m {
			modules = append(modules, map[string]interface{}{
				"module_id": i,
				"comps":     m[i],
			})
		}
		api.CompImportPath(map[string]interface{}{
			"task_id": taskId,
			"modules": modules,
		})
	}
}

func scanJavaImport(dir string, compNames []string) map[string][]string {
	// todo: 临时方案，groupId 前缀匹配全限定名，artifactId用'-'分割，检查每一个词按先后顺序出现在全限定名去除groupId的后半部分中
	compNames = utils.DistinctStringSlice(compNames)
	// processCompNameSet
	type id struct {
		name           string
		groupId        string
		artifactIdList []string
	}
	var compInfo []id
	for _, it := range compNames {
		g, a := parseMavenCompName(it)
		if g == "" || a == "" {
			continue
		}
		compInfo = append(compInfo, id{
			name:           it,
			groupId:        g,
			artifactIdList: strings.Split(a, "."),
		})
	}

	jf := make(chan java_import_scanner.JavaFileImportItem, 100)
	go func() {
		java_import_scanner.JavaImportScan(dir, jf)
		close(jf)
	}()
	// map: compName -> setOf filePath
	rs := map[string]utils.StringSet{}
	file := <-jf
	for file.FilePath != "" {
		// for each file
		// todo: 匹配算法可能需要优化
		for _, qualName := range file.Imports {
			for _, info := range compInfo {
				if matchImportQualName(info.groupId, info.artifactIdList, qualName) {
					if rs[info.name] == nil {
						rs[info.name] = utils.NewStringSet()
					}
					rs[info.name].Put(must.String(filepath.Rel(dir, file.FilePath)))
				}
			}
		}
	}
	// convert map to slice
	r := map[string][]string{}
	for compName, files := range rs {
		r[compName] = files.ToSlice()
	}
	return r
}

var mavenCompPattern = regexp.MustCompile("^([A-Za-z0-9._-]+):([A-Za-z0-9._-]+)$")

func parseMavenCompName(name string) (groupId, artifactId string) {
	if m := mavenCompPattern.FindStringSubmatch(name); m != nil {
		return m[1], m[2]
	}
	return "", ""
}

func matchImportQualName(groupId string, artifactIdList []string, qualName string) bool {
	p := strings.TrimPrefix(qualName, groupId)
	if len(p) < len(qualName) {
		list := strings.Split(p, ".")
		pPtr := 0    // ptr to artifactIdList
		listPtr := 0 // ptr to list
		for {
			if pPtr >= len(artifactIdList) {
				return true
			}
			if listPtr >= len(list) {
				return false
			}
			if artifactIdList[pPtr] == list[listPtr] {
				pPtr++
			}
			listPtr++
		}
	}
	return false
}
