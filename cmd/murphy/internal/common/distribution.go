package common

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

type DistributionFlag struct {
	name string
}

func (d *DistributionFlag) String() string {
	if d == nil {
		return ""
	}
	return d.name
}

func (d *DistributionFlag) Set(s string) error {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "external", "internal", "saas", "open_source", "":
		d.name = s
	default:
		return fmt.Errorf("invalid distribution: %s", strconv.Quote(s))
	}
	return nil
}

func (d *DistributionFlag) Type() string {
	return "DistributionFlag"
}

func (d *DistributionFlag) IsExternal() bool {
	return d.name == "external"
}

func (d *DistributionFlag) IsInternal() bool {
	return d.name == "internal"
}

func (d *DistributionFlag) IsSaaS() bool {
	return d.name == "saas"
}

func (d *DistributionFlag) IsOpenSource() bool {
	return d.name == "open_source"
}

func (d DistributionFlag) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.name)
}

func (d *DistributionFlag) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return d.Set(s)
}

var _ pflag.Value = (*DistributionFlag)(nil)
var _ json.Marshaler = (*DistributionFlag)(nil)
var _ json.Unmarshaler = (*DistributionFlag)(nil)
