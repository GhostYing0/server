package models

type EnrollForm struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	Name     string `json:"name"`
	TeamID   string `json:"team_id"`
	Contest  string `json:"contest"`
	//CreateTime string `json:"create_time"`
	School string `json:"school"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	State  int    `json:"state"`
}

type EnrollInformationForm struct {
	ID         int64  `json:"id"`
	UserName   string `json:"username"`
	Name       string `json:"name"`
	TeamID     string `json:"team_id"`
	Contest    string `json:"contest"`
	CreateTime string `json:"create_time"`
	School     string `json:"school"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	State      int    `json:"state"`
}

type EnrollInformation struct {
	ID        int64  `json:"id" xorm:"id"`
	StudentID string `json:"student_id" xorm:"student_id"`
	TeamID    string `json:"team_id" xorm:"team_id"`
	//ContestID    string    `json:"contest_id" xorm:"contest_id"`
	ContestID    int64     `json:"contest_id" xorm:"contest_id"`
	CreateTime   string    `json:"create_time" xorm:"create_time"`
	SchoolID     int64     `json:"school_id" xorm:"school_id"`
	Phone        string    `json:"phone" xorm:"phone"`
	Email        string    `json:"email" xorm:"email"`
	State        int       `json:"state" xorm:"state"`
	RejectReason string    `json:"reject_reason" xorm:"reject_reason"`
	Deleted      OftenTime `json:"deleted" xorm:"deleted"`
}

type NewEnroll struct {
	StudentID  string    `json:"student_id" xorm:"student_id"`
	TeamID     string    `json:"team_id" xorm:"team_id"`
	ContestID  int64     `json:"contest_id" xorm:"contest_id"`
	CreateTime OftenTime `json:"create_time" xorm:"create_time"`
	SchoolID   int64     `json:"school_id" xorm:"school_id"`
	Phone      string    `json:"phone" xorm:"phone"`
	Email      string    `json:"email" xorm:"email"`
	State      int       `json:"state" xorm:"state"`
}

func (NewEnroll) TableName() string {
	return "enroll_information"
}

type EnrollContestStudent struct {
	EnrollInformation `xorm:"extends"`
	Username          string    `json:"username" xorm:"username"`
	Contest           string    `json:"contest" xorm:"contest"`
	ContestState      int       `json:"contest_state" xorm:"contest_state"`
	ContestType       string    `json:"type" xorm:"type"`
	CreateTime        string    `json:"create_time" xorm:"create_time"`
	StartTime         string    `json:"start_time" xorm:"start_time"`
	Deadline          string    `json:"deadline" xorm:"deadline"`
	State             int       `json:"state" xorm:"state"`
	Describe          string    `json:"desc" xorm:"describe"`
	Deleted           OftenTime `json:"deleted" xorm:"deleted"`
	Student           `xorm:"extends"`
	//Account           `xorm:"extends"`
	School   string `xorm:"school"`
	College  string `xorm:"college"`
	Semester string `xorm:"semester"`
}

type EnrollContest struct {
	EnrollInformation `xorm:"extends"`
	Contest           `xorm:"extends"`
	Account           `xorm:"extends"`
}

func (EnrollContest) TableName() string {
	return "enroll_information"
}

func (EnrollContestStudent) TableName() string {
	return "enroll_information"
}

type ContestInfoAccount struct {
	Account     `xorm:"extends"`
	ContestInfo `xorm:"extends"`
}

type EnrollInformationReturn struct {
	ID           int64  `json:"id" xorm:"id"`
	Username     string `json:"username" xorm:"username"`
	StudentID    string `json:"student_id" xorm:"student_id"`
	TeamID       string `json:"team_id" xorm:"team_id"`
	ContestType  string `json:"contest_type" xorm:"contest_type"`
	CreateTime   string `json:"create_time" xorm:"create_time"`
	StartTime    string `json:"start_time" xorm:"start_time"`
	School       string `json:"school" xorm:"school"`
	College      string `json:"college" xorm:"college"`
	Semester     string `json:"semester" xorm:"semester"`
	Class        string `json:"student_class" xorm:"student_class"`
	Phone        string `json:"phone" xorm:"phone"`
	Email        string `json:"email" xorm:"email"`
	State        int    `json:"state" xorm:"state"`
	Name         string `json:"name" xorm:"name"`
	Contest      string `json:"contest" xorm:"contest"`
	RejectReason string `json:"reject_reason" xorm:"reject_reason"`
	DoUpload     bool   `json:"do_upload"`
}

func (EnrollInformationReturn) TableName() string {
	return "enroll_information"
}

func (EnrollInformation) TableName() string {
	return "enroll_information"
}

type EnrollDeleteId struct {
	ID []int64 `json:"ids"`
}

type EnrollIds EnrollDeleteId
type GradeIds EnrollDeleteId
