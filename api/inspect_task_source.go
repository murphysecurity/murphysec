package api

import (
	"encoding/json"
)

type InspectTaskSource int

const (
	TaskSourceIdea InspectTaskSource = iota + 1
	TaskSourceCli
)

func (receiver InspectTaskSource) String() string {
	switch receiver {
	case TaskSourceCli:
		return "client"
	case TaskSourceIdea:
		return "plugin"
	}
	panic(int(receiver))
}

func (receiver InspectTaskSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(receiver.String())
}
