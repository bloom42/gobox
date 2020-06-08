package log_test

import "gitlab.com/bloom42/gobox/log"

func ExampleNew() {
	log := log.New(log.SetFields(log.Timestamp(false)))

	log.Info("hello world")
	// Output: {"level":"info","message":"hello world"}
}

func ExampleFields() {
	log := log.New(log.SetFields(log.Timestamp(false), log.String("foo", "bar")))

	log.Info("hello world")

	// Output: {"level":"info","foo":"bar","message":"hello world"}
}

// func ExampleLogger_Level() {
// 	log := log.New(os.Stdout).Level(log.WarnLevel)

// 	log.Info().Msg("filtered out message")
// 	log.Error().Msg("kept message")

// 	// Output: {"level":"error","message":"kept message"}
// }

// func ExampleLogger_Sample() {
// 	log := log.New(os.Stdout).Sample(&log.BasicSampler{N: 2})

// 	log.Info().Msg("message 1")
// 	log.Info().Msg("message 2")
// 	log.Info().Msg("message 3")
// 	log.Info().Msg("message 4")

// 	// Output: {"level":"info","message":"message 1"}
// 	// {"level":"info","message":"message 3"}
// }

// type LevelNameHook struct{}

// func (h LevelNameHook) Run(e *log.Event, l log.Level, msg string) {
// 	if l != log.NoLevel {
// 		e.Str("level_name", l.String())
// 	} else {
// 		e.Str("level_name", "NoLevel")
// 	}
// }

// type MessageHook string

// func (h MessageHook) Run(e *log.Event, l log.Level, msg string) {
// 	e.Str("the_message", msg)
// }

// func ExampleLogger_Hook() {
// 	var levelNameHook LevelNameHook
// 	var messageHook MessageHook = "The message"

// 	log := log.New(os.Stdout).Hook(levelNameHook).Hook(messageHook)

// 	log.Info().Msg("hello world")

// 	// Output: {"level":"info","level_name":"info","the_message":"hello world","message":"hello world"}
// }

// func ExampleLogger_Print() {
// 	log := log.New(os.Stdout)

// 	log.Print("hello world")

// 	// Output: {"level":"debug","message":"hello world"}
// }

// func ExampleLogger_Printf() {
// 	log := log.New(os.Stdout)

// 	log.Printf("hello %s", "world")

// 	// Output: {"level":"debug","message":"hello world"}
// }

// func ExampleLogger_Debug() {
// 	log := log.New(os.Stdout)

// 	log.Debug().
// 		Str("foo", "bar").
// 		Int("n", 123).
// 		Msg("hello world")

// 	// Output: {"level":"debug","foo":"bar","n":123,"message":"hello world"}
// }

// func ExampleLogger_Info() {
// 	log := log.New(os.Stdout)

// 	log.Info().
// 		Str("foo", "bar").
// 		Int("n", 123).
// 		Msg("hello world")

// 	// Output: {"level":"info","foo":"bar","n":123,"message":"hello world"}
// }

// func ExampleLogger_Warn() {
// 	log := log.New(os.Stdout)

// 	log.Warn().
// 		Str("foo", "bar").
// 		Msg("a warning message")

// 	// Output: {"level":"warning","foo":"bar","message":"a warning message"}
// }

// func ExampleLogger_Error() {
// 	log := log.New(os.Stdout)

// 	log.Error().
// 		Err(errors.New("some error")).
// 		Msg("error doing something")

// 	// Output: {"level":"error","error":"some error","message":"error doing something"}
// }

// func ExampleLogger_WithLevel() {
// 	log := log.New(os.Stdout)

// 	log.WithLevel(log.InfoLevel).
// 		Msg("hello world")

// 	// Output: {"level":"info","message":"hello world"}
// }

// func ExampleLogger_Write() {
// 	log := log.New(os.Stdout).With().
// 		Str("foo", "bar").
// 		Logger()

// 	stdlog.SetFlags(0)
// 	stdlog.SetOutput(log)

// 	stdlog.Print("hello world")

// 	// Output: {"foo":"bar","message":"hello world"}
// }

// func ExampleLogger_Log() {
// 	log := log.New(os.Stdout)

// 	log.Log().
// 		Str("foo", "bar").
// 		Str("bar", "baz").
// 		Msg("")

// 	// Output: {"foo":"bar","bar":"baz"}
// }

