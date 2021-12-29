package logger

import (
	"encoding/json"
	"io"
	"murphysec-cli-simple/util/must"
	"murphysec-cli-simple/version"
	"os"
	"sync"
)

var outputMutex = sync.Mutex{}
var stdout = os.Stdout
var wg = sync.WaitGroup{}

func writeOutput(m ManagedOutput) {
	outputMutex.Lock()
	defer outputMutex.Unlock()
	must.Int(stdout.Write(must.Byte(json.Marshal(m))))
	must.Int(stdout.Write([]byte("\n")))
}

func InitManagedMode() {
	// 接管 stdout
	fakeStdoutReader, fakeStdoutWriter, e := os.Pipe()
	must.Must(e)
	os.Stdout = fakeStdoutWriter
	// sending version info
	writeOutput(ManagedOutput{VersionInfo: map[string]interface{}{"version": version.Version()}})
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			buf := make([]byte, 1024)
			n, e := fakeStdoutReader.Read(buf)
			if e == nil {
				writeOutput(ManagedOutput{Stdout: buf[0:n]})
				continue
			}
			if e == io.EOF {
				writeOutput(ManagedOutput{StdoutEOF: true})
				_ = stdout.Close()
				break
			}
			panic(e)
		}
	}()
}

type ManagedOutput struct {
	Stdout      []byte                 `json:"stdout,omitempty"`
	StdoutEOF   bool                   `json:"stdout_eof,omitempty"`
	VersionInfo map[string]interface{} `json:"version_info,omitempty"`
}

func CloseAndWait() {
	must.Must(os.Stdout.Close())
	wg.Wait()
}
