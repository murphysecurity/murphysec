package model

import "fmt"

type InspectError struct {
	Language string `json:"language"`
	Message  string `json:"message"`
}

func (i InspectError) Error() string {
	return fmt.Sprintf("%s: %s", i.Language, i.Message)
}

func NewInspectError(language Language, message string) error {
	return &InspectError{
		Language: string(language),
		Message:  message,
	}
}
