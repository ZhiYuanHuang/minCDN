package cmd

import "time"

func UTCNow() time.Time {
	return time.Now().UTC()
}
