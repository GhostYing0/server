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
	g.GET("/getStudent", self.GetStudent)           // 查看学生用户
	g.POST("/addStudent", self.AddStudent)          // 添加学生用户
	g.POST("/updateStudent", self.UpdateUser)       // 更新学生用户
	g.DELETE("/deleteStudent", self.DeleteStudent)  // 删除学生用户
	g.GET("/getStudentCount", self.GetStudentCount) // 获取学生用户数量

	g.GET("/getTeacher", self.GetTeacher)           // 查看教师用户
	g.POST("/addTeacher", self.AddTeacher)          // 添加教师用户
	g.POST("/updateTeacher", self.UpdateTeacher)    // 更新教师用户
	g.DELETE("/deleteTeacher", self.DeleteTeacher)  // 删除教师用户
	g.GET("/getTeacherCount", self.GetTeacherCount) // 获取教师用户数量

	g.GET("/getManager", self.GetManager)           // 查看管理员用户
	g.POST("/addManager", self.AddManager)          // 添加管理员用户
	g.POST("/updateManager", self.UpdateManager)    // 修改管理员用户
	g.DELETE("/deleteManager", self.DeleteManager)  // 删除管理员用户
	g.GET("/getManagerCount", self.GetManagerCount) // 获取管理员用户数量
}

