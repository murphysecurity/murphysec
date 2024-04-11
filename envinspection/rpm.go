package envinspection

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"os/exec"
	"strings"
)

func listRPMPackage(ctx context.Context) ([]model.DependencyItem, error) {
	LOG := logctx.Use(ctx)
	cmd := exec.Command("rpm", "-qa", "--qf", "%{NAME} %{VERSION}-%{RELEASE}\\n")
	LOG.Sugar().Infof("Execute: %s", cmd.String())
	data, e := cmd.Output()
	if e != nil {
		return nil, e
	}
	var rs []model.DependencyItem
	for _, s := range strings.Split(string(data), "\n") {
		s = strings.TrimSpace(s)
		chunks := strings.Split(s, " ")
		if len(chunks) != 2 {
			continue
		}
		rs = append(rs, model.DependencyItem{
			Component: model.Component{
				CompName:    chunks[0],
				CompVersion: chunks[1],
			},
		})
	}
	return rs, nil
}
