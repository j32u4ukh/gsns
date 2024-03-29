package utils

import (
	"internal/pbgo"
	"time"

	"github.com/pkg/errors"
)

// 定義時間格式
var TIME_LAYOUT string = "2006-01-02 15:04:05"

func StringToTime(s string) (time.Time, error) {
	// 使用 time.Parse() 進行轉換
	t, err := time.Parse(TIME_LAYOUT, s)
	if err != nil {
		var null time.Time
		return null, errors.Wrapf(err, "Failed to parse time from %s", s)
	}
	return t, nil
}

func TimeToUtc(t time.Time) int64 {
	return t.Unix()
}

func UtcToTime(utc int64) time.Time {
	t := time.Unix(utc, 0).UTC()
	return t
}

func UtcToTimestamp(utc int64) *pbgo.TimeStamp {
	return TimeToTimestamp(UtcToTime(utc))
}

func TimestampToUtc(ts *pbgo.TimeStamp) int64 {
	t := TimestampToTime(ts)
	return TimeToUtc(t)
}

func TimeToTimestamp(t time.Time) *pbgo.TimeStamp {
	return &pbgo.TimeStamp{
		Year:   int32(t.Year()),
		Month:  int32(t.Month()),
		Day:    int32(t.Day()),
		Hour:   int32(t.Hour()),
		Minute: int32(t.Minute()),
		Second: int32(t.Second()),
	}
}

func TimestampToTime(t *pbgo.TimeStamp) time.Time {
	if t == nil {
		return time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	} else {
		return time.Date(int(t.Year), time.Month(t.Month), int(t.Day), int(t.Hour), int(t.Minute), int(t.Second), 0, time.UTC)
	}
}

func TimeToString(t time.Time) string {
	return t.Format(TIME_LAYOUT)
}
