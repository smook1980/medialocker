package medialocker

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// Log represents an logging sink
type Log interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

// Logger represents an abstracted structured logging implementation. It
// provides methods to trigger log messages at various alert levels and a
// WithField method to set keys for a structured log message.
type Logger struct {
	*logrus.Logger
}

type Level struct {
	logrus.Level
}

func (l *Level) UnmarshalText(text []byte) (err error) {
	l.Level, err = logrus.ParseLevel(string(text))
	return
}

type loggerSettings struct {
	ConsoleLogging, ForceColor, DisableColor bool
	LogLevel                                 logrus.Level
	LogPath                                  string
}

// Must takes a message and returns a Do func with takes the func who's error must be nil.
// If error is not nil, panic, showing the given message and error info before crashing!
// Must("Oh noos... Stupid %s", failWhale).Do(ImportFunc())
func (l *Logger) Must(msg string, args ...interface{}) struct{ Do func(args ...interface{}) } {
	return struct{ Do func(args ...interface{}) }{
		Do: func(args ...interface{}) {
			switch err := args[len(args)-1].(type) {
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
	logger.Formatter = &prefixed.TextFormatter{ShortTimestamp: true, ForceColors: true}
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

	logFile, err := fs.OpenFile(settings.LogPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)

	if err != nil {
		l.Errorf("Failed to open log file %s: %s", settings.LogPath, err)
		return
	}

	// go func(fn io.Closer, c AppContext) {
	//	<-c.Done()
	//	fn.Close()
	// }(logFile, ctx)

	// runtime.SetFinalizer(*l, logFile.Close)

	l.Out = logFile
	// l.Out = io.MultiWriter(os.Stderr, logFile)
}

func NewLoggerWith(config Config, fs *FileSystem) *Logger {
	logger := logrus.New()
	logFormatter := new(prefixed.TextFormatter)
	logFormatter.ForceColors = true
	logger.Formatter = logFormatter
	l := &Logger{Logger: logger}
	if config.DebugLogging {
		logger.Level = logrus.DebugLevel
	}

	if config.ConsoleLog {
		logger.Out = os.Stderr

		// logger.Formatter = &prefixed.TextFormatter{
		//	ForceColors: config.ForceColor,
		//	// DisableColors:  config.DisableColors,
		//	ShortTimestamp: true,
		// }

		return l
	}

	fs.EnsureFileDirectory(config.LogPath)

	logFile, err := fs.OpenFile(config.LogPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)

	if err != nil {
		logger.Errorf("Failed to open log file %s: %s", config.LogPath, err)
		return l
	}

	logger.Out = logFile
	runtime.SetFinalizer(l, func(interface{}) { logFile.Close() })

	return l
}
