package api

import (
	"bytes"
	"encoding/json"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"net/http"
)

func CompImportPath(data interface{}) {
	uri := serverAddress() + "/message/v1/access/detect/import_path"
	logger.Info.Println("Call API:", uri)
	logger.Debug.Println("API request body:")
	requestData := must.Byte(json.Marshal(data))
	logger.Debug.Println(string(requestData))
	body := bytes.NewReader(requestData)
	//body := new(bytes.Buffer)
	//g := gzip.NewWriter(body)
	//must.Int(g.Write(requestData))
	//must.Close(g)
	req, e := http.NewRequest(http.MethodPost, uri, body)
	//req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	must.Must(e)
	res, e := client.Do(req)
	if e != nil {
		logger.Err.Println("API request, failed.", e.Error())
	}
	if res.StatusCode != http.StatusOK {
		logger.Err.Println("API status code!=OK, code:", res.StatusCode)
	} else {
		logger.Info.Println("API request succeed")
	}
}
