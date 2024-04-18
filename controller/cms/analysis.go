package cms

import (
	"errors"
	. "server/database"
	_ "server/database"
	logic "server/logic/cms"
	"server/utils/app"
	"server/utils/logging"

	"github.com/gin-gonic/gin"
)

type AnalysisController struct{}

// RegisterRoutes
func (self AnalysisController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/totalEnrollCountOfPerYear", self.TotalEnrollCountOfPerYear)     // 查看最近五年总报名数量
	g.GET("/preTypeEnrollCountOfPerYear", self.PreTypeEnrollCountOfPerYear) // 查看今年各类竞赛报名数
	g.GET("/compareEnrollCount", self.CompareEnrollCount)                   // 今年与往年报名数对比
	g.GET("/schoolEnrollCount", self.SchoolEnrollCount)
}

func CheckManager(id int64) error {
	exist, err := MasterDB.Table("cms_account").Where("id = ?", id).Exist()
	if !exist {
		return errors.New("管理员不存在")
	}
	return err
}

func (self AnalysisController) TotalEnrollCountOfPerYear(c *gin.Context) {
	appG := app.Gin{C: c}

	Manager, exist := c.Get("user_id")
	if !exist {
		logging.L.Error("管理员不存在")
		appG.ResponseErr("管理员不存在")
		return
	}

	err := CheckManager(Manager.(int64))
	if err != nil {
		logging.L.Error(err)
		appG.ResponseErr(err.Error())
		return
	}

	data, err := logic.DefaultCmsAnalysis.GetTotalEnrollCountOfPerYear()
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg(data)
}

func (self AnalysisController) PreTypeEnrollCountOfPerYear(c *gin.Context) {
	appG := app.Gin{C: c}

	Manager, exist := c.Get("user_id")
	if !exist {
		logging.L.Error("管理员不存在")
		appG.ResponseErr("管理员不存在")
		return
	}

	err := CheckManager(Manager.(int64))
	if err != nil {
		logging.L.Error(err)
		appG.ResponseErr(err.Error())
		return
	}

	data, err := logic.DefaultCmsAnalysis.GetPreTypeEnrollCountOfPerYear()
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg(data)
}

func (self AnalysisController) CompareEnrollCount(c *gin.Context) {
	appG := app.Gin{C: c}

	Manager, exist := c.Get("user_id")
	if !exist {
		logging.L.Error("管理员不存在")
		appG.ResponseErr("管理员不存在")
		return
	}

	err := CheckManager(Manager.(int64))
	if err != nil {
		logging.L.Error(err)
		appG.ResponseErr(err.Error())
		return
	}

	data, err := logic.DefaultCmsAnalysis.CompareEnrollCount()
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg(data)
}

func (self AnalysisController) SchoolEnrollCount(c *gin.Context) {
	appG := app.Gin{C: c}

	Manager, exist := c.Get("user_id")
	if !exist {
		logging.L.Error("管理员不存在")
		appG.ResponseErr("管理员不存在")
		return
	}

	err := CheckManager(Manager.(int64))
	if err != nil {
		logging.L.Error(err)
		appG.ResponseErr(err.Error())
		return
	}

	data, err := logic.DefaultCmsAnalysis.SchoolEnrollCount()
	if err != nil {
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg(data)
}
