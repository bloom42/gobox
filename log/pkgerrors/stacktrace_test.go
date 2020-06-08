package pkgerrors

// import (
// 	"bytes"
// 	"regexp"
// 	"testing"

// 	"github.com/pkg/errors"
// 	"gitlab.com/bloom42/gobox/log"
// )

// func TestLogStack(t *testing.T) {
// 	log.ErrorStackMarshaler = MarshalStack

// 	out := &bytes.Buffer{}
// 	logger := log.New(log.SetWriter(out), log.SetFields(log.Timestamp(false)))

// 	err := errors.Wrap(errors.New("error message"), "from error")
// 	logger.Log("", log.Stack(true), log.Err("error", err))

// 	got := out.String()
// 	want := `\{"stack":\[\{"func":"TestLogStack","line":"18","source":"stacktrace_test.go"\},.*\],"error":"from error: error message"\}\n`
// 	if ok, _ := regexp.MatchString(want, got); !ok {
// 		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
// 	}
// }

// func TestContextStack(t *testing.T) {
// 	log.ErrorStackMarshaler = MarshalStack

// 	out := &bytes.Buffer{}
// 	logger := log.New(
// 		log.SetWriter(out),
// 		log.SetFields(log.Stack(true), log.Timestamp(false)),
// 	)

// 	err := errors.Wrap(errors.New("error message"), "from error")
// 	logger.Log("", log.Err("error", err))

// 	got := out.String()
// 	want := `\{"stack":\[\{"func":"TestContextStack","line":"37","source":"stacktrace_test.go"\},.*\],"error":"from error: error message"\}\n`
// 	if ok, _ := regexp.MatchString(want, got); !ok {
// 		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
// 	}
// }

// func BenchmarkLogStack(b *testing.B) {
// 	log.ErrorStackMarshaler = MarshalStack
// 	out := &bytes.Buffer{}
// 	logger := log.New(log.SetWriter(out))
// 	err := errors.Wrap(errors.New("error message"), "from error")
// 	b.ReportAllocs()

// 	for i := 0; i < b.N; i++ {
// 		logger.Log("", log.Stack(true), log.Err("error", err))
// 		out.Reset()
// 	}
// }
