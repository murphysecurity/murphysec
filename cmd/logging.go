package cmd

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/version"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var cliLogFilePathOverride string
var disableLogFile bool
var consoleLogLevelOverride string
var enableNetworkLog bool

var loggerInitialized = false

func mustInitLogger() {
	if e := initLogger(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		exitcode.Set(1)
		exitcode.Exit()
	}
}

func initLogger() error {
	if loggerInitialized {
		panic("loggerInitialized == true")
	}

	consoleCore := zapcore.NewNopCore()
	jsonCore := zapcore.NewNopCore()

	// 如果日志文件没被禁用
	if !disableLogFile {
		// 创建日志文件
		logFile, e := logger.CreateLogFile(cliLogFilePathOverride)
		if e != nil {
			return e
		}
		// 绑定日志core
		jsonCore = zapcore.NewCore(logger.ZapConsoleEncoder, logFile, zapcore.DebugLevel)
	}

	// 有关标准错误流的日志输出
	var stderr = zapcore.Lock(os.Stderr)
	if logLevel > logger.LevelSilent {
		consoleCore = zapcore.NewCore(logger.ZapConsoleEncoder, stderr, logLevel.ZapLevel())
	}

	loggerCore := zapcore.NewTee(consoleCore, jsonCore)
	_logger := zap.New(loggerCore, zap.AddCaller())

	//inspector.Logger = _logger       // todo: legacy
	logger.InitLegacyLogger(_logger) // todo: 是时候去掉了

	_logger.Sugar().Infof("Log start: %s, %s", time.Now().Format(time.RFC3339), version.UserAgent())
	_logger.Sugar().Infof("Args: %s", utils.DesensitizedArgs)
	for _, s := range os.Environ() {
		_logger.Sugar().Debugf("Env: %s", s)
	}
	_logger.Sugar().Infof("Machine id: %s", version.MachineId())

	rootCtx = logctx.With(context.TODO(), _logger)
	loggerInitialized = true
	return nil
}
