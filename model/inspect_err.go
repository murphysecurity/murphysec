package model

import "fmt"

type InspectError struct {
	Language string `json:"language"`
	Message  string `json:"message"`
}

func (i InspectError) Error() string {
	return fmt.Sprintf("%s: %s", i.Language, i.Message)
}
