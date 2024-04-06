package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"server/logic"
	. "server/logic"
	"server/models"
	"server/utils/app"
	. "server/utils/mydebug"
	"strconv"
)

type EnrollController struct{}

func (self EnrollController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/viewContest", self.ViewContest)                // 查看竞赛信息
	g.POST("/enrollContest", self.EnrollContest)           // 报名竞赛
	g.GET("/searchEnrollResult", self.DisplayEnrollResult) // 用户查看报名结果
	//g.GET("/teacherSearchEnroll", self.TeacherSearchEnroll)    // 教师查看报名结果
	g.POST("/processPassEnroll", self.ProcessPassEnroll)       // 教师审核通过
	g.POST("/processRejectEnroll", self.ProcessRejectEnroll)   // 教师审核驳回
	g.POST("/processRecoverEnroll", self.ProcessRecoverEnroll) // 教师审核恢复
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

	err = logic.DefaultEnrollLogic.InsertEnrollInformation(param.UserName, param.TeamID, param.ContestName, param.CreateTime, param.School, param.Phone, param.Email)
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

	contest := c.DefaultQuery("contest_name", "")
	startTime := c.DefaultQuery("start_time", "")
	endTime := c.DefaultQuery("end_time", "")
	state := com.StrTo(c.DefaultQuery("state", "-1")).MustInt()

	fmt.Println(startTime)
	fmt.Println(endTime)

	key, exist := appG.C.Get("user_id")
	if !exist {
		DPrintf("无Token")
		appG.ResponseErr("无Token")
		return
	}
	userID, ok := key.(int64)
	if !ok {
		DPrintf("DisplayEnrollResult userID类型错误")
		appG.ResponseErr("userID类型错误")
		return
	}

	if limit < 1 || curPage < 1 {
		DPrintf("DisplayEnrollResult 查询表容量和页码应大于0")
		appG.ResponseErr("查询表容量和页码应大于0")
		return
	}

	data := make(map[string]interface{})

	paginator := logic.NewPaginator(curPage, limit)

	list, total, err := logic.DefaultEnrollLogic.Search(paginator, userID, contest, startTime, endTime, state)

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

// ProcessPassEnroll
// 审核通过
func (EnrollController) ProcessPassEnroll(c *gin.Context) {
	appG := app.Gin{C: c}

	role, exist := appG.C.Get("role")
	if !exist {
		DPrintf("无效身份")
		appG.ResponseErr("无效身份")
		return
	}
	if role != 2 {
		DPrintf("无教师权限")
		appG.ResponseErr("无教师权限")
		return
	}

	request := models.EnrollIds{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultEnrollLogic.ProcessEnroll(&request.ID, 1)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(strconv.Itoa(int(count)))
}

// ProcessRejectEnroll
// 审核驳回
func (EnrollController) ProcessRejectEnroll(c *gin.Context) {
	appG := app.Gin{C: c}

	role, exist := appG.C.Get("role")
	if !exist {
		DPrintf("无效身份")
		appG.ResponseErr("无效身份")
		return
	}
	if role != 2 {
		DPrintf("无教师权限")
		appG.ResponseErr("无教师权限")
		return
	}

	request := models.EnrollIds{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultEnrollLogic.ProcessEnroll(&request.ID, 2)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(strconv.Itoa(int(count)))
}

// ProcessRecoverEnroll
// 审核恢复
func (EnrollController) ProcessRecoverEnroll(c *gin.Context) {
	appG := app.Gin{C: c}

	role, exist := appG.C.Get("role")
	if !exist {
		DPrintf("无效身份")
		appG.ResponseErr("无效身份")
		return
	}
	if role != 2 {
		DPrintf("无教师权限")
		appG.ResponseErr("无教师权限")
		return
	}

	request := models.EnrollIds{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultEnrollLogic.ProcessEnroll(&request.ID, 3)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(strconv.Itoa(int(count)))
}