// =================================================学生用户
// GetStudentUser
func (UserController) GetStudent(c *gin.Context) {
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
	list, total, err := logic.DefaultCmsStudent.DisplayStudent(paginator, username, gender, school, semester, college, class, name)
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
func (UserController) AddStudent(c *gin.Context) {
	appG := app.Gin{C: c}
	var form models.StudentForm

	err := c.ShouldBindJSON(&form)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultCmsStudent.AddStudent(form.Username, form.Password, form.Name, form.Gender, form.School,
		form.College, form.Class, form.Semester, form.Avatar)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
}

// UpdateUser
func (UserController) UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	form := &models.StudentForm{}

	err := c.ShouldBindJSON(&form)
	if err != nil {
		DPrintf("UpdateUser c.ShouldBindJSON err:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultCmsStudent.UpdateStudent(form.ID, form.Username, form.Password, form.Name, form.Gender, form.School,
		form.College, form.Class, form.Semester, form.Avatar)
	if err != nil {
		DPrintf("UpdateUser logic.DefaultCmsUser.UpdateUser err:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("更新成功")
}

// DeleteUser
func (UserController) DeleteStudent(c *gin.Context) {
	appG := app.Gin{C: c}

	request := models.UserDeleteId{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		DPrintf("DeleteUser c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultCmsStudent.DeleteStudent(&request.ID)
	if err != nil {
		DPrintf("DeleteUser logic.DefaultCmsUser.DeleteUser 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseNumber(count)
}

// GetCount
func (UserController) GetStudentCount(c *gin.Context) {
	appG := app.Gin{C: c}

	count, err := logic.DefaultCmsStudent.GetStudentCount()

	if err != nil {
		DPrintf("GetUserCount 出错:", err)
		appG.ResponseErr(err.Error())
		return
	}
	data := make(map[string]interface{})
	data["count"] = count

	appG.ResponseSucMsg(data, "查询成功")
}

// =================================================教师用户
// GetTeacher
func (UserController) GetTeacher(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()
	username := c.DefaultQuery("searchUser", "")
	gender := c.DefaultQuery("gender", "")
	school := c.DefaultQuery("school", "")
	college := c.DefaultQuery("college", "")
	name := c.DefaultQuery("name", "")

	fmt.Println("limit:", limit, " curPage:", curPage)
	fmt.Println(c.Request.URL)
	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultCmsTeacher.DisplayTeacher(paginator, username, gender, school, college, name)
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

// AddTeacher
func (UserController) AddTeacher(c *gin.Context) {
	appG := app.Gin{C: c}
	var form models.TeacherForm

	err := c.ShouldBindJSON(&form)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultCmsTeacher.AddTeacher(form.Username, form.Password, form.Name, form.Gender, form.School, form.College, form.Avatar)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
}

// UpdateTeacher
func (UserController) UpdateTeacher(c *gin.Context) {
	appG := app.Gin{C: c}

	form := &models.TeacherForm{}

	err := c.ShouldBindJSON(&form)
	if err != nil {
		DPrintf("UpdateUser c.ShouldBindJSON err:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultCmsTeacher.UpdateTeacher(form.ID, form.Username, form.Password, form.Name, form.Gender, form.School, form.College, form.Avatar)
	if err != nil {
		DPrintf("UpdateUser logic.DefaultCmsUser.UpdateUser err:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("更新成功")
}

// DeleteTeacher
func (UserController) DeleteTeacher(c *gin.Context) {
	appG := app.Gin{C: c}

	request := models.UserDeleteId{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		DPrintf("DeleteUser c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultCmsTeacher.DeleteTeacher(&request.ID)
	if err != nil {
		DPrintf("DeleteUser logic.DefaultCmsUser.DeleteUser 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseNumber(count)
}

// GetTeacherCount
func (UserController) GetTeacherCount(c *gin.Context) {
	appG := app.Gin{C: c}

	count, err := logic.DefaultCmsTeacher.GetTeacherCount()

	if err != nil {
		DPrintf("GetUserCount 出错:", err)
		appG.ResponseErr(err.Error())
		return
	}
	data := make(map[string]interface{})
	data["count"] = count

	appG.ResponseSucMsg(data, "查询成功")
}

// =================================================管理员用户
// GetStudentUser
func (UserController) GetManager(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.DefaultQuery("page_size", "10")).MustInt()
	curPage := com.StrTo(c.DefaultQuery("page_number", "1")).MustInt()
	username := c.DefaultQuery("username", "")

	fmt.Println("limit:", limit, " curPage:", curPage)
	fmt.Println(c.Request.URL)
	paginator := NewPaginator(curPage, limit)

	data := make(map[string]interface{})
	list, total, err := logic.DefaultCmsManager.DisplayManager(paginator, username)
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
func (UserController) AddManager(c *gin.Context) {
	appG := app.Gin{C: c}
	var form models.NewManager

	err := c.ShouldBindJSON(&form)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultCmsManager.AddManager(form.Username, form.Password, form.ConfirmPassword)
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc()
}

// UpdateUser
func (UserController) UpdateManager(c *gin.Context) {
	appG := app.Gin{C: c}

	form := &models.ManagerUpdate{}

	err := c.ShouldBindJSON(&form)
	if err != nil {
		DPrintf("UpdateUser c.ShouldBindJSON err:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = logic.DefaultCmsManager.UpdateManager(form.ID, form.Username, form.Password, form.ConfirmPassword)
	if err != nil {
		DPrintf("UpdateUser logic.DefaultCmsUser.UpdateUser err:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSuc("更新成功")
}

// DeleteUser
func (UserController) DeleteManager(c *gin.Context) {
	appG := app.Gin{C: c}

	request := models.UserDeleteId{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		DPrintf("DeleteUser c.ShouldBindJSON 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	count, err := logic.DefaultCmsManager.DeleteManager(&request.ID)
	if err != nil {
		DPrintf("DeleteUser logic.DefaultCmsUser.DeleteUser 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseNumber(count)
}

// GetCount
func (UserController) GetManagerCount(c *gin.Context) {
	appG := app.Gin{C: c}

	count, err := logic.DefaultCmsManager.GetManagerCount()

	if err != nil {
		DPrintf("GetUserCount 出错:", err)
		appG.ResponseErr(err.Error())
		return
	}
	data := make(map[string]interface{})
	data["count"] = count

	appG.ResponseSucMsg(data, "查询成功")
}
