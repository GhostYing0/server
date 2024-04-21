package cms

import (
	"fmt"
	. "server/database"
	"server/logic/public"
	"server/models"
	. "server/utils/e"
	"server/utils/logging"
	"sort"
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
		Where("create_time > ? and create_time < ? and state = ?", startYear, endYear, Pass).
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
		Where("enroll_information.create_time > ? and enroll_information.create_time < ? and enroll_information.state = ?", startYear, endYear, Pass).
		FindAndCount(&list)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	data := &models.PreTypeEnrollCountOfPerYear{}
	data.EnrollData = make(map[string]map[string]int64, 5)

	for i := int64(0); i < 5; i++ {
		data.EnrollData[yearMap[i]] = make(map[string]int64)
		for _, value := range typeMap {
			data.EnrollData[yearMap[i]][value] = 0
		}
	}

	for i := 0; i < int(len(list)); i++ {
		data.EnrollData[strconv.Itoa(int(list[i].Date.Year()))][typeMap[list[i].ContestType]] += 1
	}
	return data, err
}

func (CmsAnalysisLogic) CompareEnrollCount() (*models.CompareEnrollCount, error) {
	nowDate := time.Now()

	curYear := time.Date(nowDate.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
	lastYear := time.Date(nowDate.Year()-1, time.January, 1, 0, 0, 0, 0, time.Local)

	allContestType := &[]models.ContestType{}
	_, err := MasterDB.Table("contest_type").FindAndCount(allContestType)

	typeMap := make(map[int64]string)
	for i := 0; i < len(*allContestType); i++ {
		typeMap[(*allContestType)[i].ContestTypeID] = (*allContestType)[i].ContestType
	}

	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	compareData := make([][]int64, 2)
	for i := 0; i < len(compareData); i++ {
		compareData[i] = make([]int64, len(typeMap))
	}

	index := 0
	for key, _ := range typeMap {
		// 统计今年的
		compareData[0][index], err = MasterDB.
			Table("enroll_information").
			Join("LEFT", "contest", "contest.id = enroll_information.contest_id").
			Where("enroll_information.create_time > ? and contest.contest_type_id = ? and enroll_information.state = ?", curYear, key, Pass).
			Count()
		if err != nil {
			logging.L.Error(err)
			return nil, err
		}

		// 统计去年的
		compareData[1][index], err = MasterDB.
			Table("enroll_information").
			Join("LEFT", "contest", "contest.id = enroll_information.contest_id").
			Where("enroll_information.create_time > ? and enroll_information.create_time < ? and contest.contest_type_id = ?", lastYear, curYear, key).
			Count()
		if err != nil {
			logging.L.Error(err)
			return nil, err
		}
		index++
	}

	data := &models.CompareEnrollCount{}
	data.EnrollCompare = make(map[string]float64)

	curr_sum := float64(0)
	prev_sum := float64(0)
	for i := 0; i < len(compareData[0]); i++ {
		curr := float64(compareData[0][i])
		prev := float64(compareData[1][i])
		rate := float64(0)
		if prev != 0 {
			rate = (curr - prev) / prev
		}
		data.EnrollCompare[typeMap[int64(i)+1]], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", rate), 64)
		curr_sum += curr
		prev_sum += prev
	}
	if prev_sum != 0 {
		rate := (curr_sum - prev_sum) / prev_sum
		data.EnrollCompare["总共"], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", rate), 64)
	} else {
		data.EnrollCompare["总共"] = 0
	}
	return data, err
}

func (CmsAnalysisLogic) SchoolEnrollCount(year int) (*models.SchoolEnrollSortedCount, error) {

	curYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	nextYear := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local)

	allSchool := &[]models.School{}
	_, err := MasterDB.Table("school").FindAndCount(allSchool)

	schoolMap := make(map[int64]string)
	for i := 0; i < len(*allSchool); i++ {
		schoolMap[(*allSchool)[i].SchoolID] = (*allSchool)[i].School
	}

	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	list := &[]models.SchoolEnrollCount{}

	// 统计今年的
	_, err = MasterDB.
		Table("enroll_information").
		Where("create_time > ? and create_time < ? and state = ?", curYear, nextYear, Pass).
		FindAndCount(list)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	temp := make(map[string]int64)
	for _, value := range schoolMap {
		temp[value] = 0
	}
	for i := 0; i < len(*list); i++ {
		temp[schoolMap[(*list)[i].SchoolID]] += 1
	}

	array := make([]models.SchoolEnroll, 0)
	for key, value := range temp {
		array = append(array, models.SchoolEnroll{key, value})
	}
	sort.Slice(array, func(i, j int) bool {
		return array[i].EnrollCount > array[j].EnrollCount
	})

	data := &models.SchoolEnrollSortedCount{SchoolEnrollData: array[:10]}
	return data, err
}

func (CmsAnalysisLogic) StudentContestSemester(contest string) (*models.EnrollSemesterArray, error) {
	contestSearch, err := public.SearchContestByName(contest)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	semesterMap := make(map[int64]string)

	semester := &[]models.Semester{}
	_, err = MasterDB.Table("semester").FindAndCount(semester)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	for i := 0; i < len(*semester); i++ {
		semesterMap[(*semester)[i].SemesterID] = (*semester)[i].Semester
	}

	list := &[]models.EnrollAndSemester{}
	total, err := MasterDB.
		Table("enroll_information").
		Where("contest_id = ? and state = ?", contestSearch.ID, Pass).
		Join("LEFT", "student", "enroll_information.student_id = student.student_id").
		FindAndCount(list)

	temp := make(map[int64]int64) // k:semester_id v:count
	for i := 0; i < len(*list); i++ {
		temp[(*list)[i].Semester] += 1
	}

	data := models.EnrollSemesterArray{}
	data.Total = total
	for k, v := range temp {
		data.Data = append(data.Data, models.EnrollSemester{semesterMap[k], v})
	}

	return &data, nil
}

func (CmsAnalysisLogic) StudentRewardRate(contest string) (*models.RewardRate, error) {
	contestSearch, err := public.SearchContestByName(contest)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	rewardCount, err := MasterDB.
		Table("grade").
		Where("contest_id = ? and state = ?", contestSearch.ID, Pass).
		Count()
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	fmt.Println(rewardCount)

	enrollCount, err := MasterDB.
		Table("enroll_information").
		Where("contest_id = ? and state = ?", contestSearch.ID, Pass).
		Count()
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	data := models.RewardRate{
		RewardCount: rewardCount,
		EnrollCount: enrollCount,
	}

	data.Rate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(rewardCount)/float64(enrollCount)), 64)
	return &data, err
}
