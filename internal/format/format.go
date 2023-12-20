package format

import (
	"time"
)

const DATE_FORMAT = "01-02-2006"

func ParseDate(dateString string) (time.Time, error) {
	date, err := time.Parse(DATE_FORMAT, dateString)
	if err != nil {
		return time.Now(), err
	}
	return date, nil
}
