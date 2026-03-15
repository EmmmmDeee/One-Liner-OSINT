package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus logger
type Logger struct {
	*logrus.Logger
}

var defaultLogger *Logger

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	if defaultLogger != nil {
		return defaultLogger
	}

	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
	l.SetLevel(logrus.InfoLevel)

	defaultLogger = &Logger{l}
	return defaultLogger
}

// SetLevel sets the logger level
func (l *Logger) SetLevel(level string) {
	switch level {
	case "debug":
		l.Logger.SetLevel(logrus.DebugLevel)
	case "info":
		l.Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		l.Logger.SetLevel(logrus.WarnLevel)
	case "error":
		l.Logger.SetLevel(logrus.ErrorLevel)
	default:
		l.Logger.SetLevel(logrus.InfoLevel)
	}
}

// SetOutput sets the logger output
func (l *Logger) SetOutput(w io.Writer) {
	l.Logger.SetOutput(w)
}

// SetVerbose enables verbose logging
func (l *Logger) SetVerbose(verbose bool) {
	if verbose {
		l.Logger.SetLevel(logrus.DebugLevel)
	}
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}
