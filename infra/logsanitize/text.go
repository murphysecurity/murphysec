package logsanitize

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ForLog normalizes text for log transport:
// 1) invalid UTF-8 is replaced with U+FFFD
// 2) control characters are escaped
// 3) NUL is escaped as \0
func ForLog(s string) string {
	s = strings.ToValidUTF8(s, "\uFFFD")
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r == 0:
			b.WriteString(`\0`)
		case r == '\t':
			b.WriteString(`\t`)
		case r < 0x20 || r == 0x7f:
			_, _ = fmt.Fprintf(&b, `\x%02x`, r)
		default:
			if r == utf8.RuneError {
				b.WriteRune('\uFFFD')
			} else {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}

