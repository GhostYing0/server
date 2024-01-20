package cms

import (
	"fmt"
	. "server/database"
	. "server/logic"
	"server/models"
)

type CmsUserLogic struct{}

var DefaultCmsUser = CmsUserLogic{}

func (self CmsUserLogic) Display(paginator *Paginator) (*[]models.DisplayUserForm, int64, error) {
	tx := MasterDB.NewSession()
	var total int64
	var err error

	var List []models.DisplayUserForm

	if paginator.PerPage() > 0 {
		total, err = tx.Table("account").Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(&List)
		if err != nil {
			tx.Rollback()
			return nil, 0, err
		}
	} else {
		total, err = tx.Table("account").Limit(10, 10*(paginator.CurPage()-1)).FindAndCount(&List)
		if err != nil {
			tx.Rollback()
			return nil, 0, err
		}
	}

	tx.Commit()
	return &List, total, err
}

func (self CmsUserLogic) AddUser(username string, password string) (string, error) {
	if len(username) == 0 || len(password) == 0 {
		return "账号和密码不能为空 ", nil
	}

	tx := MasterDB.NewSession()

	exist, err := tx.Table("account").Where("username = ?", username).Exist()
	if exist {
		return "用户已存在", err
	}

	param := &models.UserInfo{
		Username: username,
		Password: password,
	}

	_, err = tx.Table("account").Insert(param)
	if err != nil {
		tx.Rollback()
		return "操作出错", err
	}

	tx.Commit()
	return "操作成功", err
}

func (self CmsUserLogic) UpdateUser(ID int64, NewUsername string, NewPassword string) (string, error) {
	tx := MasterDB.NewSession()

	exist, err := tx.Table("account").Where("id = ?", ID).Exist()
	if !exist {
		return "用户不存在", err
	}
	if err != err {
		tx.Rollback()
		return "操作错误", err
	}

	param := &models.UserParam{
		Username: NewUsername,
		Password: NewPassword,
	}

	if len(NewUsername) > 0 {
		_, err = tx.Table("account").Where("id = ?", ID).Cols("username").Update(param)
		if err != nil {
			tx.Rollback()
			return "更新出错", err
		}
	}
	if len(NewPassword) > 0 {
		_, err = tx.Table("account").Where("id = ?", ID).Cols("password").Update(param)
		if err != nil {
			tx.Rollback()
			return "更新出错", err
		}
	}

	tx.Commit()
	return "操作成功", err
}

func (self CmsUserLogic) DeleteUser(ids *[]int) (string, error, int64) {
	tx := MasterDB.NewSession()
	var count int64

	for _, id := range *ids {
		var contest models.ContestInfo
		if id < 1 {
			fmt.Println("非法id")
			continue
		}
		contest.ID = int64(id)
		affected, err := tx.Table("account").Delete(&contest)
		if err != nil {
			tx.Rollback()
			return "操作出错", err, 0
		}
		if affected > 0 {
			count += affected
		}
	}

	tx.Commit()
	return "操作成功", nil, count
}
