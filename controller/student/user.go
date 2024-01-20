package student

import (
	"github.com/gin-gonic/gin"
)

type UserController struct{}

// RegisterRoutes
func (self UserController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/personal_center", self.GetUserInformation) // 查看学生用户个人信息

	g.POST("/update_personal_information", self.UpdateUserInformation) // 更新学生用户个人信息
}

// GetUserInformation
func (UserController) GetUserInformation(c *gin.Context) {

}

// UpdateUserInformation
func (UserController) UpdateUserInformation(c *gin.Context) {

}
