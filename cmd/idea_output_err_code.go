package cmd

type IdeaErrCode int

const (
	IdeaUnknownErr IdeaErrCode = iota - 2
	IdeaInspectErr
	IdeaSucceed
	IdeaEngineScanFailed
	IdeaServerRequestFailed
	IdeaUnknownProject
	IdeaTokenInvalid
	IdeaApiTimeout
	IdeaScanDirInvalid
	IdeaLogFileCreateErr
)

func (code IdeaErrCode) Error() string {
	switch code {
	case IdeaUnknownErr:
		return "UnknownErr"
	case IdeaInspectErr:
		return "InspectErr"
	case IdeaSucceed:
		return "Succeed"
	case IdeaEngineScanFailed:
		return "EngineScanFailed"
	case IdeaServerRequestFailed:
		return "ServerRequestFailed"
	case IdeaUnknownProject:
		return "UnknownProject"
	case IdeaTokenInvalid:
		return "TokenInvalid"
	case IdeaApiTimeout:
		return "ApiTimeout"
	case IdeaScanDirInvalid:
		return "ScanDirInvalid"
	case IdeaLogFileCreateErr:
		return "LogFileCreateErr"
	}
	panic(code)
}
