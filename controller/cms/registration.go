package cms

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	. "server/logic"
	logic "server/logic/cms"
	"server/models"
	"server/utils/app"
)

type RegistrationController struct{}

// 查看参赛者报名情况
// 以及比赛获奖情况

// RegisterRoutes
func (self RegistrationController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/get_registration", self.GetRegistration)          // 查看用户报名信息
	g.POST("/add_registration", self.AddRegistration)         // 添加用户报名信息
	g.POST("/update_registration", self.UpdateRegistration)   // 更改用户报名信息
	g.DELETE("/delete_registration", self.DeleteRegistration) // 删除用户报名信息

	g.GET("/get_grade", self.GetGradeByUser)        // 查看学生竞赛分数
	g.POST("/update_grade", self.UpdateGradeByUser) // 更改学生竞赛分数
}

// GetRegisteredContestByUser
func (RegistrationController) GetRegistration(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	var ret string

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultRegistrationContest.Display(paginator)
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

// AddRegisteredContestByUser
func (RegistrationController) AddRegistration(c *gin.Context) {
	appG := app.Gin{C: c}
	var Param models.EntryContestParam
	var err error
	var ret string

	err = c.ShouldBindJSON(&Param)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	ret, err = logic.DefaultRegistrationContest.AddRegistration(Param.Contestant, Param.Contest)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(ret)
}

// UpdateRegisteredContestByUser
func (RegistrationController) UpdateRegistration(c *gin.Context) {

}

// DeleteRegisteredContestByUser
func (RegistrationController) DeleteRegistration(c *gin.Context) {

}

// GetGradeByUser
func (RegistrationController) GetGradeByUser(c *gin.Context) {

}

// UpdateGradeByUser
func (RegistrationController) UpdateGradeByUser(c *gin.Context) {

}
