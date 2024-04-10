package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"server/logic"
	"server/models"
	"server/utils/app"
	. "server/utils/mydebug"
	"strconv"
)

type GradeController struct{}

func (self GradeController) RegisterRoutes(g *gin.RouterGroup) {
	g.POST("/uploadGrade", self.UploadGrade)                 // 上传成绩
	g.GET("/searchGrade", self.DisplayGrade)                 // 查看成绩
	g.POST("/processPassGrade", self.ProcessPassGrade)       // 教师审核成绩通过
	g.POST("/processRejectGrade", self.ProcessRejectGrade)   // 教师审核成绩驳回
	g.POST("/processRecoverGrade", self.ProcessRecoverGrade) // 教师审核成绩恢复
}

func (self GradeController) UploadGrade(c *gin.Context) {
	appG := app.Gin{C: c}

	form := &models.GradeForm{}

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("用户不存在")
		return
	}

	err := c.ShouldBindJSON(form)
	if err != nil {
		DPrintf("UploadGrade c.ShouldBindJSON() 发生错误:", err)
		appG.ResponseErr("上传成绩失败", err.Error())
		return
	}

	err = logic.DefaultGradeLogic.InsertGradeInformation(userID.(int64), form.Contest, form.Grade, form.Certificate)
	if err != nil {
		DPrintf("EnrollContest logic.DefaultEnrollLogic.InsertEnrollInformation() 发生错误:", err)
		appG.ResponseErr("上传成绩失败", err.Error())
		return
	}

	//file, err := c.FormFile("file")
	//
	//if err != nil {
	//	DPrintf("grade c.FormFile 发生错误:", err)
	//	appG.ResponseErr("获取文件失败", err.Error())
	//	return
	//}
	//saveDir, err := public.DefaultPublic.UploadImg(file)
	//if err != nil {
	//	DPrintf("logic.DefaultPublic.Upload 发生错误:", err)
	//	appG.ResponseErr(err.Error())
	//	return
	//}
	//
	//err = c.SaveUploadedFile(file, saveDir)
	//if err != nil {
	//	DPrintf("上传错误")
	//	// 返回值
	//	appG.ResponseErr("文件保存失败")
	//	return
	//}
	//
	//imageurl := strings.Replace(saveDir, "D:/GDesign/code/server/picture/img", "http://127.0.0.1:9006/api/public/picture", -1)

	//appG.ResponseSucMsg(gin.H{"imageurl": imageurl}, "上传成功")
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
	if role != 2 {
		DPrintf("无教师权限")
		appG.ResponseErr("无教师权限")
		return
	}

	request := models.GradeIds{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultGradeLogic.ProcessGrade(&request.ID, 1)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(strconv.Itoa(int(count)))
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
	if role != 2 {
		DPrintf("无教师权限")
		appG.ResponseErr("无教师权限")
		return
	}

	request := models.GradeIds{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultGradeLogic.ProcessGrade(&request.ID, 2)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(strconv.Itoa(int(count)))
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
	if role != 2 {
		DPrintf("无教师权限")
		appG.ResponseErr("无教师权限")
		return
	}

	request := models.GradeIds{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		DPrintf("ProcessEnroll c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultGradeLogic.ProcessGrade(&request.ID, 3)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Process 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(strconv.Itoa(int(count)))
}
