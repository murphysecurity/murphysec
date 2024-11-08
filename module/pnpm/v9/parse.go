package v9

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/pnpm/shared"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io"
	"sort"
	"strconv"
	"strings"
)

type importerDepItem struct {
	Specifier string `yaml:"specifier"`
	Version   string `yaml:"version"`
}

type importerItem struct {
	DevDependencies map[string]importerDepItem `yaml:"devDependencies"`
	Dependencies    map[string]importerDepItem `yaml:"dependencies"`
}

type snapshotItem struct {
	Dependencies map[string]string `yaml:"dependencies"`
}

func Parse(ctx context.Context, reader io.Reader) (trees []shared.DepTree, e error) {
	var d = yaml.NewDecoder(reader)
	var doc struct {
		Importers map[string]importerItem `yaml:"importers"`
		Snapshots map[string]snapshotItem `yaml:"snapshots"`
	}
	e = d.Decode(&doc)
	if e != nil {
		return
	}
	for path, importer := range doc.Importers {
		var c = importerHandlingCtx{
			Logger:            logctx.Use(ctx).Sugar(),
			Snapshot:          doc.Snapshots,
			Handled:           make(map[[2]string]struct{}),
			CircularDetectMap: make(map[[2]string]struct{}),
			CircularPath:      make([][2]string, 0),
		}
		var deps []model.DependencyItem
		for name, obj := range importer.Dependencies {
			var r model.DependencyItem
			r, e = c.handle(name, obj.Version, true)
			if e != nil {
				return
			}
			deps = append(deps, r)
		}
		c.Handled = make(map[[2]string]struct{})
		c.CircularDetectMap = make(map[[2]string]struct{})
		c.CircularPath = make([][2]string, 0)
		for name, obj := range importer.DevDependencies {
			var r model.DependencyItem
			r, e = c.handle(name, obj.Version, false)
			if e != nil {
				return
			}
			deps = append(deps, r)
		}
		trees = append(trees, shared.DepTree{
			Name:         path,
			Dependencies: deps,
		})
	}
	return
}

type importerHandlingCtx struct {
	Logger            *zap.SugaredLogger
	Snapshot          map[string]snapshotItem
	Handled           map[[2]string]struct{}
	CircularDetectMap map[[2]string]struct{}
	CircularPath      [][2]string
}

func (i *importerHandlingCtx) handle(name, version string, online bool) (d model.DependencyItem, e error) {
	var rKey = [2]string{name, version}
	var _, circleDetected = i.CircularDetectMap[rKey]
	defer func() {
		i.CircularPath = i.CircularPath[:len(i.CircularPath)-1]
		delete(i.CircularDetectMap, rKey)
	}()
	i.CircularDetectMap[rKey] = struct{}{}
	i.CircularPath = append(i.CircularPath, rKey)
	realVersion := splitRealVersionInVersionString(version)
	if len(realVersion) > 1 && strings.Contains(realVersion[1:], "@") {
		var suffix string
		name, version, suffix = splitNameVersion(realVersion)
		realVersion = version
		version += suffix
	}
	var searchKey = name + "@" + version
	f, ok := i.Snapshot[searchKey]
	if !ok {
		i.Logger.Warnf("dependency %s not found in snapshot", strconv.Quote(searchKey))
		return
	}
	d.CompName = name
	d.CompVersion = splitRealVersionInVersionString(version)
	d.Ecosystem = "npm"
	d.IsOnline.SetOnline(online)
	if _, ok := i.Handled[rKey]; ok || circleDetected {
		return
	}
	i.Handled[rKey] = struct{}{}
	var pairs = lo.ToPairs(f.Dependencies)
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Key < pairs[j].Key })
	for _, pair := range pairs {
		depName := pair.Key
		depVersion := pair.Value
		var r model.DependencyItem
		r, e = i.handle(depName, depVersion, online)
		if e != nil {
			return
		}
		d.Dependencies = append(d.Dependencies, r)
	}
	return
}

func splitRealVersionInVersionString(input string) string {
	var i = strings.Index(input, "(")
	if i < 0 {
		return input
	}
	return input[:i]
}

func splitNameVersion(input string) (name, version, suffix string) {
	if len(input) < 2 {
		input = name
		return
	}
	var i = strings.Index(input, "(")
	if i >= 0 {
		input = input[:i]
		suffix = input[i:]
	}
	var parts = strings.Split(input, "@")
	name = parts[0]
	if len(parts) > 1 {
		version = parts[1]
	}
	return
}
