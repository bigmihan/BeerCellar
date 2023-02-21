package logger

import (
	"github.com/sirupsen/logrus"
	"io"
)

type Logger struct {
	Out io.Writer
}

func NewLogger(out io.Writer) *Logger {
	log := Logger{Out: out}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(log.Out)
	//logrus.SetLevel(logrus.InfoLevel)

	return &log
}

func (log *Logger) WriteLog(description string, URL string, err error) {

	logWithField := logrus.WithFields(logrus.Fields{
		"URL": URL,
	})

	if err != nil {
		logWithField.Error(err)
	} else {
		logWithField.Info(description)
	}
}
