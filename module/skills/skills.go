package skills

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/murphysecurity/murphysec/infra/pathignore"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
)

const (
	skillManifestName   = "SKILL.md"
	skillsModuleName    = "Skills"
	skillsModuleUUID    = "7f1d8b7f-fd00-4d61-a4f0-f1c1165bb7aa"
	skillsPackageManger = "skills"
)

type Inspector struct{}

func (i Inspector) String() string {
	return "Skills"
}

func (i Inspector) CheckDir(ctx context.Context, dir string) bool {
	task := model.UseScanTask(ctx)
	return task != nil && dir == task.ProjectPath
}

func (i Inspector) InspectProject(ctx context.Context) error {
	inspTask := model.UseInspectionTask(ctx)
	module, ok, err := ScanDir(inspTask.Dir())
	if err != nil || !ok {
		return err
	}
	inspTask.AddModule(module)
	return nil
}

func (i Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func ScanDir(root string) (model.Module, bool, error) {
	skillDirs, err := findSkillDirs(root)
	if err != nil {
		return model.Module{}, false, err
	}
	if len(skillDirs) == 0 {
		return model.Module{}, false, nil
	}

	deps := make([]model.DependencyItem, 0, len(skillDirs))
	for _, skillDir := range skillDirs {
		dep, err := scanSkill(root, skillDir)
		if err != nil {
			return model.Module{}, false, err
		}
		deps = append(deps, dep)
	}

	return model.Module{
		ModuleName:     skillsModuleName,
		ModuleVersion:  skillsModuleUUID,
		ModulePath:     root,
		PackageManager: skillsPackageManger,
		Dependencies:   deps,
		ScanStrategy:   model.ScanStrategyNormal,
	}, true, nil
}

func findSkillDirs(root string) ([]string, error) {
	var skillDirs []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root && shouldSkipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.EqualFold(d.Name(), skillManifestName) {
			skillDirs = append(skillDirs, filepath.Dir(path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(skillDirs)
	return skillDirs, nil
}

func scanSkill(root, skillDir string) (model.DependencyItem, error) {
	skillFiles, err := collectSkillFiles(skillDir)
	if err != nil {
		return model.DependencyItem{}, err
	}
	skillFilesValue := model.SkillFiles(skillFiles)

	relDir, err := filepath.Rel(root, skillDir)
	if err != nil {
		return model.DependencyItem{}, err
	}

	return model.DependencyItem{
		Component: model.Component{
			CompName:    filepath.Base(skillDir),
			CompVersion: "",
			SkillFiles:  &skillFilesValue,
			EcoRepo: model.EcoRepo{
				Ecosystem:  skillsPackageManger,
				Repository: filepath.ToSlash(relDir),
			},
		},
		DependencyRelation: model.DependencyRelationDirect,
		IsOnline:           model.IsOnlineFalse(),
	}, nil
}

func collectSkillFiles(skillDir string) ([]model.SkillFile, error) {
	var files []string
	err := filepath.WalkDir(skillDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != skillDir && shouldSkipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.EqualFold(filepath.Ext(d.Name()), ".md") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)

	result := make([]model.SkillFile, 0, len(files))
	for _, path := range files {
		relPath, err := filepath.Rel(skillDir, path)
		if err != nil {
			return nil, err
		}
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		raw, normalized, err := utils.ComputeSHA256AndLF(file)
		_ = file.Close()
		if err != nil {
			return nil, err
		}
		result = append(result, model.SkillFile{
			RelativePath: filepath.ToSlash(relPath),
			SHA256Hashes: []model.SHA256Hash{raw, normalized},
		})
	}
	return result, nil
}

func shouldSkipDir(name string) bool {
	if name == ".agent" {
		return false
	}
	return utils.HasHiddenFilePrefix(name) || pathignore.DirName(name)
}

var _ model.Inspector = Inspector{}
