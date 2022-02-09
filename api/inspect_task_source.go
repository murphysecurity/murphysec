package api

import "encoding/json"

type InspectTaskSource int

const (
	TaskSourceIdea InspectTaskSource = iota + 1
	TaskSourceCli
	TaskSourceCI
)

func (receiver InspectTaskSource) String() string {
	switch receiver {
	case TaskSourceCI:
		return "CI"
	case TaskSourceCli:
		return "CLI"
	case TaskSourceIdea:
		return "IDEA"
	}
	panic(receiver)
}

func (receiver InspectTaskSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(receiver.String())
}
