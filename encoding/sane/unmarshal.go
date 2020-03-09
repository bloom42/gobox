package sane

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"gitlab.com/bloom42/lily/encoding/sane/parser"
)

var marshalerType = reflect.TypeOf(new(Marshaler)).Elem()

type Options struct {
	Name      string
	Include   bool
	OmitEmpty bool
}

// Decoder reads and decodes E values from an input stream.
type Decoder struct {
	r    io.Reader
	tree *parser.Tree
}

// Unmarshaler is the interface implemented by types that can unmarshal a SANE description of
// themselves.
// The input can be assumed to be a valid encoding of a SANE value. UnmarshalSANE must copy the
// SANE data if it wishes to retain the data after returning.
type Unmarshaler interface {
	UnmarshalSANE([]byte) error
}

// format error utility
func e(format string, args ...interface{}) error {
	return fmt.Errorf("sane: "+format, args...)
}

// Unmarshal decodes the contents of `p` in SANE format into a pointer `v`.
func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return e("unmarshalling of non-pointer %s", reflect.TypeOf(v))
	}
	if rv.IsNil() {
		return e("unmarshalling of nil %s", reflect.TypeOf(v))
	}
	if rv.Elem().Kind() != reflect.Struct {
		return e("only a pointer to struct can be unmarshaled from SANE")
	}

	tree, err := parser.Parse(data)
	if err != nil {
		return err
	}

	typ := reflect.TypeOf(v)

	treeVal, err := valueFromTree(typ.Elem(), tree)
	if err != nil {
		return err
	}

	reflect.ValueOf(v).Elem().Set(treeVal)
	return nil
}

// Convert sane tree to marshal struct or map, using marshal type
func valueFromTree(mtype reflect.Type, treeVal *parser.Tree) (reflect.Value, error) {
	if mtype.Kind() == reflect.Ptr {
		return unwrapPointer(mtype, treeVal)
	}
	var mval reflect.Value
	switch mtype.Kind() {
	case reflect.Struct:
		mval = reflect.New(mtype).Elem()
		for i := 0; i < mtype.NumField(); i++ {
			mtypef := mtype.Field(i)
			opts := sanOptions(mtypef)
			if opts.Include {
				baseKey := opts.Name
				keysToTry := []string{baseKey, strings.ToLower(baseKey), strings.ToTitle(baseKey)}
				for _, key := range keysToTry {
					exists := treeVal.Has(key)
					if !exists {
						continue
					}
					val := treeVal.Get(key)
					mvalf, err := valueFromSan(mtypef.Type, val)
					if err != nil {
						return mval, e(err.Error()) //formatError(err, tval.GetPosition(key))
					}
					mval.Field(i).Set(mvalf)
					break
				}
			}
		}
	case reflect.Map:
		mval = reflect.MakeMap(mtype)
		for _, key := range treeVal.Keys() {
			val := treeVal.GetPath([]string{key})
			mvalf, err := valueFromSan(mtype.Elem(), val)
			if err != nil {
				return mval, e(err.Error()) //formatError(err, tval.GetPosition(key))
			}
			mval.SetMapIndex(reflect.ValueOf(key), mvalf)
		}
	}
	return mval, nil
}

func unwrapPointer(mtype reflect.Type, tval interface{}) (reflect.Value, error) {
	val, err := valueFromSan(mtype.Elem(), tval)
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	mval := reflect.New(mtype.Elem())
	mval.Elem().Set(val)
	return mval, nil
}

