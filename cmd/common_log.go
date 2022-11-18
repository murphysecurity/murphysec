package cmd

import (
	"fmt"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/version"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

var cliLogFilePathOverride string
var disableLogFile bool
var consoleLogLevelOverride string
var enableNetworkLog bool

var LOG = zap.NewNop()
var SLOG = zap.NewNop().Sugar()

func initConsoleLoggerOrExit() {
	e := initLogger()
	if e == nil {
		return
	}
	fmt.Println("Error:", e.Error())
	os.Exit(1)
}

func initLogger() error {
	consoleCore := zapcore.NewNopCore()
	jsonCore := zapcore.NewNopCore()

	if !disableLogFile {
		logFile, e := logger.CreateLogFile(cliLogFilePathOverride)
		if e != nil {
			return e
		}
		jsonCore = zapcore.NewCore(logger.ZapConsoleEncoder, logFile, zapcore.DebugLevel)
	}

	var stderr = zapcore.Lock(os.Stderr)
	switch strings.ToLower(strings.TrimSpace(consoleLogLevelOverride)) {
	case "error":
		consoleCore = zapcore.NewCore(logger.ZapConsoleEncoder, stderr, zapcore.ErrorLevel)
	case "warn":
		consoleCore = zapcore.NewCore(logger.ZapConsoleEncoder, stderr, zapcore.WarnLevel)
	case "info":
		consoleCore = zapcore.NewCore(logger.ZapConsoleEncoder, stderr, zapcore.InfoLevel)
	case "debug":
		consoleCore = zapcore.NewCore(logger.ZapConsoleEncoder, stderr, zapcore.DebugLevel)
	}

	loggerCore := zapcore.NewTee(consoleCore, jsonCore)
	_logger := zap.New(loggerCore, zap.AddCaller())

	if enableNetworkLog {
		api.NetworkLogger = _logger.Named("Net").WithOptions(zap.WithCaller(false))
	}
	LOG = _logger
	SLOG = LOG.Sugar()
	api.Logger = _logger
	inspector.Logger = _logger
	model.Logger = _logger
	logger.InitLegacyLogger(_logger)

	_logger.Sugar().Infof("Log start: %s, %s", time.Now().Format(time.RFC3339), version.UserAgent())
	_logger.Sugar().Infof("Args: %s", utils.DesensitizedArgs)
	for _, s := range os.Environ() {
		_logger.Sugar().Debugf("Env: %s", s)
	}
	_logger.Sugar().Infof("Machine id: %s", version.MachineId())
	return nil
}
