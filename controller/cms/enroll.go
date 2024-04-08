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

type RegistrationController struct{}

// 查看参赛者报名情况
// 以及比赛获奖情况

// RegisterRoutes
func (self RegistrationController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/getEnrollInformation", self.GetEnrollInformation)          // 查看用户报名信息
	g.POST("/addEnrollInformation", self.AddEnrollInformation)         // 添加用户报名信息
	g.POST("/updateEnrollInformation", self.UpdateEnrollInformation)   // 更改用户报名信息
	g.DELETE("/deleteEnrollInformation", self.DeleteEnrollInformation) // 删除用户报名信息

	g.GET("/getEnrollCount", self.GetEnrollCount)
}

// GetRegisteredContestByUser
func (RegistrationController) GetEnrollInformation(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	name := c.DefaultQuery("name", "")
	//userID := com.StrTo(c.DefaultQuery("user_id", "-1")).MustInt64()
	//teamID := c.DefaultQuery("team_id", "")
	contest := c.DefaultQuery("contest", "")
	startTime := c.DefaultQuery("start_time", "")
	endTime := c.DefaultQuery("end_time", "")
	//school := c.DefaultQuery("school", "")
	//phone := c.DefaultQuery("phone", "")
	//email := c.DefaultQuery("email", "")
	state := com.StrTo(c.DefaultQuery("state", "-1")).MustInt()

	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultEnrollContest.Display(paginator, name, contest, startTime, endTime, state)
	if err != nil {
		DPrintf("logic.DefaultRegistrationContest.Display 发生错误:", err)
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
func (RegistrationController) AddEnrollInformation(c *gin.Context) {
	appG := app.Gin{C: c}
	param := &models.NewEnrollInformation{}

	err := c.ShouldBindJSON(&param)
	if err != nil {
		DPrintf("AddEnrollInformation c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	fmt.Println("asd:", param)
	err = logic.DefaultEnrollContest.Add(param.UserName, param.TeamID, param.ContestName, param.CreateTime, param.School, param.Phone, param.Email, param.State)
	if err != nil {
		DPrintf("AddEnrollInformation 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("添加成功")
}

// UpdateRegisteredContestByUser
func (RegistrationController) UpdateEnrollInformation(c *gin.Context) {
	appG := app.Gin{C: c}
	param := &models.EnrollForm{}

	err := c.ShouldBindJSON(param)
	if err != nil {
		DPrintf("UpdateEnrollInformation c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultEnrollContest.Update(param.ID, param.UserName, param.TeamID, param.ContestName, param.CreateTime, param.School, param.Phone, param.Email, param.State)
	if err != nil {
		DPrintf("UpdateEnrollInformation 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("操作成功")
}

// DeleteRegisteredContestByUser
func (RegistrationController) DeleteEnrollInformation(c *gin.Context) {
	appG := app.Gin{C: c}

	request := models.EnrollDeleteId{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		DPrintf("DeleteEnrollInformation c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultEnrollContest.Delete(&request.ID)
	if err != nil {
		DPrintf("DeleteEnrollInformation 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("删除", strconv.Itoa(int(count)), "条报名信息成功")
}

// GetCount
func (RegistrationController) GetEnrollCount(c *gin.Context) {
	appG := app.Gin{C: c}

	count, err := logic.DefaultEnrollContest.GetEnrollCount()

	if err != nil {
		DPrintf("GetEnrollCount 出错:", err)
		appG.ResponseErr(err.Error())
		return
	}
	data := make(map[string]interface{})
	data["count"] = count

	appG.ResponseSucMsg(data, "查询成功")
}