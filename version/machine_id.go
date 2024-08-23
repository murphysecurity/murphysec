package version

import (
	"github.com/iseki0/machineid"
)

func MachineId() string {
	s, e := machineid.ProtectedID("murphysec")
	if e != nil {
		return "<NoMachineID>"
	}
	return s
}
