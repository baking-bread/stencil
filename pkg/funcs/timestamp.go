package funcs

import "time"

func Timestamp(format string) string {

	var currentTime = time.Now()

	// uses the standard time library
	return currentTime.Format(format)
}
