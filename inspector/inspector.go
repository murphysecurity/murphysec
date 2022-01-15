package inspector

import (
	"murphysec-cli-simple/logger"
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

func tryMatchInspector(dir string) base.Inspector {
	for _, it := range engines {
		logger.Debug.Println("Try match project by inspector:", it.String(), "...")
		if it.CheckDir(dir) {
			logger.Info.Println("Matched.")
			return it
		}
	}
	logger.Debug.Println("Match failed")
	return nil
}
