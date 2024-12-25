package buildout

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"

	list "github.com/bahlo/generic-list-go"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/scanerr"
	"github.com/murphysecurity/murphysec/utils"
)

func doBuildout(ctx context.Context, dir string) (errorText string, e error) {
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
		errorText = e.Error()
		e = fmt.Errorf("failed to start buildout process: %w", e)
		return
	}
	var locker sync.Mutex
	var errLines = list.New[string]()
	var errAppender = func(prefix string, s string) {
		locker.Lock()
		defer locker.Unlock()
		if s == "" {
			return
		}
		if errLines.Len() > 32 {
			errLines.Remove(errLines.Front())
		}
		errLines.PushBack(prefix + ": " + s)
	}
	var wg sync.WaitGroup
	launchPipeForward(ctx, "o", stdout, errAppender, &wg)
	launchPipeForward(ctx, "e", stderr, errAppender, &wg)
	e = cmd.Wait()
	wg.Wait()
	var errLineI = errLines.Front()
	for errLineI != nil {
		errorText += errLineI.Value + "\n"
		errLineI = errLineI.Next()
	}
	if e != nil {
		e = fmt.Errorf("buildout process failed: %w", e)
		return
	}
	return
}

func launchPipeForward(ctx context.Context, prefix string, reader io.ReadCloser, lineAppender func(prefix, s string), wg *sync.WaitGroup) {
	var log = logctx.Use(ctx).Sugar()
	wg.Add(1)
	go func() {
		defer func() { wg.Done() }()
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
			var text = scanner.Text()
			if lineAppender != nil {
				lineAppender(prefix, text)
			}
			log.Debugf("%s: %s", prefix, text)
		}
	}()
}

type Inspector struct{}

func DirHasBuildout(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "buildout.cfg"))
}

func InspectProject(ctx context.Context, dir string) error {
	var log = logctx.Use(ctx).Sugar()
	var task = model.UseInspectionTask(ctx)
	var errText, e = doBuildout(ctx, dir)
	if e != nil {
		log.Warnf("failed to run buildout: %s", e.Error())
		if errText != "" {
			scanerr.Add(ctx, scanerr.Param{
				Kind:    "auto_build_error",
				Content: errText,
			})
		}
	}
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
			task.AddModule(model.Module{
				ModuleName:     filepath.Base(path),
				ModulePath:     filepath.Join(dir, "METADATA"),
				PackageManager: "Buildout",
				Dependencies: []model.DependencyItem{
					{
						Component: model.Component{
							CompName:    n,
							CompVersion: v,
							EcoRepo: model.EcoRepo{
								Ecosystem:  "pypi",
								Repository: "",
							},
						},
						IsOnline: model.IsOnlineTrue(),
					},
				},
				ScanStrategy: model.ScanStrategyNormal,
			})
		}
		if d.Name() == "buildout.cfg" {
			if err := base(ctx, path); err != nil {
				return err
			}
		}
		return nil
	})

	return nil
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
