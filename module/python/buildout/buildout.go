package buildout

import (
	"bufio"
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/repeale/fp-go"
	"golang.org/x/exp/maps"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func doBuildout(ctx context.Context, dir string) (e error) {
	var cmd = exec.CommandContext(ctx, "buildout")
	cmd.Dir = dir
	cmd.Stdin = nil
	stdout, e := cmd.StdoutPipe()
	if e != nil {
		e = fmt.Errorf("failed to create stdout pipe: %w", e)
		return
	}
	stderr, e := cmd.StderrPipe()
	if e != nil {
		e = fmt.Errorf("failed to create stderr pipe: %w", e)
		return
	}
	e = cmd.Start()
	if e != nil {
		e = fmt.Errorf("failed to start buildout process: %w", e)
		return
	}
	launchPipeForward(ctx, "o", stdout)
	launchPipeForward(ctx, "e", stderr)
	e = cmd.Wait()
	if e != nil {
		e = fmt.Errorf("buildout process failed: %w", e)
		return
	}
	return
}

func launchPipeForward(ctx context.Context, prefix string, reader io.ReadCloser) {
	var log = logctx.Use(ctx).Sugar()
	go func() {
		var scanner = bufio.NewScanner(reader)
		scanner.Split(bufio.ScanLines)
		scanner.Buffer(nil, 4096)
		for scanner.Scan() {
			var e = scanner.Err()
			if e != nil {
				log.Errorf("error during read lines, prefix %s: %s", strconv.Quote(prefix), e.Error())
				e = reader.Close()
				if e != nil {
					log.Errorf("failed to close reader, prefix %s: %s", strconv.Quote(prefix), e.Error())
				}
				return
			}
			log.Debugf("%s: %s", prefix, scanner.Text())
		}
	}()
}

type Inspector struct{}

func DirHasBuildout(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "buildout.cfg"))
}

func InspectProject(ctx context.Context, dir string) (*model.Module, error) {
	var log = logctx.Use(ctx).Sugar()

	var e = doBuildout(ctx, dir)
	if e != nil {
		log.Warnf("failed to run buildout: %s", e.Error())
		return nil, e
	}

	var comps = make(map[[2]string]struct{})
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, e error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "METADATA" {
			log.Debugf("inspecting file: %s", path)
			n, v, e := parseMetadataFile(ctx, path)
			if e != nil || n == "" {
				return nil
			}
			comps[[2]string{n, v}] = struct{}{}
		}
		return nil
	})
	var compList = maps.Keys(comps)
	if len(compList) == 0 {
		return nil, nil
	}
	var module = model.Module{
		ModuleName:     filepath.Dir(dir),
		ModulePath:     filepath.Join(dir, "buildout.cfg"),
		PackageManager: "Buildout",
		Dependencies: fp.Map(func(it [2]string) model.DependencyItem {
			return model.DependencyItem{
				Component: model.Component{
					CompName:    it[0],
					CompVersion: it[1],
					EcoRepo: model.EcoRepo{
						Ecosystem:  "pypi",
						Repository: "",
					},
				},
				IsOnline: model.IsOnlineTrue(),
			}
		})(compList),
		ScanStrategy: model.ScanStrategyNormal,
	}

	return &module, nil
}

func parseMetadataFile(ctx context.Context, path string) (name, version string, e error) {
	var log = logctx.Use(ctx).Sugar()
	var file *os.File
	file, e = os.Open(path)
	if e != nil {
		e = fmt.Errorf("failed to open metadata file: %w", e)
		return
	}
	defer func() { _ = file.Close() }()
	r, e := ParseMetadata(file)
	if e != nil {
		e = fmt.Errorf("failed to parse metadata file: %w", e)
		return
	}
	name = getFieldFromResult(r, "Name")
	version = getFieldFromResult(r, "Version")
	log.Debugf("metadata: %s %s", name, version)
	return
}

func getFieldFromResult(r map[string][]string, field string) string {
	if v, ok := r[field]; ok {
		return v[0]
	}
	return ""
}
