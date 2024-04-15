package cms

import (
	. "server/database"
	"server/models"
	"server/utils/logging"
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

func (CmsAnalysisLogic) GetPreTypeEnrollCountOfPerYear() (*models.PreTypeEnrollCountOfPerYear, error) {
	nowDate := time.Now()

	endYear := time.Date(nowDate.Year(), nowDate.Month(), nowDate.Day(), 0, 0, 0, 0, time.Local)
	startYear := time.Date(nowDate.Year()-6, time.December, 31, 0, 0, 0, 0, time.Local)

	allContestType := &[]models.ContestType{}
	_, err := MasterDB.Table("contest_type").FindAndCount(allContestType)

	yearMap := make([]string, 5)
	for i := 0; i < 5; i++ {
		yearMap[i] = strconv.Itoa(startYear.Year() + 2 + i)
	}

	typeMap := make(map[int64]string)
	for i := 0; i < len(*allContestType); i++ {
		typeMap[(*allContestType)[i].ContestTypeID] = (*allContestType)[i].ContestType
	}

	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	list := []models.MysqlSelectEnrollYearAndContestType{}
	_, err = MasterDB.
		Table("enroll_information").
		Join("LEFT", "contest", "contest.id = enroll_information.contest_id").
		Where("enroll_information.create_time > ? and enroll_information.create_time < ?", startYear, endYear).
		FindAndCount(&list)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	data := &models.PreTypeEnrollCountOfPerYear{}
	data.EnrollData = make(map[string]map[string]int64, 5)

	data.NowYear = strconv.Itoa(startYear.Year())
	for i := int64(0); i < 5; i++ {
		data.EnrollData[yearMap[i]] = make(map[string]int64)
		for _, value := range typeMap {
			data.EnrollData[yearMap[i]][value] = 0
		}
	}

	for i := 0; i < int(len(list)); i++ {
		//data.EnrollData[strconv.Itoa(int(list[i].Date.Year()))[typeMap[list[i].ContestType]] += 1
	}
	return data, err
}
