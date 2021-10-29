package utils

import "time"

func StampToTime(s int64) *time.Time {
	t := time.Unix(s,0)
	return &t
}
func TimeToStamp(t *time.Time) int64 {
	return  t.Unix()
}
