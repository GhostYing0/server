package models

import (
	"server/utils/mydebug"
	"time"
)

type OftenTime time.Time

func NewOftenTime() OftenTime {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)
	return OftenTime(t)
}

func FormatString2OftenTime(str string) OftenTime {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	if err != nil {
		mydebug.DPrintf(err)
		return NewOftenTime()
	}
	return OftenTime(t)
}

func (self OftenTime) String() string {
	t := time.Time(self)
	if t.IsZero() {
		return "0000-00-00 00:00:00"
	}
	return t.Format("2006-01-02 15:04:05")
}

func (this *OftenTime) IsZero() bool {
	t := time.Time(*this)
	return t.IsZero()
}

func (this *OftenTime) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	if str == "null" {
		return nil
	}

	if str == `"0001-01-01 08:00:00"` {
		ft := NewOftenTime()
		this = &ft
		return nil
	}

	var t time.Time
	t, err = time.ParseInLocation(`"2006-01-02 15:04:05"`, str, time.Local)
	*this = OftenTime(t)

	return
}

func MysqlFormatString2String(str string) string {
	parsedTimeMysql, err := time.Parse(time.RFC3339, str)
	if err != nil {
		mydebug.DPrintf("解析时间出错")
		return ""
	}
	stringTime := parsedTimeMysql.String()

	loc, err := time.LoadLocation("Asia/Shanghai")
	parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05 MST", stringTime[:len(stringTime)-6], loc)
	if err != nil {
		mydebug.DPrintf("解析时间出错")
		return ""
	}
	formattedTime := parsedTime.Format("2006-01-02 15:04:05")
	return formattedTime
}
