package utils

import "time"

func GetParseTime(alertTime string) string {

	ts, err := time.Parse(time.RFC3339, alertTime)
	if err != nil {
		return err.Error()
	}

	return ts.In(&time.Location{}).Format("2006-01-02 15:04:05")
}
