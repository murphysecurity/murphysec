//go:build !noprint

package output

import (
	"fmt"
	"github.com/fatih/color"
	"murphysec-cli-simple/util/must"
	"sync"
)

// Error 输出
func Error(text string) {
	syn.Lock()
	defer syn.Unlock()
	text = wrapStr(text)
	writeToFile(fmt.Sprintf("[ERROR] %s\n", text))
	if Colorful {
		must.Int(color.New(color.Bold, color.FgRed).Printf("[ERROR] %s\n", text))
	} else {
		fmt.Printf("[ERROR] %s\n", text)
	}
}

var syn sync.Mutex

// Debug 输出，当 Verbose == false 时不输出
func Debug(text string) {
	syn.Lock()
	defer syn.Unlock()
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
	syn.Lock()
	defer syn.Unlock()
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
	syn.Lock()
	defer syn.Unlock()
	text = wrapStr(text)
	writeToFile(fmt.Sprintf("[WARN] %s\n", text))
	if Colorful {
		must.Int(color.New(color.Bold, color.FgRed).Printf("[WARN] %s\n", text))
	} else {
		fmt.Printf("[WARN] %s\n", text)
	}
}
