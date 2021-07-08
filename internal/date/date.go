package date

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

func Parse(date string) (time.Time, error) {
	re := regexp.MustCompile(`(?m)(\d{4})(\d{2})(\d{2})`)
	matches := re.FindStringSubmatch(date)

	if len(matches) == 4 {
		var Y, m, d int

		if val, err := strconv.Atoi(matches[1]); err == nil {
			Y = val
		}

		if val, err := strconv.Atoi(matches[2]); err == nil {
			m = val
		}

		if val, err := strconv.Atoi(matches[3]); err == nil {
			d = val
		}

		return time.Date(Y, time.Month(m), d, 0, 0, 0, 0, time.Local), nil
	}

	return time.Now(), errors.New("invalid date format")
}
