package pnpm

import (
	"fmt"
	"gopkg.in/yaml.v3"
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
