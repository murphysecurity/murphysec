package pathignore

import (
	_ "embed"
	"github.com/murphysecurity/murphysec/infra/sl"
	"github.com/repeale/fp-go"
	"path/filepath"
	"strings"
)

//go:embed dirskip
var _dirskipData []byte

var commonDirSkip []string

func init() {
	commonDirSkip = fp.Pipe2(fp.Map(strings.TrimSpace), fp.Filter(sl.NotF1(lineShouldSkip)))(strings.Split(string(_dirskipData), "\n"))
}

func DirName(s string) bool {
	for _, it := range commonDirSkip {
		if m, _ := filepath.Match(it, s); m {
			return true
		}
	}
	return false
}

func lineShouldSkip(s string) bool {
	return s == "" || s[0] == '#'
}
