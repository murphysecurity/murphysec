package inspector

import (
	"fmt"
	"os"
	"testing"
)

func TestFileScan(t *testing.T) {
	iter := FileScan(os.TempDir())
	for iter.Next() {
		if iter.Err() != nil {
			fmt.Println("err:", iter.Err().Error())
		} else {
			fmt.Println("path:", iter.Path())
		}
	}
}

func TestFileHashInspectScan(t *testing.T) {
	FileHashInspectScan(".")
}
