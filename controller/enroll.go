package controller

import (
	"github.com/gin-gonic/gin"
	"server/logic"
	"server/models"
	"server/utils/app"
	. "server/utils/mydebug"
)

type EnrollController struct{}

func (self EnrollController) RegisterRoutes(g *gin.RouterGroup) {
	g.POST("/enrollContest", self.EnrollContest)           // 报名竞赛
	g.POST("/searchEnrollResult", self.SearchEnrollResult) // 查看报名结果
	g.POST("/processEnroll", self.ProcessEnroll)           // 审核竞赛
}

// EnrollContest
// 报名竞赛
func (EnrollController) EnrollContest(c *gin.Context) {
	appG := app.Gin{C: c}

	param := &models.EnrollForm{}

	err := c.ShouldBindJSON(param)
	if err != nil {
		DPrintf("EnrollContest c.ShouldBindJSON() 发生错误:", err)
		appG.ResponseErr("报名失败", err.Error())
		return
	}

	err = logic.DefaultEnrollLogic.InsertEnrollInformation(param.Members, param.Contest, param.CreateTime, param.School, param.Phone, param.Email)
	if err != nil {
		DPrintf("EnrollContest logic.DefaultEnrollLogic.InsertEnrollInformation() 发生错误:", err)
		appG.ResponseErr("报名失败 ", err.Error())
		return
	}

	appG.ResponseSuc("报名成功")
}

// SearchEnrollResult
// 查看报名结果
func (EnrollController) SearchEnrollResult(c *gin.Context) {

}

// ProcessEnroll
// 审核竞赛
func (EnrollController) ProcessEnroll(c *gin.Context) {

}
