package common

import (
	"context"
	"errors"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/config"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/version"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var loggerInitialized = false

var (
	CliTokenOverride         string
	CliServerAddressOverride string
	EnableNetworkLogging     bool
	AllowInsecure            bool
	NoLogFile                bool
	LogFileOverride          string
	LogLevel                 logger.Level
)

func GetToken(ctx context.Context) (string, error) {
	if CliTokenOverride != "" {
		return CliTokenOverride, nil
	}
	if env.APITokenOverride != "" {
		return env.APITokenOverride, nil
	}
	return config.ReadTokenFile(ctx)
}

func InitAPIClient(ctx context.Context) error {
	if !loggerInitialized {
		panic("logger not initialized")
	}
	token, e := GetToken(ctx)
	if errors.Is(e, config.ErrNoToken) {
		displayTokenNotSet(ctx)
		return e
	}
	if e != nil {
		displayGetTokenErr(ctx, e)
		return e
	}
	var cf = &api.Config{
		Logger:             logctx.Use(ctx),
		EnableNetworkDebug: EnableNetworkLogging,
		Token:              token,
		AllowInsecure:      AllowInsecure,
	}
	if env.ServerURLOverride != "" {
		cf.ServerURL = env.ServerURLOverride
	}
	if CliServerAddressOverride != "" {
		cf.ServerURL = CliServerAddressOverride
	}
	e = api.InitDefaultClient(cf)
	if e != nil {
		displayInitializeFailed(ctx, e)
		return e
	}
	return nil
}

func InitLogger(ctx context.Context) (context.Context, error) {
	if loggerInitialized {
		panic("loggerInitialized == true")
	}

	consoleCore := zapcore.NewNopCore()
	jsonCore := zapcore.NewNopCore()

	// 如果日志文件没被禁用
	if !NoLogFile {
		// 创建日志文件
		logFile, e := logger.CreateLogFile(LogFileOverride)
		if e != nil {
			return nil, e
		}
		// 绑定日志core
		jsonCore = zapcore.NewCore(logger.ZapConsoleEncoder, logFile, zapcore.DebugLevel)
	}

	// 有关标准错误流的日志输出
	var stderr = zapcore.Lock(os.Stderr)
	if LogLevel > logger.LevelSilent {
		consoleCore = zapcore.NewCore(logger.ZapConsoleEncoder, stderr, LogLevel.ZapLevel())
	}

	loggerCore := zapcore.NewTee(consoleCore, jsonCore)
	_logger := zap.New(loggerCore, zap.AddCaller())

	_logger.Sugar().Infof("Log start: %s, %s", time.Now().Format(time.RFC3339), version.UserAgent())
	_logger.Sugar().Infof("Args: %s", utils.DesensitizedArgs)
	for _, s := range os.Environ() {
		_logger.Sugar().Debugf("Env: %s", s)
	}
	_logger.Sugar().Infof("Machine id: %s", version.MachineId())

	loggerInitialized = true
	return logctx.With(ctx, _logger), nil
}
