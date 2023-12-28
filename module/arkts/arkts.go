package arkts

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/samber/lo"
	"github.com/titanous/json5"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	packageFile = "oh-package.json5"
	lockFile    = "oh-package-lock.json5"
)

var ecoRepo = model.EcoRepo{
	Ecosystem: "npm",
}

type Inspector struct{}

func (Inspector) String() string {
	return "ArkTS"
}

func (Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, lockFile))
}

func (Inspector) InspectProject(ctx context.Context) error {
	return analyze(ctx)
}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

var _ model.Inspector = Inspector{}

type lockRoot struct {
	Specifiers map[string]string  `json:"specifiers"`
	Packages   map[string]lockPkg `json:"packages"`
}

type lockPkg struct {
	Dependencies map[string]string `json:"dependencies"`
}

type packageRoot struct {
	Dependencies map[string]string `json:"dependencies"`
}

func (p packageRoot) getRootsWithLock(lock lockRoot) (r [][2]string) {
	r = make([][2]string, 0, len(p.Dependencies))
	for n, v := range p.Dependencies {
		vv := lock.Specifiers[n]
		if vv == "" {
			vv = v
		}
		r = append(r, [2]string{n, vv})
	}
	sortPairOfString(r)
	return
}
func _buildDepTreeVisit(visited map[[2]string]struct{}, next [2]string, root *lockRoot) model.DependencyItem {
	var dep = model.DependencyItem{
		Component: model.Component{
			CompName:    next[0],
			CompVersion: next[1],
			EcoRepo:     ecoRepo,
		},
		Dependencies:       nil,
		IsDirectDependency: false,
	}

	if _, ok := visited[next]; ok {
		return dep
	}
	visited[next] = struct{}{}
	defer func() { delete(visited, next) }()
	var key = next[0] + "@" + next[1]
	if root.Packages[key].Dependencies == nil {
		return dep
	}
	for n, v := range root.Packages[key].Dependencies {
		vv := versionOfSpecifier(root.Specifiers[n+"@"+v])
		if vv == "" {
			vv = v
		}
		dep.Dependencies = append(dep.Dependencies, _buildDepTreeVisit(visited, [2]string{n, vv}, root))
	}
	return dep
}

func (l lockRoot) findRoots() [][2]string {
	var p = make(map[[2]string]struct{}, len(l.Packages))
	for s := range l.Packages {
		var n = nameOfSpecifier(s)
		if n == "" {
			continue
		}
		var v = versionOfSpecifier(s)
		if v == "" {
			continue
		}
		p[[2]string{n, v}] = struct{}{}
	}
	for _, pkg := range l.Packages {
		for n, v := range pkg.Dependencies {
			vv := versionOfSpecifier(l.Specifiers[n+"@"+v])
			if vv == "" {
				vv = v // I don't know if it's correct, just do it.
			}
			delete(p, [2]string{n, vv})
		}
	}
	var r = lo.Keys(p)
	sortPairOfString(r)
	return r
}

func sortPairOfString(a [][2]string) {
	sort.Slice(a, func(i, j int) bool {
		if n := strings.Compare(a[i][0], a[j][0]); n != 0 {
			return n < 0
		}
		return a[i][1] < a[j][1]
	})
}

func nameOfSpecifier(input string) string {
	var i = strings.LastIndex(input, "@")
	if i == -1 {
		return ""
	}
	return input[:i]
}

func versionOfSpecifier(input string) string {
	var i = strings.LastIndex(input, "@")
	if i == -1 || i+1 >= len(input) {
		return ""
	}
	return input[i+1:]
}

func analyze(ctx context.Context) (e error) {
	task := model.UseInspectionTask(ctx)
	dir := task.Dir()
	var logger = logctx.Use(ctx).Sugar()
	logger.Debugf("analyzing dir %s", dir)
	pkg, pkgReadErr := readFileHelper(ctx, filepath.Join(dir, packageFile), readPackageFile)
	lock, lockReadErr := readFileHelper(ctx, filepath.Join(dir, lockFile), readLockfile)
	if lockReadErr != nil {
		return lockReadErr
	}
	var roots [][2]string
	if pkgReadErr != nil {
		roots = pkg.getRootsWithLock(lock)
	} else {
		roots = lock.findRoots()
		if len(roots) == 0 {
			return ErrRootNotFound
		}
	}
	var m = model.Module{
		ModulePath:     filepath.Join(dir, lockFile),
		PackageManager: "arkts",
		Dependencies:   make([]model.DependencyItem, 0),
		ScanStrategy:   model.ScanStrategyNormal,
	}
	for _, it := range roots {
		m.Dependencies = append(m.Dependencies, _buildDepTreeVisit(map[[2]string]struct{}{}, it, &lock))
	}
	for i := range m.Dependencies {
		m.Dependencies[i].IsDirectDependency = true
	}
	task.AddModule(m)
	return
}

// readPackageFile, close input is caller's responsibility
func readPackageFile(ctx context.Context, input io.Reader) (pkg packageRoot, e error) {
	var decoder = json5.NewDecoder(input)
	e = decoder.Decode(&pkg)
	if e != nil {
		e = ErrReadPackage.withCause(e)
		return
	}
	if pkg.Dependencies == nil {
		pkg.Dependencies = make(map[string]string)
	}
	return
}

// readLockfile, close input is caller's responsibility
func readLockfile(ctx context.Context, input io.Reader) (root lockRoot, e error) {
	var decoder = json5.NewDecoder(input)
	e = decoder.Decode(&root)
	if e != nil {
		e = ErrReadPackageLock.withCause(e)
		return
	}
	if root.Specifiers == nil {
		root.Specifiers = make(map[string]string)
	}
	if root.Packages == nil {
		root.Packages = make(map[string]lockPkg)
	}
	return
}

func readFileHelper[T any](ctx context.Context, fp string, parseFn func(context.Context, io.Reader) (T, error)) (r T, e error) {
	var logger = logctx.Use(ctx).Sugar()
	var f *os.File
	f, e = os.Open(fp)
	var fn = filepath.Base(fp)
	if e != nil {
		logger.Warnf("open %s failed: %s", fn, e)
		return
	}
	defer func() {
		if e := f.Close(); e != nil {
			logger.Errorf("close %s failed: %s", fn, e)
		}
	}()
	r, e = parseFn(ctx, f)
	if e != nil {
		logger.Errorf("read %s failed: %s", fn, e)
		return
	}
	return
}
