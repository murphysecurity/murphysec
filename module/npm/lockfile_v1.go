package npm

import (
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/npm/shared"
	v1 "github.com/murphysecurity/murphysec/module/npm/v1"
	"github.com/samber/lo"
)

func processV1Lockfile(data []byte, pkg *pkgFile) ([]model.DependencyItem, error) {
	lf, e := v1.ParseLockfile(data)
	if e != nil {
		return nil, e
	}
	entries := pkg.DependenciesEntries()
	entries = append(entries, pkg.DevDependenciesEntries()...)
	lo.Uniq(entries)
	nodes, e := lf.Build(pkg.DependenciesEntries(), false)
	if e != nil {
		return nil, fmt.Errorf("build dependencies tree: %w", e)
	}
	return shared.ConvNodes(nodes), nil
}
