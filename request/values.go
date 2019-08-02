package request

import (
	"strconv"
	"time"
)

// Values interface for the add values
type Values interface {
	Add(key, val string)
}

// RValue configure values
type RValue func(Values)

// StringValue add string values
func StringValue(key string, val string) RValue {
	return func(values Values) {
		values.Add(key, val)
	}
}

// Int64Value add int64
func Int64Value(key string, val int64) RValue {
	return func(values Values) {
		values.Add(key, strconv.FormatInt(val, 10))
	}
}

// TimeValue add time by format
func TimeValue(key string, val time.Time, format string) RValue {
	return func(values Values) {
		if !val.IsZero() {
			values.Add(key, val.Format(format))
		}
	}
}
