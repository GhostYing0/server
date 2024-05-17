package cms

import (
	"errors"
	"fmt"
	. "server/database"
	. "server/logic"
	"server/models"
	"server/utils/logging"
	. "server/utils/mydebug"
	"server/utils/util"
)

type CmsManagerLogic struct{}

var DefaultCmsManager = CmsManagerLogic{}

func (self CmsManagerLogic) DisplayManager(paginator *Paginator, username string) (*[]models.Manager, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		logging.L.Error("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser Display session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser Display session.Close() 发生错误:", err)
		}
	}()

	//查询管理员
	session.Table("cms_account")

	if username != "" {
		session.Table("cms_account").Where("username like ?", "%"+username+"%")
	}

	data := &[]models.Manager{}

	total, err := session.
		//Join("LEFT", "student", "account.user_id = student.student_id").
		Limit(paginator.PerPage(), paginator.Offset()).
		FindAndCount(data)
	if err != nil {
		DPrintf("Search 查找成绩信息失败:", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	for i := 0; i < len(*data); i++ {
		(*data)[i].CreateTime = models.MysqlFormatString2String((*data)[i].CreateTime)
		(*data)[i].UpdateTime = models.MysqlFormatString2String((*data)[i].UpdateTime)
		(*data)[i].LastLoginTime = models.MysqlFormatString2String((*data)[i].LastLoginTime)
	}

	return data, total, session.Commit()
}

func (self CmsManagerLogic) AddManager(username, password, confirmPassword string) error {
	if len(username) == 0 || len(password) == 0 {
		logging.L.Error("账号和密码不能为空")
		return nil
	}
	if password != confirmPassword {
		logging.L.Error("两次密码输入不相同")
		return nil
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser Display session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser Display session.Close() 发生错误:", err)
		}
	}()

	exist, err := session.Table("cms_account").Where("username = ?", username).Exist()
	if exist {
		logging.L.Error("用户已存在")
		return errors.New("用户已存在")
	}

	if err != nil {
		logging.L.Error(err)
		return err
	}

	newManager := &models.ManagerInfo{
		Username:   username,
		Password:   util.EncodeMD5(password),
		Role:       0,
		CreateTime: models.NewOftenTime(),
		UpdateTime: models.NewOftenTime(),
	}

	_, err = session.Insert(newManager)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(err)
			return fail
		}
		logging.L.Error(err)
		return err
	}

	return session.Commit()
}

func (self CmsManagerLogic) UpdateManager(id int64, username, password, confirmPassword string) error {
	session := MasterDB.NewSession()
	if password != confirmPassword {
		logging.L.Error("密码两次输入不同")
		return errors.New("密码两次输入不同")
	}
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser UpdateUser session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser UpdateUser session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	exist, err := session.Table("cms_account").Where("id = ?", id).Exist()
	if !exist {
		logging.L.Error("用户不存在")
		return errors.New("用户不存在")
	}
	if err != nil {
		logging.L.Error(err)
		DPrintf("UpdateUser 查询用户失败:", err)
		return err
	}

	exist, err = session.Table("cms_account").Where("username = ?", username).Exist()
	if exist {
		logging.L.Error("已有同名用户")
		return errors.New("已有同名用户")
	}
	if err != nil {
		logging.L.Error(err)
		return err
	}

	_, err = session.Where("id = ?", id).
		Update(&models.ManagerInfo{
			Username:   username,
			Password:   util.EncodeMD5(password),
			UpdateTime: models.NewOftenTime(),
		})
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(fail)
			return fail
		}
		logging.L.Error(err)
		DPrintf("UpdateUser Update 更新用户失败:", err)
		return err
	}

	return session.Commit()
}

func (self CmsManagerLogic) DeleteManager(ids *[]int64) (int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser DeleteManager session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser DeleteManager session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	var count int64

	for _, id := range *ids {
		account := &models.Manager{}
		if id < 1 {
			fmt.Println("非法id")
			continue
		}

		_, err := session.Where("id = ?", id).Get(account)
		if err != nil {
			logging.L.Error(err)
			return 0, err
		}

		affected, err := session.Delete(account)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				logging.L.Error(fail)
				return 0, fail
			}
			logging.L.Error(err)
			return 0, err
		}

		if affected > 0 {
			count += affected
		}
	}

	return count, session.Commit()
}

func (self CmsManagerLogic) GetManagerCount() (int64, error) {
	session := MasterDB.NewSession()

	count, err := session.Table("cms_account").Count()
	if err != nil {
		DPrintf("GetUserCount Count 发生错误:", err)
		logging.L.Error(err)
		return count, err
	}
	return count, err
}
