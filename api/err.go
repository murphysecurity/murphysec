package api

import "github.com/murphysecurity/murphysec/model"

type err string

func (e err) Error() string {
	return string(e)
}

var ErrTokenInvalid = model.WrapIdeaErr(err("Token invalid"), model.IdeaTokenInvalid)
var ErrServerRequest = model.WrapIdeaErr(err("Send request failed"), model.IdeaServerRequestFailed)
var UnprocessableResponse = model.WrapIdeaErr(err("Unprocessable response"), model.IdeaServerRequestFailed)
var ErrTimeout = model.WrapIdeaErr(err("API request timeout"), model.IdeaApiTimeout)
