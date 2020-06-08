package log

// Formatter can be used to log to another format than JSON
type Formatter func(ev *Event) ([]byte, error)
