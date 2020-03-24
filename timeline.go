package utils

import (
	"fmt"
	"math"
	"strings"
	"time"
)

var (
	TimeLocal, _          = time.LoadLocation("Asia/Shanghai")
	DefaultDatetimeFormat = "2006-01-02 15:04:05"
	DefaultDateFormat     = "2006-01-02"
)

func GetNow() time.Time {
	return time.Now().In(TimeLocal)
}

//返回指定时间的年月日
func GetYmd(timestamp int64) (y, m, d int) {
	var mm time.Month
	if timestamp > 0 {
		y, mm, d = time.Unix(timestamp, 0).Date()
	} else {
		y, mm, d = GetNow().Date()
	}
	m = int(mm)
	return
}

//返回指定时间的年月日时
func GetYmdHms(timestamp int64) (y, month, d, h, minute, s int) {
	var mm time.Month
	var t time.Time
	if timestamp > 0 {
		t = time.Unix(timestamp, 0)
	} else {
		t = GetNow()
	}
	y, mm, d = t.Date()
	month = int(mm)
	h = t.Hour()
	minute = t.Minute()
	s = t.Second()
	return
}

//返回指定时间的时间戳
func GetToday() (start, end int64) {
	year, month, day := GetNow().Date()
	todayStart := time.Date(year, month, day, 0, 0, 0, 0, TimeLocal)
	start = todayStart.Unix()
	end = start + 24*3600 - 1
	return
}

//返回昨天的时间戳
func GetYesterday() (start, end int64) {
	year, month, day := GetNow().Date()
	yesterdayStart := time.Date(year, month, day-1, 0, 0, 0, 0, TimeLocal)
	start = yesterdayStart.Unix()
	end = start + 24*3600 - 1
	return
}

//返回明天的时间戳
func GetTomorrow() (start, end int64) {
	year, month, day := GetNow().Date()
	tomorrowStart := time.Date(year, month, day+1, 0, 0, 0, 0, TimeLocal)
	start = tomorrowStart.Unix()
	end = start + 24*3600 - 1
	return
}

//返回本周一到本周日的时间戳
func GetWeek() (start, end int64) {
	year, month, day := GetNow().Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, TimeLocal)
	//循环到减少一天直至星期一！
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}
	start = date.Unix()
	end = start + 7*24*3600 - 1
	return
}

//返回本月1日至月末的时间戳
func GetMonth() (start, end int64) {
	year, month, _ := GetNow().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, TimeLocal)
	nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, TimeLocal)
	start = thisMonth.Unix()
	end = nextMonth.Unix() - 1
	return
}

//返回本年初到年末的时间戳
func GetYear() (start, end int64) {
	year, _, _ := GetNow().Date()
	thisYear := time.Date(year, 1, 1, 0, 0, 0, 0, TimeLocal)
	nextYear := time.Date(year+1, 1, 1, 0, 0, 0, 0, TimeLocal)
	start = thisYear.Unix()
	end = nextYear.Unix() - 1
	return
}

//返回从现在开始到N年后的时间戳, now 是否为当前时间还是23:59:59
func GetLaterYear(num int, now bool) int64 {
	t := GetNow()
	year, month, day := t.Date()
	if now {
		return time.Date(year+num, month, day, t.Hour(), t.Minute(), t.Second(), 0, TimeLocal).Unix()
	} else {
		return time.Date(year+num, month, day, 23, 59, 59, 0, TimeLocal).Unix()
	}
}

//返回上周一至上周末的时间戳
func GetLastWeek() (start, end int64) {
	start, end = GetWeek()
	start -= 7 * 24 * 3600
	end -= 7 * 24 * 3600
	return
}

//返回上个月1号到月末的时间戳。
func GetLastMonth() (start, end int64) {
	year, month, _ := GetNow().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, TimeLocal)
	startTime := thisMonth.AddDate(0, -1, 0)
	start = startTime.Unix()
	end = thisMonth.Unix() - 1 //本月第一天00:00:00减一秒。就是上个月的最后一秒。
	return
}
func GetLastYear() (start, end int64) {
	year, _, _ := GetNow().Date()
	lastYear := time.Date(year-1, 1, 1, 0, 0, 0, 0, TimeLocal)
	thisYear := time.Date(year, 1, 1, 0, 0, 0, 0, TimeLocal)
	start = lastYear.Unix()
	end = thisYear.Unix() - 1 //本年第一天00:00:00减一秒。就是去年的最后一秒。
	return
}

