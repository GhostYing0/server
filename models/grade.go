package models

type GradeForm struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Contest     string `json:"contest_name"`
	Grade       string `json:"grade"`
	Certificate string `json:"certificate"`
	CreateTime  string `json:"create_time"`
	State       int    `json:"state" xorm:"state"`
}

type GradeInformation struct {
	ID          int64     `json:"id" xorm:"id"`
	Username    string    `json:"username" xorm:"username" `
	Contest     string    `json:"contest" xorm:"contest"`
	Grade       string    `json:"grade" xorm:"grade"`
	Certificate string    `json:"certificate" xorm:"certificate"`
	State       int       `json:"state" xorm:"state"`
	CreateTime  OftenTime `json:"create_time" xorm:"create_time"`
	Deleted     OftenTime `json:"deleted" xorm:"deleted"`
}

type ReturnGradeInformation struct {
	ID          int64  `json:"id" xorm:"id"`
	Username    string `json:"username" xorm:"username" `
	Contest     string `json:"contest" xorm:"contest"`
	Grade       string `json:"grade" xorm:"grade"`
	Certificate string `json:"certificate" xorm:"certificate"`
	State       int    `json:"state" xorm:"state"`
	CreateTime  string `json:"create_time" xorm:"create_time"`
}

func (GradeInformation) TableName() string {
	return "grade"
}
