package version

import (
	"github.com/denisbrodbeck/machineid"
)

func MachineId() string {
	s, e := machineid.ProtectedID("murphysec")
	if e != nil {
		return "<NoMachineID>"
	}
	return s
}
