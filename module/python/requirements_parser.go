package python

import (
	"bufio"
	"context"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"io"
	"os"
	"regexp"
	"strings"
)

func parsePythonRequirements(ctx context.Context, p string) (rs []model.Dependency) {
	logger := utils.UseLogger(ctx)
	logger.Sugar().Debugf("Parsing python requirements: %s", p)
	var pyRequirementsPattern = regexp.MustCompile("^([\\w-]+) *.?= *([^= \\n\\r]+)$")
	f, e := os.Open(p)
	if e != nil {
		logger.Sugar().Warnf("Open file failed: %s", e.Error())
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(io.LimitReader(f, 4*1024*1024))
	for scanner.Scan() {
		if scanner.Err() != nil {
			logger.Sugar().Warnf("Scan requirements failed: %s", e.Error())
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
