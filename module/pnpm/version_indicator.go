package pnpm

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
)

type lockfileVersionIndicator struct {
	LockfileVersion string `json:"lockfile_version" yaml:"lockfileVersion"`
}

func parseLockfileVersion(data []byte) (string, error) {
	var indicator lockfileVersionIndicator
	if e := yaml.Unmarshal(data, &indicator); e != nil {
		return "", fmt.Errorf("parseLockfileVersion: %w", e)
	}
	return indicator.LockfileVersion, nil
}

func matchLockfileVersion(s string) int {
	if regexp.MustCompile(`^v?5\.`).MatchString(s) {
		return 5
	}
	if regexp.MustCompile(`^v?6\.`).MatchString(s) {
		return 6
	}
	return 0
}
