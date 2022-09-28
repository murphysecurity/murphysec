package gradle

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
)

//goland:noinspection GoNameStartsWithPackageName
type GradleWrapperStatus string

//goland:noinspection GoNameStartsWithPackageName
const (
	// GradleWrapperStatusNotDetected Gradle wrapper is not detected
	GradleWrapperStatusNotDetected GradleWrapperStatus = "gradle_wrapper_not_detected"
	// GradleWrapperStatusUsed Gradle wrapper is detected and used
	GradleWrapperStatusUsed GradleWrapperStatus = "gradle_wrapper_used"
	// GradleWrapperStatusError Gradle wrapper is detected but some errors happened
	GradleWrapperStatusError GradleWrapperStatus = "gradle_wrapper_error"
)

// check gradlew script exists and grant permission is required
//
// returns gradle wrapper script if exists
func prepareGradleWrapperScriptFile(ctx context.Context, dir string) string {
	var isWindows = runtime.GOOS == "windows"

	var wrapperScriptPath string
	if isWindows {
		wrapperScriptPath = filepath.Join(dir, "gradlew.bat")
	} else {
		wrapperScriptPath = filepath.Join(dir, "gradlew")
	}
	stat, e := os.Stat(wrapperScriptPath)
	if e != nil || stat.IsDir() {
		return ""
	}
	if isWindows {
		return wrapperScriptPath
	}
	// on other platform, check executable permission
	mode := stat.Mode()
	if !mode.IsRegular() {
		return ""
	}
	if mode.Perm()&0111 == 0 {
		_ = os.Chmod(wrapperScriptPath, 0755)
	}
	return wrapperScriptPath
}
