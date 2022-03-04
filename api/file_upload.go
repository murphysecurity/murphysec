package api

import (
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func UploadFile(taskId string, filePath string, baseDir string) error {
	f, e := os.Open(filePath)
	if e != nil {
		return e
	}
	defer f.Close()
	relPath := filepath.ToSlash(must.String(filepath.Rel(baseDir, filePath)))
	u := must.Url(url.Parse(serverAddress() + "/message/v2/access/client/upload_check_files"))
	v := u.Query()
	v.Set("task_info", taskId)
	v.Set("path", relPath)
	u.RawQuery = v.Encode()
	logger.Info.Println("Upload file to:", u.String())
	resp, e := http.Post(u.String(), "application/octet-stream", f)
	if e != nil {
		return ErrSendRequest
	}
	if resp.StatusCode == 200 {
		return nil
	}
	data, e := readHttpBody(resp)
	if e != nil {
		return e
	}
	return readCommonErr(data, resp.StatusCode)
}
