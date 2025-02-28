package lib

import "time"

func GetParsedTime(date string) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
