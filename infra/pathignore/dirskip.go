package pathignore

import (
	_ "embed"
	"path/filepath"
	"strings"
)

//go:embed dirskip
var _dirskipData []byte

var commonDirSkip []string

func init() {
	for _, s := range strings.Split(string(_dirskipData), "\n") {
		s = strings.TrimSpace(s)
		if s == "" || s[0] == '#' {
			continue
		}
		commonDirSkip = append(commonDirSkip, s)
	}
}

func DirName(s string) bool {
	for _, it := range commonDirSkip {
		if m, _ := filepath.Match(it, s); m {
			return true
		}
	}
	return false
}