// func ExampleEvent_Dict() {
// 	log := log.New(os.Stdout)

// 	log.Log().
// 		Str("foo", "bar").
// 		Dict("dict", log.Dict().
// 			Str("bar", "baz").
// 			Int("n", 1),
// 		).
// 		Msg("hello world")

// 	// Output: {"foo":"bar","dict":{"bar":"baz","n":1},"message":"hello world"}
// }

// type User struct {
// 	Name    string
// 	Age     int
// 	Created time.Time
// }

// func (u User) MarshalLogObject(e *log.Event) {
// 	e.Str("name", u.Name).
// 		Int("age", u.Age).
// 		Time("created", u.Created)
// }

// type Price struct {
// 	val  uint64
// 	prec int
// 	unit string
// }

// func (p Price) MarshalLogObject(e *log.Event) {
// 	denom := uint64(1)
// 	for i := 0; i < p.prec; i++ {
// 		denom *= 10
// 	}
// 	result := []byte(p.unit)
// 	result = append(result, fmt.Sprintf("%d.%d", p.val/denom, p.val%denom)...)
// 	e.Str("price", string(result))
// }

// type Users []User

// func (uu Users) MarshalLogArray(a *log.Array) {
// 	for _, u := range uu {
// 		a.Object(u)
// 	}
// }

// func ExampleEvent_Array() {
// 	log := log.New(os.Stdout)

// 	log.Log().
// 		Str("foo", "bar").
// 		Array("array", log.Arr().
// 			Str("baz").
// 			Int(1),
// 		).
// 		Msg("hello world")

// 	// Output: {"foo":"bar","array":["baz",1],"message":"hello world"}
// }

// func ExampleEvent_Array_object() {
// 	log := log.New(os.Stdout)

// 	// Users implements log.LogArrayMarshaler
// 	u := Users{
// 		User{"John", 35, time.Time{}},
// 		User{"Bob", 55, time.Time{}},
// 	}

// 	log.Log().
// 		Str("foo", "bar").
// 		Array("users", u).
// 		Msg("hello world")

// 	// Output: {"foo":"bar","users":[{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},{"name":"Bob","age":55,"created":"0001-01-01T00:00:00Z"}],"message":"hello world"}
// }

// func ExampleEvent_Object() {
// 	log := log.New(os.Stdout)

// 	// User implements log.LogObjectMarshaler
// 	u := User{"John", 35, time.Time{}}

// 	log.Log().
// 		Str("foo", "bar").
// 		Object("user", u).
// 		Msg("hello world")

// 	// Output: {"foo":"bar","user":{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},"message":"hello world"}
// }

// func ExampleEvent_EmbedObject() {
// 	log := log.New(os.Stdout)

// 	price := Price{val: 6449, prec: 2, unit: "$"}

// 	log.Log().
// 		Str("foo", "bar").
// 		EmbedObject(price).
// 		Msg("hello world")

// 	// Output: {"foo":"bar","price":"$64.49","message":"hello world"}
// }

// func ExampleEvent_Interface() {
// 	log := log.New(os.Stdout)

// 	obj := struct {
// 		Name string `json:"name"`
// 	}{
// 		Name: "john",
// 	}

// 	log.Log().
// 		Str("foo", "bar").
// 		Interface("obj", obj).
// 		Msg("hello world")

// 	// Output: {"foo":"bar","obj":{"name":"john"},"message":"hello world"}
// }

// func ExampleEvent_Dur() {
// 	d := time.Duration(10 * time.Second)

// 	log := log.New(os.Stdout)

// 	log.Log().
// 		Str("foo", "bar").
// 		Dur("dur", d).
// 		Msg("hello world")

// 	// Output: {"foo":"bar","dur":10000,"message":"hello world"}
// }

// func ExampleEvent_Durs() {
// 	d := []time.Duration{
// 		time.Duration(10 * time.Second),
// 		time.Duration(20 * time.Second),
// 	}

// 	log := log.New(os.Stdout)

// 	log.Log().
// 		Str("foo", "bar").
// 		Durs("durs", d).
// 		Msg("hello world")

// 	// Output: {"foo":"bar","durs":[10000,20000],"message":"hello world"}
// }

// func ExampleContext_Dict() {
// 	log := log.New(os.Stdout).With().
// 		Str("foo", "bar").
// 		Dict("dict", log.Dict().
// 			Str("bar", "baz").
// 			Int("n", 1),
// 		).Logger()

