package pnpm

import (
	"fmt"
	"github.com/murphysecurity/murphysec/module/pnpm/shared"
	"regexp"
	"strconv"
)

type lockfileVersionIndicator struct {
	LockfileVersion string `json:"lockfile_version" yaml:"lockfileVersion"`
}

func parseLockfileVersion(data []byte) (string, error) {
	var indicator lockfileVersionIndicator
	if e := shared.ParseYaml(data, &indicator); e != nil {
		return "", fmt.Errorf("parseLockfileVersion: %w", e)
	}
	return indicator.LockfileVersion, nil
}

func matchLockfileVersion(s string) int {
	d := regexp.MustCompile(`^v?(\d+)(\.|$)`).FindStringSubmatch(s)
	if d != nil {
		n, _ := strconv.Atoi(d[1])
		return n
	}
	return 0
}
