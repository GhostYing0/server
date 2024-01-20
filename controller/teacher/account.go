package teacher

import (
	"github.com/gin-gonic/gin"
)

type AccountController struct{}

// RegisterRoutes
// 注册学生用户路由
func (self AccountController) RegisterRoutes(g *gin.RouterGroup) {
	g.POST("/login", self.Login)                // 指导教师用户登录
	g.POST("/register", self.Register)          // 指导教师用户注册
	g.POST("/update_passwd", self.UpdatePasswd) // 指导教师用户修改密码
}

// Login
func (AccountController) Login(c *gin.Context) {

}

// Register
func (AccountController) Register(c *gin.Context) {

}

// UpdatePasswd
func (AccountController) UpdatePasswd(c *gin.Context) {

}
