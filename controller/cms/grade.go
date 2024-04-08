package cms

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	. "server/logic"
	logic "server/logic/cms"
	"server/models"
	"server/utils/app"
	. "server/utils/mydebug"
	"strconv"
)

type GradeController struct{}

// 查看参赛者报名情况
// 以及比赛获奖情况

// GradeRoutes
func (self GradeController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/getGrade", self.GetGrade)          // 查看成绩信息
	g.POST("/addGrade", self.AddGrade)         // 添加成绩信息
	g.POST("/updateGrade", self.UpdateGrade)   // 更改成绩信息
	g.DELETE("/deleteGrade", self.DeleteGrade) // 删除成绩信息

	g.GET("/getGradeCount", self.GetCount)
}

// GetRegisteredContestByUser
func (GradeController) GetGrade(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	username := c.DefaultQuery("username", "")
	contest := c.DefaultQuery("contest_name", "")
	startTime := c.DefaultQuery("start_time", "")
	endTime := c.DefaultQuery("end_time", "")
	grade := c.DefaultQuery("grade", "")
	state := com.StrTo(c.DefaultQuery("state", "-1")).MustInt()

	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultGradeContest.Display(paginator, username, contest, startTime, endTime, grade, state)
	if err != nil {
		DPrintf("logic.DefaultGradeContest.Display 发生错误:", err)
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

	appG.ResponseSucMsg(data, "查询成功")
}

// AddRegisteredContestByUser
func (GradeController) AddGrade(c *gin.Context) {
	appG := app.Gin{C: c}
	param := &models.GradeForm{}

	err := c.ShouldBindJSON(&param)
	if err != nil {
		DPrintf("AddGradeInformation c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	fmt.Println("asd:", param)
	err = logic.DefaultGradeContest.Add(param.Username, param.Contest, param.Grade, param.CreateTime, param.Certificate, param.State)
	if err != nil {
		DPrintf("AddGradeInformation 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("添加成功")
}

// UpdateRegisteredContestByUser
func (GradeController) UpdateGrade(c *gin.Context) {
	appG := app.Gin{C: c}
	param := &models.GradeForm{}

	err := c.ShouldBindJSON(param)
	if err != nil {
		DPrintf("UpdateGradeInformation c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultGradeContest.Update(param.ID, param.Username, param.Contest, param.Grade, param.CreateTime, param.Certificate, param.State)
	if err != nil {
		DPrintf("UpdateGradeInformation 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("操作成功")
}

// DeleteRegisteredContestByUser
func (GradeController) DeleteGrade(c *gin.Context) {
	appG := app.Gin{C: c}

	request := models.GradeIds{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		DPrintf("DeleteGradeInformation c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultGradeContest.Delete(&request.ID)
	if err != nil {
		DPrintf("DeleteGradeInformation 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("删除", strconv.Itoa(int(count)), "条成绩信息成功")
}

// GetCount
func (GradeController) GetCount(c *gin.Context) {
	appG := app.Gin{C: c}

	count, err := logic.DefaultGradeContest.GetGradeCount()

	if err != nil {
		DPrintf("GetUserCount 出错:", err)
		appG.ResponseErr(err.Error())
		return
	}
	data := make(map[string]interface{})
	data["count"] = count

	appG.ResponseSucMsg(data, "查询成功")
}