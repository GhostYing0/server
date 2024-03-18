package models

type GradeForm struct {
	Username    string `json:"username"`
	Contest     string `json:"contest"`
	Grade       string `json:"grade"`
	Certificate string `json:"certificate"`
	CreateTime  string `json:"create_time"`
}

type GradeInformation struct {
	UserID      int64     `json:"user_id" xorm:"user_id"`
	Username    string    `json:"username" xorm:"username" `
	Contest     string    `json:"contest" xorm:"contest"`
	Grade       string    `json:"grade" xorm:"grade"`
	Certificate string    `json:"certificate" xorm:"certificate"`
	State       int       `json:"state" xorm:"state"`
	CreateTime  OftenTime `json:"create_time" xorm:"create_time"`
	Deleted     OftenTime `json:"deleted" xorm:"deleted"`
}

func (GradeInformation) TableName() string {
	return "grade"
}
