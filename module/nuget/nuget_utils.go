package nuget

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/murphysecurity/murphysec/infra/pathignore"
	"github.com/murphysecurity/murphysec/utils"
)

// 判断项目中是否有packages.config文件
func checkPackagesIsExistence(fileName string) bool {
	return utils.IsFile(fileName)
}

func findCLNList(dir string) (filePath []string, err error) {

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if pathignore.DirName(info.Name()) {
				fmt.Println(info.Name())
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".sln" {
			// 找到.sln文件所在的目录
			filePath = append(filePath, filepath.Dir(path))
			return nil
		}
		return nil // 继续搜索
	})

	return filePath, err
}
