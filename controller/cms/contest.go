package cms

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	. "server/logic"
	logic "server/logic/cms"
	"server/models"
	"server/utils/app"
	"strconv"
)

type ContestController struct{}

// RegisterRoutes
func (self ContestController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/getContest", self.GetContest)          // 查看竞赛信息
	g.POST("/addContest", self.AddContest)         // 添加竞赛信息
	g.POST("/updateContest", self.UpdateContest)   // 更改竞赛信息
	g.DELETE("/deleteContest", self.DeleteContest) // 删除竞赛信息
	g.POST("/processContest", self.ProcessContest) //审核竞赛
}

// GetContest
func (ContestController) GetContest(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	var ret string

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultCmsContest.Display(paginator)
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
	var Param models.ContestParam
	var err error
	var ret string

	err = c.ShouldBindJSON(&Param)
	if err != nil {
		fmt.Println("ShouldBindJSON error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	ret, err = logic.DefaultCmsContest.InsertContest(Param.Name, Param.Type, Param.StartDate, Param.Deadline)
	if err != nil {
		fmt.Println("logic.InsertContestInfo error:", err)
		appG.ResponseErr(ret)
		return
	}

	appG.ResponseSuc(ret)
}

// UpdateContest
func (ContestController) UpdateContest(c *gin.Context) {
	appG := app.Gin{C: c}
	var Param models.UpdateContestParam
	var err error
	var ret string

	err = c.ShouldBindJSON(&Param)
	if err != nil {
		fmt.Println("ShouldBindJSON error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	ret, err = logic.DefaultCmsContest.UpdateContest(Param.ID, Param.Name, Param.Type, Param.StartDate, Param.Deadline)
	if err != nil {
		fmt.Println("logic.UpdateContestInfo error:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(ret)
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

func (ContestController) ProcessContest(c *gin.Context) {
}
