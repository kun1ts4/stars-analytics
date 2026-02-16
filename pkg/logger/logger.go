// Package logger предоставляет структурированное логирование с использованием logrus.
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log — глобальный экземпляр логгера.
var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	Log.SetLevel(logrus.InfoLevel)
}

// SetLevel устанавливает уровень логирования.
func SetLevel(level logrus.Level) {
	Log.SetLevel(level)
}

// SetTextFormatter устанавливает текстовый форматтер для логов.
func SetTextFormatter() {
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

// WithField создает запись с одним полем.
func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

// WithFields создает запись с несколькими полями.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}

// WithError создает запись с ошибкой.
func WithError(err error) *logrus.Entry {
	return Log.WithError(err)
}

// Info логирует информационное сообщение.
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Infof логирует форматированное информационное сообщение.
func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

// Warn логирует предупреждение.
func Warn(args ...interface{}) {
	Log.Warn(args...)
}

// Warnf логирует форматированное предупреждение.
func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

// Error логирует ошибку.
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Errorf логирует форматированную ошибку.
func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

// Fatal логирует критическую ошибку и завершает приложение.
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

// Fatalf логирует форматированную критическую ошибку и завершает приложение.
func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

// Debug логирует отладочное сообщение.
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

// Debugf логирует форматированное отладочное сообщение.
func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}
