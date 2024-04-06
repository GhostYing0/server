package logic

import (
	"errors"
	"github.com/unknwon/com"
	. "server/database"
	"server/models"
	"server/utils/gredis"
	. "server/utils/mydebug"
	"server/utils/util"
	"time"
)

type UserAccountLogic struct{}

var DefaultUserAccount = UserAccountLogic{}

func (self UserAccountLogic) Login(username string, password string, role int) (string, error) {
	if username == "" || password == "" {
		DPrintf("用户名和密码不能为空")
		return "", errors.New("用户名和密码不能为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("InsertEnrollInformation session.Begin() 发生错误:", err)
		return "", err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("InsertEnrollInformation session.Close() 发生错误:", err)
		}
	}()

	account := &models.Account{}
	has, err := session.Where("username = ?", username).And("role = ?", role).Get(account)
	if err != nil {
		DPrintf("Login 查询用户失败:", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return "", fail
		}
		return "", err
	}
	if !has {
		DPrintf("Login 用户不存在:")
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return "", fail
		}
		return "", errors.New("用户不存在")
	}

	if username != account.Username || util.EncodeMD5(password) != account.Password {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return "", fail
		}
		return "", errors.New("用户名或密码错误")
	}

	token, err := util.GenerateToken(com.ToStr(account.ID), account.Username, account.Role)
	if err != nil {
		DPrintf("生成token失败")
		return "", err
	}

	err = gredis.Set(token, 1, time.Minute*60)
	if err != nil {
		DPrintf("login Set 失败:", err)
		return "", err
	}

	return token, session.Commit()
}

func (self UserAccountLogic) Register(username string, password string, confirmPassword string, role int) error {
	if username == "" || password == "" {
		return errors.New("用户名或密码不能为空")
	}
	if password != confirmPassword {
		return errors.New("密码两次输入不一致")
	}
	if role != 1 && role != 2 {
		return errors.New("用户不是学生或老师")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Register session.Begin() 发生错误:", err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Register session.Close() 发生错误:", err)
		}
	}()

	newAccount := models.LoginForm{
		Username: username,
		Password: password,
		Role:     role,
	}

	has, err := session.Table("account").Where("username = ?", newAccount.Username).And("role = ?", newAccount.Role).Exist()
	if err != nil {
		DPrintf("Register 查询用户发生错误:", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}
	if has {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return errors.New("用户已存在")
	}

	newAccount.Password = util.EncodeMD5(newAccount.Password)

	_, err = session.Insert(newAccount)
	if err != nil {
		DPrintf("Login 添加新用户失败:", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}

	return session.Commit()
}

func (self UserAccountLogic) UpdatePassword(username string, newPassword string, confirmPassword string, role int) error {
	if username == "" || newPassword == "" {
		return errors.New("用户名或密码不能为空")
	}
	if newPassword != confirmPassword {
		return errors.New("密码两次输入不一致")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Register session.Begin() 发生错误:", err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Register session.Close() 发生错误:", err)
		}
	}()

	has, err := session.Table("account").Where("username = ?", username).And("role = ?", role).Exist()
	if err != nil {
		DPrintf("UpdatePassword 查找用户发生错误:", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}
	if !has {
		DPrintf("用户不存在")
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return errors.New("用户不存在")
	}

	_, err = session.Table("account").Where("username = ?", username).Cols("password").Update(&models.Account{Password: newPassword})
	if err != nil {
		DPrintf("UpdatePassword 修改密码发生错误:", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}

	return session.Commit()
}
