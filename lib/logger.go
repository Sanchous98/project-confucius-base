package lib

import (
	"github.com/Sanchous98/project-confucius-base/utils"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	loggerConfigPath       = "config/logger.yaml"
	debugLevel       level = 1 << iota
	infoLevel
	noticeLevel
	warningLevel
	errorLevel
	criticalLevel
	alertLevel
	emergencyLevel
)

type (
	level uint16

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
	}

	Log struct {
		config  *loggerConfig
		channel io.Writer
		recover func(...interface{})
	}

	loggerConfig struct {
		Channels map[string]io.Writer
	}
)

func (c *loggerConfig) Unmarshall() error {
	absPath, _ := filepath.Abs(loggerConfigPath)
	content, err := ioutil.ReadFile(absPath)
	cfg, err := utils.HydrateConfig(c, content, yaml.Unmarshal)

	if err != nil {
		return err
	}

	c = cfg.(*loggerConfig)

	return nil
}

func (l *Log) Construct() {
	l.config = new(loggerConfig)
	_ = l.config.Unmarshall()
}

func (l Log) Channel(channel string) *Log {
	log := new(Log)
	log.Construct()
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
	l.recover(context)
}

func (l *Log) Emergency(message error, context ...interface{}) {
	l.Log(emergencyLevel, message, context)
	os.Exit(1)
}
