package cpphasher

import (
	_ "embed"
	"github.com/murphysecurity/murphysec/infra/predata"
	"github.com/murphysecurity/murphysec/utils"
	"path/filepath"
)

//go:embed dirignore
var _dirIgnoreText string

var ignoredDirMap = predata.StringsToMapBool(predata.ParseString(_dirIgnoreText))

func dirShouldIgnore(name string) bool {
	return utils.HasHiddenFilePrefix(name) || ignoredDirMap[name] || ignoredDirMap[filepath.Base(name)]
}
