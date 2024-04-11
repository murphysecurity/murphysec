package shared

import (
	"bytes"
	"gopkg.in/yaml.v3"
)

func ParseYaml(data []byte, target any) error {
	return yaml.Unmarshal(bytes.ReplaceAll(data, []byte{'?'}, []byte("(QuestionMark)")), target)
}
