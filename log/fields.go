package log

import (
	"net"
	"sort"
	"time"
)

func (e *Event) appendFields(dst []byte, fields map[string]interface{}) []byte {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		dst = enc.AppendKey(dst, key)
		val := fields[key]
		if val, ok := val.(ObjectMarshaler); ok {
			e := newEvent(nil, 0)
			e.buf = e.buf[:0]
			e.appendObject(val)
			dst = append(dst, e.buf...)
			putEvent(e)
			continue
		}
		switch val := val.(type) {
		case string:
			dst = enc.AppendString(dst, val)
		case []byte:
			dst = enc.AppendBytes(dst, val)
		case error:
			marshaled := ErrorMarshalFunc(val)
			switch m := marshaled.(type) {
			case ObjectMarshaler:
				e := newEvent(nil, 0)
				e.buf = e.buf[:0]
				e.appendObject(m)
				dst = append(dst, e.buf...)
				putEvent(e)
			case error:
				dst = enc.AppendString(dst, m.Error())
			case string:
				dst = enc.AppendString(dst, m)
			default:
				dst = enc.AppendInterface(dst, m)
			}
		case []error:
			dst = enc.AppendArrayStart(dst)
			for i, err := range val {
				marshaled := ErrorMarshalFunc(err)
				switch m := marshaled.(type) {
				case ObjectMarshaler:
					e := newEvent(nil, 0)
					e.buf = e.buf[:0]
					e.appendObject(m)
					dst = append(dst, e.buf...)
					putEvent(e)
				case error:
					dst = enc.AppendString(dst, m.Error())
				case string:
					dst = enc.AppendString(dst, m)
				default:
					dst = enc.AppendInterface(dst, m)
				}

				if i < (len(val) - 1) {
					dst = enc.AppendArrayDelim(dst)
				}
			}
			dst = enc.AppendArrayEnd(dst)
		case bool:
			dst = enc.AppendBool(dst, val)
		case int:
			dst = enc.AppendInt(dst, val)
		case int8:
			dst = enc.AppendInt8(dst, val)
		case int16:
			dst = enc.AppendInt16(dst, val)
		case int32:
			dst = enc.AppendInt32(dst, val)
		case int64:
			dst = enc.AppendInt64(dst, val)
		case uint:
			dst = enc.AppendUint(dst, val)
		case uint8:
			dst = enc.AppendUint8(dst, val)
		case uint16:
			dst = enc.AppendUint16(dst, val)
		case uint32:
			dst = enc.AppendUint32(dst, val)
		case uint64:
			dst = enc.AppendUint64(dst, val)
		case float32:
			dst = enc.AppendFloat32(dst, val)
		case float64:
			dst = enc.AppendFloat64(dst, val)
		case time.Time:
			dst = enc.AppendTime(dst, val, DefaultTimeFieldFormat)
		case time.Duration:
			dst = enc.AppendDuration(dst, val, DurationFieldUnit, DurationFieldInteger)
		case *string:
			if val != nil {
				dst = enc.AppendString(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *bool:
			if val != nil {
				dst = enc.AppendBool(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *int:
			if val != nil {
				dst = enc.AppendInt(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *int8:
			if val != nil {
				dst = enc.AppendInt8(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *int16:
			if val != nil {
				dst = enc.AppendInt16(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *int32:
			if val != nil {
				dst = enc.AppendInt32(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *int64:
			if val != nil {
				dst = enc.AppendInt64(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *uint:
			if val != nil {
				dst = enc.AppendUint(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *uint8:
			if val != nil {
				dst = enc.AppendUint8(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *uint16:
			if val != nil {
				dst = enc.AppendUint16(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *uint32:
			if val != nil {
				dst = enc.AppendUint32(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *uint64:
			if val != nil {
				dst = enc.AppendUint64(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *float32:
			if val != nil {
				dst = enc.AppendFloat32(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *float64:
			if val != nil {
				dst = enc.AppendFloat64(dst, *val)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *time.Time:
			if val != nil {
				dst = enc.AppendTime(dst, *val, DefaultTimeFieldFormat)
			} else {
				dst = enc.AppendNil(dst)
			}
		case *time.Duration:
			if val != nil {
				dst = enc.AppendDuration(dst, *val, DurationFieldUnit, DurationFieldInteger)
			} else {
				dst = enc.AppendNil(dst)
			}
		case []string:
			dst = enc.AppendStrings(dst, val)
		case []bool:
			dst = enc.AppendBools(dst, val)
		case []int:
			dst = enc.AppendInts(dst, val)
		case []int8:
			dst = enc.AppendInts8(dst, val)
		case []int16:
			dst = enc.AppendInts16(dst, val)
		case []int32:
			dst = enc.AppendInts32(dst, val)
		case []int64:
			dst = enc.AppendInts64(dst, val)
		case []uint:
			dst = enc.AppendUints(dst, val)
		// case []uint8:
		// 	dst = enc.AppendUints8(dst, val)
		case []uint16:
			dst = enc.AppendUints16(dst, val)
		case []uint32:
			dst = enc.AppendUints32(dst, val)
		case []uint64:
			dst = enc.AppendUints64(dst, val)
		case []float32:
			dst = enc.AppendFloats32(dst, val)
		case []float64:
			dst = enc.AppendFloats64(dst, val)
		case []time.Time:
			dst = enc.AppendTimes(dst, val, DefaultTimeFieldFormat)
		case []time.Duration:
			dst = enc.AppendDurations(dst, val, DurationFieldUnit, DurationFieldInteger)
		case nil:
			dst = enc.AppendNil(dst)
		case net.IP:
			dst = enc.AppendIPAddr(dst, val)
		case net.IPNet:
			dst = enc.AppendIPPrefix(dst, val)
		case net.HardwareAddr:
			dst = enc.AppendMACAddr(dst, val)
		default:
			dst = enc.AppendInterface(dst, val)
		}
	}
	return dst
}
