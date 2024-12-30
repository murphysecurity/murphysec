package python

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const pipConf = `[global]
index-url = https://%s/simple/
trusted-host = %s`

type PipdeptreeStruct struct {
	Key              string             `json:"key"`
	PackageName      string             `json:"package_name"`
	InstalledVersion string             `json:"installed_version"`
	RequiredVersion  string             `json:"required_version"`
	Dependencies     []PipdeptreeStruct `json:"dependencies"`
}

func getVenvPath(basePath string) string {
	goos := runtime.GOOS
	switch goos {
	case "darwin":
		path := filepath.Join(basePath, "/virtual_venv/bin")
		return path
	case "linux":
		path := filepath.Join(basePath, "/virtual_venv/bin")
		return path
	case "windows":
		path := filepath.Join(basePath, "/virtual_venv/Scripts")
		return path
	}

	return ""
}
func newVenv(dir string, logger *zap.SugaredLogger) error {
	var out bytes.Buffer
	var errout bytes.Buffer
	env := os.Environ()
	logger.Debug(zap.Any("env", env))
	cmd := exec.Command("bash", "-c", "/usr/local/python3.10/bin/python3.10 -m venv virtual_venv")
	cmd.Dir = dir
	cmd.Stdout = &out
	cmd.Stderr = &errout
	if err := cmd.Run(); err != nil {
		logger.Error("new venv error :", zap.String("venv", errout.String()))
		return err
	}
	logger.Debug("new venv success ")
	return nil
}
func newPipConf(basePath string, privateAddr string) error {
	goos := runtime.GOOS
	switch goos {
	case "darwin":
		path := filepath.Join(basePath, "/virtual_venv/pip.conf")
		if err := os.WriteFile(path, []byte(fmt.Sprintf(pipConf, privateAddr, privateAddr)), 0777); err != nil {
			return err
		}
		return nil
	case "linux":
		path := filepath.Join(basePath, "/virtual_venv/pip.conf")
		if err := os.WriteFile(path, []byte(fmt.Sprintf(pipConf, privateAddr, privateAddr)), 0777); err != nil {
			return err
		}
		return nil
	case "windows":
		path := filepath.Join(basePath, "/virtual_venv/pip.ini")
		if err := os.WriteFile(path, []byte(fmt.Sprintf(pipConf, privateAddr, privateAddr)), 0777); err != nil {
			return err
		}
		return nil
	}
	return nil
}
func updatePip(dir string, logger *zap.SugaredLogger) error {
	var out bytes.Buffer
	var errout bytes.Buffer
	cmd := exec.Command("./python3.10", "-m", "pip", "install", "--upgrade", "pip")
	cmd.Stdout = &out
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		logger.Error("pip update error :", zap.String("pip", errout.String()))
		return err
	}
	logger.Debug("pip update success ")
	return nil
}

