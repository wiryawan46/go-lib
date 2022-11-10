package general

import "time"

func GetUnixTimestamp() int64 {
	return time.Now().Local().Unix()
}

func GetCurrentTime() time.Time {
	local := time.Now().UTC()
	timezone, err := time.LoadLocation("Asia/Jakarta")
	if err == nil {
		local = local.In(timezone)
	}

	currentTime := local.Format("2006-01-02 15:04:05")
	t, _ := time.Parse("2006-01-02 15:04:05", currentTime)

	return t
}
