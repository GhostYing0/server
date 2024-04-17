package models

type TotalEnrollCountOfPerYear struct {
	EnrollData map[string]int64 `json:"enroll_data"` //key:年份 value:报名数
}

type PreTypeEnrollCountOfPerYear struct {
	EnrollData map[string]map[string]int64 `json:"contest_type_with_enroll_data"` //KEY:年份 VALUE:{key:竞赛类型 value:数量}
}

type MysqlSelectEnrollYear struct {
	Date OftenTime `xorm:"create_time"`
}

type MysqlSelectEnrollYearAndContestType struct {
	Date        OftenTime `xorm:"create_time"`
	ContestType int64     `xorm:"contest_type_id"`
}

type CompareEnrollCount struct {
	EnrollCompare map[string]float64 `json:"enroll_compare"`
}
