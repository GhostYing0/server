package models

type ContestForm struct {
	Contest     string `json:"contest"`
	Username    string `json:"username"`
	ContestType string `json:"contest_type"`
	StartTime   string `json:"start_time"`
	Deadline    string `json:"deadline"`
	Describe    string `json:"desc"`
	State       int    `json:"state"`
}

type ContestInfo struct {
	ID           int64     `json:"id" xorm:"id"`
	TeacherID    string    `json:"teacher_id" xorm:"teacher_id"`
	SchoolID     int64     `json:"school_id" xorm:"school_id"`
	CollegeID    int64     `json:"college_id" xorm:"college_id"`
	ContestState int       `json:"contest_state" xorm:"contest_state"`
	Contest      string    `json:"contest" xorm:"contest"`
	ContestType  int64     `json:"contest_type_id" xorm:"contest_type_id"`
	CreateTime   OftenTime `json:"create_time" xorm:"create_time"`
	StartTime    OftenTime `json:"start_time" xorm:"start_time"`
	Deadline     OftenTime `json:"deadline" xorm:"deadline"`
	State        int       `json:"state" xorm:"state"`
	Describe     string    `json:"describe" xorm:"describe"`
	Deleted      OftenTime `json:"deleted" xorm:"deleted"`
}

type ContestInfoType struct {
	ContestInfo `xorm:"extends"`
	ContestType `xorm:"extends"`
}

func (ContestInfoType) TableName() string {
	return "contest"
}

type ContestReturn struct {
	ID           int64     `json:"id" xorm:"id"`
	ContestState int       `json:"contest_state" xorm:"contest_state"`
	Username     string    `json:"username" xorm:"username"`
	Name         string    `json:"name" xorm:"name"`
	School       string    `json:"school" xorm:"school"`
	College      string    `json:"college" xorm:"college"`
	Contest      string    `json:"contest" xorm:"contest"`
	ContestType  string    `json:"contest_type" xorm:"contest_type"`
	CreateTime   string    `json:"create_time" xorm:"create_time"`
	StartTime    string    `json:"start_time" xorm:"start_time"`
	Deadline     string    `json:"deadline" xorm:"deadline"`
	State        int       `json:"state" xorm:"state"`
	Describe     string    `json:"desc" xorm:"describe"`
	Deleted      OftenTime `json:"deleted" xorm:"deleted"`
}

type Contest struct {
	ID           int64     `json:"id" xorm:"id"`
	Username     string    `json:"username" xorm:"username"`
	Contest      string    `json:"contest" xorm:"contest"`
	ContestState int       `json:"contest_state" xorm:"contest_state"`
	ContestType  string    `json:"type" xorm:"type"`
	CreateTime   string    `json:"create_time" xorm:"create_time"`
	StartTime    string    `json:"start_time" xorm:"start_time"`
	Deadline     string    `json:"deadline" xorm:"deadline"`
	State        int       `json:"state" xorm:"state"`
	Describe     string    `json:"desc" xorm:"describe"`
	Deleted      OftenTime `json:"deleted" xorm:"deleted"`
}

type ContestContestTypeTeacher struct {
	Contest `xorm:"extends"`
	Name    string `xorm:"name"`
	School  string `xorm:"school"`
	College string `xorm:"college"`
}

func (ContestContestTypeTeacher) TableName() string {
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

type ContestAndType struct {
	Contest string `json:"contest" xrom:"contest"`
	Type    string `json:"type" xorm:"type"`
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
