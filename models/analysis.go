package models

type TotalEnrollCountOfPerYear struct {
	EnrollData map[string]int64 `json:"enroll_data"` //key:年份 value:报名数
}

type PreTypeEnrollCountOfPerYear struct {
	EnrollData map[string]map[string]int64 `json:"contest_type_with_enroll_data"` //KEY:年份 VALUE:{key:竞赛类型 value:数量}
}

type PreLevelEnrollCountOfPerYear struct {
	EnrollData map[string]map[string]int64 `json:"contest_level_with_enroll_data"` //KEY:年份 VALUE:{key:LEVEL value:数量}
}

type SchoolEnrollSortedCount struct {
	SchoolEnrollData []SchoolEnroll
}

type EnrollSemesterArray struct {
	Data  []EnrollSemester `json:"data"`
	Total int64            `json:"total"`
}

type EnrollSemester struct {
	Semester    string `json:"semester"`
	EnrollCount int64  `json:"enroll_count"`
}

type SchoolEnroll struct {
	School      string `json:"school"`
	EnrollCount int64  `json:"enroll_count"`
}

type SchoolEnrollCount struct {
	SchoolID int64 `json:"school_id" xorm:"school_id"`
}

type MysqlSelectEnrollYear struct {
	Date OftenTime `xorm:"create_time"`
}

type RewardRate struct {
	Rate        float64 `json:"rate"`
	RewardCount int64   `json:"reward_count"`
	EnrollCount int64   `json:"enroll_count"`
	Prize1      int64   `json:"prize1"`
	Prize2      int64   `json:"prize2"`
	Prize3      int64   `json:"prize3"`
	Prize4      int64   `json:"prize4"`
}

type MysqlSelectEnrollYearAndContestType struct {
	Date        OftenTime `xorm:"create_time"`
	ContestType int64     `xorm:"contest_type_id"`
}

type MysqlSelectEnrollYearAndContestLevel struct {
	Date         OftenTime `xorm:"create_time"`
	ContestLevel int64     `xorm:"contest_level_id"`
}

type CompareEnrollCount struct {
	EnrollCompare map[string]float64 `json:"enroll_compare"`
}

type EnrollAndSemester struct {
	Semester int64 `json:"semester_id" xorm:"semester_id"`
}
