package cpphasher

import (
	_ "embed"
	"github.com/murphysecurity/murphysec/infra/predata"
)

//go:embed cxx_file_ext
var __cxxFileExtData string

var cppFileExtSet = predata.StringsToMapBool(predata.ParseString(__cxxFileExtData))
