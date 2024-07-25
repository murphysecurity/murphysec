package nuget

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/murphysecurity/murphysec/utils"
)

// 判断项目中是否有packages.config文件
func checkPackagesIsExistence(fileName string) bool {
	return utils.IsFile(fileName)
}

// 判断nuget命令是否存在
func checkNugetCommand() bool {
	_, err := exec.LookPath("nuget")
	return err == nil
}
func findCLN(dir string) (filePath string) {

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".sln" {
			// 找到.sln文件所在的目录
			filePath = filepath.Dir(path)
			return nil
		}
		return nil // 继续搜索
	})

	return filePath
}
func findCLNList(dir string) (filePath []string) {

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".sln" {
			// 找到.sln文件所在的目录
			filePath = append(filePath, filepath.Dir(path))
			return nil
		}
		return nil // 继续搜索
	})

	return filePath
}
