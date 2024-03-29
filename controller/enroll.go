package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"server/logic"
	. "server/logic"
	"server/models"
	"server/utils/app"
	. "server/utils/mydebug"
)

type EnrollController struct{}

func (self EnrollController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/viewContest", self.ViewContest)                // 查看竞赛信息
	g.POST("/enrollContest", self.EnrollContest)           // 报名竞赛
	g.GET("/searchEnrollResult", self.DisplayEnrollResult) // 查看报名结果
	g.POST("/processEnroll", self.ProcessEnroll)           // 审核竞赛
}

// Viewcontest
// 查看竞赛
func (EnrollController) ViewContest(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	var ret string

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	if limit < 0 || curPage < 0 {
		DPrintf("分页器参数错误")
		appG.ResponseErr("分页器参数错误")
		return
	}

	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultEnrollLogic.DisplayContest(paginator)
	if err != nil {
		DPrintf(" logic.DefaultEnrollLogic.DisplayContest 错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	paginator.SetTotalPage(total)

	if list != nil {
		data["list"] = list
		data["total"] = total
		data["page_size"] = limit
		data["page_number"] = curPage
		data["total_page"] = paginator.GetTotalPage()
	}

	appG.ResponseSucMsg(data, ret)
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

	err = logic.DefaultEnrollLogic.InsertEnrollInformation(param.UserName, param.TeamID, param.ContestID, param.CreateTime, param.School, param.Phone, param.Email)
	if err != nil {
		DPrintf("EnrollContest logic.DefaultEnrollLogic.InsertEnrollInformation() 发生错误:", err)
		appG.ResponseErr("报名失败", err.Error())
		return
	}

	appG.ResponseSuc("报名成功")
}

// SearchEnrollResult
// 查看报名结果
func (EnrollController) DisplayEnrollResult(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	username := c.DefaultQuery("username", "")
	userID := int64(com.StrTo(c.DefaultQuery("user_id", "0")).MustInt())
	contest := c.DefaultQuery("contest", "")
	startTime := c.DefaultQuery("startTime", "")
	endTime := c.DefaultQuery("endTime", "")
	school := c.DefaultQuery("school", "")
	phone := c.DefaultQuery("phone", "")
	email := c.DefaultQuery("email", "")
	state := com.StrTo(c.DefaultQuery("state", "-1")).MustInt()

	if limit < 1 || curPage < 1 {
		DPrintf("DisplayEnrollResult 查询表容量和页码应大于0")
		appG.ResponseErr("查询表容量和页码应大于0")
		return
	}

	data := make(map[string]interface{})

	paginator := logic.NewPaginator(curPage, limit)

	user_id, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("user_id不存在")
		return
	}
	role, exist := c.Get("role")
	if !exist {
		appG.ResponseErr("role不存在")
		return
	}

	list, total, err := logic.DefaultEnrollLogic.Search(paginator, username, userID, contest, startTime, endTime, school, phone, email, state, user_id.(int64), role.(int))

	if err != nil {
		DPrintf("DisplayEnrollResult 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	paginator.SetTotalPage(total)

	data["list"] = list
	data["pageNumber"] = curPage
	data["perSize"] = limit
	data["total"] = total
	data["totalPage"] = paginator.GetTotalPage()

	appG.ResponseSucMsg(data, "查询成功")
}

// ProcessEnroll
// 审核竞赛
func (EnrollController) ProcessEnroll(c *gin.Context) {

}
