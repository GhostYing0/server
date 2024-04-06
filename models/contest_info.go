package models

type ContestForm struct {
	Contest     string `json:"contest"`
	ContestType string `json:"contest_type"`
	StartTime   string `json:"start_time"`
	Deadline    string `json:"deadline"`
}

type ContestInfo struct {
	ID          int64     `json:"id" xorm:"id"`
	Contest     string    `json:"contest" xorm:"contest"`
	ContestType string    `json:"contest_type" xorm:"contest_type"`
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

type UpdateContestParam struct {
	ID        int64  `json:"id" xorm:"id"`
	Name      string `json:"name" xorm:"name"`
	Type      string `json:"type" xorm:"type"`
	StartDate string `json:"start_date" xorm:"start_date"`
	Deadline  string `json:"deadline" xorm:"deadline"`
}

type ContestDeleteId struct {
	ID []int64 `json:"id_number"`
}

type DisplayContestForm ContestInfo

type NewContest ContestInfo

func (ContestInfo) TableName() string {
	return "contest"
}

func (DisplayContestForm) TableName() string {
	return "contest"
}
