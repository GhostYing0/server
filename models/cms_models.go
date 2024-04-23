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
	Phone           string `json:"phone" xorm:"phone"`
	Email           string `json:"email" xorm:"email"`
	Class           string `json:"class" xorm:"class"`
}

type UpdatePasswordForm struct {
	Username        string `json:"username" xorm:"username"`
	NewPassword     string `json:"new_password" xorm:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Role            int    `json:"role" xorm:"role"`
}

type UpdateUserPassword struct {
	Password        string `json:"password" xorm:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginReturn struct {
	ID       int64  `xorm:"id"`
	Password string `xorm:"password"`
}

type Account struct {
	ID            int64     `json:"id" xorm:"id"`
	Username      string    `json:"username" xorm:"username"`
	Password      string    `json:"password" xorm:"password"`
	Role          int       `json:"role" xorm:"role"`
	UserID        string    `json:"user_id" xorm:"user_id"`
	Phone         string    `json:"phone" xorm:"phone"`
	Email         string    `json:"email" xorm:"email"`
	CreateTime    OftenTime `json:"create_time" xorm:"create_time"`
	LastLoginTime OftenTime `json:"last_login_time" xorm:"last_login_time"`
	Deleted       OftenTime `json:"deleted" xorm:"deleted"`
}

type Avatar struct {
	Avatar string `json:"avatar" xorm:"avatar"`
}

type NewAccount LoginForm

func (LoginForm) TableName() string {
	return "account"
}

func (Account) TableName() string {
	return "account"
}

type Manager struct {
	ID            int64     `json:"id" xorm:"id"`
	Username      string    `json:"username" xorm:"username"`
	Password      string    `json:"password" xorm:"password"`
	Role          int       `json:"role" xorm:"role"`
	CreateTime    string    `json:"create_time" xorm:"create_time"`
	LastLoginTime string    `json:"last_login_time" xorm:"last_login_time"`
	UpdateTime    string    `json:"update_time" xorm:"update_time"`
	Deleted       OftenTime `json:"deleted" xorm:"deleted"`
}

func (Manager) TableName() string {
	return "cms_account"
}

type NewManager struct {
	Username        string `json:"username" xorm:"username"`
	Password        string `json:"password" xorm:"password"`
	ConfirmPassword string `json:"confirm_password" xorm:"confirm_password"`
}
