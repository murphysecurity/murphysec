package common

import (
	"github.com/murphysecurity/murphysec/logger"
	"github.com/spf13/pflag"
	"strings"
)

type LogLevelFlag struct {
	Level logger.Level
	Valid bool
}

var _ pflag.Value

func (l *LogLevelFlag) String() string {
	if !l.Valid {
		return "unset"
	}
	return l.Level.String()
}
func (l *LogLevelFlag) Set(s string) error {
	var ll logger.Level
	if e := ll.Of(strings.ToLower(s)); e != nil {
		return e
	}
	l.Level = ll
	l.Valid = true
	return nil
}
func (l *LogLevelFlag) Type() string {
	return "logLevelFlag"
}
