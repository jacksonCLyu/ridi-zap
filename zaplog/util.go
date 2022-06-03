package zaplog

import "time"

// GetInitialDelay returns the initial logRotateInitialDelay millisecond before the first event is logged.
// hour: hour of the day
// minute: minute of the hour
// second: second of the minute
// return: logRotateInitialDelay time.Duration of time.Now() and the input parse time.
//if the input parse time if Before the current time, return the Duration that between the current time and the input parse time add one day.
func GetInitialDelay(hour int, minute int, second int) (time.Duration, error) {
	if hour < 0 {
		hour = 0
	}
	if hour > 23 {
		hour = 23
	}
	if minute < 0 {
		minute = 0
	}
	if minute > 59 {
		minute = 59
	}
	if second < 0 {
		second = 0
	}
	if second > 59 {
		second = 59
	}
	localTime := time.Now().Local()
	parseTime, err := time.ParseInLocation(localTimeLayout, getFormatTimeValue(localTime, hour, minute, second), time.Local)
	if err != nil {
		return 0, err
	}
	if parseTime.Before(localTime) {
		parseTime = parseTime.AddDate(0, 0, 1)
	}
	return parseTime.Sub(localTime), nil
}
