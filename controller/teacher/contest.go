package teacher

import (
	"github.com/gin-gonic/gin"
)

type ContestController struct{}

// RegisterRoutes
// 注册学生用户路由
func (self ContestController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/competition_view", self.ViewContest) // 查看竞赛信息
	g.GET("/grade_view", self.ViewGrade)         // 成绩查询
}

// ViewCompetition
func (ContestController) ViewContest(c *gin.Context) {

}

// ViewGrade
func (ContestController) ViewGrade(c *gin.Context) {

}
