package controller

import (
	"server/logic"
	"server/models"
	"server/utils/app"
	. "server/utils/e"
	"server/utils/logging"
	. "server/utils/mydebug"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type GradeController struct{}

func (self GradeController) RegisterRoutes(g *gin.RouterGroup) {
	g.POST("/uploadGrade", self.UploadGrade) // 上传成绩
	g.GET("/searchGrade", self.DisplayGrade) // 查看成绩

	g.GET("/teacherSearchGrade", self.TeacherDisplayGrade) // 教师查看自身竞赛成绩
	//g.POST("/processPassGrade", self.ProcessPassGrade)       // 教师审核成绩通过
	//g.POST("/processRejectGrade", self.ProcessRejectGrade)   // 教师审核成绩驳回
	//g.POST("/processRecoverGrade", self.ProcessRecoverGrade) // 教师审核成绩恢复

	g.GET("/departmentManagerSearchGrade", self.DepartmentManagerSearchGrade) // 系部管理查看自身竞赛成绩
	g.POST("/processPassGrade", self.ProcessPassGrade)                        // 系部管理员审核成绩通过
	g.POST("/processRejectGrade", self.ProcessRejectGrade)                    // 系部管理员审核成绩驳回
	g.POST("/processRecoverGrade", self.ProcessRecoverGrade)                  // 系部管理员审核成绩恢复

	g.POST("/revokeGrade", self.RevokeGrade)
	g.POST("/studentUpdateGrade", self.StudentUpdateGrade)
}

func (self GradeController) UploadGrade(c *gin.Context) {
	appG := app.Gin{C: c}

	form := &models.UploadGradeForm{}

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("用户不存在")
		return
	}

	role, exist := c.Get("role")
	if !exist {
		appG.ResponseErr("权限错误")
		return
	}
	if role.(int) != TeacherRole {
		appG.ResponseErr("无权限")
		return
	}

	err := c.ShouldBindJSON(form)
	if err != nil {
		DPrintf("UploadGrade c.ShouldBindJSON() 发生错误:", err)
		appG.ResponseErr("上传成绩失败", err.Error())
		return
	}

	err = logic.DefaultGradeLogic.InsertGradeInformation(userID.(int64), form.EnrollID, form.Grade, form.RewardTime, form.Certificate,
		form.GuidanceTeacher, form.TeacherDepartment, form.TeacherTitle)
	if err != nil {
		DPrintf("EnrollContest logic.DefaultEnrollLogic.InsertEnrollInformation() 发生错误:", err)
		appG.ResponseErr("上传成绩失败", err.Error())
		return
	}

	appG.ResponseSucMsg("上传成功")
}

func (self GradeController) DisplayGrade(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	grade := c.DefaultQuery("grade", "")
	contest := c.DefaultQuery("contest", "")
	startTime := c.DefaultQuery("startTime", "")
	endTime := c.DefaultQuery("endTime", "")
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

	list, total, err := logic.DefaultGradeLogic.Search(paginator, grade, contest, startTime, endTime, state, user_id.(int64), role.(int))

	if err != nil {
		DPrintf("DisplayGrade 发生错误:", err)
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

func (self GradeController) TeacherDisplayGrade(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	contestID := com.StrTo(c.DefaultQuery("id", "0")).MustInt64()
	grade := c.DefaultQuery("grade", "")
	contest := c.DefaultQuery("contest", "")
	startTime := c.DefaultQuery("startTime", "")
	endTime := c.DefaultQuery("endTime", "")
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

	if role != TeacherRole {
		appG.ResponseErr("无权限")
		logging.L.Error("无权限")
		return
	}

	list, total, err := logic.DefaultGradeLogic.TeacherSearch(paginator, grade, contest, startTime, endTime, state, contestID, user_id.(int64), role.(int))

	if err != nil {
		DPrintf("DisplayGrade 发生错误:", err)
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

// ProcessPassGrade
// 审核通过
func (GradeController) ProcessPassGrade(c *gin.Context) {
	appG := app.Gin{C: c}

	role, exist := appG.C.Get("role")
	if !exist {
		DPrintf("无效身份")
		appG.ResponseErr("无效身份")
		return
	}
	if role != DepartmentRole {
		DPrintf("无教师权限")
		appG.ResponseErr("无教师权限")
		return
	}

	form := models.PassGradeID{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultGradeLogic.PassGrade(&form.IDS, Pass)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("审核通过" + strconv.Itoa(int(count)) + "个成功")
}

// ProcessRejectGrade
// 审核驳回
func (GradeController) ProcessRejectGrade(c *gin.Context) {
	appG := app.Gin{C: c}

	role, exist := appG.C.Get("role")
	if !exist {
		DPrintf("无效身份")
		appG.ResponseErr("无效身份")
		return
	}
	if role != DepartmentRole {
		DPrintf("无教师权限")
		appG.ResponseErr("无教师权限")
		return
	}

	form := models.GradeForm{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultGradeLogic.ProcessGrade(form.ID, Reject, form.RejectReason)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
}

// ProcessRecoverGrade
// 审核恢复
func (GradeController) ProcessRecoverGrade(c *gin.Context) {
	appG := app.Gin{C: c}

	role, exist := appG.C.Get("role")
	if !exist {
		DPrintf("无效身份")
		appG.ResponseErr("无效身份")
		return
	}
	if role != DepartmentRole {
		DPrintf("无教师权限")
		appG.ResponseErr("无教师权限")
		return
	}

	form := models.GradeForm{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultGradeLogic.ProcessGrade(form.ID, Processing, form.RejectReason)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
}

func (GradeController) RevokeGrade(c *gin.Context) {
	appG := app.Gin{C: c}

	role, exist := appG.C.Get("role")
	if !exist {
		DPrintf("无效身份")
		appG.ResponseErr("无效身份")
		return
	}
	if role != TeacherRole {
		DPrintf("无学生权限")
		appG.ResponseErr("无学生权限")
		return
	}

	form := models.GradeForm{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultGradeLogic.ProcessGrade(form.ID, Revoked, form.RejectReason)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
}

func (GradeController) StudentUpdateGrade(c *gin.Context) {
	appG := app.Gin{C: c}
	form := &models.GradeForm{}

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

	err := c.ShouldBindJSON(form)
	if err != nil {
		DPrintf("UpdateGradeInformation c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultGradeLogic.Update(form.ID, form.Grade, form.Certificate)
	if err != nil {
		DPrintf("UpdateGradeInformation 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("操作成功")
}

func (self GradeController) DepartmentManagerSearchGrade(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	grade := c.DefaultQuery("grade", "")
	contest := c.DefaultQuery("contest", "")
	startTime := c.DefaultQuery("startTime", "")
	endTime := c.DefaultQuery("endTime", "")
	contestID := com.StrTo(c.DefaultQuery("id", "0")).MustInt64()
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

	if role != DepartmentRole {
		appG.ResponseErr("无权限")
		logging.L.Error("无权限")
		return
	}

	list, total, err := logic.DefaultGradeLogic.DepartmentManagerSearchGrade(paginator, grade, contest, startTime, endTime, state, contestID, user_id.(int64), role.(int))

	if err != nil {
		DPrintf("DisplayGrade 发生错误:", err)
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
