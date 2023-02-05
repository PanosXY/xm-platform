// Package logger is a structured logger that tries to ensure a unified approach for logs.
package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

const (
	durationFieldUnit = time.Millisecond
	messageFieldName  = "message"
	timeFieldName     = "timestamp"
	timeFieldFormat   = time.RFC3339Nano
)

// Logger is the default logger struct.
type Logger struct {
	logger   *zerolog.Logger
	reserved map[string]struct{}
	defaults map[string]interface{}
}

// Fields type to be used along with any additional metadata for each log.
type Fields map[string]interface{}

// NewLogger creates a new logger.
func NewLogger(isDebug, isEngineerFriendly bool) *Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimestampFieldName = timeFieldName
	zerolog.MessageFieldName = messageFieldName
	zerolog.TimeFieldFormat = timeFieldFormat
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if isEngineerFriendly {
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: timeFieldFormat,
		}).With().Timestamp().Logger()
	}

	return &Logger{
		logger:   &logger,
		reserved: reservedWords(),
		defaults: make(map[string]interface{}),
	}
}

// SetDefaultField sets a default field in the log output.
func (l *Logger) SetDefaultField(key string, value interface{}) error {
	if _, ok := l.reserved[key]; ok {
		return fmt.Errorf("field %s is reserved", key)
	}

	l.defaults[key] = value

	return nil
}

// reservedWords returns a map with reserved words to avoid any duplications / overrides.
func reservedWords() map[string]struct{} {
	reserved := make(map[string]struct{})
	reserved["message"] = struct{}{}
	reserved["component"] = struct{}{}
	reserved["context"] = struct{}{}
	reserved["level"] = struct{}{}
	reserved["timestamp"] = struct{}{}
	reserved["error"] = struct{}{}
	reserved["request_id"] = struct{}{}

	return reserved
}

// Trace logs a new message at level TRACE.
func (l *Logger) Trace(component, msg string) {
	l.logger.Trace().Fields(l.formatFields("", component, nil)).Msg(msg)
}

// Tracef logs a new message with additional fields at level TRACE.
func (l *Logger) Tracef(requestID, component, msg string, fields Fields) {
	l.logger.Trace().Fields(l.formatFields(requestID, component, fields)).Msg(msg)
}

// Debug logs a new message at level DEBUG.
func (l *Logger) Debug(component, msg string) {
	l.logger.Debug().Fields(l.formatFields("", component, nil)).Msg(msg)
}

// Debugf logs a new message with additional fields at level DEBUG.
func (l *Logger) Debugf(requestID, component, msg string, fields Fields) {
	l.logger.Debug().Fields(l.formatFields(requestID, component, fields)).Msg(msg)
}

// Info logs a new message at level INFO.
func (l *Logger) Info(component, msg string) {
	l.logger.Info().Fields(l.formatFields("", component, nil)).Msg(msg)
}

// Infof logs a new message with additional fields at level INFO combined.
func (l *Logger) Infof(requestID, component, msg string, fields Fields) {
	l.logger.Info().Fields(l.formatFields(requestID, component, fields)).Msg(msg)
}

// Error logs a new message at level ERROR.
func (l *Logger) Error(component, msg string, err error) {
	l.logger.Error().Err(err).Fields(l.formatFields("", component, nil)).Msg(msg)
}

// Errorf logs a new message with additional fields at level ERROR.
func (l *Logger) Errorf(requestID, component, msg string, err error, fields Fields) {
	l.logger.Error().Err(err).Fields(l.formatFields(requestID, component, fields)).Msg(msg)
}

// Warn logs a new message with additional fields at level WARN.
func (l *Logger) Warn(component, msg string, err error) {
	l.logger.Warn().Err(err).Fields(l.formatFields("", component, nil)).Msg(msg)
}

// Warnf logs a new message with additional fields at level WARN.
func (l *Logger) Warnf(requestID, component, msg string, err error, fields Fields) {
	l.logger.Warn().Err(err).Fields(l.formatFields(requestID, component, fields)).Msg(msg)
}

// Fatal logs a new message at level Fatal then the process will exit with status set to 1.
func (l *Logger) Fatal(component, msg string, err error) {
	l.logger.Fatal().Err(err).Fields(l.formatFields("", component, nil)).Msg(msg)
}

// Fatalf logs a new message with additional fields at level Fatal then the process will exit with status set to 1.
func (l *Logger) Fatalf(requestID, component, msg string, err error, fields Fields) {
	l.logger.Fatal().Err(err).Fields(l.formatFields(requestID, component, fields)).Msg(msg)
}

// FormatDuration formats a duration in the durationFieldUnit.
func (l *Logger) FormatDuration(duration time.Duration) float64 {
	return float64(duration) / float64(durationFieldUnit)
}

// FormatTimestamp format a timestamp in the timeFieldFormat.
func (l *Logger) FormatTimestamp(timestamp time.Time) string {
	return timestamp.Format(timeFieldFormat)
}

// formatFields converts Fields to an interface, so that it can be asserted correctly.
// Also, it checks if there is a reserved word and bypass it.
func (l *Logger) formatFields(requestID, component string, fields Fields) interface{} {
	fieldsMap := make(map[string]interface{}, 0)

	for key, value := range fields {
		// Check for reserved words
		if _, found := l.reserved[key]; found {
			continue
		}

		fieldsMap[key] = value
	}

	logCtx := make(map[string]interface{}, 0)

	for key, value := range l.defaults {
		logCtx[key] = value
	}
	if requestID != "" {
		logCtx["request_id"] = requestID
	}
	logCtx["component"] = component
	logCtx["context"] = fieldsMap

	return logCtx
}
