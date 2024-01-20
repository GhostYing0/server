package cms

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "server/database"
	logic "server/logic/cms"
	"server/models"
	"server/utils/app"
)

type AccountController struct{}

func (self AccountController) RegisterRoutes(g *gin.RouterGroup) {
	g.POST("/login", self.Login)                // 管理员登录
	g.POST("/register", self.Register)          // 管理员注册
	g.POST("/update_passwd", self.UpdatePasswd) // 管理员修改密码
}

// Login
func (AccountController) Login(c *gin.Context) {
	appG := app.Gin{C: c}
	var Param models.LoginParam
	var token string
	var err error
	var ret string
	data := make(map[string]interface{})

	err = c.ShouldBindJSON(&Param)

	if err != nil {
		fmt.Println("ShouldBindJSON error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	ret, token, err = logic.DefaultCmsAccount.Login(&Param)
	if err != nil {
		fmt.Println("logic.Login error:", err)
		appG.ResponseErr(ret, err.Error())
		return
	}

	data["token"] = token
	appG.ResponseSucMsg(data, ret)
}

// Register
func (AccountController) Register(c *gin.Context) {
	appG := app.Gin{C: c}
	var Param models.RegisterParam
	var err error
	var ret string

	err = c.ShouldBindJSON(&Param)

	if err != nil {
		fmt.Println("ShouldBindJSON error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	ret, err = logic.DefaultCmsAccount.Register(&Param)
	if err != nil {
		fmt.Println("logic.Register error:", err)
		appG.ResponseErr(ret, err.Error())
		return
	}
	appG.ResponseSuc(ret)
}

// UpdatePasswd
func (AccountController) UpdatePasswd(c *gin.Context) {
	appG := app.Gin{C: c}
	var Param models.UpdatePasswordParam
	var err error
	var ret string

	err = c.ShouldBindJSON(&Param)
	if err != nil {
		fmt.Println("ShouldBindJSON error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	ret, err = logic.DefaultCmsAccount.UpdatePassword(&Param)
	if err != nil {
		fmt.Println("logic.UpdatePassword error:", err)
		appG.ResponseErr(ret, err.Error())
		return
	}

	appG.ResponseSuc(ret)
}
