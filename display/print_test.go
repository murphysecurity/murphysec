package display

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {
	CLI.Display(MsgError, "Error message")
	CLI.Display(MsgWarn, "Warn message")
	CLI.Display(MsgNotice, "Notice message")
	CLI.Display(MsgInfo, "Info message")
	var l MsgLevel
	assert.NoError(t, json.Unmarshal([]byte("\"info\""), &l))
	fmt.Println(l)
}
