package module

import (
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/arkts"
	"github.com/murphysecurity/murphysec/module/bundler"
	"github.com/murphysecurity/murphysec/module/cargo"
	"github.com/murphysecurity/murphysec/module/cocoapods"
	"github.com/murphysecurity/murphysec/module/composer"
	"github.com/murphysecurity/murphysec/module/conan"
	"github.com/murphysecurity/murphysec/module/go_mod"
	"github.com/murphysecurity/murphysec/module/gradle"
	"github.com/murphysecurity/murphysec/module/ivy"
	"github.com/murphysecurity/murphysec/module/maven"
	"github.com/murphysecurity/murphysec/module/npm"
	"github.com/murphysecurity/murphysec/module/nuget"
	"github.com/murphysecurity/murphysec/module/perl"
	"github.com/murphysecurity/murphysec/module/pnpm"
	"github.com/murphysecurity/murphysec/module/poetry"
	"github.com/murphysecurity/murphysec/module/python"
	"github.com/murphysecurity/murphysec/module/rebar3"
	"github.com/murphysecurity/murphysec/module/renv"
	"github.com/murphysecurity/murphysec/module/sbt"
	"github.com/murphysecurity/murphysec/module/yarn"
	"os"
	"sort"
	"strconv"
	"strings"
)

var Inspectors []model.Inspector

func GetSupportedModuleList() []string {
	var r []string
	for _, it := range Inspectors {
		r = append(r, it.String())
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return r
}

func init() {
	Inspectors = append(Inspectors, arkts.Inspector{})
	Inspectors = append(Inspectors, bundler.Inspector{})
	Inspectors = append(Inspectors, cargo.Inspector{})
	Inspectors = append(Inspectors, cocoapods.Inspector{})
	Inspectors = append(Inspectors, composer.Inspector{})
	Inspectors = append(Inspectors, conan.Inspector{})
	Inspectors = append(Inspectors, go_mod.Inspector{})
	if enableScan("MAVEN") {
		Inspectors = append(Inspectors, maven.Inspector{})
	}
	if enableScan("GRADLE") {
		Inspectors = append(Inspectors, gradle.Inspector{})
	}
	Inspectors = append(Inspectors, ivy.Inspector{})
	Inspectors = append(Inspectors, npm.Inspector{})
	Inspectors = append(Inspectors, nuget.Inspector{})
	Inspectors = append(Inspectors, perl.Inspector{})
	Inspectors = append(Inspectors, pnpm.Inspector{})
	Inspectors = append(Inspectors, poetry.Inspector{})
	Inspectors = append(Inspectors, python.Inspector{})
	Inspectors = append(Inspectors, rebar3.Inspector{})
	Inspectors = append(Inspectors, renv.Inspector{})
	Inspectors = append(Inspectors, sbt.Inspector{})
	Inspectors = append(Inspectors, yarn.Inspector{})
}

func enableScan(name string) bool {
	return !boolEnv("DO_NOT_SCAN_" + strings.ToUpper(name))
}

func boolEnv(name string) bool {
	b, _ := strconv.ParseBool(os.Getenv(name))
	return b
}
