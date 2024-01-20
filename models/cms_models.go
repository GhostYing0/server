package models

type LoginParam struct {
	Username string `json:"username" xorm:"username"`
	Password string `json:"password" xorm:"password"`
	Role     int    `json:"role" xorm:"role"`
}

type RegisterParam struct {
	Username        string `json:"username" xorm:"username"`
	Password        string `json:"password" xorm:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Role            int    `json:"role" xorm:"role"`
}

type UpdatePasswordParam struct {
	Username        string `json:"username" xorm:"username"`
	NewPassword     string `json:"new_password" xorm:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginReturn struct {
	ID       int64  `xorm:"id"`
	Password string `xorm:"password"`
}

type NewAccount LoginParam
