package cms

import (
	"fmt"
	. "server/logic"
	logic "server/logic/cms"
	"server/models"
	"server/utils/app"
	. "server/utils/e"
	. "server/utils/mydebug"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type ContestController struct{}

// RegisterRoutes
func (self ContestController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/getContest", self.GetContest)          // 查看竞赛信息
	g.POST("/addContest", self.AddContest)         // 添加竞赛信息
	g.POST("/updateContest", self.UpdateContest)   // 更改竞赛信息
	g.DELETE("/deleteContest", self.DeleteContest) // 删除竞赛信息

	g.POST("/processPassContest", self.ProcessPassContest)       // 管理员审核竞赛通过
	g.POST("/processRejectContest", self.ProcessRejectContest)   // 管理员审核竞赛驳回
	g.POST("/processRecoverContest", self.ProcessRecoverContest) // 管理员审核竞赛恢复

	g.GET("/getContestCount", self.GetCount)
}

// GetContest
func (ContestController) GetContest(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	var ret string

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()
	contest := c.DefaultQuery("contest", "")
	contestType := c.DefaultQuery("contest_type", "")
	state := com.StrTo(c.DefaultQuery("state", "-1")).MustInt()
	contestLevel := com.StrTo(c.DefaultQuery("state", "0")).MustInt()

	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultCmsContest.Display(paginator, contest, contestType, state, contestLevel)
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

// AddContest
func (ContestController) AddContest(c *gin.Context) {
	appG := app.Gin{C: c}
	var form models.ContestForm

	err := c.ShouldBindJSON(&form)
	if err != nil {
		fmt.Println("ShouldBindJSON error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultCmsContest.InsertContest(form.Username, form.Contest, form.ContestType, form.StartTime, form.Deadline, form.State)
	if err != nil {
		fmt.Println("logic.InsertContestInfo error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
}

// UpdateContest
func (ContestController) UpdateContest(c *gin.Context) {
	appG := app.Gin{C: c}
	var form models.UpdateContestForm

	err := c.ShouldBindJSON(&form)
	if err != nil {
		fmt.Println("ShouldBindJSON error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultCmsContest.UpdateContest(form.ID, form.Username, form.Contest, form.ContestType, form.StartTime, form.Deadline, form.State)
	if err != nil {
		fmt.Println("logic.UpdateContestInfo error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
}

// DeleteContest
func (ContestController) DeleteContest(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	var ret string
	var count int64

	request := models.ContestDeleteId{}
	err = c.ShouldBindJSON(&request)
	if err != nil {
		appG.ResponseErr(err.Error())
	}

	ret, count, err = logic.DefaultCmsContest.DeleteContest(&request.ID)
	if err != nil {
		appG.ResponseErr(err.Error())
	}

	appG.ResponseSuc(ret, "删除", strconv.Itoa(int(count)), "个竞赛成功")
}

// GetCount
func (ContestController) GetCount(c *gin.Context) {
	appG := app.Gin{C: c}

	count, err := logic.DefaultCmsContest.GetContestCount()

	if err != nil {
		DPrintf("GetEnrollCount 出错:", err)
		appG.ResponseErr(err.Error())
		return
	}
	data := make(map[string]interface{})
	data["count"] = count

	appG.ResponseSucMsg(data, "查询成功")
}

func (ContestController) ProcessPassContest(c *gin.Context) {
	appG := app.Gin{C: c}

	form := &models.Contest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultCmsContest.ProcessContest(form.ID, Pass)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg("success")
}

func (ContestController) ProcessRejectContest(c *gin.Context) {
	appG := app.Gin{C: c}

	form := &models.Contest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultCmsContest.ProcessContest(form.ID, Reject)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg("success")
}

func (ContestController) ProcessRecoverContest(c *gin.Context) {
	appG := app.Gin{C: c}

	form := &models.Contest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultCmsContest.ProcessContest(form.ID, Processing)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg("success")
}
