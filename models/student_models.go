package models

type StudentEntryParam struct {
	ContestantID int64 `json:"contestant_id" xorm:"contestant_id"`
	ContestID    int64 `json:"contest_id" xorm:"contest_id"`
}

// 竞赛报名信息结构体
type ContestantInfo struct {
	ID           int64     `json:"id" xorm:"id"`
	ContestantID int64     `json:"contestant_id" xorm:"contestant_id"`
	ContestID    int64     `json:"contest_id" xorm:"contest_id"`
	EntryTime    OftenTime `json:"entry_time" xorm:"entry_time"`
	Deleted      OftenTime `xorm:"deleted"`
}

// 查询竞赛信息结构体
type ContestGrade struct {
	ID           int64     `xorm:"id"`
	ContestantID int64     `xorm:"contestant_id"`
	ContestID    int64     `xorm:"contest_id"`
	EntryTime    OftenTime `xorm:"entry_time"`
	Awards       string    `xorm:"awards"`
	Deleted      OftenTime `xorm:"deleted"`
}

// 显示竞赛信息结构体
type RegistrationInfo ContestGrade

type RegistrationDeleteId struct {
	ID []int64 `json:"id_number"`
}

type EntryContestParam struct {
	Contestant string `json:"contestant"`
	Contest    string `json:"contest"`
}

func (*StudentEntryParam) TableName() string {
	return "registration"
}

func (*ContestantInfo) TableName() string {
	return "registration"
}

func (*ContestGrade) TableName() string {
	return "registration"
}

func (*RegistrationInfo) TableName() string {
	return "registration"
}

func (*EntryContestParam) TableName() string {
	return "registration"
}
