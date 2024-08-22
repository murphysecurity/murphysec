package go_mod

import (
	"bufio"
	"context"
	"encoding/json"
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
	"io"
	"os/exec"
	"path/filepath"
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
	}
	if len(dependencies) == 0 {
		if !env.DoNotBuild {
			logger.Warn("no dependencies found, backup")
		}
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
	var cmd = exec.CommandContext(ctx, "go", "list", "-json", "all")
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
		defer func() { _ = stderr.Close() }()
		var scanner = bufio.NewScanner(stderr)
		scanner.Buffer(nil, 1024*4)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			logger.Warn("go: " + scanner.Text())
		}
	}()
	logger.Sugar().Infof("executing command: %s", cmd)
	var decoder = json.NewDecoder(stdout)
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
	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			logger.Debug("process killed, waiting...")

			_, _ = cmd.Process.Wait()
			logger.Debug("after wait.")
		}
	}()
	logger.Debug("start scanning...")
	var m struct {
		Path     string `json:"path"`
		Version  string `json:"version"`
		Indirect bool   `json:"indirect"`
	}
	for {
		e = decoder.Decode(&m)
		if e != nil {
			break
		}
		result = append(result, model.DependencyItem{
			Component: model.Component{
				CompName:    m.Path,
				CompVersion: m.Version,
				EcoRepo:     EcoRepo,
			},
			IsDirectDependency: !m.Indirect,
		})
	}
	if e != io.EOF {
		e = fmt.Errorf("decode json failed: %w", e)
		logger.Error(e.Error())
		return
	}
	logger.Debug("done.")
	_ = stdout.Close()
	return
}
