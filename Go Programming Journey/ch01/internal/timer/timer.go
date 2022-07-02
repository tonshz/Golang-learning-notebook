package timer

import (
	"time"
)

// 封装返回当前本地时间的 Time 对象
func GetNowTime() time.Time {
	//return time.Now()
	// 设置当前时区为 Asia/Shanghai
	location, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(location)
}

// 在当前时间上加上 duration 获得最终时间
func GetCalculateTime(currentTime time.Time, d string) (time.Time, error) {
	duration, err := time.ParseDuration(d)
	if err != nil {
		return time.Time{}, err
	}
	return currentTime.Add(duration), nil
}
