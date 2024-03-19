package models

type EnrollForm struct {
	UserName   string `json:"username"`
	TeamID     int64  `json:"team_id"`
	ContestID  int64  `json:"contest_id"`
	CreateTime string `json:"create_time"`
	School     string `json:"school"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type EnrollInformation struct {
	ID         int64     `json:"id" xorm:"id"`
	Username   string    `json:"username" xorm:"username"`
	UserID     int64     `json:"user_id" xorm:"user_id"`
	TeamID     int64     `json:"team_id" xorm:"team_id"`
	Contest    string    `json:"contest" xorm:"contest"`
	CreateTime OftenTime `json:"create_time" xorm:"create_time"`
	School     string    `json:"school" xorm:"school"`
	Phone      string    `json:"phone" xorm:"phone"`
	Email      string    `json:"email" xorm:"email"`
	State      int       `json:"state" xorm:"state"`
	Deleted    OftenTime `json:"deleted" xorm:"deleted"`
}

func (EnrollInformation) TableName() string {
	return "enroll_information"
}

type EnrollDeleteId struct {
	ID []int64 `json:"id_number"`
}
