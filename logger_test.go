package medialocker

import (
	"github.com/Sirupsen/logrus"
	"io"
)

func NewTestLogger(out io.Writer) *Logger {
	log := &logrus.Logger{
		Out:       out,
		Formatter: &logrus.TextFormatter{ForceColors: true},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	return &Logger{Logger: log}
}
