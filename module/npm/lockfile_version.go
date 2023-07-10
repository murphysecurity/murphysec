package npm

import (
	"encoding/json"
	"fmt"
)

func parseLockfileVersion(data []byte) (int, error) {
	type unknownVersionLockfile struct {
		LockfileVersion int `json:"lockfileVersion"`
	}
	var u unknownVersionLockfile
	if e := json.Unmarshal(data, &u); e != nil {
		return 0, fmt.Errorf("parsing lockfile version: %w", e)
	}
	return u.LockfileVersion, nil
}
