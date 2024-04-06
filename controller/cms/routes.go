package cms

import "github.com/gin-gonic/gin"

func RegisterRoutes(g *gin.RouterGroup) {
	new(AccountController).RegisterRoutes(g)
	new(UserController).RegisterRoutes(g)
	new(ContestController).RegisterRoutes(g)
	new(RegistrationController).RegisterRoutes(g)
	new(GradeController).RegisterRoutes(g)
}
