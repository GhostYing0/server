package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"server/logic"
	"server/models"
	"server/utils/app"
	"server/utils/e"
	"server/utils/logging"
	. "server/utils/mydebug"
)

type ContestController struct{}

func (self ContestController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/viewContest", self.ViewContest) // 查看竞赛信息

	g.GET("/viewTeacherContest", self.ViewTeacherContest)     // 教师查看自身上传的竞赛信息
	g.GET("/getContestForTeacher", self.GetContestForTeacher) // 获取自身竞赛，给选择框用
	g.GET("getDepartmentContest", self.GetDepartmentContest)  // 系部管理员获取同校同院系竞赛
	g.POST("/uploadContest", self.UploadContest)              //教师上传竞赛信息
	g.POST("/updateContest", self.UpdateContest)              //教师更改竞赛信息
	g.POST("/transformState", self.TransformState)            //教师开关竞赛报名
	g.DELETE("/deleteContest", self.DeleteContest)            //教师删除竞赛信息(暂时没用)
	g.POST("/cancelContest", self.CancelContest)              //教师撤回竞赛

	g.POST("/processPassContest", self.ProcessPassContest)       // 系部管理员审核竞赛通过
	g.POST("/processRejectContest", self.ProcessRejectContest)   // 系部管理员审核竞赛驳回
	g.POST("/processRecoverContest", self.ProcessRecoverContest) // 系部管理员审核竞赛恢复
}

// Viewcontest
// 查看竞赛
func (ContestController) ViewContest(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()
	contestType := c.DefaultQuery("type", "")
	contest := c.DefaultQuery("contest", "")

	if limit < 0 || curPage < 0 {
		DPrintf("分页器参数错误")
		appG.ResponseErr("分页器参数错误")
		return
	}

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("请登录")
		return
	}

	paginator := logic.NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultContestLogic.DisplayContest(paginator, contest, contestType, userID.(int64))
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

	appG.ResponseSucMsg(data)
}

// ViewTeacherContest
// 查看自身竞赛
func (ContestController) ViewTeacherContest(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()
	contestType := c.DefaultQuery("type", "")
	contest := c.DefaultQuery("contest", "")
	state := com.StrTo(c.DefaultQuery("", "-1")).MustInt()

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("分页器参数错误")
		logging.L.Error("用户不存在")
		return
	}

	if limit < 0 || curPage < 0 {
		DPrintf("分页器参数错误")
		appG.ResponseErr("分页器参数错误")
		return
	}

	paginator := logic.NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultContestLogic.ViewTeacherContest(paginator, userID.(int64), contest, contestType, state)
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

	appG.ResponseSucMsg(data)
}

// UploadContest
func (ContestController) UploadContest(c *gin.Context) {
	appG := app.Gin{C: c}
	form := &models.ContestForm{}

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("用户不存在")
		return
	}

	err := c.ShouldBindJSON(form)
	if err != nil {
		fmt.Println("ShouldBindJSON error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultContestLogic.UploadContest(userID.(int64), form.Contest, form.ContestType, form.StartTime, form.Deadline, &form.Describe)
	if err != nil {
		fmt.Println("logic.UploadContest error:", err)
		appG.ResponseErr(err.Error())
		return
	}
	appG.ResponseSuc()
}

// UpdateContest
func (ContestController) UpdateContest(c *gin.Context) {
	appG := app.Gin{C: c}
	form := &models.UpdateContestForm{}

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("用户不存在")
		return
	}

	err := c.ShouldBindJSON(form)
	if err != nil {
		fmt.Println("ShouldBindJSON error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultContestLogic.UpdateContest(form.ID, userID.(int64), form.Contest, form.ContestType, form.StartTime, form.Deadline, form.ContestState, form.State)
	if err != nil {
		fmt.Println("logic.UpdateContestInfo error:", err)
		appG.ResponseErr(err.Error())
		return
	}
	appG.ResponseSuc()
}

// DeleteContest
func (ContestController) DeleteContest(c *gin.Context) {}

// GetContest
func (ContestController) GetContestForTeacher(c *gin.Context) {
	appG := app.Gin{C: c}

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

	if role != e.TeacherRole {
		appG.ResponseErr("无权限")
		return
	}

	data, err := logic.DefaultContestLogic.GetContestForTeacher(userID.(int64))
	if err != nil {
		DPrintf("教师获取自身竞赛出错:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg(data)
}

func (ContestController) TransformState(c *gin.Context) {
	appG := app.Gin{C: c}
	form := &models.UpdateContestForm{}

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

	if role != e.TeacherRole {
		appG.ResponseErr("无权限")
		return
	}

	err := c.ShouldBindJSON(form)
	if err != nil {
		fmt.Println("ShouldBindJSON error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultContestLogic.TransformState(userID.(int64), form.ID, form.ContestState)
	if err != nil {
		DPrintf("教师获取自身竞赛出错:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
}

func (ContestController) CancelContest(c *gin.Context) {
	appG := app.Gin{C: c}
	form := &models.UpdateContestForm{}

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("用户不存在")
		return
	}

	err := c.ShouldBindJSON(form)
	if err != nil {
		fmt.Println("ShouldBindJSON error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultContestLogic.CancelContest(form.ID, userID.(int64))
	if err != nil {
		fmt.Println("logic.UpdateContestInfo error:", err)
		appG.ResponseErr(err.Error())
		return
	}
	appG.ResponseSuc()
}

func (ContestController) GetDepartmentContest(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	var ret string

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

	if role != e.DepartmentRole {
		appG.ResponseErr("无权限")
		return
	}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()
	contest := c.DefaultQuery("contest", "")
	contestType := c.DefaultQuery("contest_type", "")
	state := com.StrTo(c.DefaultQuery("state", "-1")).MustInt()

	paginator := logic.NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultContestLogic.DepartmentManagerGetContest(paginator, contest, contestType, state, userID.(int64))
	if err != nil {
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

func (ContestController) ProcessPassContest(c *gin.Context) {
	appG := app.Gin{C: c}

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

	if role != e.DepartmentRole {
		appG.ResponseErr("无权限")
		return
	}

	form := &models.ProcessContest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultContestLogic.ProcessContest(form.ID, e.Pass, userID.(int64), form.RejectReason)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg("success")
}

func (ContestController) ProcessRejectContest(c *gin.Context) {
	appG := app.Gin{C: c}

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

	if role != e.DepartmentRole {
		appG.ResponseErr("无权限")
		return
	}

	form := &models.ProcessContest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultContestLogic.ProcessContest(form.ID, e.Reject, userID.(int64), form.RejectReason)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg("success")
}

func (ContestController) ProcessRecoverContest(c *gin.Context) {
	appG := app.Gin{C: c}

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

	if role != e.DepartmentRole {
		appG.ResponseErr("无权限")
		return
	}

	form := &models.ProcessContest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultContestLogic.ProcessContest(form.ID, e.Processing, userID.(int64), form.RejectReason)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg("success")
}
