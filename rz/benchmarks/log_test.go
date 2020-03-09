package bench

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"gitlab.com/bloom42/lily/rz"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func newDisabledLogrus() *logrus.Logger {
	logger := newLogrus()
	logger.Level = logrus.ErrorLevel
	return logger
}

func newLogrus() *logrus.Logger {
	return &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
}

func newRz() rz.Logger {
	return rz.New(
		rz.Writer(ioutil.Discard),
		rz.Level(rz.DebugLevel),
		rz.Fields(rz.Timestamp(true)),
		rz.TimeFieldFormat(""),
		rz.TimestampFunc(time.Now),
	)
}

func newDisabledRz() rz.Logger {
	return newRz().With(rz.Level(rz.Disabled))
}

func newZerolog() zerolog.Logger {
	zerolog.TimeFieldFormat = ""
	return zerolog.New(ioutil.Discard).With().Timestamp().Logger().Level(zerolog.DebugLevel)
}

func newDisabledZerolog() zerolog.Logger {
	zerolog.TimeFieldFormat = ""
	return newZerolog().Level(zerolog.Disabled)
}

func zapMillisecondDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}

func zapEncoder() zapcore.Encoder {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeDuration = zapMillisecondDurationEncoder
	ec.EncodeTime = zapcore.EpochTimeEncoder
	ec.TimeKey = "timestamp"
	ec.MessageKey = "message"
	return zapcore.NewJSONEncoder(ec)
}

func newZap() *zap.Logger {
	lvl := zap.DebugLevel
	var w zapcore.WriteSyncer = &zaptest.Discarder{}
	return zap.New(zapcore.NewCore(zapEncoder(), w, lvl))
}

func newDisabledZap() *zap.Logger {
	lvl := zap.FatalLevel
	var w zapcore.WriteSyncer = &zaptest.Discarder{}
	return zap.New(zapcore.NewCore(zapEncoder(), w, lvl))
}

var _tenInts = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var _tenStrings = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
var _tenTimes = []time.Time{time.Now()}

type _testUser struct {
	Username string
	Name     string
	Phone    string
}

var _oneUser = _testUser{
	Username: "lol",
	Name:     "lol2",
	Phone:    "lollol",
}

var _tenUsers = []_testUser{{}, {}, {}, {}, {},
	{}, {}, {}, {}, {}}
var errExample = errors.New("lolerror")

func logrus10Fields() logrus.Fields {
	return logrus.Fields{
		"int":     _tenInts[0],
		"ints":    _tenInts,
		"string":  _tenStrings[0],
		"strings": _tenStrings,
		"time":    _tenTimes[0],
		"times":   _tenTimes,
		"user1":   _oneUser,
		"user2":   _oneUser,
		"users":   _tenUsers,
		"error":   errExample,
	}
}

func zerolog10Fields(e *zerolog.Event) *zerolog.Event {
	return e.
		Int("int", _tenInts[0]).
		Ints("ints", _tenInts).
		Str("string", _tenStrings[0]).
		Strs("strings", _tenStrings).
		Time("time", _tenTimes[0]).
		Times("times", _tenTimes).
		Interface("user1", _oneUser).
		Interface("user2", _oneUser).
		Interface("users", _tenUsers).
		Err(errExample)
}

func zerolog10Context(c zerolog.Context) zerolog.Context {
	return c.
		Int("int", _tenInts[0]).
		Ints("ints", _tenInts).
		Str("string", _tenStrings[0]).
		Strs("strings", _tenStrings).
		Time("time", _tenTimes[0]).
		Times("times", _tenTimes).
		Interface("user1", _oneUser).
		Interface("user2", _oneUser).
		Interface("users", _tenUsers).
		Err(errExample)
}

func rz10Fields() []rz.Field {
	return []rz.Field{
		rz.Int("int", _tenInts[0]),
		rz.Ints("ints", _tenInts),
		rz.String("string", _tenStrings[0]),
		rz.Strings("strings", _tenStrings),
		rz.Time("time", _tenTimes[0]),
		rz.Times("times", _tenTimes),
		rz.Any("user1", _oneUser),
		rz.Any("user2", _oneUser),
		rz.Any("users", _tenUsers),
		rz.Err(errExample),
	}
}

func zap10Fields() []zap.Field {
	return []zap.Field{
		zap.Int("int", _tenInts[0]),
		zap.Ints("ints", _tenInts),
		zap.String("string", _tenStrings[0]),
		zap.Strings("strings", _tenStrings),
		zap.Time("time", _tenTimes[0]),
		zap.Times("times", _tenTimes),
		zap.Any("user1", _oneUser),
		zap.Any("user2", _oneUser),
		zap.Any("users", _tenUsers),
		zap.Error(errExample),
	}
}

var _testMessage = "hello world"

func BenchmarkDisabledWithoutFields(b *testing.B) {
	b.Logf("Logging without any structured context and at a disabled level.")
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newDisabledLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage)
			}
		})
	})
	b.Run("uber-go/zap", func(b *testing.B) {
		logger := newDisabledZap()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage)
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newDisabledZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(_testMessage)
			}
		})
	})
	b.Run("z0mbie42/rz-go", func(b *testing.B) {
		logger := newDisabledRz()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage)
			}
		})
	})
}

func BenchmarkWithoutFields(b *testing.B) {
	b.Logf("Logging without any structured context.")
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage)
			}
		})
	})
	b.Run("uber-go/zap", func(b *testing.B) {
		logger := newZap()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage)
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(_testMessage)
			}
		})
	})
	b.Run("z0mbie42/rz-go", func(b *testing.B) {
		logger := newRz()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage)
			}
		})
	})
}

func Benchmark10Context(b *testing.B) {
	b.Logf("Logging with 10 fields in context")
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		fields := logrus10Fields()
		l := logger.WithFields(fields)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Info(_testMessage)
			}
		})
	})
	b.Run("uber-go/zap", func(b *testing.B) {
		fields := zap10Fields()
		logger := newZap().With(fields...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage)
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := zerolog10Context(newZerolog().With()).Logger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(_testMessage)
			}
		})
	})
	b.Run("z0mbie42/rz-go", func(b *testing.B) {
		fields := rz10Fields()
		logger := newRz().With(rz.Fields(fields...))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage)
			}
		})
	})
}

func Benchmark10Fields(b *testing.B) {
	b.Logf("Logging without context and 10 fields")
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		fields := logrus10Fields()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fields).Info(_testMessage)
			}
		})
	})
	b.Run("uber-go/zap", func(b *testing.B) {
		logger := newZap()
		fields := zap10Fields()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage, fields...)
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				zerolog10Fields(logger.Info()).Msg(_testMessage)
			}
		})
	})
	b.Run("z0mbie42/rz-go", func(b *testing.B) {
		fields := rz10Fields()
		logger := newRz()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage, fields...)
			}
		})
	})
}

func Benchmark10Fields10Context(b *testing.B) {
	b.Logf("Logging without context and 10 fields")
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		fields := logrus10Fields()
		l := logger.WithFields(fields)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.WithFields(fields).Info(_testMessage)
			}
		})
	})
	b.Run("uber-go/zap", func(b *testing.B) {
		fields := zap10Fields()
		logger := newZap().With(fields...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage, fields...)
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := zerolog10Context(newZerolog().With()).Logger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				zerolog10Fields(logger.Info()).Msg(_testMessage)
			}
		})
	})
	b.Run("z0mbie42/rz-go", func(b *testing.B) {
		fields := rz10Fields()
		logger := newRz().With(rz.Fields(fields...))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_testMessage, fields...)
			}
		})
	})
}