// 	log.Log().Msg("hello world")

// 	// Output: {"foo":"bar","dict":{"bar":"baz","n":1},"message":"hello world"}
// }

// func ExampleContext_Array() {
// 	log := log.New(os.Stdout).With().
// 		Str("foo", "bar").
// 		Array("array", log.Arr().
// 			Str("baz").
// 			Int(1),
// 		).Logger()

// 	log.Log().Msg("hello world")

// 	// Output: {"foo":"bar","array":["baz",1],"message":"hello world"}
// }

// func ExampleContext_Array_object() {
// 	// Users implements log.LogArrayMarshaler
// 	u := Users{
// 		User{"John", 35, time.Time{}},
// 		User{"Bob", 55, time.Time{}},
// 	}

// 	log := log.New(os.Stdout).With().
// 		Str("foo", "bar").
// 		Array("users", u).
// 		Logger()

// 	log.Log().Msg("hello world")

// 	// Output: {"foo":"bar","users":[{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},{"name":"Bob","age":55,"created":"0001-01-01T00:00:00Z"}],"message":"hello world"}
// }

// func ExampleContext_Object() {
// 	// User implements log.LogObjectMarshaler
// 	u := User{"John", 35, time.Time{}}

// 	log := log.New(os.Stdout).With().
// 		Str("foo", "bar").
// 		Object("user", u).
// 		Logger()

// 	log.Log().Msg("hello world")

// 	// Output: {"foo":"bar","user":{"name":"John","age":35,"created":"0001-01-01T00:00:00Z"},"message":"hello world"}
// }

// func ExampleContext_EmbedObject() {

// 	price := Price{val: 6449, prec: 2, unit: "$"}

// 	log := log.New(os.Stdout).With().
// 		Str("foo", "bar").
// 		EmbedObject(price).
// 		Logger()

// 	log.Log().Msg("hello world")

// 	// Output: {"foo":"bar","price":"$64.49","message":"hello world"}
// }

// func ExampleContext_Interface() {
// 	obj := struct {
// 		Name string `json:"name"`
// 	}{
// 		Name: "john",
// 	}

// 	log := log.New(os.Stdout).With().
// 		Str("foo", "bar").
// 		Interface("obj", obj).
// 		Logger()

// 	log.Log().Msg("hello world")

// 	// Output: {"foo":"bar","obj":{"name":"john"},"message":"hello world"}
// }

// func ExampleContext_Dur() {
// 	d := time.Duration(10 * time.Second)

// 	log := log.New(os.Stdout).With().
// 		Str("foo", "bar").
// 		Dur("dur", d).
// 		Logger()

// 	log.Log().Msg("hello world")

// 	// Output: {"foo":"bar","dur":10000,"message":"hello world"}
// }

// func ExampleContext_Durs() {
// 	d := []time.Duration{
// 		time.Duration(10 * time.Second),
// 		time.Duration(20 * time.Second),
// 	}

// 	log := log.New(os.Stdout).With().
// 		Str("foo", "bar").
// 		Durs("durs", d).
// 		Logger()

// 	log.Log().Msg("hello world")

// 	// Output: {"foo":"bar","durs":[10000,20000],"message":"hello world"}
// }

// func ExampleContext_IPAddr() {
// 	hostIP := net.IP{192, 168, 0, 100}
// 	log := log.New(os.Stdout).With().
// 		IPAddr("HostIP", hostIP).
// 		Logger()

// 	log.Log().Msg("hello world")

// 	// Output: {"HostIP":"192.168.0.100","message":"hello world"}
// }

// func ExampleContext_IPPrefix() {
// 	route := net.IPNet{IP: net.IP{192, 168, 0, 0}, Mask: net.CIDRMask(24, 32)}
// 	log := log.New(os.Stdout).With().
// 		IPPrefix("Route", route).
// 		Logger()

// 	log.Log().Msg("hello world")

// 	// Output: {"Route":"192.168.0.0/24","message":"hello world"}
// }

// func ExampleContext_MacAddr() {
// 	mac := net.HardwareAddr{0x00, 0x14, 0x22, 0x01, 0x23, 0x45}
// 	log := log.New(os.Stdout).With().
// 		MACAddr("hostMAC", mac).
// 		Logger()

// 	log.Log().Msg("hello world")

// 	// Output: {"hostMAC":"00:14:22:01:23:45","message":"hello world"}
// }
