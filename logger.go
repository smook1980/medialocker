package medialocker

import (
	"github.com/Sirupsen/logrus"
	"os"
	"path"
)

// Logger represents an abstracted structured logging implementation. It
// provides methods to trigger log messages at various alert levels and a
// WithField method to set keys for a structured log message.
type Logger struct {
	*logrus.Logger
}

func NewLoggerWith(config Config) *Logger {
	logger := logrus.New()
	if config.DebugLogging {
		logger.Level = logrus.DebugLevel
	}

	if LocalFileExists(config.LogPath) {
		logFile, err := LocalFileSystem().OpenFile(config.LogPath, os.O_WRONLY, os.ModePerm)
		if err != nil {
			logger.Errorf("Couldn't open log file %s: %v!", config.LogPath, err)
		} else {
			logger.Out = logFile
		}
	} else {
		logDir := path.Dir(config.LogPath)
		err := LocalFileSystem().MkdirAll(logDir, os.ModeDir|os.ModePerm)
		if err != nil {
			logger.Errorf("Couldn't create log file directory %s: %v!", logDir, err)
		}
		logFile, err := LocalFileSystem().Create(config.LogPath)

		if err != nil {
			logger.Errorf("Couldn't open log file %s: %v!", config.LogPath, err)
		} else {
			logger.Out = logFile
		}
	}

	return &Logger{Logger: logger}
}