func GetNextMonth(timestamp int64) int64 {
	var year, day int
	var month time.Month
	if timestamp > 0 {
		year, month, day = time.Unix(timestamp, 0).Date()
	} else {
		year, month, day = GetNow().Date()
	}
	thisMonth := time.Date(year, month, day, 0, 0, 0, 0, TimeLocal)
	nextMonth := thisMonth.AddDate(0, +1, 0)
	endTime := nextMonth.Unix() + 86399
	return endTime
}

func GetNextYear(timestamp int64) (int64, int64) {
	var year, day int
	var month time.Month
	if timestamp > 0 {
		year, month, day = time.Unix(timestamp, 0).Date()
	} else {
		year, month, day = GetNow().Date()
	}
	thisYear := time.Date(year, month, day, 0, 0, 0, 0, TimeLocal)
	nextYear := thisYear.AddDate(+1, 0, 0)
	startTime := nextYear.Unix()
	endTime := startTime + 86399
	return startTime, endTime
}

//返回当前时间戳，而不是0点。
func GetDaysAgo(days int) (start, end int64) {
	end = GetNow().Unix()
	startTime := GetNow().AddDate(0, 0, -days)
	start = startTime.Unix()
	return
}

// Format 跟 PHP 中 date 类似的使用方式，如果 ts 没传递，则使用当前时间
func TimeFormat(format string, timestamp ...int64) string {
	patterns := []string{
		// 年
		"Y", "2006", // 4 位数字完整表示的年份
		"y", "06", // 2 位数字表示的年份

		// 月
		"m", "01", // 数字表示的月份，有前导零
		"n", "1", // 数字表示的月份，没有前导零
		"M", "Jan", // 三个字母缩写表示的月份
		"F", "January", // 月份，完整的文本格式，例如 January 或者 March

		// 日
		"d", "02", // 月份中的第几天，有前导零的 2 位数字
		"j", "2", // 月份中的第几天，没有前导零

		"D", "Mon", // 星期几，文本表示，3 个字母
		"l", "Monday", // 星期几，完整的文本格式;L的小写字母

		// 时间
		"g", "3", // 小时，12 小时格式，没有前导零
		"G", "15", // 小时，24 小时格式，没有前导零
		"h", "03", // 小时，12 小时格式，有前导零
		"H", "15", // 小时，24 小时格式，有前导零

		"a", "pm", // 小写的上午和下午值
		"A", "PM", // 小写的上午和下午值

		"i", "04", // 有前导零的分钟数
		"s", "05", // 秒数，有前导零
	}
	replacer := strings.NewReplacer(patterns...)
	format = replacer.Replace(format)

	t := GetNow()
	if len(timestamp) > 0 {
		t = time.Unix(timestamp[0], 0).In(TimeLocal)
	}
	return t.Format(format)
}

func StrToLocalTimeBak(value string) time.Time {
	if value == "" {
		return GetNow()
	}
	zoneName, offset := time.Now().Zone()

	zoneValue := offset / 3600 * 100
	if zoneValue > 0 {
		value += fmt.Sprintf(" +%04d", zoneValue)
	} else {
		value += fmt.Sprintf(" -%04d", zoneValue)
	}

	if zoneName != "" {
		value += " " + zoneName
	}
	return StrToLocalTime(value)
}