// Convert sane value to marshal value, using marshal type
func valueFromSan(mtype reflect.Type, tval interface{}) (reflect.Value, error) {
	if mtype.Kind() == reflect.Ptr {
		return unwrapPointer(mtype, tval)
	}

	switch tval.(type) {
	case *parser.Tree:
		if isTree(mtype) {
			return valueFromTree(mtype, tval.(*parser.Tree))
		}
		return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to a tree", tval, tval)
	case []*parser.Tree:
		if isTreeSlice(mtype) {
			return valueFromTreeSlice(mtype, tval.([]*parser.Tree))
		}
		return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to trees", tval, tval)
	case []interface{}:
		if mtype.Kind() == reflect.Interface {
			mtype = reflect.SliceOf(mtype)
		}
		if isOtherSlice(mtype) {
			return valueFromOtherSlice(mtype, tval.([]interface{}))
		}
		return reflect.ValueOf(nil), fmt.Errorf("Can't  convert %v(%T) to a slice", tval, tval)
	default:
		switch mtype.Kind() {
		case reflect.Bool, reflect.Struct:
			val := reflect.ValueOf(tval)
			// if this passes for when mtype is reflect.Struct, tval is a time.Time
			if !val.Type().ConvertibleTo(mtype) {
				return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v", tval, tval, mtype.String())
			}

			return val.Convert(mtype), nil
		case reflect.String:
			val := reflect.ValueOf(tval)
			if !val.Type().ConvertibleTo(mtype) {
				return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v", tval, tval, mtype.String())
			}

			return val.Convert(mtype), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val := reflect.ValueOf(tval)
			if !val.Type().ConvertibleTo(mtype) {
				return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v", tval, tval, mtype.String())
			}
			if reflect.Indirect(reflect.New(mtype)).OverflowInt(val.Int()) {
				return reflect.ValueOf(nil), fmt.Errorf("%v(%T) would overflow %v", tval, tval, mtype.String())
			}

			return val.Convert(mtype), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			val := reflect.ValueOf(tval)
			if !val.Type().ConvertibleTo(mtype) {
				return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v", tval, tval, mtype.String())
			}
			if val.Int() < 0 {
				return reflect.ValueOf(nil), fmt.Errorf("%v(%T) is negative so does not fit in %v", tval, tval, mtype.String())
			}
			if reflect.Indirect(reflect.New(mtype)).OverflowUint(uint64(val.Int())) {
				return reflect.ValueOf(nil), fmt.Errorf("%v(%T) would overflow %v", tval, tval, mtype.String())
			}

			return val.Convert(mtype), nil
		case reflect.Float32, reflect.Float64:
			val := reflect.ValueOf(tval)
			if !val.Type().ConvertibleTo(mtype) {
				return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v", tval, tval, mtype.String())
			}
			if reflect.Indirect(reflect.New(mtype)).OverflowFloat(val.Float()) {
				return reflect.ValueOf(nil), fmt.Errorf("%v(%T) would overflow %v", tval, tval, mtype.String())
			}

			return val.Convert(mtype), nil
		case reflect.Interface:
			val := reflect.ValueOf(tval)
			if !val.Type().ConvertibleTo(mtype) {
				return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v", tval, tval, mtype.String())
			}
			return val.Convert(mtype), nil
		default:
			return reflect.ValueOf(nil), fmt.Errorf("Can't convert %v(%T) to %v(%v)", tval, tval, mtype, mtype.Kind())
		}
	}
}

// Convert sane value to marshal struct/map slice, using marshal type
func valueFromTreeSlice(mtype reflect.Type, tval []*parser.Tree) (reflect.Value, error) {
	mval := reflect.MakeSlice(mtype, len(tval), len(tval))
	for i := 0; i < len(tval); i++ {
		val, err := valueFromTree(mtype.Elem(), tval[i])
		if err != nil {
			return mval, err
		}
		mval.Index(i).Set(val)
	}
	return mval, nil
}

// Convert sane value to marshal primitive slice, using marshal type
func valueFromOtherSlice(mtype reflect.Type, tval []interface{}) (reflect.Value, error) {
	mval := reflect.MakeSlice(mtype, len(tval), len(tval))
	for i := 0; i < len(tval); i++ {
		val, err := valueFromSan(mtype.Elem(), tval[i])
		if err != nil {
			return mval, err
		}
		mval.Index(i).Set(val)
	}
	return mval, nil
}

func sanOptions(vf reflect.StructField) Options {
	tag := vf.Tag.Get("sane")
	parse := strings.Split(tag, ",")
	result := Options{Name: vf.Name, Include: true, OmitEmpty: false}
	if parse[0] != "" {
		if parse[0] == "-" && len(parse) == 1 {
			result.Include = false
		} else {
			result.Name = strings.Trim(parse[0], " ")
		}
	}
	if vf.PkgPath != "" {
		result.Include = false
	}
	if len(parse) > 1 && strings.Trim(parse[1], " ") == "omitempty" {
		result.OmitEmpty = true
	}
	if vf.Type.Kind() == reflect.Ptr {
		result.OmitEmpty = true
	}
	return result
}

// Check if the given marshall type maps to a Tree primitive
func isPrimitive(mtype reflect.Type) bool {
	switch mtype.Kind() {
	case reflect.Ptr:
		return isPrimitive(mtype.Elem())
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		return true
	case reflect.Struct:
		return isCustomMarshaler(mtype)
	case reflect.Interface:
		return true
	default:
		return false
	}
}

// Check if the given marshall type maps to a Tree
func isTree(mtype reflect.Type) bool {
	switch mtype.Kind() {
	case reflect.Map:
		return true
	case reflect.Struct:
		return !isPrimitive(mtype)
	default:
		return false
	}
}

// Check if the given marshall type maps to a Tree slice
func isTreeSlice(mtype reflect.Type) bool {
	switch mtype.Kind() {
	case reflect.Slice:
		return !isOtherSlice(mtype)
	default:
		return false
	}
}

// Check if the given marshall type maps to a non-Tree slice
func isOtherSlice(mtype reflect.Type) bool {
	switch mtype.Kind() {
	case reflect.Ptr:
		return isOtherSlice(mtype.Elem())
	case reflect.Slice:
		return isPrimitive(mtype.Elem()) || isOtherSlice(mtype.Elem())
	default:
		return false
	}
}

func isCustomMarshaler(mtype reflect.Type) bool {
	return mtype.Implements(marshalerType)
}
