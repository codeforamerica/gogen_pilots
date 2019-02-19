package data

import (
	"time"
)

func parseDate(format string, date string) time.Time {
	//TODO: handle dates that have format YYYY0000
	t, err := time.Parse(format, date)
	if err != nil {
		t = time.Time{}
	}
	//if time.Now().Before(t) {
	//	t = t.AddDate(-100, 0, 0)
	//}
	return t
}
