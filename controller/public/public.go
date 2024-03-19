package public

import (
	"github.com/gin-gonic/gin"
	logic "server/logic/public"
	"server/utils/app"
	. "server/utils/mydebug"
)

type PublicController struct{}

func (self PublicController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/get_info", self.GetInfo) // 从Redis里读token
	g.GET("/logout", self.Logout)
}

func (PublicController) GetInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	data := make(map[string]interface{})
	token := c.Query("token")

	id, username, role, err := logic.DefaultPublic.GetInfo(token)
	if err != nil {
		DPrintf("User GetInfo 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	data["id"] = id
	data["username"] = username
	data["role"] = role

	appG.ResponseSucMsg(data)
}

// Logout
func (PublicController) Logout(c *gin.Context) {
	appG := app.Gin{C: c}

	token := c.Query("token")

	if len(token) <= 0 {
		DPrintf("token为空")
		appG.ResponseErr("token为空")
		return
	}

	err := logic.DefaultPublic.Logout(token)
	if err != nil {
		DPrintf("登出发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg("登出成功")
}
