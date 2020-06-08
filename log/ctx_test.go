package log

import (
	"context"
	"io/ioutil"
	"testing"
)

func TestFromCtx(t *testing.T) {
	log := NewLogger(SetWriter(ioutil.Discard))
	ctx := log.ToCtx(context.Background())
	log2 := FromCtx(ctx)
	if &log != log2 {
		t.Error("FromCtx did not return the expected logger")
	}

	// update
	log = log.Clone(SetLevel(InfoLevel))
	ctx = log.ToCtx(ctx)
	log2 = FromCtx(ctx)
	if &log != log2 {
		t.Error("FromCtx did not return the expected logger")
	}

	log2 = FromCtx(context.Background())
	if log2 == nil || log2 == &log {
		t.Error("FromCtx did not return the expected logger")
	}
}

func TestFromCtxDisabled(t *testing.T) {
	dl := NewLogger(SetWriter(ioutil.Discard), SetLevel(Disabled))
	ctx := dl.ToCtx(context.Background())
	if ctx != context.Background() {
		t.Error("ToCtx stored a disabled logger")
	}

	l := NewLogger(
		SetWriter(ioutil.Discard),
		SetFields(String("foo", "bar")),
	)
	ctx = l.ToCtx(ctx)
	if FromCtx(ctx) != &l {
		t.Error("Clone(Context did not store logger")
	}

	l = l.Clone(SetLevel(DebugLevel))
	ctx = l.ToCtx(ctx)
	if FromCtx(ctx) != &l {
		t.Error("ToCtx did not store copied logger")
	}

	ctx = dl.ToCtx(ctx)
	if FromCtx(ctx) != &dl {
		t.Error("ToCtx did not overide logger Clone( a disabled logger")
	}
}
