// Package log provides a lightweight logging library dedicated to JSON logging.
//
// A global Logger can be use for simple logging:
//
//     import "gitlab.com/bloom42/gobox/log"
//
//     log.Info("hello world")
//
//     // Output: {"timestamp":"2019-02-07T09:30:07Z","level":"info","message":"hello world"}
//
// Fields can be added to log messages:
//
//     log.Info("hello world", log.String("foo", "bar"))
//
//     // Output: {"timestamp":"2019-02-07T09:30:07Z","level":"info","message":"hello world","foo":"bar"}
//
// Create logger instance to manage different outputs:
//
//     logger := log.New()
//     log.Info("hello world",log.String("foo", "bar"))
//
//     // Output: {"timestamp":"2019-02-07T09:30:07Z","level":"info","message":"hello world","foo":"bar"}
//
// Sub-loggers let you chain loggers with additional context:
//
//     sublogger := log.Clone(log.String("component", "foo"))
//     sublogger.Info("hello world", nil)
//
//     // Output: {"timestamp":"2019-02-07T09:30:07Z","level":"info","message":"hello world","component":"foo"}
package log
