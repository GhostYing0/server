package cms

import (
	"fmt"
	"github.com/unknwon/com"
	. "server/database"
	"server/models"
	"server/utils/gredis"
	"server/utils/logging"
	"server/utils/util"
	"time"
)

type CmsAccountLogic struct{}

var DefaultCmsAccount = CmsAccountLogic{}

func (self CmsAccountLogic) Login(param *models.LoginForm) (string, string, error) {
	if param.Username == "" || param.Password == "" {
		return "用户名或密码不能为空", "", nil
	}

	if param.Role != 0 {
		return "", "用户不是管理员", nil
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

	if util.EncodeMD5(param.Password) != loginReturn.Password {
		return "密码错误", "", err
	}

	token, err := util.GenerateToken(com.ToStr(loginReturn.ID), param.Username, param.Role)
	if err != nil {
		return "token创建出错", token, err
	}

	err = gredis.Set(token, 1, time.Minute*60)
	if err != nil {
		fmt.Println("login Set err:", err)
	}

	_, err = tx.Where("id = ?", loginReturn.ID).Update(models.ManagerInfo{LastLoginTime: models.NewOftenTime()})
	if err != nil {
		logging.L.Error(err)
		return "登录失败", "", err
	}
	tx.Commit()
	return "登陆成功", token, err
}

func (self CmsAccountLogic) Register(param *models.RegisterForm) (string, error) {
	if param.Username == "" || param.Password == "" {
		return "用户名或密码不能为空", nil
	}

	if param.Password != param.ConfirmPassword {
		return "密码两次输入不一致", nil
	}

	if param.Role != 0 {
		return "用户不为管理员", nil
	}

	newAccount := models.LoginForm{
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

	newAccount.Password = util.EncodeMD5(newAccount.Password)

	_, err = tx.Table("cms_account").Insert(newAccount)
	if err != nil {
		fmt.Println("Register Insert error:", err)
		tx.Rollback()
		return "操作错误", err
	}

	tx.Commit()
	return "操作成功", err
}

func (self CmsAccountLogic) UpdatePassword(param *models.UpdatePasswordForm) (string, error) {
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

//func (self CmsAccountLogic) GetInfo(token string) (int64, string, int, error) {
//	tokenHasExpired, err := gredis.Get(token)
//	if err != nil {
//		fmt.Println("")
//		return 0, "", 0, err
//	}
//	if tokenHasExpired != "1" {
//		return 0, "", 0, errors.New("token已过期, 请重新登录")
//	}
//
//	claims, err := util.ParseToken(token)
//	if err != nil {
//		fmt.Println("")
//		return 0, "", 0, err
//	}
//
//	id, err := strconv.Atoi(claims.ID)
//	if err != nil {
//		fmt.Println("")
//		return 0, "", 0, err
//	}
//
//	return int64(id), claims.Username, claims.Role, err
//}
