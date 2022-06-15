package python

import (
	"bufio"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"io"
	"os"
	"regexp"
	"strings"
)

var pyImportPattern1 = regexp.MustCompile("import\\s+(?:[A-Za-z_-][\\w.-]*)(?:\\s*,\\s*(?:[A-Za-z_-][\\w.-]*))")
var pyImportPattern2 = regexp.MustCompile("from\\s+([A-Za-z_-][\\w-]*)")

func parsePythonRequirements(p string) (rs []model.Dependency) {
	var pyRequirementsPattern = regexp.MustCompile("^([\\w-]+) *.?= *([^= \\n\\r]+)$")
	f, e := os.Open(p)
	if e != nil {
		logger.Warn.Println("Open file failed.", e.Error(), p)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(io.LimitReader(f, 4*1024*1024))
	for scanner.Scan() {
		if scanner.Err() != nil {
			logger.Warn.Println("read file failed.", e.Error(), p)
			return
		}
		t := strings.TrimSpace(scanner.Text())
		m := pyRequirementsPattern.FindStringSubmatch(t)
		if m == nil {
			continue
		}
		rs = append(rs, model.Dependency{
			Name:    m[1],
			Version: m[2],
		})
	}
	return
}
