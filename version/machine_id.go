package version

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/murphysecurity/murphysec/utils/must"
)

func MachineId() string {
	return must.A(machineid.ProtectedID("murphysec"))
}
