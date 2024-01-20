package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
)

var MasterDB *xorm.Engine

type DataBase struct {
	Name     string
	Password string
	User     string
	Host     string
}

func Setup() {
	var err error
	MasterDB, err = xorm.NewEngine("mysql", "root:gu112233..@tcp(127.0.0.1:3306)/contest_manager?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("xorm.NewEngine error:", err)
	}
}
