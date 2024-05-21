package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"server/logic"
	"server/utils/app"
	"server/utils/e"
	"strconv"
	"time"
)

type StatisticController struct{}

func (self StatisticController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/student_statistic", self.StudentStatistic)       //  学生首页
	g.GET("/teacher_statistic", self.TeacherStatistic)       // 教师首页首页
	g.GET("/department_statistic", self.DepartmentStatistic) // 系部管理员首页
	g.GET("/manager_statistic", self.ManagerStatistic)       // 系统管理员首页
	g.GET("/statistic_slice", self.StatisticSlice)           // 系部管理员数据统计
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

func (self StatisticController) StatisticSlice(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, exist := c.Get("user_id")
	if !exist {
		appG.ResponseErr("用户不存在")
		return
	}

	year := com.StrTo(c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))).MustInt()

	role, exist := c.Get("role")
	if !exist {
		appG.ResponseErr("权限错误")
		return
	}
	if role.(int) != e.CmsManagerRole && role.(int) != e.DepartmentRole {
		appG.ResponseErr("无权限")
		return
	}

	data, err := logic.DefaultStatisticLogic.StatisticSlice(userID.(int64), year)
	if err != nil {
		appG.ResponseErr("获取信息失败", err.Error())
		return
	}

	appG.ResponseSucMsg(data)
}
