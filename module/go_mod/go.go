package go_mod

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"

	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/mod/modfile"
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
	logger := logctx.Use(ctx)
	if privatePath, ok := ctx.Value("privateSourceAddr").(string); ok {
		logger.Debug("Use private path", zap.String("path", privatePath))
		if err := setPrivatePath(privatePath, logger); err != nil {
			return err
		}
	}
	if proxyPath, ok := ctx.Value("proxyAddr").(string); ok {
		logger.Debug("Use proxy path", zap.String("path", proxyPath))
		if err := setProxyPath(proxyPath, logger); err != nil {
			return err
		}
	}
	if err := buildScan(ctx); err != nil {
		if err := baseScan(ctx); err != nil {
			return err
		}
	}
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
	var cmd = exec.CommandContext(ctx, "go", "list", "-json", "-m", "all")
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
