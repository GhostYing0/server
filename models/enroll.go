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

type NewEnrollInformation struct {
	UserName    string `json:"username"`
	TeamID      string `json:"team_id"`
	ContestName string `json:"contest_name"`
	CreateTime  string `json:"create_time"`
	School      string `json:"school"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	State       int    `json:"state"`
}

type EnrollInformation struct {
	ID         int64     `json:"id" xorm:"id"`
	Username   string    `json:"username" xorm:"username"`
	UserID     int64     `json:"user_id" xorm:"user_id"`
	TeamID     string    `json:"team_id" xorm:"team_id"`
	Contest    string    `json:"contest" xorm:"contest"`
	CreateTime OftenTime `json:"create_time" xorm:"create_time"`
	School     string    `json:"school" xorm:"school"`
	Phone      string    `json:"phone" xorm:"phone"`
	Email      string    `json:"email" xorm:"email"`
	State      int       `json:"state" xorm:"state"`
	Deleted    OftenTime `json:"deleted" xorm:"deleted"`
}

type ReturnEnrollInformation struct {
	ID         int64  `json:"id" xorm:"id"`
	Username   string `json:"username" xorm:"username"`
	UserID     int64  `json:"user_id" xorm:"user_id"`
	TeamID     string `json:"team_id" xorm:"team_id"`
	Contest    string `json:"contest" xorm:"contest"`
	CreateTime string `json:"create_time" xorm:"create_time"`
	School     string `json:"school" xorm:"school"`
	Phone      string `json:"phone" xorm:"phone"`
	Email      string `json:"email" xorm:"email"`
	State      int    `json:"state" xorm:"state"`
}

func (ReturnEnrollInformation) TableName() string {
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
