package output

import (
	"fmt"
	"github.com/fatih/color"
	"murphysec-cli-simple/util/must"
	"strings"
)

// Colorful 控制 Error Debug Info 输出是否带颜色
var Colorful = true

// Verbose 控制是否输出 Debug
var Verbose = false

// Error 输出
func Error(text string) {
	text = wrapStr(text)
	if Colorful {
		must.Int(color.New(color.Bold, color.FgRed).Printf("[ERROR] %s\n", text))
	} else {
		fmt.Printf("[ERROR] %s\n", text)
	}
}

// Debug 输出，当 Verbose == false 时不输出
func Debug(text string) {
	text = wrapStr(text)
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
	if Colorful {
		must.Int(color.New(color.Bold, color.FgCyan).Printf("[INFO] %s\n", text))
	} else {
		fmt.Printf("[INFO] %s\n", text)
	}
}

func Warn(text string) {
	text = wrapStr(text)
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
