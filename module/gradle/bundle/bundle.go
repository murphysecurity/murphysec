package bundle

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
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
		var addPath = func(it os.DirEntry) {
			if !it.IsDir() {
				return
			}
			var name = it.Name()
			ver, e := semver.NewVersion(strings.TrimPrefix(name, "gradle-"))
			if e != nil {
				return
			}
			var item = Item{Version: ver, Path: filepath.Join(defaultGradlePath, name)}
			list = append(list, item)
		}
		dir, e := os.ReadDir(defaultGradlePath)
		if e != nil {
			return
		}
		for _, it := range dir {
			addPath(it)
		}
		dir, e = os.ReadDir("/opt")
		if e != nil {
			return
		}
		for _, it := range dir {
			if strings.HasPrefix(it.Name(), "gradle-") {
				addPath(it)
			}
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
	selectedJavaHome = SelectJavaHome(gradleVersion)
	return
}

var java8Paths = []string{"/opt/java/8", "/opt/openjdk/jdk8u402-b06"}
var java11Paths = []string{"/opt/java/11", "/opt/openjdk/jdk-11.0.22+7"}
var java17Paths = []string{"/opt/java/17", "/opt/openjdk/jdk-17.0.10+7"}
var java21Paths = []string{"/opt/java/21", "/opt/openjdk/jdk-21.0.2+13"}
var java24Paths = []string{"/opt/java/24", "/opt/openjdk/jdk-24.0.1"}

func SelectJavaHome(gradleVersion string) string {
	parsedVersion, e := semver.NewVersion(gradleVersion)
	if e != nil {
		return ""
	}
	var javaHome []string
	if parsedVersion.Compare(_GradleVersionFirstTimeAllowJavaWith2XVersionNumber) < 0 {
		javaHome = java8Paths
	} else if parsedVersion.Compare(_GradleVersionFirstTimeSupportJdk17) < 0 {
		javaHome = java11Paths
	} else if parsedVersion.Compare(_GradleVersionFirstTimeSupportJdk21) < 0 {
		javaHome = java17Paths
	} else if parsedVersion.Compare(_GradleVersionFirstTimeSupportJdk24) < 0 {
		javaHome = java21Paths
	} else {
		java24Exists := false
		for _, it := range java24Paths {
			if utils.IsDir(it) {
				java24Exists = true
				break
			}
		}
		if java24Exists {
			javaHome = java24Paths
		} else {
			javaHome = java21Paths
		}
	}
	for _, p := range javaHome {
		if utils.IsDir(p) {
			return p
		}
	}
	return ""
}

var _GradleVersionFirstTimeAllowJavaWith2XVersionNumber = semver.MustParse("4.7")
var _GradleVersionFirstTimeSupportJdk17 = semver.MustParse("7.3.3")
var _GradleVersionFirstTimeSupportJdk21 = semver.MustParse("8.5")
var _GradleVersionFirstTimeSupportJdk24 = semver.MustParse("8.14")
