package models

type EnrollForm struct {
	ID          int64  `json:"id"`
	UserName    string `json:"username"`
	TeamID      string `json:"team_id"`
	ContestName string `json:"contest_name"`
	CreateTime  string `json:"create_time"`
	School      string `json:"school"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	State       int    `json:"state"`
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
	ID         int64     `json:"id" xorm:"id"`
	StudentID  string    `json:"student_id" xorm:"student_id"`
	TeamID     string    `json:"team_id" xorm:"team_id"`
	ContestID  string    `json:"contest_id" xorm:"contest_id"`
	CreateTime string    `json:"create_time" xorm:"create_time"`
	School     string    `json:"school" xorm:"school"`
	Phone      string    `json:"phone" xorm:"phone"`
	Email      string    `json:"email" xorm:"email"`
	State      int       `json:"state" xorm:"state"`
	Deleted    OftenTime `json:"deleted" xorm:"deleted"`
}

type NewEnroll struct {
	StudentID  string    `json:"student_id" xorm:"student_id"`
	TeamID     string    `json:"team_id" xorm:"team_id"`
	ContestID  int64     `json:"contest_id" xorm:"contest_id"`
	CreateTime OftenTime `json:"create_time" xorm:"create_time"`
	//School     string    `json:"school" xorm:"school"`
	Phone string `json:"phone" xorm:"phone"`
	Email string `json:"email" xorm:"email"`
	State int    `json:"state" xorm:"state"`
}

func (NewEnroll) TableName() string {
	return "enroll_information"
}

type EnrollContestStudent struct {
	EnrollInformation `xorm:"extends"`
	Contest           `xorm:"extends"`
	Student           `xorm:"extends"`
	Account           `xorm:"extends"`
}

func (EnrollContestStudent) TableName() string {
	return "enroll_information"
}

type ContestInfoAccount struct {
	Account     `xorm:"extends"`
	ContestInfo `xorm:"extends"`
}

type EnrollInformationReturn struct {
	ID         int64  `json:"id" xorm:"id"`
	Username   string `json:"username" xorm:"username"`
	StudentID  string `json:"student_id" xorm:"student_id"`
	TeamID     string `json:"team_id" xorm:"team_id"`
	ContestID  string `json:"contest_id" xorm:"contest_id"`
	CreateTime string `json:"create_time" xorm:"create_time"`
	School     string `json:"school" xorm:"school"`
	Phone      string `json:"phone" xorm:"phone"`
	Email      string `json:"email" xorm:"email"`
	State      int    `json:"state" xorm:"state"`
	Name       string `json:"name" xorm:"name"`
	Contest    string `json:"contest" xorm:"contest"`
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
