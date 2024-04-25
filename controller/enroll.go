package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"server/logic"
	"server/models"
	"server/utils/app"
	. "server/utils/e"
	. "server/utils/mydebug"
)

type EnrollController struct{}

func (self EnrollController) RegisterRoutes(g *gin.RouterGroup) {
	g.POST("/enrollContest", self.EnrollContest)            // 学生上传报名信息
	g.GET("/searchEnrollResult", self.DisplayEnrollResult)  // 用户查看报名结果
	g.GET("/teacherSearchEnroll", self.TeacherSearchEnroll) // 教师查看报名结果

	g.POST("/processPassEnroll", self.ProcessPassEnroll)       // 教师审核报名通过
	g.POST("/processRejectEnroll", self.ProcessRejectEnroll)   // 教师审核报名驳回
	g.POST("/processRecoverEnroll", self.ProcessRecoverEnroll) // 教师审核报名恢复

	g.POST("/revokeEnroll", self.RevokeEnroll) //学生用户撤回报名信息
}

func (EnrollController) TeacherSearchEnroll(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	contest := c.DefaultQuery("contest", "")
	contestType := c.DefaultQuery("type", "")
	startTime := c.DefaultQuery("start_time", "")
	endTime := c.DefaultQuery("end_time", "")
	state := com.StrTo(c.DefaultQuery("state", "-1")).MustInt()

	key, exist := appG.C.Get("user_id")
	if !exist {
		appG.ResponseErr("用户不存在")
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

	list, total, err := logic.DefaultEnrollLogic.TeacherSearch(paginator, userID, contest, startTime, endTime, state, contestType)

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

// EnrollContest
// 报名竞赛
func (EnrollController) EnrollContest(c *gin.Context) {
	appG := app.Gin{C: c}

	form := &models.EnrollForm{}

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("用户不存在")
		return
	}

	err := c.ShouldBindJSON(form)
	if err != nil {
		DPrintf("EnrollContest c.ShouldBindJSON() 发生错误:", err)
		appG.ResponseErr("报名失败", err.Error())
		return
	}

	err = logic.DefaultEnrollLogic.InsertEnrollInformation(userID.(int64), form.Name, form.TeamID, form.Contest, form.School, form.Phone, form.Email)
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
	if role != TeacherRole {
		DPrintf("无教师权限")
		appG.ResponseErr("无教师权限")
		return
	}

	form := models.EnrollInformation{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultEnrollLogic.ProcessEnroll(form.ID, Pass, form.RejectReason)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
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
	if role != TeacherRole {
		DPrintf("无教师权限")
		appG.ResponseErr("无教师权限")
		return
	}

	form := models.EnrollInformation{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultEnrollLogic.ProcessEnroll(form.ID, Reject, form.RejectReason)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
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
	if role != TeacherRole {
		DPrintf("无教师权限")
		appG.ResponseErr("无教师权限")
		return
	}

	form := models.EnrollInformation{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultEnrollLogic.ProcessEnroll(form.ID, Processing, form.RejectReason)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
}

func (EnrollController) RevokeEnroll(c *gin.Context) {
	appG := app.Gin{C: c}

	role, exist := appG.C.Get("role")
	if !exist {
		DPrintf("无效身份")
		appG.ResponseErr("无效身份")
		return
	}
	if role != StudentRole {
		DPrintf("无学生权限")
		appG.ResponseErr("无学生权限")
		return
	}

	form := models.EnrollInformation{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultEnrollLogic.ProcessEnroll(form.ID, Revoked, form.RejectReason)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
}
