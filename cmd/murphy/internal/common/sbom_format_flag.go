package common

import (
	"fmt"
	"github.com/spf13/pflag"
)


type SBOMFormatFlag struct {
	Valid bool
}

func (t *SBOMFormatFlag) String() string {
	if !t.Valid {
		return "invalid"
	}
	return "murphysec1.1+json"
}

func (t *SBOMFormatFlag) Set(s string) error {
	// temporary implementation
	if s!="murphysec1.1+json" {
		return fmt.Errorf("unsupported format: %s", s)
	}
	t.Valid = true
	return nil
}

func (t *SBOMFormatFlag) Type() string {
	return "sbomFormatFlag"
}

var _ pflag.Value = (*SBOMFormatFlag)(nil)
