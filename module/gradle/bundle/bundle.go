package bundle

import (
	"github.com/Masterminds/semver"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Item struct {
	// Version
	Version *semver.Version
	// Path
	Path string
}

var once sync.Once
var list []Item

const defaultGradlePath = "/opt/gradle"

func List() []Item {
	once.Do(func() {
		dir, e := os.ReadDir(defaultGradlePath)
		if e != nil {
			return
		}
		for _, it := range dir {
			if !it.IsDir() {
				continue
			}
			var name = it.Name()
			ver, e := semver.NewVersion(strings.TrimPrefix(name, "gradle-"))
			if e != nil {
				continue
			}
			var item = Item{Version: ver, Path: filepath.Join(defaultGradlePath, name)}
			list = append(list, item)
		}
		slices.SortFunc(list, func(a, b Item) int { return a.Version.Compare(b.Version) })
	})
	var _list = make([]Item, len(list))
	copy(_list, list)
	return _list
}

func FindOkVersion(gradleVersion string) (selectedGradleBinPath, selectedJavaHome string) {
	gradleParsedVersion, e := semver.NewVersion(gradleVersion)
	if e != nil {
		return
	}
	r, _, ok := lo.FindLastIndexOf(List(), func(it Item) bool {
		return it.Version.Major() == gradleParsedVersion.Major() &&
			it.Version.Minor() == gradleParsedVersion.Minor()
	})
	if !ok {
		return
	}
	selectedGradleBinPath = filepath.Join(r.Path, "bin", "gradle")
	selectedJavaHome = selectJavaVersionOnly(gradleParsedVersion)
	selectedJavaHome = filepath.Join("/opt/openjdk", selectedJavaHome)
	return
}

func SelectJavaHome(gradleVersion string) string {
	var parsed, e = semver.NewVersion(gradleVersion)
	if e != nil {
		return ""
	}
	return filepath.Join(selectJavaVersionOnly(parsed))
}

const (
	jdk8  = "jdk8u402-b06"
	jdk11 = "jdk-11.0.22+7"
	jdk17 = "jdk-17.0.10+7"
	jdk21 = "jdk-21.0.2+13"
	jdk24 = "jdk-24.0.1"
)

var _GradleVersionFirstTimeAllowJavaWith2XVersionNumber = semver.MustParse("4.7")
var _GradleVersionFirstTimeSupportJdk17 = semver.MustParse("7.3.3")
var _GradleVersionFirstTimeSupportJdk21 = semver.MustParse("8.5")
var _GradleVersionFirstTimeSupportJdk24 = semver.MustParse("8.14")

func selectJavaVersionOnly(parsedVersion *semver.Version) string {
	if parsedVersion.Compare(_GradleVersionFirstTimeAllowJavaWith2XVersionNumber) < 0 {
		return jdk8
	} else if parsedVersion.Compare(_GradleVersionFirstTimeSupportJdk17) < 0 {
		return jdk11
	} else if parsedVersion.Compare(_GradleVersionFirstTimeSupportJdk21) < 0 {
		return jdk17
	} else if parsedVersion.Compare(_GradleVersionFirstTimeSupportJdk24) < 0 {
		return jdk21
	}
	return jdk24
}
