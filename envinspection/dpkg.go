package envinspection

import (
	"context"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"os/exec"
	"strings"
)

func inspectDpkgPackage(ctx context.Context) ([]model.Dependency, error) {
	LOG := utils.UseLogger(ctx)
	cmd := exec.Command("dpkg-query", "-W", "-f", "${binary:Package} ${Version}\\n")
	LOG.Sugar().Infof("Execute: %s", cmd.String())
	data, e := cmd.Output()
	if e != nil {
		return nil, e
	}
	var rs []model.Dependency
	for _, s := range strings.Split(string(data), "\n") {
		s = strings.TrimSpace(s)
		chunks := strings.Split(s, " ")
		if len(chunks) != 2 {
			continue
		}
		rs = append(rs, model.Dependency{
			Name:    chunks[0],
			Version: chunks[1],
		})
	}
	return rs, nil
}
