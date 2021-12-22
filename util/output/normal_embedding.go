//go:build embedding

package output

import "fmt"

// Error 输出
func Error(text string) {
	text = wrapStr(text)
	writeToFile(fmt.Sprintf("[ERROR] %s\n", text))
}

// Debug 输出，当 Verbose == false 时不输出
func Debug(text string) {
	text = wrapStr(text)
	writeToFile(fmt.Sprintf("[DEBUG] %s\n", text))
}

// Info 输出
func Info(text string) {
	text = wrapStr(text)
	writeToFile(fmt.Sprintf("[INFO] %s\n", text))
}

func Warn(text string) {
	text = wrapStr(text)
	writeToFile(fmt.Sprintf("[WARN] %s\n", text))
}
