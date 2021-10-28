package utils

import "time"

func GetParseTime(alertTime string) string {

	ts, err := time.Parse(time.RFC3339, alertTime)
	if err != nil {
		return err.Error()
	}
	// 设置时区
	cstSH, _ := time.LoadLocation("Asia/Shanghai")
	// if err != nil {
	// 	return err.Error()
	// }

	return ts.In(cstSH).Format("2006-01-02 15:04:05")
}
