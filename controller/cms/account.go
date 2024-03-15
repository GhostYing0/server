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
	g.GET("/get_info", self.GetInfo)            // 从Redis里读token
}

// Login
func (AccountController) Login(c *gin.Context) {
	appG := app.Gin{C: c}
	var Param models.LoginForm
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
	var Param models.RegisterForm
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
	var Param models.UpdatePasswordForm
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

func (AccountController) GetInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	data := make(map[string]interface{})
	token := c.Query("token")

	id, username, role, err := logic.DefaultCmsAccount.GetInfo(token)
	if err != nil {
		appG.ResponseErr(err.Error())
	}

	data["id"] = id
	data["username"] = username
	data["role"] = role

	appG.ResponseSucMsg(data)
}
