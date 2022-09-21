package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/conf"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/murphysecurity/murphysec/version"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var defaultLogFile = filepath.Join(must.A(homedir.Dir()), ".murphysec", "logs", fmt.Sprintf("%d.log", time.Now().UnixMilli()))
var cliLogFilePathOverride string
var disableLogFile bool
var consoleLogLevelOverride string
var enableNetworkLog bool

func initConsoleLoggerOrExit() {
	e := initLogger()
	if e == nil {
		return
	}
	fmt.Println("Error:", e.Error())
	os.Exit(1)
}

func initLogger() error {
	var loggerEncoderConfig = zapcore.EncoderConfig{
		MessageKey:          "message",
		LevelKey:            "level",
		TimeKey:             "time",
		NameKey:             "name",
		CallerKey:           "caller",
		FunctionKey:         "",
		StacktraceKey:       "stacktrace",
		SkipLineEnding:      false,
		LineEnding:          zapcore.DefaultLineEnding,
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          zapcore.RFC3339TimeEncoder,
		EncodeDuration:      zapcore.StringDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		EncodeName:          zapcore.FullNameEncoder,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    " ",
	}
	var consoleEncoder = zapcore.NewConsoleEncoder(loggerEncoderConfig)

	consoleCore := zapcore.NewNopCore()
	jsonCore := zapcore.NewNopCore()

	logFile, e := createLogFile()
	if e != nil && !errors.Is(e, ErrLogFileDisabled) {
		return e
	}
	if e == nil {
		jsonCore = zapcore.NewCore(consoleEncoder, logFile, zapcore.DebugLevel)
	}

	var stderr = zapcore.Lock(os.Stderr)
	switch strings.ToLower(strings.TrimSpace(consoleLogLevelOverride)) {
	case "error":
		consoleCore = zapcore.NewCore(consoleEncoder, stderr, zapcore.ErrorLevel)
	case "warn":
		consoleCore = zapcore.NewCore(consoleEncoder, stderr, zapcore.WarnLevel)
	case "info":
		consoleCore = zapcore.NewCore(consoleEncoder, stderr, zapcore.InfoLevel)
	case "debug":
		consoleCore = zapcore.NewCore(consoleEncoder, stderr, zapcore.DebugLevel)
	}

	loggerCore := zapcore.NewTee(consoleCore, jsonCore)
	_logger := zap.New(loggerCore, zap.AddCaller())

	if enableNetworkLog {
		api.NetworkLogger = _logger.Named("Net").WithOptions(zap.WithCaller(false))
	}
	api.Logger = _logger
	inspector.Logger = _logger
	conf.Logger = _logger
	model.Logger = _logger
	logger.InitLegacyLogger(_logger)

	_logger.Sugar().Infof("Log start: %s, %s", time.Now().Format(time.RFC3339), version.UserAgent())
	_logger.Sugar().Infof("Args: %s", os.Args)
	_logger.Sugar().Infof("Machine id: %s", version.MachineId())
	return nil
}

func createLogFile() (*os.File, error) {
	if disableLogFile {
		return nil, ErrLogFileDisabled
	}
	logFilePath := defaultLogFile
	if cliLogFilePathOverride != "" {
		logFilePath = cliLogFilePathOverride
	}
	// ensure log dir created
	if e := os.MkdirAll(filepath.Dir(logFilePath), 0755); e != nil {
		return nil, wrapLogErr(ErrCreateLogFileFailed, e)
	}
	if f, e := os.OpenFile(logFilePath, os.O_CREATE+os.O_RDWR+os.O_APPEND, 0644); e != nil {
		return nil, wrapLogErr(ErrCreateLogFileFailed, e)
	} else {
		return f, nil
	}
}

// LogFileCleanup auto remove log files which created between staticRefTime and 7 days ago
func logFileCleanup() {
	// file before staticRefTime will be ignored
	var staticRefTime = must.A(time.Parse(time.RFC3339, "2020-01-01T00:00:00Z"))

	logFilePattern := regexp.MustCompile("^(\\d+)\\.log$")
	basePath := filepath.Dir(defaultLogFile)
	if basePath == "" {
		return
	}
	d, e := os.ReadDir(basePath)
	if e != nil {
		return
	}
	for _, entry := range d {
		if entry.IsDir() || !entry.Type().IsRegular() {
			continue
		}
		if m := logFilePattern.FindStringSubmatch(entry.Name()); m != nil {
			ts, e := strconv.Atoi(m[1])
			if e != nil {
				continue
			}
			lt := time.UnixMilli(int64(ts))
			if lt.Before(staticRefTime) {
				continue
			}
			if time.Now().Sub(time.UnixMilli(int64(ts))) > time.Hour*24*7 {
				_ = os.Remove(filepath.Join(basePath, entry.Name()))
			}
		}
	}
}
