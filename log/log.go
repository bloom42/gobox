package log

// global logger

// logger is the global globalLogger.
var globalLogger = New(CallerSkipFrameCount(4))

// SetGlobalLogger update the global logger
func SetGlobalLogger(logger Logger) {
	globalLogger = logger.With(CallerSkipFrameCount(4))
}

// GlobalLogger returns the global logger
func GlobalLogger() Logger {
	return globalLogger.With(CallerSkipFrameCount(3))
}

// With duplicates the global logger and update it's configuration.
func With(options ...LoggerOption) Logger {
	options = append([]LoggerOption{CallerSkipFrameCount(3)}, options...)
	return globalLogger.With(options...)
}

// LogWithLevel logs a new message with the given level.
func LogWithLevel(level LogLevel, message string, fields ...Field) {
	globalLogger.LogWithLevel(level, message, fields...)
}

// Debug starts a new message with debug level.
func Debug(message string, fields ...Field) {
	globalLogger.Debug(message, fields...)
}

// Info logs a new message with info level.
func Info(message string, fields ...Field) {
	globalLogger.Info(message, fields...)
}

// Warn logs a new message with warn level.
func Warn(message string, fields ...Field) {
	globalLogger.Warn(message, fields...)
}

// Error logs a message with error level.
func Error(message string, fields ...Field) {
	globalLogger.Error(message, fields...)
}

// Fatal logs a new message with fatal level. The os.Exit(1) function
// is then called, which terminates the program immediately.
func Fatal(message string, fields ...Field) {
	globalLogger.Fatal(message, fields...)
}

// Panic logs a new message with panic level. The panic() function
// is then called, which stops the ordinary flow of a goroutine.
func Panic(message string, fields ...Field) {
	globalLogger.Panic(message, fields...)
}

// Log logs a new message with no level. Setting GlobalLevel to Disabled
// will still disable events produced by this method.
func Log(message string, fields ...Field) {
	globalLogger.Log(message, fields...)
}

// Append the fields to the internal logger's context.
// It does not create a noew copy of the logger and rely on a mutex to enable thread safety,
// so `Config(With(fields...))` often is preferable.
func Append(fields ...Field) {
	globalLogger.Append(fields...)
}

// NewDict create a new Dict with the logger's configuration
func NewDict(fields ...Field) *Event {
	return globalLogger.NewDict(fields...)
}
