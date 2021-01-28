package zaplog

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger encapsulates log level handler with zap fields
type Logger interface {
	Debug(message string, fields ...zapcore.Field)
	Info(message string, fields ...zapcore.Field)
	Warn(message string, fields ...zapcore.Field)
	Error(message string, fields ...zapcore.Field)
	Fatal(message string, fields ...zapcore.Field)
	CheckErr(message string, err error, fields ...zapcore.Field)
	SafeClose(closer io.Closer, message string)
}

// zaplogger delegates all calls to the underlying zaplog.Logger
type zaplogger struct {
	logger *zap.Logger
}

// Debug logs an debug message with fields
func (l *zaplogger) Debug(message string, fields ...zapcore.Field) {
	l.logger.Debug(message, fields...)
}

// Info logs an info message with fields
func (l *zaplogger) Info(message string, fields ...zapcore.Field) {
	l.logger.Info(message, fields...)
}

// Error logs an error message with fields
func (l *zaplogger) Error(message string, fields ...zapcore.Field) {
	l.logger.Error(message, fields...)
}

// Warn logs a warning with fields
func (l *zaplogger) Warn(message string, fields ...zapcore.Field) {
	l.logger.Warn(message, fields...)
}

// Fatal logs a fatal error message with fields
func (l *zaplogger) Fatal(message string, fields ...zapcore.Field) {
	l.logger.Fatal(message, fields...)
}

// SafeClose closes a Closer and log a message of error in case it happened
func (l *zaplogger) SafeClose(closer io.Closer, message string) {
	err := closer.Close()
	l.CheckErr(message, err)
}

// CheckErr handles error correctly
func (l *zaplogger) CheckErr(message string, err error, fields ...zapcore.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.logger.Error(message, fields...)
	}
}
