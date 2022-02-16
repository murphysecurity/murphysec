package api

import (
	"encoding/json"
)

type InspectTaskSource int

const (
	TaskSourceIdea InspectTaskSource = iota + 1
	TaskSourceCli
	TaskSourceJenkins
)

func (receiver InspectTaskSource) String() string {
	switch receiver {
	case TaskSourceCli:
		return "client"
	case TaskSourceIdea:
		return "plugin"
	case TaskSourceJenkins:
		return "jenkins"
	}
	panic(int(receiver))
}

func (receiver InspectTaskSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(receiver.String())
}
