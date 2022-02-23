package api

import (
	"encoding/json"
)

type InspectTaskType int

const (
	TaskTypeIdea InspectTaskType = iota + 1
	TaskTypeCli
	TaskTypeJenkins
)

func (receiver InspectTaskType) String() string {
	switch receiver {
	case TaskTypeCli:
		return "client"
	case TaskTypeIdea:
		return "plugin"
	case TaskTypeJenkins:
		return "jenkins"
	}
	panic(int(receiver))
}

func (receiver InspectTaskType) MarshalJSON() ([]byte, error) {
	return json.Marshal(receiver.String())
}
