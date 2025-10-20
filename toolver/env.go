package toolver

import (
	"encoding/json"
	"os"
	"sync"
)

var envOnce sync.Once
var envConfig *Config

func getFromEnv() *Config {
	envOnce.Do(func() {
		var s = os.Getenv("AUTO_BUILD_TOOLCHAIN_CONFIG")
		if s != "" {
			var o Config
			if json.Unmarshal([]byte(s), &o) != nil {
				envConfig = &o
			}
		}
	})
	return envConfig
}
