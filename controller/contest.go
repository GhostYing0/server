package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"server/logic"
	"server/models"
	"server/utils/app"
	"server/utils/logging"
	. "server/utils/mydebug"
	"strconv"
)

type ContestController struct{}

func (self ContestController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/viewContest", self.ViewContest) // 查看竞赛信息

	g.GET("/viewTeacherContest", self.ViewTeacherContest) // 教师查看自身上传的竞赛信息
	g.POST("/uploadContest", self.UploadContest)          //教师上传竞赛信息
	g.POST("/updateContest", self.UpdateContest)          //教师更改竞赛信息
	g.DELETE("/deleteContest", self.DeleteContest)        //教师删除竞赛信息(暂时没用)

	g.POST("/processPassEnroll", self.ProcessPassEnroll)       // 教师审核报名通过
	g.POST("/processRejectEnroll", self.ProcessRejectEnroll)   // 教师审核报名驳回
	g.POST("/processRecoverEnroll", self.ProcessRecoverEnroll) // 教师审核报名恢复
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

	paginator := logic.NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultContestLogic.DisplayContest(paginator, contest, contestType)
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
func (ContestController) UploadContest(c *gin.Context) {}

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

// ProcessPassEnroll
// 审核通过
func (ContestController) ProcessPassEnroll(c *gin.Context) {
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

	count, err := logic.DefaultContestLogic.ProcessEnroll(&request.ID, 1)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(strconv.Itoa(int(count)))
}

// ProcessRejectEnroll
// 审核驳回
func (ContestController) ProcessRejectEnroll(c *gin.Context) {
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

	count, err := logic.DefaultContestLogic.ProcessEnroll(&request.ID, 2)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(strconv.Itoa(int(count)))
}

// ProcessRecoverEnroll
// 审核恢复
func (ContestController) ProcessRecoverEnroll(c *gin.Context) {
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

	count, err := logic.DefaultContestLogic.ProcessEnroll(&request.ID, 3)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(strconv.Itoa(int(count)))
}
