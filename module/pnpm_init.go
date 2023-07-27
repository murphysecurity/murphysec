package module

import (
	"github.com/murphysecurity/murphysec/module/pnpm"
)

func init() {
	Inspectors = append(Inspectors, &pnpm.Inspector{})
}
