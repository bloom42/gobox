package log

import (
	"io"
	"os"
	"time"
)

// LoggerOption is used to configure a logger.
type LoggerOption func(logger *Logger)

var (
	// DurationFieldUnit defines the unit for time.Duration type fields added
	// using the Duration method.
	DurationFieldUnit = time.Millisecond

	// DurationFieldInteger renders Dur fields as integer instead of float if
	// set to true.
	DurationFieldInteger = false

	// ErrorHandler is called whenever log fails to write an event on its
	// output. If not set, an error is printed on the stderr. This handler must
	// be thread safe and non-blocking.
	ErrorHandler func(err error)

	// ErrorStackMarshaler extract the stack from err if any.
	ErrorStackMarshaler func(err error) interface{}

	// ErrorMarshalFunc allows customization of global error marshaling
	ErrorMarshalFunc = func(err error) interface{} {
		return err
	}
)

// SetWriter update logger's writer.
func SetWriter(writer io.Writer) LoggerOption {
	return func(logger *Logger) {
		if writer == nil {
			writer = os.Stdout
		}
		lw, ok := writer.(LevelWriter)
		if !ok {
			lw = levelWriterAdapter{writer}
		}
		logger.writer = lw
	}
}

// SetLevel update logger's level.
func SetLevel(lvl Level) LoggerOption {
	return func(logger *Logger) {
		logger.level = lvl
	}
}

// SetSampler update logger's sampler.
func SetSampler(sampler Sampler) LoggerOption {
	return func(logger *Logger) {
		logger.sampler = sampler
	}
}

// AddHook appends hook to logger's hook
func AddHook(hook Hook) LoggerOption {
	return func(logger *Logger) {
		logger.hooks = append(logger.hooks, hook)
	}
}

// SetHooks replaces logger's hooks
func SetHooks(hooks ...Hook) LoggerOption {
	return func(logger *Logger) {
		logger.hooks = hooks
	}
}

// SetFields update logger's context fields
func SetFields(fields ...Field) LoggerOption {
	return func(logger *Logger) {
		e := newEvent(logger.writer, logger.level)
		e.buf = nil
		copyInternalLoggerFieldsToEvent(logger, e)
		for i := range fields {
			fields[i](e)
		}
		if e.stack != logger.stack {
			logger.stack = e.stack
		}
		if e.caller != logger.caller {
			logger.caller = e.caller
		}
		if e.timestamp != logger.timestamp {
			logger.timestamp = e.timestamp
		}
		if e.buf != nil {
			logger.context = enc.AppendObjectData(logger.context, e.buf)
		}
	}
}

// SetFormatter update logger's formatter.
func SetFormatter(formatter Formatter) LoggerOption {
	return func(logger *Logger) {
		logger.formatter = formatter
	}
}

// SetTimestampFieldName update logger's timestampFieldName.
func SetTimestampFieldName(timestampFieldName string) LoggerOption {
	return func(logger *Logger) {
		logger.timestampFieldName = timestampFieldName
	}
}

// SetLevelFieldName update logger's levelFieldName.
func SetLevelFieldName(levelFieldName string) LoggerOption {
	return func(logger *Logger) {
		logger.levelFieldName = levelFieldName
	}
}

// SetMessageFieldName update logger's messageFieldName.
func SetMessageFieldName(messageFieldName string) LoggerOption {
	return func(logger *Logger) {
		logger.messageFieldName = messageFieldName
	}
}

// SetCallerFieldName update logger's callerFieldName.
func SetCallerFieldName(callerFieldName string) LoggerOption {
	return func(logger *Logger) {
		logger.callerFieldName = callerFieldName
	}
}

// SetCallerSkipFrameCount update logger's callerSkipFrameCount.
func SetCallerSkipFrameCount(callerSkipFrameCount int) LoggerOption {
	return func(logger *Logger) {
		logger.callerSkipFrameCount = callerSkipFrameCount
	}
}

// SetErrorStackFieldName update logger's errorStackFieldName.
func SetErrorStackFieldName(errorStackFieldName string) LoggerOption {
	return func(logger *Logger) {
		logger.errorStackFieldName = errorStackFieldName
	}
}

// SetTimeFieldFormat update logger's timeFieldFormat.
func SetTimeFieldFormat(timeFieldFormat string) LoggerOption {
	return func(logger *Logger) {
		logger.timeFieldFormat = timeFieldFormat
	}
}

// SetTimestampFunc update logger's timestampFunc.
func SetTimestampFunc(timestampFunc func() time.Time) LoggerOption {
	return func(logger *Logger) {
		logger.timestampFunc = timestampFunc
	}
}
