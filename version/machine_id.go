package version

import (
	"github.com/denisbrodbeck/machineid"
	"murphysec-cli-simple/utils/must"
)

func MachineId() string {
	return must.A(machineid.ProtectedID("murphysec"))
}
