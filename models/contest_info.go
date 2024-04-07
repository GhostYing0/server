package models

type ContestForm struct {
	Contest     string `json:"contest"`
	ContestType string `json:"contest_type"`
	StartTime   string `json:"start_time"`
	Deadline    string `json:"deadline"`
	State       int    `json:"state"`
}

type ContestInfo struct {
	ID          int64     `json:"id" xorm:"id"`
	Contest     string    `json:"contest" xorm:"contest"`
	ContestType int64     `json:"contest_type_id" xorm:"contest_type_id"`
	CreateTime  OftenTime `json:"create_time" xorm:"create_time"`
	StartTime   OftenTime `json:"start_time" xorm:"start_time"`
	Deadline    OftenTime `json:"deadline" xorm:"deadline"`
	State       int       `json:"state" xorm:"state"`
	Deleted     OftenTime `json:"deleted" xorm:"deleted"`
}

type ContestReturn struct {
	ID          int64     `json:"id" xorm:"id"`
	Contest     string    `json:"contest" xorm:"contest"`
	ContestType string    `json:"contest_type" xorm:"contest_type"`
	CreateTime  string    `json:"create_time" xorm:"create_time"`
	StartTime   string    `json:"start_time" xorm:"start_time"`
	Deadline    string    `json:"deadline" xorm:"deadline"`
	State       int       `json:"state" xorm:"state"`
	Deleted     OftenTime `json:"deleted" xorm:"deleted"`
}

type Contest struct {
	ID          int64     `json:"id" xorm:"id"`
	Contest     string    `json:"contest" xorm:"contest"`
	ContestType string    `json:"contest_type" xorm:"contest_type"`
	CreateTime  string    `json:"create_time" xorm:"create_time"`
	StartTime   string    `json:"start_time" xorm:"start_time"`
	Deadline    string    `json:"deadline" xorm:"deadline"`
	State       int       `json:"state" xorm:"state"`
	Deleted     OftenTime `json:"deleted" xorm:"deleted"`
}

type ContestContestType struct {
	Contest     `xorm:"extends"`
	ContestType string `xorm:"type"`
}

func (ContestContestType) TableName() string {
	return "contest"
}

type UpdateContestForm Contest

type ContestDeleteId struct {
	ID []int64 `json:"ids"`
}

type ContestType struct {
	ContestTypeID int64  `json:"id" xorm:"id"`
	ContestType   string `json:"type" xorm:"type"`
}

type DisplayContestForm ContestInfo

type NewContest ContestInfo

func (ContestInfo) TableName() string {
	return "contest"
}

func (DisplayContestForm) TableName() string {
	return "contest"
}

func (ContestReturn) TableName() string {
	return "contest"
}

func (ContestType) TableName() string {
	return "contest_type"
}
