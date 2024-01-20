package models

import "time"

type UserParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID       int64     `xorm:"id"`
	Username string    `xorm:"username"`
	Password string    `xorm:"password"`
	Deleted  time.Time `xorm:"deleted"`
}

type UpdateUserInfo UserInfo
type DisplayUserForm UserInfo

type UserDeleteId struct {
	ID []int `json:"id_number"`
}
