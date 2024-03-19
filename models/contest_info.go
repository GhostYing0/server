package models

type ContestParam struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	StartDate string `json:"start_date"`
	Deadline  string `json:"deadline"`
}

type ContestInfo struct {
	ID        int64     `xorm:"id"`
	Name      string    `xorm:"name"`
	Type      string    `xorm:"type"`
	StartDate OftenTime `xorm:"start_date"`
	Deadline  OftenTime `xorm:"deadline"`
	Deleted   OftenTime `xorm:"deleted"`
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
