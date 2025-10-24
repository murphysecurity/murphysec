package toolver

import (
	"encoding/json"
	"os"
	"sync"
)

var envOnce sync.Once
var envConfig Config

func getFromEnv() *Config {
	envOnce.Do(func() {
		var s = os.Getenv("AUTO_BUILD_TOOLCHAIN_CONFIG")
		if s != "" {
			if e := json.Unmarshal([]byte(s), &envConfig); e != nil {
				panic(e)
			}
		}
	})
	return &envConfig
}
