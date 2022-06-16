package model

import (
	"fmt"
	"github.com/pkg/errors"
)

type IdeaErrCode int

var _errCodeList = map[IdeaErrCode]string{
	IdeaSucceed:             "Succeeded",
	IdeaUnknownErr:          "UnknownError",
	IdeaInspectErr:          "InspectError",
	IdeaEngineScanFailed:    "EngineScanFailed",
	IdeaServerRequestFailed: "ServerRequestFailed",
	IdeaTokenInvalid:        "TokenInvalid",
	IdeaApiTimeout:          "APITimeout",
	IdeaScanDirInvalid:      "ScanDirInvalid",
	IdeaLogFileCreateFailed: "LogFileCreateFailed",
}

var _errCodeRMap = map[string]IdeaErrCode{}

func init() {
	for code, s := range _errCodeList {
		_errCodeRMap[s] = code
	}
}

const (
	IdeaSucceed IdeaErrCode = iota + 100
	IdeaUnknownErr
	IdeaInspectErr
	IdeaEngineScanFailed
	IdeaServerRequestFailed
	IdeaTokenInvalid
	IdeaApiTimeout
	IdeaScanDirInvalid
	IdeaLogFileCreateFailed
)

func (code IdeaErrCode) Error() string {
	s, _ := _errCodeList[code]
	return s
}

func (code IdeaErrCode) String() string {
	return code.Error()
}

func (code *IdeaErrCode) UnmarshalText(data []byte) error {
	c, ok := _errCodeRMap[string(data)]
	if !ok {
		return fmt.Errorf("bad IdeaErrCode: %s", string(data))
	}
	*code = c
	return nil
}

type ideaErr struct {
	error
	Code IdeaErrCode
}

func (e *ideaErr) Error() string {
	if e.error != nil {
		return e.error.Error()
	}
	return e.Code.Error()
}

func (e *ideaErr) Unwrap() error {
	return e.error
}

func (e *ideaErr) Is(target error) bool {
	return e.Code == target
}

func WrapIdeaErr(e error, ideaCode IdeaErrCode) error {
	return &ideaErr{e, ideaCode}
}

func GetIdeaErrCode(e error) IdeaErrCode {
	var t *ideaErr
	if errors.As(e, &t) {
		return t.Code
	}
	return IdeaUnknownErr
}
