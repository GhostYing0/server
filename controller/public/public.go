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
