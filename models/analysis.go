package models

type TotalEnrollCountOfPerYear struct {
	EnrollData map[string]int64 `json:"enroll_data"`
}

type MysqlSelectEnrollYear struct {
	Date OftenTime `xorm:"create_time"`
}
