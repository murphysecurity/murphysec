//go:build embedding

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/plugin/plugin_base"
	"murphysec-cli-simple/util/must"
	"murphysec-cli-simple/util/output"
	"path/filepath"
)

func scanByPlugin(p plugin_base.Plugin, dir string) error {
	dir = must.String(filepath.Abs(dir))
	output.Info(fmt.Sprintf("Scan dir: %s", dir))
	if !p.MatchPath(dir) {
		return errors.New(fmt.Sprintf("The project can't be processed by plugin %s.", p.Info().Name))
	}
	packageInfo, err := p.DoScan(dir)
	if err != nil {
		return err
	}
	fmt.Println(string(must.Byte(json.Marshal(packageInfo))))
	return nil
}
