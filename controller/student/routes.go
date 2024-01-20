package student

import "github.com/gin-gonic/gin"

// RegisterRoutes
// 注册学生用户路由
func RegisterRoutes(g *gin.RouterGroup) {
	new(AccountController).RegisterRoutes(g)
	new(UserController).RegisterRoutes(g)
	new(ContestController).RegisterRoutes(g)
}
