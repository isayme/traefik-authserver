package util

import "time"

func Millisecond(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func Second(t time.Time) int64 {
	return t.UnixNano() / int64(time.Second)
}
