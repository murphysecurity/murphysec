package cmd

type IdeaErrCode int

const (
	IdeaUnknownErr          IdeaErrCode = -2
	IdeaInspectErr          IdeaErrCode = -1
	IdeaSucceed             IdeaErrCode = 0
	IdeaEngineScanFailed    IdeaErrCode = 1
	IdeaServerRequestFailed IdeaErrCode = 2
	IdeaUnknownProject      IdeaErrCode = 3
	IdeaTokenInvalid        IdeaErrCode = 4
	IdeaApiTimeout          IdeaErrCode = 5
	IdeaScanDirInvalid      IdeaErrCode = 6
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
	}
	panic(code)
}
