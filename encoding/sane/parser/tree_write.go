package parser

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Write write a value to w
func (t *Tree) Write(w io.Writer, indent, keyspace string, bytesCount int64, arraysOneElementPerLine bool) (int64, error) {
	simpleValuesKeys := make([]string, 0)
	complexValuesKeys := make([]string, 0)

	for k := range t.Values {
		v := t.Values[k]
		switch v.(type) {
		case *Tree, []*Tree:
			complexValuesKeys = append(complexValuesKeys, k)
		default:
			simpleValuesKeys = append(simpleValuesKeys, k)
		}
	}

	sort.Strings(simpleValuesKeys)
	sort.Strings(complexValuesKeys)

	for _, k := range simpleValuesKeys {
		v, ok := t.Values[k].(*Value)
		if !ok {
			return bytesCount, fmt.Errorf("invalid value type at %s: %T", k, t.Values[k])
		}

		stringValue, err := sanValueStringRepresentation(v, indent, arraysOneElementPerLine)
		if err != nil {
			return bytesCount, err
		}

		writtenBytesCount, err := writeStrings(w, indent, "", k, " = ", stringValue, "\n")
		bytesCount += int64(writtenBytesCount)
		if err != nil {
			return bytesCount, err
		}

	}

	for _, k := range complexValuesKeys {
		v := t.Values[k]

		switch node := v.(type) {
		// node has to be of those two types given how keys are sorted above
		case *Tree:
			_, ok := t.Values[k].(*Tree)
			if !ok {
				return bytesCount, fmt.Errorf("invalid value type at %s: %T", k, t.Values[k])
			}
			writtenBytesCount, err := writeStrings(w, "", indent, "", "", k, " = {\n")
			bytesCount += int64(writtenBytesCount)
			if err != nil {
				return bytesCount, err
			}
			bytesCount, err := node.Write(w, indent+"  ", k, bytesCount, arraysOneElementPerLine)
			if err != nil {
				return bytesCount, err
			}
			writtenBytesCount, err = writeStrings(w, indent, "}\n")
			bytesCount += int64(writtenBytesCount)
			if err != nil {
				return bytesCount, err
			}
		case []*Tree:
			writtenBytesCount, err := writeStrings(w, "\n", indent, k, " = [")
			bytesCount += int64(writtenBytesCount)
			if err != nil {
				return bytesCount, err
			}
			for _, subTree := range node {
				writtenBytesCount, err := writeStrings(w, "\n", indent+"  ", "{\n")
				bytesCount += int64(writtenBytesCount)
				if err != nil {
					return bytesCount, err
				}
				bytesCount, err = subTree.Write(w, indent+"    ", k, bytesCount, arraysOneElementPerLine)
				if err != nil {
					return bytesCount, err
				}
				writtenBytesCount, err = writeStrings(w, indent+"  ", "},")
				bytesCount += int64(writtenBytesCount)
				if err != nil {
					return bytesCount, err
				}
			}
			writtenBytesCount, err = writeStrings(w, "\n", indent, "]\n")
			bytesCount += int64(writtenBytesCount)
			if err != nil {
				return bytesCount, err
			}
		}
	}

	return bytesCount, nil
}

func sanValueStringRepresentation(v interface{}, indent string, arraysOneElementPerLine bool) (string, error) {
	// this interface check is added to dereference the change made in the Write function.
	// That change was made to allow this function to see formatting options.
	tv, ok := v.(*Value)
	if ok {
		v = tv.Val
	} else {
		tv = &Value{}
	}

	switch value := v.(type) {
	case uint64:
		return strconv.FormatUint(value, 10), nil
	case int64:
		return strconv.FormatInt(value, 10), nil
	case float64:
		// Ensure a round float does contain a decimal point. Otherwise feeding
		// the output back to the parser would convert to an integer.
		if math.Trunc(value) == value {
			return strings.ToLower(strconv.FormatFloat(value, 'f', 1, 32)), nil
		}
		return strings.ToLower(strconv.FormatFloat(value, 'f', -1, 32)), nil
	case string:
		return "\"" + encodeSanString(value) + "\"", nil
	case []byte:
		b, _ := v.([]byte)
		return sanValueStringRepresentation(string(b), indent, arraysOneElementPerLine)
	case bool:
		if value {
			return "true", nil
		}
		return "false", nil
	case time.Time:
		return value.Format(time.RFC3339), nil
	case nil:
		return "", nil
	}

	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Slice {
		var values []string
		for i := 0; i < rv.Len(); i++ {
			item := rv.Index(i).Interface()
			itemRepr, err := sanValueStringRepresentation(item, indent, arraysOneElementPerLine)
			if err != nil {
				return "", err
			}
			values = append(values, itemRepr)
		}
		if arraysOneElementPerLine && len(values) > 1 {
			stringBuffer := bytes.Buffer{}
			valueIndent := indent + "  "

			stringBuffer.WriteString("[\n")

			for _, value := range values {
				stringBuffer.WriteString(valueIndent)
				stringBuffer.WriteString(value)
				stringBuffer.WriteString(`,`)
				stringBuffer.WriteString("\n")
			}

			stringBuffer.WriteString(indent + "]")

			return stringBuffer.String(), nil
		}
		return "[" + strings.Join(values, ",") + "]", nil
	}
	return "", fmt.Errorf("unsupported value type %T: %v", v, v)
}

func writeStrings(w io.Writer, s ...string) (int, error) {
	var n int
	for i := range s {
		b, err := io.WriteString(w, s[i])
		n += b
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func encodeSanString(value string) string {
	var b bytes.Buffer

	for _, rr := range value {
		switch rr {
		case '\b':
			b.WriteString(`\b`)
		case '\t':
			b.WriteString(`\t`)
		case '\n':
			b.WriteString(`\n`)
		case '\f':
			b.WriteString(`\f`)
		case '\r':
			b.WriteString(`\r`)
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		default:
			intRr := uint16(rr)
			if intRr < 0x001F {
				b.WriteString(fmt.Sprintf("\\u%0.4X", intRr))
			} else {
				b.WriteRune(rr)
			}
		}
	}
	return b.String()
}
