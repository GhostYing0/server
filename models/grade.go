package models

type GradeForm struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Contest     string `json:"contest"`
	School      string `json:"school"`
	Grade       string `json:"grade"`
	Certificate string `json:"certificate"`
	CreateTime  string `json:"create_time"`
	State       int    `json:"state" xorm:"state"`
}

type GradeInformation struct {
	ID          int64     `json:"id" xorm:"id"`
	StudentID   string    `json:"student_id" xorm:"student_id" `
	ContestID   int64     `json:"contest_id" xorm:"contest_id"`
	SchoolID    int64     `json:"school_id" xorm:"school_id"`
	Grade       string    `json:"grade" xorm:"grade"`
	Certificate string    `json:"certificate" xorm:"certificate"`
	State       int       `json:"state" xorm:"state"`
	CreateTime  string    `json:"create_time" xorm:"create_time"`
	Deleted     OftenTime `json:"deleted" xorm:"deleted"`
}

type GradeStudentSchoolContestAccount struct {
	GradeInformation `xorm:"extends"`
	Student          `xorm:"extends"`
	School           `xorm:"extends"`
	Contest          `xorm:"extends"`
	Account          `xorm:"extends"`
}

func (GradeStudentSchoolContestAccount) TableName() string {
	return "grade"
}

type ReturnGradeInformation struct {
	ID          int64  `json:"id" xorm:"id"`
	Username    string `json:"username" xorm:"username"`
	Name        string `json:"name" xorm:"name" `
	Contest     string `json:"contest" xorm:"contest"`
	School      string `json:"school" xorm:"school"`
	Grade       string `json:"grade" xorm:"grade"`
	Certificate string `json:"certificate" xorm:"certificate"`
	State       int    `json:"state" xorm:"state"`
	CreateTime  string `json:"create_time" xorm:"create_time"`
}

func (GradeInformation) TableName() string {
	return "grade"
}
