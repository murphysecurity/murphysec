package output

import (
	"fmt"
	"github.com/fatih/color"
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

// Error 输出
func Error(text string) {
	text = wrapStr(text)
	writeToFile(fmt.Sprintf("[ERROR] %s\n", text))
	if Colorful {
		must.Int(color.New(color.Bold, color.FgRed).Printf("[ERROR] %s\n", text))
	} else {
		fmt.Printf("[ERROR] %s\n", text)
	}
}

// Debug 输出，当 Verbose == false 时不输出
func Debug(text string) {
	text = wrapStr(text)
	writeToFile(fmt.Sprintf("[DEBUG] %s\n", text))
	if !Verbose {
		return
	}
	if Colorful {
		must.Int(color.New(color.Bold, color.FgCyan).Printf("[DEBUG] %s\n", text))
	} else {
		fmt.Printf("[DEBUG] %s\n", text)
	}
}

// Info 输出
func Info(text string) {
	text = wrapStr(text)
	writeToFile(fmt.Sprintf("[INFO] %s\n", text))
	if !Verbose {
		return
	}
	if Colorful {
		must.Int(color.New(color.Bold, color.FgCyan).Printf("[INFO] %s\n", text))
	} else {
		fmt.Printf("[INFO] %s\n", text)
	}
}

func Warn(text string) {
	text = wrapStr(text)
	writeToFile(fmt.Sprintf("[WARN] %s\n", text))
	if Colorful {
		must.Int(color.New(color.Bold, color.FgRed).Printf("[WARN] %s\n", text))
	} else {
		fmt.Printf("[WARN] %s\n", text)
	}
}

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
	must.Must(os.MkdirAll(must.String(homedir.Expand("~/.murphysec/logs")), 755))
	path := must.String(homedir.Expand(fmt.Sprintf("~/.murphysec/logs/%s.log", time.Now().Format(logFileName))))
	fmt.Println(path)
	f, e := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 644)
	must.Must(e)
	return func() *os.File {
		return f
	}
}()

func writeToFile(s string) {
	_, _ = logFile().WriteString(time.Now().Format(time.RFC3339) + " " + s)
}
