package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//GOOS=windows GOARCH=amd64 go build -v -ldflags "-s -w -X github.com/murphysecurity/murphysec/infra/buildinfo.version=$CI_COMMIT_TAG -buildid=" -trimpath -o out/bin/302-murphysec-windows-amd64.exe
//GOOS=linux GOARCH=amd64 go build -v -ldflags "-s -w -X github.com/murphysecurity/murphysec/infra/buildinfo.version=$CI_COMMIT_TAG -buildid=" -trimpath -o out/bin/302-murphysec-linux-amd64
//GOOS=linux GOARCH=arm64 go build -v -ldflags "-s -w -X github.com/murphysecurity/murphysec/infra/buildinfo.version=$CI_COMMIT_TAG -buildid=" -trimpath -o out/bin/302-murphysec-linux-arm64
//GOOS=darwin GOARCH=amd64 go build -v -ldflags "-s -w -X github.com/murphysecurity/murphysec/infra/buildinfo.version=$CI_COMMIT_TAG -buildid=" -trimpath -o out/bin/302-murphysec-darwin-amd64
//GOOS=darwin GOARCH=arm64 go build -v -ldflags "-s -w -X github.com/murphysecurity/murphysec/infra/buildinfo.version=$CI_COMMIT_TAG -buildid=" -trimpath -o out/bin/302-murphysec-darwin-arm64

func TestHash(t *testing.T) {
	dir, _ := os.Getwd()
	all := strings.ReplaceAll(dir, "utils", filepath.Join("out", "bin"))
	fileSystem := os.DirFS(all)
	if fileSystem != nil {
		tt := TT{
			PluginVersion: "3.0.6",
		}
		// nolint:all
		_ = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
			if d == nil {
				return nil
			}
			if d.IsDir() {
				return nil
			}

			file, err := os.Open(filepath.Join(all, d.Name()))
			hash := md5.New()
			if _, err := io.Copy(hash, file); err != nil {
				fmt.Println(err)
				return nil
			}

			// 计算哈希值并将其转换为十六进制字符串
			md5Hash := hex.EncodeToString(hash.Sum(nil))

			typ := ""
			if strings.Contains(d.Name(), "murphysec-windows") {
				typ = "win-amd64"
			} else if strings.Contains(d.Name(), "murphysec-darwin-amd64") {
				typ = "mac-amd64"
			} else if strings.Contains(d.Name(), "murphysec-darwin-arm64") {
				typ = "mac-arm64"

			} else if strings.Contains(d.Name(), "murphysec-linux-amd64") {
				typ = "linux-amd64"

			} else if strings.Contains(d.Name(), "murphysec-linux-arm64") {
				typ = "linux-arm64"

			}

			detail := Detail{
				Type:     typ,
				FileName: d.Name(),
				Hash:     md5Hash,
			}
			tt.Details = append(tt.Details, detail)

			return nil
		})

		marshal, _ := json.Marshal(tt)
		println(string(marshal))
	}

}

type TT struct {
	PluginVersion string   `json:"plugin_version"`
	Details       []Detail `json:"details"`
}
type Detail struct {
	Type     string `json:"type"`
	FileName string `json:"file_name"`
	Hash     string `json:"hash"`
}
