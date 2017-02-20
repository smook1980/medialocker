package medialocker

import (
	"github.com/Sirupsen/logrus"
	"os"
	"path"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"fmt"
	"io"
	"runtime"
)

// Logger represents an abstracted structured logging implementation. It
// provides methods to trigger log messages at various alert levels and a
// WithField method to set keys for a structured log message.
type Logger struct {
	*logrus.Logger
}

type Level struct {
	logrus.Level
}

func (l *Level) UnmarshalText(text []byte) (err error){
	l.Level, err = logrus.ParseLevel(string(text))
	return
}

type loggerSettings struct {
	ConsoleLogging, ForceColor, DisableColor bool
	LogLevel logrus.Level
	LogPath string
}

// Must takes a message and returns a Do func with takes the func who's error must be nil.
// If error is not nil, panic, showing the given message and error info before crashing!
// Must("Oh noos... Stupid %s", failWhale).Do(ImportFunc())
func (l *Logger) Must(msg string, args ...interface{}) struct { Do func(args ...interface{}) } {
	return struct { Do func(args ...interface{}) }{
		Do: func(args ...interface{}) {
			switch err := args[len(args) - 1].(type) {
			case nil:
				break
			case error:
				message := fmt.Sprintf(msg, args...)
				l.Panicf("%s\nCaused By Error:\n%s", message, err)
			}
		},
	}
}

func NewDefaultLogger() *Logger {
	logger := &Logger{Logger: logrus.New()}
	logger.Level = logrus.WarnLevel
	logger.Out = os.Stderr
	logger.Formatter = &prefixed.TextFormatter{ShortTimestamp: true}
	return logger
}

func (l *Logger) Configure(ctx AppContext, settings loggerSettings) {
	l.Level = settings.LogLevel
	if settings.ConsoleLogging {
		l.Out = os.Stderr

		l.Formatter = &prefixed.TextFormatter{
			ForceColors:    settings.ForceColor,
			DisableColors:  settings.DisableColor,
			ShortTimestamp: true,
		}
	}

	fs := ctx.FileSystem()

	fs.EnsureFileDirectory(settings.LogPath)
	logFile, err := fs.OpenFile(settings.LogPath, os.O_APPEND, os.ModePerm)

	if err != nil {
		l.Errorf("Failed to open log file %s: %s", settings.LogPath, err)
		return
	}

	logFileCloser := func(c io.Closer) error {
		return c.Close()
	}(logFile)
	runtime.SetFinalizer(l, logFileCloser)

	l.Out = io.MultiWriter(os.Stderr, logFile)
}

func NewLoggerWith(config Config) *Logger {
	logger := logrus.New()
	if config.DebugLogging {
		logger.Level = logrus.DebugLevel
	}

	if config.ConsoleLog {
		logger.Out = os.Stderr

		logger.Formatter = &prefixed.TextFormatter{
			ForceColors:    config.ForceColor,
			// DisableColors:  config.DisableColors,
			ShortTimestamp: true,
		}

		return &Logger{Logger: logger}
	}

	if LocalFileExists(config.LogPath) {
		logFile, err := LocalFileSystem().OpenFile(config.LogPath, os.O_APPEND, os.ModePerm)
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
