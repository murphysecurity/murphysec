package buildout

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const MetadataMaxLine = 4096

func ParseMetadata(input io.Reader) (result map[string][]string, err error) {
	var scanner = bufio.NewScanner(bufio.NewReader(input))
	scanner.Buffer(nil, MetadataMaxLine)
	scanner.Split(bufio.ScanLines)
	result = make(map[string][]string)
	for scanner.Scan() {
		var e = scanner.Err()
		if e != nil {
			err = fmt.Errorf("error during read lines: %w", e)
			return
		}
		var text = scanner.Text()
		if strings.TrimSpace(text) == "" {
			break
		}
		var i = strings.Index(text, ":")
		if i == -1 || i == 0 || i == len(text)-1 {
			break
		}
		var key = strings.TrimSpace(text[:i])
		var value = strings.TrimSpace(text[i+1:])
		if key == "" || value == "" {
			break
		}
		result[key] = append(result[key], value)
	}
	return
}
