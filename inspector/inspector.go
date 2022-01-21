package inspector

import (
	"github.com/pkg/errors"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/module/go_mod"
	"murphysec-cli-simple/module/maven"
	"murphysec-cli-simple/module/npm"
	"murphysec-cli-simple/utils/must"
	"time"
)

var engines = []base.Inspector{
	go_mod.New(),
	maven.New(),
	npm.New(),
}

func tryMatchInspector(dir string) base.Inspector {
	for _, it := range engines {
		logger.Debug.Println("Try match project by inspector:", it.String(), "...")
		if it.CheckDir(dir) {
			logger.Info.Println("Matched.")
			return it
		}
	}
	logger.Debug.Println("Match failed")
	return nil
}

var ErrNoEngineMatched = errors.New("ErrNoEngineMatched")

func autoInspectDir(dir string) ([]base.Module, error) {
	startTime := time.Now()
	logger.Info.Println("Auto scan dir:", dir)
	engine := tryMatchInspector(dir)
	if engine == nil {
		return nil, ErrNoEngineMatched
	}
	logger.Info.Println("Engine matched.", engine.String())
	modules, e := engine.Inspect(dir)
	if e != nil {
		logger.Warn.Println("Engine report some error.", e.Error())
		return nil, e
	}
	endTime := time.Now()
	must.True(!endTime.Before(startTime))
	logger.Info.Println("Scan terminated. Cost time:", endTime.Sub(startTime))
	return modules, nil
}