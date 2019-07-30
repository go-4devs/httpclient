package request

import (
	"net/url"
	"strconv"
	"time"
)

// RValue configure values
type RValue func(url.Values)

// StringValue add string values
func StringValue(key string, val string) RValue {
	return func(values url.Values) {
		values.Add(key, val)
	}
}

// Int64Value add int64
func Int64Value(key string, val int64) RValue {
	return func(values url.Values) {
		values.Add(key, strconv.FormatInt(val, 10))
	}
}

// TimeValue add time by format
func TimeValue(key string, val time.Time, format string) RValue {
	return func(values url.Values) {
		if !val.IsZero() {
			values.Add(key, val.Format(format))
		}
	}
}
