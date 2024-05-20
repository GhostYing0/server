package controller

import (
	"github.com/gin-gonic/gin"
	"server/logic"
	"server/utils/app"
	"server/utils/e"
)

type StatisticController struct{}

func (self StatisticController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/student_statistic", self.StudentStatistic)       // 查看成绩
	g.GET("/teacher_statistic", self.TeacherStatistic)       // 上传成绩
	g.GET("/department_statistic", self.DepartmentStatistic) // 上传成绩
	g.GET("/manager_statistic", self.ManagerStatistic)       // 上传成绩
}

func (self StatisticController) StudentStatistic(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("用户不存在")
		return
	}

	role, exist := c.Get("role")
	if !exist {
		appG.ResponseErr("权限错误")
		return
	}
	if role.(int) != e.StudentRole {
		appG.ResponseErr("无权限")
		return
	}

	data, err := logic.DefaultStatisticLogic.StudentStatistic(userID.(int64))
	if err != nil {
		appG.ResponseErr("获取信息失败", err.Error())
		return
	}

	appG.ResponseSucMsg(data)
}

func (self StatisticController) TeacherStatistic(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("用户不存在")
		return
	}

	role, exist := c.Get("role")
	if !exist {
		appG.ResponseErr("权限错误")
		return
	}
	if role.(int) != e.TeacherRole {
		appG.ResponseErr("无权限")
		return
	}

	data, err := logic.DefaultStatisticLogic.TeacherStatistic(userID.(int64))
	if err != nil {
		appG.ResponseErr("获取信息失败", err.Error())
		return
	}

	appG.ResponseSucMsg(data)
}

func (self StatisticController) DepartmentStatistic(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("用户不存在")
		return
	}

	role, exist := c.Get("role")
	if !exist {
		appG.ResponseErr("权限错误")
		return
	}
	if role.(int) != e.DepartmentRole {
		appG.ResponseErr("无权限")
		return
	}

	data, err := logic.DefaultStatisticLogic.DepartmentStatistic(userID.(int64))
	if err != nil {
		appG.ResponseErr("获取信息失败", err.Error())
		return
	}

	appG.ResponseSucMsg(data)
}

func (self StatisticController) ManagerStatistic(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("用户不存在")
		return
	}

	role, exist := c.Get("role")
	if !exist {
		appG.ResponseErr("权限错误")
		return
	}
	if role.(int) != e.CmsManagerRole {
		appG.ResponseErr("无权限")
		return
	}

	data, err := logic.DefaultStatisticLogic.ManagerStatistic(userID.(int64))
	if err != nil {
		appG.ResponseErr("获取信息失败", err.Error())
		return
	}

	appG.ResponseSucMsg(data)
}
