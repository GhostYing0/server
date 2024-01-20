package teacher

import (
	"github.com/gin-gonic/gin"
)

type UserController struct{}

// RegisterRoutes
func (self UserController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/personal_center", self.GetUserInformation) // 查看指导教师个人信息

	g.POST("/update_personal_information", self.UpdateUserInformation) // 更新指导教师个人信息
}

// GetUserInformation
func (UserController) GetUserInformation(c *gin.Context) {

}

// UpdateUserInformation
func (UserController) UpdateUserInformation(c *gin.Context) {

}
