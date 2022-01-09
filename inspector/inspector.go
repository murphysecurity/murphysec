package inspector

import (
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/module/go_mod"
	"murphysec-cli-simple/module/maven"
	"murphysec-cli-simple/module/npm"
)

var engines = []base.Inspector{
	go_mod.New(),
	maven.New(),
	npm.New(),
}

func getInspectorSupportPkgManagerType(pkgType base.PackageManagerType) base.Inspector {
	for _, it := range engines {
		if it.PackageManagerType() == pkgType {
			return it
		}
	}
	return nil
}
