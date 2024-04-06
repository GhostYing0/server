package models

type LoginForm struct {
	Username string `json:"username" xorm:"username"`
	Password string `json:"password" xorm:"password"`
	Role     int    `json:"role" xorm:"role"`
}

type RegisterForm struct {
	Username        string `json:"username" xorm:"username"`
	Password        string `json:"password" xorm:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Role            int    `json:"role" xorm:"role"`
	Name            string `json:"name" xorm:"name"`
	Gender          string `json:"gender" xorm:"gender"`
	School          string `json:"school" xorm:"school"`
	College         string `json:"college" xorm:"college"`
	Semester        string `json:"semester" xorm:"semester"`
	Class           string `json:"class" xorm:"class"`
}

type UpdatePasswordForm struct {
	Username        string `json:"username" xorm:"username"`
	NewPassword     string `json:"new_password" xorm:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Role            int    `json:"role" xorm:"role"`
}

type LoginReturn struct {
	ID       int64  `xorm:"id"`
	Password string `xorm:"password"`
}

type Account struct {
	ID       int64     `json:"id" xorm:"id"`
	Username string    `json:"username" xorm:"username"`
	Password string    `json:"password" xorm:"password"`
	Role     int       `json:"role" xorm:"role"`
	UserID   string    `json:"user_id" xorm:"user_id"`
	Deleted  OftenTime `json:"deleted" xorm:"deleted"`
}

type NewAccount LoginForm

func (LoginForm) TableName() string {
	return "account"
}

func (Account) TableName() string {
	return "account"
}
