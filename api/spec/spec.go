package spec

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils/must"
	"io"
	"net/http"
	"sync"
)

//go:embed api.json.gz
var data []byte
var once sync.Once
var _oas *openapi3.T

func GetSpec() *openapi3.T {
	once.Do(func() {
		d := must.A(io.ReadAll(must.A(gzip.NewReader(bytes.NewReader(data)))))
		_oas = must.A(openapi3.NewLoader().LoadFromData(d))
	})
	return _oas
}

func Validate(ctx context.Context, request *http.Request, response *http.Response, respBody []byte) error {
	router := must.A(gorillamux.NewRouter(GetSpec()))
	var e error
	route, m, _ := router.FindRoute(request)
	if route == nil {
		return nil
	}
	options := &openapi3filter.Options{
		ExcludeRequestBody:    true,
		ExcludeResponseBody:   false,
		IncludeResponseStatus: true,
		MultiError:            true,
		AuthenticationFunc:    nil,
	}
	reqV := &openapi3filter.RequestValidationInput{
		Request:      request,
		PathParams:   m,
		QueryParams:  request.URL.Query(),
		Route:        route,
		Options:      options,
		ParamDecoder: nil,
	}
	resV := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: reqV,
		Status:                 response.StatusCode,
		Header:                 response.Header,
		Body:                   io.NopCloser(bytes.NewReader(respBody)),
		Options:                options,
	}
	e = openapi3filter.ValidateRequest(ctx, reqV)
	if e != nil {
		logctx.Use(ctx).Sugar().Warnf("request error: %s %v", request.URL.String(), e)
		return e
	}
	e = openapi3filter.ValidateResponse(ctx, resV)
	if e != nil {
		logctx.Use(ctx).Sugar().Warnf("response error: %s %v", request.URL.String(), e)
		return e
	}
	return nil
}
