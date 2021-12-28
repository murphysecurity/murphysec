package output

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"murphysec-cli-simple/util/must"
	"os"
	"strings"
	"time"
)

// Colorful 控制 Error Debug Info 输出是否带颜色
var Colorful = true

// Verbose 控制是否输出 Debug
var Verbose = false

func wrapStr(input string) string {
	lines := strings.Split(input, "\n")
	for i := range lines {
		if i > 0 {
			lines[i] = "    " + strings.Trim(lines[i], "\r\n")
		}
	}
	return strings.Join(lines, "\n")
}

const logFileName string = "20060102-150405"

var logFile = func() func() *os.File {
	must.Must(os.MkdirAll(must.String(homedir.Expand("~/.murphysec/logs")), 0755))
	path := must.String(homedir.Expand(fmt.Sprintf("~/.murphysec/logs/%s.log", time.Now().Format(logFileName))))
	f, e := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	must.Must(e)
	return func() *os.File {
		return f
	}
}()

func writeToFile(s string) {
	_, _ = logFile().WriteString(time.Now().Format(time.RFC3339) + " " + s)
}
