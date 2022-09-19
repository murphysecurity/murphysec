//go:build murphydev

package module

import "github.com/murphysecurity/murphysec/module/renv"

func init() {
	Inspectors = append(Inspectors, renv.Instance)
}
