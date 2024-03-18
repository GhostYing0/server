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
	Deleted  OftenTime `json:"deleted" xorm:"deleted"`
}

type NewAccount LoginForm

func (LoginForm) TableName() string {
	return "account"
}

func (Account) TableName() string {
	return "account"
}
