package utils

import "github.com/sirupsen/logrus"

var LoggerLevel = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"error": logrus.ErrorLevel,
	"trace": logrus.TraceLevel,
}
