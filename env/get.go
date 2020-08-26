package env

import (
	"os"
	"strconv"
)

func Get(key string) string {
	return os.Getenv(key)
}

func GetBool(key string) (ret bool, err error) {
	val := os.Getenv(key)
	ret, err = strconv.ParseBool(val)
	return
}

func GetFloat64(key string) (ret float64, err error) {
	val := os.Getenv(key)
	ret, err = strconv.ParseFloat(val, 64)
	return
}

func GetFloat32(key string) (ret float32, err error) {
	strVal := os.Getenv(key)
	val, err := strconv.ParseFloat(strVal, 32)
	if err != nil {
		return
	}
	ret = float32(val)
	return
}

func GetInt64(key string) (ret int64, err error) {
	strVal := os.Getenv(key)
	ret, err = strconv.ParseInt(strVal, 10, 64)
	return
}

func GetInt32(key string) (ret int32, err error) {
	strVal := os.Getenv(key)
	val, err := strconv.ParseInt(strVal, 10, 32)
	if err != nil {
		return
	}
	ret = int32(val)
	return
}

func GetInt16(key string) (ret int16, err error) {
	strVal := os.Getenv(key)
	val, err := strconv.ParseInt(strVal, 10, 16)
	if err != nil {
		return
	}
	ret = int16(val)
	return
}

func GetUint64(key string) (ret uint64, err error) {
	strVal := os.Getenv(key)
	ret, err = strconv.ParseUint(strVal, 10, 64)
	return
}

func GetUint32(key string) (ret uint32, err error) {
	strVal := os.Getenv(key)
	val, err := strconv.ParseUint(strVal, 10, 32)
	if err != nil {
		return
	}
	ret = uint32(val)
	return
}

func GetUint16(key string) (ret uint16, err error) {
	strVal := os.Getenv(key)
	val, err := strconv.ParseUint(strVal, 10, 16)
	if err != nil {
		return
	}
	ret = uint16(val)
	return
}
