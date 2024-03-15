package controller

import "github.com/gin-gonic/gin"

// RegisterRoutes
// 注册普通用户路由
func RegisterRoutes(g *gin.RouterGroup) {
	new(AccountController).RegisterRoutes(g)
	new(EnrollController).RegisterRoutes(g)
}
