package sane

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"time"

	"gitlab.com/bloom42/lily/encoding/sane/parser"
)

// Marshaler is the interface implemented by types that
// can marshal themselves into valid SANE.
type Marshaler interface {
	MarshalSANE() ([]byte, error)
}

// Marshal returns the SANE encoding of sata.  Behavior is similar to the Go json
// encoder, except that there is no concept of a Marshaler interface or MarshalSAN
// function for sub-structs, and currently only definite types can be marshaled
// (i.e. no `interface{}`).
// The following struct annotations are supported:
//   sane:"Field"      Overrides the field's name to output.
//   omitempty         When set, empty values and groups are not emitted.
// Note that pointers are automatically assigned the "omitempty" option, as SANE
// explicitly does not handle null values (saying instead the label should be
// dropped).
// Tree structural types and corresponding marshal types:
//   *Tree                            (*)struct, (*)map[string]interface{}
//   []*Tree                          (*)[](*)struct, (*)[](*)map[string]interface{}
//   []interface{} (as interface{})   (*)[]primitive, (*)[]([]interface{})
//   interface{}                      (*)primitive
// Tree primitive types and corresponding marshal types:
//   uint64     uint, uint8-uint64, pointers to same
//   int64      int, int8-uint64, pointers to same
//   float64    float32, float64, pointers to same
//   string     string, pointers to same
//   bool       bool, pointers to same
func Marshal(data interface{}) ([]byte, error) {
	mtype := reflect.TypeOf(data)
	if mtype.Kind() != reflect.Struct && mtype.Kind() != reflect.Map {
		return []byte{}, errors.New("Only a struct can be marshaled to SANE")
	}
	sval := reflect.ValueOf(data)
	if isCustomMarshaler(mtype) {
		return callCustomMarshaler(sval)
	}
	t, err := valueToTree(mtype, sval)
	if err != nil {
		return []byte{}, err
	}

	var buf bytes.Buffer
	_, err = t.Write(&buf, "", "", 0, true)

	return buf.Bytes(), err
}

// Convert given marshal struct or map value to sane tree
func valueToTree(mtype reflect.Type, mval reflect.Value) (*parser.Tree, error) {
	if mtype.Kind() == reflect.Ptr {
		return valueToTree(mtype.Elem(), mval.Elem())
	}
	tval := parser.NewTree()
	switch mtype.Kind() {
	case reflect.Struct:
		for i := 0; i < mtype.NumField(); i++ {
			mtypef, mvalf := mtype.Field(i), mval.Field(i)
			opts := sanOptions(mtypef)
			if opts.Include && (!opts.OmitEmpty || !isZero(mvalf)) {
				val, err := valueToSan(mtypef.Type, mvalf)
				if err != nil {
					return nil, err
				}

				tval.SetWithOptions(opts.Name, val, true)
			}
		}
	case reflect.Map:
		for _, key := range mval.MapKeys() {
			mvalf := mval.MapIndex(key)
			val, err := valueToSan(mtype.Elem(), mvalf)
			if err != nil {
				return nil, err
			}
			//if e.quoteMapKeys {
			//	keyStr, err := sanValueStringRepresentation(key.String(), "", true)
			//	if err != nil {
			//		return nil, err
			//	}
			//	tval.SetPath([]string{keyStr}, val)
			tval.Set(key.String(), val)
		}
	}
	return tval, nil
}

// Convert given marshal value to sane value
func valueToSan(mtype reflect.Type, mval reflect.Value) (interface{}, error) {
	if mtype.Kind() == reflect.Ptr {
		return valueToSan(mtype.Elem(), mval.Elem())
	}
	switch {
	case isCustomMarshaler(mtype):
		return callCustomMarshaler(mval)
	case isTree(mtype):
		return valueToTree(mtype, mval)
	case isTreeSlice(mtype):
		return valueToTreeSlice(mtype, mval)
	case isOtherSlice(mtype):
		return valueToOtherSlice(mtype, mval)
	default:
		switch mtype.Kind() {
		case reflect.Bool:
			return mval.Bool(), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return mval.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return mval.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return mval.Float(), nil
		case reflect.String:
			return mval.String(), nil
		case reflect.Struct:
			return mval.Interface().(time.Time).Format(time.RFC3339), nil
		case reflect.Interface:
			return valueToSan(mval.Elem().Type(), mval.Elem())
		default:
			return nil, fmt.Errorf("Marshal can't handle %v (%v)", mtype, mtype.Kind())
		}
	}
}

// Convert given marshal slice to slice of SANE trees
func valueToTreeSlice(mtype reflect.Type, mval reflect.Value) ([]*parser.Tree, error) {
	tval := make([]*parser.Tree, mval.Len(), mval.Len())
	for i := 0; i < mval.Len(); i++ {
		val, err := valueToTree(mtype.Elem(), mval.Index(i))
		if err != nil {
			return nil, err
		}
		tval[i] = val
	}
	return tval, nil
}

// Convert given marshal slice to slice of sane values
func valueToOtherSlice(mtype reflect.Type, mval reflect.Value) (interface{}, error) {
	tval := make([]interface{}, mval.Len(), mval.Len())
	for i := 0; i < mval.Len(); i++ {
		val, err := valueToSan(mtype.Elem(), mval.Index(i))
		if err != nil {
			return nil, err
		}
		tval[i] = val
	}
	return tval, nil
}

func callCustomMarshaler(mval reflect.Value) ([]byte, error) {
	return mval.Interface().(Marshaler).MarshalSANE()
}

func isZero(val reflect.Value) bool {
	switch val.Type().Kind() {
	case reflect.Map:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return val.Len() == 0
	default:
		return reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
	}
}
