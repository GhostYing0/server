package cms

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	. "server/logic"
	logic "server/logic/cms"
	"server/models"
	"server/utils/app"
	"strconv"
)

type UserController struct{}

// RegisterRoutes
func (self UserController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/get_user", self.GetUser)          // 查看用户
	g.POST("/add_user", self.AddUser)         // 添加用户
	g.POST("/update_user", self.UpdateUser)   // 更新用户
	g.DELETE("/delete_user", self.DeleteUser) // 删除用户
}

// GetUser
func (UserController) GetUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	var ret string

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()

	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultCmsUser.Display(paginator)
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

// AddUser
func (UserController) AddUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var Param models.UserParam
	var err error
	var ret string

	err = c.ShouldBindJSON(&Param)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	ret, err = logic.DefaultCmsUser.AddUser(Param.Username, Param.Password)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(ret)
}

// UpdateUser
func (UserController) UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	var Param models.UpdateUserInfo
	var err error
	var ret string

	err = c.ShouldBindJSON(&Param)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	ret, err = logic.DefaultCmsUser.UpdateUser(Param.ID, Param.Username, Param.Password)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(ret)
}

// DeleteUser
func (UserController) DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	var ret string
	var count int64

	request := models.UserDeleteId{}
	err = c.ShouldBindJSON(&request)
	if err != nil {
		appG.ResponseErr(err.Error())
	}

	ret, err, count = logic.DefaultCmsUser.DeleteUser(&request.ID)
	if err != nil {
		appG.ResponseErr(err.Error())
	}

	appG.ResponseSuc(ret, "删除", strconv.Itoa(int(count)), "个用户成功")
}
