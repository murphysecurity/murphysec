package maven

import (
	"fmt"
	"github.com/murphysecurity/murphysec/utils"
)

const MaxPrefixSuffix = 2048

type MvnCmdExecutionStreamHandler struct {
	suffix *utils.SuffixBuffer
}

func (m MvnCmdExecutionStreamHandler) String() string {
	return fmt.Sprintf("suffix: %s", string(m.suffix.Bytes()))
}

func (m *MvnCmdExecutionStreamHandler) Write(input []byte) (int, error) {
	_, _ = m.suffix.Write(input)
	return len(input), nil
}

func NewMvnCmdExecution() *MvnCmdExecutionStreamHandler {
	return &MvnCmdExecutionStreamHandler{
		suffix: utils.MkSuffixBuffer(MaxPrefixSuffix),
	}
}
