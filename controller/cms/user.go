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
)

type UserController struct{}

// RegisterRoutes
func (self UserController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/get_user", self.GetUser)          // 查看用户
	g.POST("/add_user", self.AddUser)         // 添加用户
	g.POST("/update_user", self.UpdateUser)   // 更新用户
	g.DELETE("/delete_user", self.DeleteUser) // 删除用户

	g.GET("/getCount", self.GetCount)
}

// GetUser
func (UserController) GetUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	var ret string

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()
	mode := com.StrTo(c.DefaultQuery("mode", "0")).MustInt()
	username := c.DefaultQuery("searchUser", "")

	fmt.Println("limit:", limit, " curPage:", curPage)
	fmt.Println(c.Request.URL)
	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultCmsUser.Display(paginator, mode, username)
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
	var Param models.User
	var err error
	var ret string

	err = c.ShouldBindJSON(&Param)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	DPrintf(Param)

	ret, err = logic.DefaultCmsUser.AddUser(Param.Username, Param.Password, Param.Role)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc(ret)
}

// UpdateUser
func (UserController) UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	param := &models.UpdateUserInfo{}

	err := c.ShouldBindJSON(&param)
	if err != nil {
		DPrintf("UpdateUser c.ShouldBindJSON err:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultCmsUser.UpdateUser(param.ID, param.Username, param.Password, param.Role)
	if err != nil {
		DPrintf("UpdateUser logic.DefaultCmsUser.UpdateUser err:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("更新成功")
}

// DeleteUser
func (UserController) DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}

	request := models.UserDeleteId{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		DPrintf("DeleteUser c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultCmsUser.DeleteUser(&request.ID)
	if err != nil {
		DPrintf("DeleteUser logic.DefaultCmsUser.DeleteUser 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseNumber(count)
}

// GetCount
func (UserController) GetCount(c *gin.Context) {
	appG := app.Gin{C: c}

	count, err := logic.DefaultCmsUser.GetUserCount()

	if err != nil {
		DPrintf("GetUserCount 出错:", err)
		appG.ResponseErr(err.Error())
		return
	}
	data := make(map[string]interface{})
	data["count"] = count

	appG.ResponseSucMsg(data, "查询成功")
}
