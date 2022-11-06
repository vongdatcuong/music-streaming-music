package time_utils

import "time"

func GetCurrentUnixTime() uint64 {
	return uint64(time.Now().Unix())
}