func pipreqs(dir string, projectPath, savePath string, logger *zap.SugaredLogger) error {
	savePath = filepath.Join(savePath, "requirements.txt")
	logger.Debug(zap.String("pipreqs Path", dir))
	logger.Debug(zap.String("pipreqs projectPath", projectPath))
	logger.Debug(zap.String("pipreqs savepath", savePath))
	cmd := exec.Command("./pipreqs", projectPath, "--savepath", savePath, "--encoding=utf-8", "--ignore=virtual_venv")
	cmd.Dir = dir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("create stdout pipe failed: %w", err)
		logger.Error(err.Error())
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		err = fmt.Errorf("create stderr pipe failed: %w", err)
		logger.Error(err.Error())
		return err
	}

	go func() {
		defer stderr.Close()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			logger.Warn("go: " + scanner.Text())
		}
	}()

	if err := cmd.Start(); err != nil {
		err = fmt.Errorf("start command failed: %w", err)
		logger.Error(err.Error())
		return err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		logger.Debug("pipreqs info: " + scanner.Text())
	}
	return nil
}
func installpipreqs(dir string, logger *zap.SugaredLogger) error {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("./pip", "install", "pipreqs")
	cmd.Dir = dir
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		logger.Error(zap.String("installpipreqs error :", stderr.String()))
		return err
	}
	logger.Debug("install pipreqs success ")
	return nil
}
func installRequirements(dir string, textDir string, logger *zap.SugaredLogger) error {
	var out bytes.Buffer
	var stderr bytes.Buffer
	by, err := os.ReadFile(textDir)
	if err != nil {
		logger.Error("read requirements.txt error :", zap.Error(err))
		return err
	}
	nvmp := parseRequirements(string(by))
	for k, v := range nvmp {
		var cmd *exec.Cmd
		if v != "" {
			cmd = exec.Command("./pip", "install", k+"=="+v)
		} else {
			cmd = exec.Command("./pip", "install", k)
		}
		cmd.Dir = dir
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			logger.Error("install requirements error :", zap.String("stderr", stderr.String()))
		}
		logger.Debug("install requirements success ")
	}
	return nil
}
func installpipdeptree(dir string, logger *zap.SugaredLogger) error {
	var out bytes.Buffer
	cmd := exec.Command("./pip", "install", "pipdeptree")
	cmd.Dir = dir
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		logger.Error("install pipdeptree error :", zap.Error(err))
		return err
	}
	logger.Debug("install pipdeptree success ")
	return nil
}
func pipdeptree(dir string, logger *zap.SugaredLogger) ([]PipdeptreeStruct, error) {
	var out bytes.Buffer
	var errout bytes.Buffer
	var result []PipdeptreeStruct
	cmd := exec.Command("./pipdeptree", "--json-tree")
	cmd.Stdout = &out
	cmd.Dir = dir
	cmd.Stderr = &errout
	if err := cmd.Run(); err != nil {
		logger.Debug("pipdeptree path ", zap.String("exec:", dir))
		logger.Error("pipdeptree error :", zap.String("pipdeptree ", errout.String()))
		return nil, err
	}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		logger.Error("pipdeptree json unmarshal error :", zap.Error(err))
		return nil, err
	}
	logger.Debug("pipdeptree success ", zap.String("exec:", out.String()))
	return result, nil
}
func updatePackage(dir string, logger *zap.SugaredLogger, k, v string) {
	var out bytes.Buffer
	version := k
	if v != "" {
		version = version + "==" + v
	}
	cmd := exec.Command("./pip", "install", version)
	cmd.Stdout = &out
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		logger.Error("update pip install error :", zap.Error(err))
	}
	logger.Debug("update pip install success :" + k + "==" + v)
}
func buildTree(pipdeptree PipdeptreeStruct, level int) model.DependencyItem {
	directDependency := false
	if level == 0 {
		directDependency = true
	}
	var mod = model.DependencyItem{
		Component: model.Component{
			CompName:    pipdeptree.Key,
			CompVersion: pipdeptree.InstalledVersion,
			EcoRepo: model.EcoRepo{
				Ecosystem:  "pip",
				Repository: "",
			},
		},
		IsDirectDependency: directDependency,
	}
	for _, i := range pipdeptree.Dependencies {
		mod.Dependencies = append(mod.Dependencies, buildTree(i, level+1))
	}
	return mod
}
func delVenv(dir string, logger *zap.SugaredLogger) {
	if err := os.RemoveAll(dir); err != nil {
		logger.Error("delete venv error :", zap.Error(err))
		return
	}
	logger.Debug("delete venv success ")
}
func directDependenceSurvival(mod *[]model.DependencyItem, nvMp map[string]string) {
	var exist = make(map[string]string)
	for _, i := range *mod {
		exist[i.CompName] = i.CompVersion
	}
	for k, v := range nvMp {
		if _, ok := exist[k]; !ok {
			*mod = append(*mod, model.DependencyItem{
				Component: model.Component{
					CompName:    k,
					CompVersion: v,
					EcoRepo: model.EcoRepo{
						Ecosystem:  "pip",
						Repository: "",
					},
				},
				IsDirectDependency: true,
			})
		}
	}
}
func pipenv() string {
	return os.Getenv("PIP_SOURCE_ADDR")
}
func Run(ctx context.Context, dir string, logger *zap.SugaredLogger, nvMp map[string]string) ([]model.DependencyItem, error) {
	var mod []model.DependencyItem
	var venvDir = filepath.Join(dir, "virtual_venv")
	venvPath := getVenvPath(dir)
	requirementsPath := filepath.Join(dir, "requirements.txt")
	venvRequirementsPath := filepath.Join(venvPath, "requirements.txt")
	if err := newVenv(dir, logger); err != nil {
		return nil, err
	}
	if privatePath, ok := ctx.Value("privateSourceAddr").(string); ok {
		logger.Debug("Use private path", zap.String("path", privatePath))
		if err := newPipConf(venvPath, privatePath); err != nil {
			return nil, err
		}
	}
	if envSource := pipenv(); envSource != "" {
		logger.Debug("Use private path", zap.String("path", envSource))
		if err := newPipConf(venvPath, envSource); err != nil {
			return nil, err
		}
	}
	if err := updatePip(venvPath, logger); err != nil {
		return nil, err
	}
	if err := installpipreqs(venvPath, logger); err != nil {
		return nil, err
	}
	if err := pipreqs(venvPath, dir, venvPath, logger); err != nil {
		return nil, err
	}
	if err := installRequirements(venvPath, venvRequirementsPath, logger); err != nil {
		return nil, err
	}
	if err := installpipdeptree(venvPath, logger); err != nil {
		return nil, err
	}
	// 读取新创建的 requirements.txt
	data, err := readTextFile(requirementsPath, 64*1024)
	if err != nil {
		logger.Warnf("read requirement: %s %v", requirementsPath, err)
	}
	//  对比原本的 requirements.txt 拿到原本包的版本
	newRequirements := parseRequirements(string(data))
	for k, v := range nvMp {
		if newV, ok := newRequirements[k]; ok && newV != v {
			updatePackage(venvPath, logger, k, v)
		}
	}
	result, err := pipdeptree(venvPath, logger)
	if err != nil {
		return nil, err
	}
	for _, j := range result {
		if j.Key != "" && j.Key != "pipdeptree" && j.Key != "pipreqs" && j.Key != "pip" {
			mod = append(mod, buildTree(j, 0))
		}
	}
	// 对于没有pip install成功的依赖 只加入pipreqs中列出的直接依赖
	directDependenceSurvival(&mod, newRequirements)
	defer delVenv(venvDir, logger)
	return mod, err
}
