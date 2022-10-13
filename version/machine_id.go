package version

import (
	"github.com/iseki0/machineid"
	"github.com/murphysecurity/murphysec/utils/must"
)

func MachineId() string {
	return must.A(machineid.ProtectedID("murphysec"))
}
