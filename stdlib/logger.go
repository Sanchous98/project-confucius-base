package stdlib

import (
	"fmt"
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/Sanchous98/project-confucius-base/utils"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"io"
	"syscall"
)

const loggerConfigPath = "config/logger.yaml"

const (
	debugLevel level = iota
	infoLevel
	noticeLevel
	warningLevel
	errorLevel
	criticalLevel
	alertLevel
	emergencyLevel
)

type level uint8

type (
	Logger interface {
		Log(level, error, ...interface{})
		Debug(error, ...interface{})
		Info(error, ...interface{})
		Notice(error, ...interface{})
		Warning(error, ...interface{})
		Error(error, ...interface{})
		Critical(error, ...interface{})
		Alert(error, ...interface{})
		Emergency(error, ...interface{})
		Channel(string) Logger
		src.Service
	}

	Log struct {
		config  *loggerConfig
		channel io.Writer
	}

	loggerConfig struct {
		Channels map[string]io.Writer
	}
)

func (c *loggerConfig) Unmarshall() error {
	return utils.Unmarshall(c, loggerConfigPath, yaml.Unmarshal)
}

func (l *Log) Constructor() {}

func (l *Log) Channel(channel string) Logger {
	log := new(Log)
	log.Constructor()
	log.channel = log.config.Channels[channel]

	return log
}

func (l *Log) Log(level level, message error, context ...interface{}) {
	switch level {
	case infoLevel:
		color.Green(message.Error(), context)
	case debugLevel:
		fallthrough
	case noticeLevel:
		color.White(message.Error(), context)
	case warningLevel:
		fallthrough
	case alertLevel:
		color.Yellow(message.Error(), context)
	case errorLevel:
		fallthrough
	case criticalLevel:
		fallthrough
	case emergencyLevel:
		color.Red(message.Error(), context)
	default:
		color.White(message.Error(), context)
	}
	fmt.Println()
}

func (l *Log) Debug(message error, context ...interface{}) {
	l.Log(debugLevel, message, context)
}

func (l *Log) Info(message error, context ...interface{}) {
	l.Log(infoLevel, message, context)
}

func (l *Log) Notice(message error, context ...interface{}) {
	l.Log(noticeLevel, message, context)
}

func (l *Log) Warning(message error, context ...interface{}) {
	l.Log(warningLevel, message, context)
}

func (l *Log) Error(message error, context ...interface{}) {
	l.Log(errorLevel, message, context)
}

func (l *Log) Critical(message error, context ...interface{}) {
	l.Log(criticalLevel, message, context)
}

func (l *Log) Alert(message error, context ...interface{}) {
	l.Log(alertLevel, message, context)
	panic(message.Error())
}

func (l *Log) Emergency(message error, context ...interface{}) {
	l.Log(emergencyLevel, message, context)
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	if err != nil {
		return
	}
}
