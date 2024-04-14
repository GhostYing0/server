package cms

import (
	. "server/database"
	"server/models"
	"strconv"
	"time"
)

type CmsAnalysisLogic struct{}

var DefaultCmsAnalysis = CmsAnalysisLogic{}

func (CmsAnalysisLogic) GetTotalEnrollCountOfPerYear() (*models.TotalEnrollCountOfPerYear, error) {
	nowDate := time.Now()

	endYear := time.Date(nowDate.Year(), nowDate.Month(), nowDate.Day(), 0, 0, 0, 0, time.Local)
	startYear := time.Date(nowDate.Year()-6, time.December, 31, 0, 0, 0, 0, time.Local)

	list := []models.MysqlSelectEnrollYear{}
	_, err := MasterDB.
		Table("enroll_information").
		Where("create_time > ? and create_time < ?", startYear, endYear).
		FindAndCount(&list)
	if err != nil {
		return nil, err
	}

	data := &models.TotalEnrollCountOfPerYear{}
	data.EnrollData = make(map[string]int64)
	for i := 0; i < int(len(list)); i++ {
		year := strconv.Itoa(list[i].Date.Year())
		data.EnrollData[year] += 1
	}
	return data, err
}
