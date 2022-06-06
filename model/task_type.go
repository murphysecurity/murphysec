package model

import (
	"encoding/json"
	"murphysec-cli-simple/display"
)

type TaskType int

const (
	TaskTypeIdea TaskType = iota + 1
	TaskTypeCli
	TaskTypeJenkins
)

func (t TaskType) String() string {
	switch t {
	case TaskTypeCli:
		return "client"
	case TaskTypeIdea:
		return "plugin"
	case TaskTypeJenkins:
		return "jenkins"
	}
	panic(int(t))
}

func (t TaskType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t TaskType) UI() display.UI {
	switch t {
	case TaskTypeCli:
		return display.CLI
	case TaskTypeJenkins, TaskTypeIdea:
		return display.NONE
	}
	panic(t)
}
