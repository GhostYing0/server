package cms

import (
	"errors"
	"fmt"
	. "server/database"
	. "server/logic"
	"server/models"
	. "server/utils/mydebug"
)

type CmsUserLogic struct{}

var DefaultCmsUser = CmsUserLogic{}

func (self CmsUserLogic) Display(paginator *Paginator, mode int, username string) (*[]models.User, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser Display session.Begin() 发生错误:", err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser Display session.Close() 发生错误:", err)
		}
	}()

	if mode == 1 {
		session.Where("role = 1")
	} else if mode == 2 {
		session.Where("role = 2")
	}

	if username != "" {
		session.Where("username = ?", username)
	}

	data := &[]models.User{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return data, 0, fail
		}
		DPrintf("Search 查找成绩信息失败:", err)
		return data, 0, err
	}

	return data, total, session.Rollback()
}

func (self CmsUserLogic) AddUser(username string, password string, role int) (string, error) {
	if len(username) == 0 || len(password) == 0 {
		return "账号和密码不能为空 ", nil
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser Display session.Begin() 发生错误:", err)
		return err.Error(), err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser Display session.Close() 发生错误:", err)
		}
	}()

	exist, err := session.Table("account").Where("username = ? and role = ?", username, role).Exist()
	if exist {
		return "用户已存在", err
	}

	param := &models.User{
		Username: username,
		Password: password,
		Role:     role,
	}

	_, err = session.Insert(param)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return err.Error(), fail
		}
		return "操作出错", err
	}

	return "操作成功", session.Commit()
}

func (self CmsUserLogic) UpdateUser(ID int64, newUsername string, newPassword string, role int) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser UpdateUser session.Begin() 发生错误:", err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser UpdateUser session.Close() 发生错误:", err)
		}
	}()

	exist, err := session.Table("account").Where("id = ?", ID).Exist()
	if !exist {
		return errors.New("用户不存在")
	}
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		DPrintf("UpdateUser 查询用户失败:", err)
		return err
	}

	param := &models.User{
		Username: newUsername,
		Password: newPassword,
		Role:     role,
	}

	_, err = session.Where("id = ?", ID).Update(param)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		DPrintf("UpdateUser Update 更新用户失败:", err)
		return err
	}

	return session.Commit()
}

func (self CmsUserLogic) DeleteUser(ids *[]int64) (int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser DeleteUser session.Begin() 发生错误:", err)
		return 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser DeleteUser session.Close() 发生错误:", err)
		}
	}()

	var count int64

	for _, id := range *ids {
		var contest models.ContestInfo
		if id < 1 {
			fmt.Println("非法id")
			continue
		}
		contest.ID = int64(id)
		affected, err := session.Table("account").Delete(&contest)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return 0, fail
			}
			return 0, err
		}
		if affected > 0 {
			count += affected
		}
	}

	return count, session.Commit()
}

func (self CmsUserLogic) GetUserCount() (int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("GetUserCount session.Begin() 发生错误:", err)
		return 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("GetUserCount session.Close() 发生错误:", err)
		}
	}()

	count, err := session.Table("account").Count()
	if err != nil {
		DPrintf("GetUserCount Count 发生错误:", err)
		return count, err
	}
	return count, err
}
