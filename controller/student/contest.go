package student

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	. "server/logic"
	logic "server/logic/student"
	"server/models"
	"server/utils/app"
)

type ContestController struct{}

// RegisterRoutes
// 注册学生用户路由
func (self ContestController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/contest_view", self.ViewContest) // 查看竞赛信息
	g.GET("/grade_view", self.ViewGrade)     // 成绩查询

	g.POST("/contest_register", self.RegisterContest) // 报名竞赛
}

// ViewCompetition
func (ContestController) ViewContest(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	var ret string

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultStudentContest.Display(paginator)
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

// ViewGrade
func (ContestController) ViewGrade(c *gin.Context) {
	appG := app.Gin{C: c}

	userid := com.StrTo(c.DefaultQuery("userid", "0")).MustInt()
	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	var err error

	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultStudentContest.FindGrade(userid, paginator)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	paginator.SetTotalPage(total)
	fmt.Println(list)

	if list != nil {
		data["list"] = list
		data["total"] = total
		data["page_size"] = limit
		data["page_number"] = curPage
		data["total_page"] = paginator.GetTotalPage()
	}

	appG.ResponseSucMsg(data, "查询成功")
}

// RegisterCompetition
func (ContestController) RegisterContest(c *gin.Context) {
	appG := app.Gin{C: c}
	var Param = models.StudentEntryParam{}
	var err error
	var ret string

	err = c.ShouldBindJSON(&Param)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	ret, err = logic.DefaultStudentContest.RegisterContest(Param.ContestantID, Param.ContestID)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(ret)
}
