//Package logger contains simle logger interface
package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Infoln(args ...interface{})
	Errorln(args ...interface{})
	Warnln(args ...interface{})
	Fatalln(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Fatal(args ...interface{})
}

func New() (Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
