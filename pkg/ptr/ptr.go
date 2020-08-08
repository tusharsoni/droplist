package ptr

import "time"

func Int64(val int64) *int64 {
	return &val
}

func Duration(val time.Duration) *time.Duration {
	return &val
}

func Time(val time.Time) *time.Time {
	return &val
}
