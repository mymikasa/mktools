package mlog

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warning(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
	Level(level string)
	OutputPath(path string) (err error)
}

var mLog Logger

type defaultLogger struct {
	logger *logrus.Logger
}

func (l *defaultLogger) Config(conf ConfigOptions) (err error) {
	l.logger.Out = conf.Logger()
	return
}

func (l *defaultLogger) OutputPath(path string) (err error) {
	config := defaultConfig()
	config.OutputPath = path

	l.logger.Out = config.Logger()
	return
}

func (l *defaultLogger) Debug(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.logger.WithFields(fields).Debug(msg)
}

func (l *defaultLogger) Info(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.logger.WithFields(fields).Info(msg)
}

func (l *defaultLogger) Warning(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.logger.WithFields(fields).Warning(msg)
}

func (l *defaultLogger) Error(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.logger.WithFields(fields).Error(msg)
}

func (l *defaultLogger) Fatal(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.logger.WithFields(fields).Fatal(msg)
}

func (l *defaultLogger) Level(level string) {
	switch strings.ToLower(level) {
	case "debug":
		l.logger.SetLevel(logrus.DebugLevel)
	case "warning":
		l.logger.SetLevel(logrus.WarnLevel)
	case "error":
		l.logger.SetLevel(logrus.ErrorLevel)
	case "fatle":
		l.logger.SetLevel(logrus.FatalLevel)
	default:
		l.logger.SetLevel(logrus.InfoLevel)
	}
}

func (c *ConfigOptions) Logger() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filepath.ToSlash(c.OutputPath),
		MaxSize:    c.MaxFileSizeMB, //MB
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAges,
		Compress:   c.Compress,
		LocalTime:  c.LocalTime,
	}
}

const defaultLogPath = "/tmp/mktools.log"

func defaultConfig() ConfigOptions {
	o := ConfigOptions{
		OutputPath:    defaultLogPath,
		MaxFileSizeMB: 10,
		MaxBackups:    5,
		MaxAges:       3,
		Compress:      false,
		LocalTime:     true,
	}

	return o
}

func SetLogger(logger Logger) {
	mLog = logger
}

func SetLogLevel(level string) {
	if level == "" {
		return
	}

	mLog.Level(level)
}

func SetOutPutPath(path string) (err error) {
	if path == "" {
		return
	}
	return mLog.OutputPath(path)
}

func Debug(msg string, fields map[string]interface{}) {
	mLog.Debug(msg, fields)
}

func Info(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	mLog.Info(msg, fields)
}

func Warning(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	mLog.Warning(msg, fields)
}

func Error(msg string, fields map[string]interface{}) {
	mLog.Error(msg, fields)
}

func Fatal(msg string, fields map[string]interface{}) {
	mLog.Fatal(msg, fields)
}

func MustSetUp(options ...Option) {
	o := &ConfigOptions{
		// OutputPath:    defaultLogPath,
		MaxFileSizeMB: 10,
		MaxBackups:    5,
		MaxAges:       3,
		Compress:      false,
		LocalTime:     true,
	}

	for _, opt := range options {
		opt(o)
	}

	r := &defaultLogger{
		logger: logrus.New(),
	}
	level := os.Getenv("ROCKETMQ_GO_LOG_LEVEL")
	switch strings.ToLower(level) {
	case "debug":
		r.logger.SetLevel(logrus.DebugLevel)
	case "warn":
		r.logger.SetLevel(logrus.WarnLevel)
	case "error":
		r.logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		r.logger.SetLevel(logrus.FatalLevel)
	default:
		r.logger.SetLevel(logrus.InfoLevel)
	}
	mLog = r

	SetOutPutPath(o.OutputPath)
}
