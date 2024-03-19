package public

import "github.com/gin-gonic/gin"

func RegisterRoutes(g *gin.RouterGroup) {
	new(PublicController).RegisterRoutes(g)
}
