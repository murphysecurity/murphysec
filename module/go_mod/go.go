package go_mod

import (
	"bufio"
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/sl"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"github.com/repeale/fp-go"
	"go.uber.org/zap"
	"golang.org/x/mod/modfile"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Inspector struct{}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return model.InspectorFeatureAllowNested&feature > 0
}

func (Inspector) String() string {
	return "GoMod"
}

func (Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "go.mod"))
}

func (Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	logger := logctx.Use(ctx)
	modFilePath := filepath.Join(task.Dir(), "go.mod")
	logger.Debug("Reading go.mod", zap.String("path", modFilePath))
	data, e := utils.ReadFileLimited(modFilePath, 1024*1024*4)
	if e != nil {
		return errors.WithMessage(e, "Open GoMod file")
	}
	logger.Debug("Parsing go.mod")
	f, e := modfile.ParseLax(filepath.Base(modFilePath), data, nil)
	if e != nil {
		return errors.WithMessage(e, "Parse go mod failed")
	}
	var dependencies []model.DependencyItem
	if !env.DoNotBuild {
		// try command go list
		dependencies, e = doGoList(ctx, task.Dir())
		if e != nil {
			if errors.Is(e, _ErrGoNotFound) {
				logger.Debug("Go not found, skip GoList")
			} else {
				// log it and go on
				logger.Warn("GoList failed", zap.Error(e))
			}
			dependencies = append(dependencies, fp.Map(mapRequireToDependencyItem)(sl.FilterNotNull(f.Require))...)
		}
	} else {
		dependencies = append(dependencies, fp.Map(mapRequireToDependencyItem)(sl.FilterNotNull(f.Require))...)
	}
	m := model.Module{
		PackageManager: "gomod",
		ModulePath:     modFilePath,
		ModuleName:     "<NoNameModule>",
		Dependencies:   dependencies,
	}
	if f.Module != nil {
		m.ModuleVersion = f.Module.Mod.Version
		m.ModuleName = f.Module.Mod.Path
	}
	task.AddModule(m)
	return nil
}

func mapRequireToDependencyItem(it *modfile.Require) model.DependencyItem {
	return model.DependencyItem{
		Component: model.Component{
			CompName:    it.Mod.Path,
			CompVersion: it.Mod.Version,
			EcoRepo:     EcoRepo,
		},
		IsDirectDependency: !it.Indirect,
	}
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "go",
	Repository: "",
}

var _ErrGoNotFound = errors.New("go not found")

func doGoList(ctx context.Context, dir string) (result []model.DependencyItem, e error) {
	var logger = logctx.Use(ctx)
	var cmd = exec.CommandContext(ctx, "go", "list", "-f", `{{with .Module}}{{.Path}} {{.Version}} {{.Indirect}}{{end}}`, "all")
	cmd.Dir = dir
	stdout, e := cmd.StdoutPipe()
	if e != nil {
		e = fmt.Errorf("create stdout pipe failed: %w", e)
		logger.Error(e.Error())
		return
	}
	stderr, e := cmd.StderrPipe()
	if e != nil {
		e = fmt.Errorf("create stderr pipe failed: %w", e)
		logger.Error(e.Error())
		return
	}
	go func() {
		var scanner = bufio.NewScanner(stderr)
		scanner.Buffer(nil, 1024*4)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			logger.Warn("go: " + scanner.Text())
		}
	}()
	logger.Sugar().Infof("executing command: %s", cmd)
	var scanner = bufio.NewScanner(stdout)
	scanner.Buffer(nil, 1024*4)
	scanner.Split(bufio.ScanLines)
	e = cmd.Start()
	if e != nil {
		// if the command is not found, we should not return error
		if errors.Is(e, exec.ErrNotFound) {
			e = _ErrGoNotFound
			return
		}
		e = fmt.Errorf("start command failed: %w", e)
		logger.Error(e.Error())
		return
	}
	logger.Debug("start scanning...")
	for scanner.Scan() {
		// the text has 3 segment, name version and indirect
		// we should split them by space
		var text = scanner.Text()
		// split
		var r = strings.Split(text, " ")
		if len(r) != 3 {
			// if the length is not 3, we should skip this line
			logger.Warn("invalid line", zap.String("line", text))
			continue
		}
		// the first segment is the name of the module
		var name = r[0]
		// the second segment is the version of the module
		var version = r[1]
		// the third segment is the indirect flag
		var indirect = r[2]
		// parse indirect flag as strconv
		var isIndirect, e = strconv.ParseBool(indirect)
		if e != nil {
			logger.Warn("parse indirect flag failed", zap.String("line", text), zap.Error(e))
			continue
		}
		// check name is not empty
		if name == "" {
			logger.Warn("name is empty", zap.String("line", text))
			continue
		}
		// append to result
		result = append(result, model.DependencyItem{
			Component: model.Component{
				CompName:    name,
				CompVersion: version,
				EcoRepo:     EcoRepo,
			},
			IsDirectDependency: !isIndirect,
		})
	}
	// Thanks to copilot
	_ = cmd.Wait()
	return
}
