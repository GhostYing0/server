package controller

import (
	"github.com/gin-gonic/gin"
	"server/logic"
	"server/models"
	"server/utils/app"
	"server/utils/logging"
	. "server/utils/mydebug"
)

type AccountController struct{}

// 有想法改为学生和老师用一套登陆注册
// 就是分为普通用户和管理员
// RegisterRoutes
func (self AccountController) RegisterRoutes(g *gin.RouterGroup) {
	g.POST("/login", self.Login)                  // 普通用户登录
	g.POST("/register", self.Register)            // 普通用户注册
	g.POST("/update_passwd", self.UpdatePasswd)   // 普通用户修改密码
	g.GET("/profileStudent", self.ProfileStudent) //学生获取个人信息
	g.GET("/profileTeacher", self.ProfileTeacher) //学生获取个人信息
	g.POST("/updateProfile", self.UpdateProfile)  //更新个人信息
}

// Login
func (AccountController) Login(c *gin.Context) {
	appG := app.Gin{C: c}
	param := &models.LoginForm{}
	data := make(map[string]interface{})

	err := c.ShouldBindJSON(param)

	if err != nil {
		DPrintf("Login c.ShouldBindJSON()发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	token, err := logic.DefaultUserAccount.Login(param.Username, param.Password, param.Role)
	if err != nil {
		DPrintf("Login 登录失败:", err)
		appG.ResponseErr(err.Error())
		return
	}

	data["token"] = token
	appG.ResponseSucMsg(data, "登陆成功")
}

// Register
func (AccountController) Register(c *gin.Context) {
	appG := app.Gin{C: c}
	form := &models.RegisterForm{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		DPrintf("Register 注册失败:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultUserAccount.Register(
		form.Username, form.Password, form.ConfirmPassword, form.Role,
		form.Name, form.Gender, form.Semester, form.College, form.School, form.Class, form.Phone, form.Email)
	if err != nil {
		DPrintf("Register 注册发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}
	appG.ResponseSuc("注册成功")
}

// UpdatePasswd
func (AccountController) UpdatePasswd(c *gin.Context) {
	appG := app.Gin{C: c}
	param := &models.UpdatePasswordForm{}

	err := c.ShouldBindJSON(&param)
	if err != nil {
		DPrintf("UpdatePasswd c.ShouldBindJSON()发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultUserAccount.UpdatePassword(param.Username, param.NewPassword, param.ConfirmPassword, param.Role)
	if err != nil {
		DPrintf("UpdatePasswd 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("修改密码成功")
}

func (AccountController) ProfileStudent(c *gin.Context) {
	app := app.Gin{C: c}

	userID, exist := c.Get("user_id")
	if !exist {
		app.ResponseErr("用户查找错误")
		logging.L.Error("用户查找错误")
		return
	}

	data, err := logic.DefaultUserAccount.GetProfileStudent(userID.(int64))
	if err != nil {
		app.ResponseErr(err.Error())
		return
	}
	app.ResponseSucMsg(data)
}

func (AccountController) ProfileTeacher(c *gin.Context) {
	app := app.Gin{C: c}

	userID, exist := c.Get("user_id")
	if !exist {
		app.ResponseErr("用户查找错误")
		logging.L.Error("用户查找错误")
		return
	}

	data, err := logic.DefaultUserAccount.GetProfileTeacher(userID.(int64))
	if err != nil {
		app.ResponseErr(err.Error())
		return
	}
	app.ResponseSucMsg(data)
}

func (AccountController) UpdateProfile(c *gin.Context) {

}
