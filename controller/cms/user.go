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
	g.GET("/getStudentUser", self.GetStudentUser) // 查看用户
	g.GET("/getTeacherUser", self.GetTeacherUser) // 查看用户
	g.POST("/addUser", self.AddUser)              // 添加用户
	g.POST("/updateUser", self.UpdateUser)        // 更新用户
	g.DELETE("/deleteUSer", self.DeleteUser)      // 删除用户

	g.GET("/getUserCount", self.GetCount)
}

// GetStudentUser
func (UserController) GetStudentUser(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()
	username := c.DefaultQuery("searchUser", "")
	gender := c.DefaultQuery("gender", "")
	school := c.DefaultQuery("school", "")
	semester := c.DefaultQuery("semester", "")
	college := c.DefaultQuery("college", "")
	class := c.DefaultQuery("class", "")
	name := c.DefaultQuery("name", "")

	fmt.Println("limit:", limit, " curPage:", curPage)
	fmt.Println(c.Request.URL)
	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultCmsUser.DisplayStudent(paginator, username, gender, school, semester, college, class, name)
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

	appG.ResponseSucMsg(data)
}

// GetTeacherUser
func (UserController) GetTeacherUser(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()
	username := c.DefaultQuery("searchUser", "")
	gender := c.DefaultQuery("gender", "")
	school := c.DefaultQuery("school", "")
	semester := c.DefaultQuery("semester", "")
	college := c.DefaultQuery("college", "")
	class := c.DefaultQuery("class", "")
	name := c.DefaultQuery("name", "")

	fmt.Println("limit:", limit, " curPage:", curPage)
	fmt.Println(c.Request.URL)
	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultCmsUser.DisplayTeacher(paginator, username, gender, school, semester, college, class, name)
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

	appG.ResponseSucMsg(data)
}

// AddUser
func (UserController) AddUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var Param models.OldUser
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
