package api

import (
	"bytes"
	"encoding/json"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"net/http"
)

func CompImportPath(data interface{}){
	uri:=serverAddress()+"/api/v1/access/detect/import_path"
	logger.Info.Println("Call API:", uri)
	logger.Debug.Println("API request body:")
	body:=must.Byte(json.Marshal(data))
	logger.Debug.Println(string(body))
	req, e:=http.NewRequest(http.MethodPost, uri, bytes.NewReader(body))
	must.Must(e)
	res, e:=client.Do(req)
	if e!=nil{
		logger.Err.Println("API request, failed.", e.Error())
	}
	if res.StatusCode!=http.StatusOK{
		logger.Err.Println("API status code!=OK, code:", res.StatusCode)
	}
}
