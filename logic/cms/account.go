package cms

import (
	"fmt"
	"github.com/unknwon/com"
	. "server/database"
	"server/models"
	"server/utils/util"
)

type CmsAccountLogic struct{}

var DefaultCmsAccount = CmsAccountLogic{}

func (self CmsAccountLogic) Login(param *models.LoginParam) (string, string, error) {
	if param.Username == "" || param.Password == "" {
		return "用户名或密码不能为空", "", nil
	}

	if param.Role != "cms" {
		return "", "无效角色", nil
	}

	tx := MasterDB.NewSession()

	var loginReturn models.LoginReturn

	has, err := tx.Table("cms_account").Where("username = ?", param.Username).Get(&loginReturn)
	if err != nil {
		fmt.Println("Login Get error:", err)
		return "操作错误", "", err
	}
	if !has {
		return "用户不存在", "", err
	}

	if param.Password != loginReturn.Password {
		return "密码错误", "", err
	}

	token, err := util.GenerateToken(com.ToStr(loginReturn.ID), param.Username, param.Role)
	if err != nil {
		return "token创建出错", token, err
	}

	tx.Commit()
	return "登陆成功", token, err
}

func (self CmsAccountLogic) Register(param *models.RegisterParam) (string, error) {
	if param.Username == "" || param.Password == "" {
		return "用户名或密码不能为空", nil
	}

	if param.Password != param.ConfirmPassword {
		return "密码两次输入不一致", nil
	}

	if param.Role != "cms" {
		return "无效角色", nil
	}

	newAccount := models.NewAccount{
		Username: param.Username,
		Password: param.Password,
		Role:     param.Role,
	}

	tx := MasterDB.NewSession()

	has, err := tx.Table("cms_account").Where("username = ?", newAccount.Username).Exist()
	if err != nil {
		fmt.Println("Register Exist error:", err)
	}
	if has {
		return "用户已存在", err
	}

	_, err = tx.Table("cms_account").Insert(newAccount)
	if err != nil {
		fmt.Println("Register Insert error:", err)
		tx.Rollback()
		return "操作错误", err
	}

	tx.Commit()
	return "操作成功", err
}

func (self CmsAccountLogic) UpdatePassword(param *models.UpdatePasswordParam) (string, error) {
	if param.Username == "" || param.NewPassword == "" {
		return "用户名或密码不能为空", nil
	}

	if param.NewPassword != param.ConfirmPassword {
		return "密码两次输入不一致", nil
	}

	tx := MasterDB.NewSession()
	has, err := tx.Table("cms_account").Where("username = ?", param.Username).Exist()
	if !has {
		fmt.Println("UpdatePassword Exist error:", err)
		return "用户不存在", err
	}

	_, err = tx.Table("cms_account").Where("username = ?", param.Username).Cols("password").Update(param)
	if err != nil {
		fmt.Println("UpdatePassword Update error:", err)
		tx.Rollback()
		return "操作错误", err
	}

	tx.Commit()
	return "操作成功", err
}
