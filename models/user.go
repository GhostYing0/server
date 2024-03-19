package models

type UserParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int64     `xorm:"id"`
	Username string    `xorm:"username"`
	Password string    `xorm:"password"`
	Role     int       `xorm:"role"`
	Deleted  OftenTime `xorm:"deleted"`
}

type UserRedis struct {
	ID       int64  `json:"id"`
	Role     int    `json:"role"`
	Username string `json:"username"`
}

type UpdateUserInfo User

type UserDeleteId struct {
	ID []int64 `json:"ids"`
}

func (User) TableName() string {
	return "account"
}
