package logger

import (
	"go.uber.org/zap"
)

func With(fields ...zap.Field) *zap.Logger {
	return elg.With(fields...)
}

func Named(name string) *zap.Logger {
	return elg.Named(name)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	if elg != nil {
		elg.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
	}
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zap.Field) {
	if elg != nil {
		elg.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
	}
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zap.Field) {
	if elg != nil {
		elg.WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
	}
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zap.Field) {
	if elg != nil {
		elg.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
	}
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zap.Field) {
	if elg != nil {
		elg.WithOptions(zap.AddCallerSkip(1)).Panic(msg, fields...)
	}
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zap.Field) {
	if elg != nil {
		elg.WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
	}
}

func access(fields ...zap.Field) {
	if alg != nil {
		alg.Debug("", fields...)
	}
}