func StrToLocalTime(value string) time.Time {
	var t = GetNow()
	if value == "" {
		return t
	}
	layouts := []string{
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05 -0700 MST",
		"2006/01/02 15:04:05 -0700",
		"2006/01/02 15:04:05",
		"2006-01-02 -0700 MST",
		"2006-01-02 -0700",
		"2006-01-02",
		"2006/01/02 -0700 MST",
		"2006/01/02 -0700",
		"2006/01/02",
		"2006-01-02 15:04:05 -0700 -0700",
		"2006/01/02 15:04:05 -0700 -0700",
		"2006-01-02 -0700 -0700",
		"2006/01/02 -0700 -0700",
		"2006-01-02T15:04:05.000",
		"2006-1-2T15:4:5.000",
		"2006-1-2T15:04:05.000",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05 Mon",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	var err error
	for _, layout := range layouts {
		t, err = time.ParseInLocation(layout, value, TimeLocal)
		if err == nil {
			return t
		}
	}
	return t
}

func GetDate(timestamp int64) string {
	if timestamp == 0 {
		return GetNow().Format("2006-01-02")
	}
	return time.Unix(timestamp, 0).Format("2006-01-02")
}
func GetDateTime(timestamp int64) string {
	if timestamp == 0 {
		return GetNow().Format("2006-01-02 15:04:05")
	}
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}
func GetTime(timestamp int64) string {
	if timestamp == 0 {
		return GetNow().Format("15:04:05")
	}
	return time.Unix(timestamp, 0).In(TimeLocal).Format("15:04:05")
}

func StrToTime(dateText, timeLayout string) (timestamp int64) {
	if dateText == "" {
		return 0
	}
	//时间模板用 2006-01-02 15:04:05 ，据说是golang的诞生时间。
	var timeFormat string
	if timeLayout == "date" {
		timeFormat = "2006-01-02"
	} else if timeLayout == "datetime" {
		timeFormat = "2006-01-02 15:04:05"
	} else {
		timeFormat = timeLayout
	}
	theTime, err := time.ParseInLocation(timeFormat, dateText, TimeLocal) //使用模板在对应时区转化为time.time类型
	if err != nil {
		panic(err.Error())
	}
	timestamp = theTime.Unix() //转化为时间戳 类型是int64
	return
}

/*
兼容常用时间格式正则
[0-9]{4,}(-|/)[0-9]{1,}(-|/)[0-9]{1,}( ||T)[0-9]{1,}(:[0-9]{1,}(:[0-9]{1,}|)|)
*/
//传毫秒
func GetDurationDesc(duration int64) string {
	var day, hour, minute int
	second := int(math.Round(float64(duration) / 1000))
	if second > 60 {
		//获取分钟，除以60取整数，得到整数分钟
		minute = int(second / 60)
		//获取秒数，秒数取佘，得到整数秒数
		second = int(second % 60)
		//如果分钟大于60，将分钟转换成小时
		if minute > 60 {
			//获取小时，获取分钟除以60，得到整数小时
			hour = int(minute / 60)
			//获取小时后取佘的分，获取分钟除以60取佘的分
			minute = int(minute % 60)
			if hour > 24 {
				//获取天，获取分钟除以24，得到整数天
				day = int(hour / 24)
				//获取天后取佘的时，获取小时除以24取佘的时
				hour = int(hour % 24)
			}
		}
	}
	var result = fmt.Sprintf("%d秒", second)
	if minute > 0 {
		result = fmt.Sprintf("%d分", minute) + result
	}
	if hour > 0 {
		result = fmt.Sprintf("%d小时", hour) + result
	}
	if day > 0 {
		result = fmt.Sprintf("%d天", day) + result
	}
	return result
}

func DaysOfThisMonth() int {
	year, month, _ := GetNow().Date()
	return DaysOfMonth(year, int(month))
}

func DaysOfMonth(year, month int) (days int) {
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30
		} else {
			days = 31
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	}
	return
}

func GetMonthLater(timestamp int64, months int) int64 {
	if months < 0 {
		return 0
	}
	var year, day int
	var month time.Month
	if timestamp > 0 {
		year, month, day = time.Unix(timestamp, 0).Date()
	} else {
		year, month, day = GetNow().Date()
	}
	thisMonth := time.Date(year, month, day, 0, 0, 0, 0, TimeLocal)
	nextMonth := thisMonth.AddDate(0, +months, 0)
	endTime := nextMonth.Unix() + 86399
	return endTime
}
