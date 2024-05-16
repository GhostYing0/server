package cms

import (
	"errors"
	"fmt"
	. "server/database"
	. "server/logic"
	"server/logic/public"
	"server/models"
	"server/utils/logging"
	. "server/utils/mydebug"
	"time"
)

type CmsContestLogic struct{}

var DefaultCmsContest = CmsContestLogic{}

func (self CmsContestLogic) Display(paginator *Paginator, contest, contestType, startTime, endTime, contestEntry string, state, contestLevel, year, isGroup int) (*[]models.ContestReturn, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("CmsContestLogic Display session.Begin() 发生错误:", err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("CmsContestLogic Display session.Close() 发生错误:", err)
		}
	}()

	startYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	endYear := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")

	session.Join("LEFT", "teacher", "contest.teacher_id = teacher.teacher_id")
	session.Where("contest.create_time > ? and contest.create_time < ?", startYear, endYear)
	session.Join("LEFT", "account", "account.user_id = teacher.teacher_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("LEFT", "contest_level", "contest.contest_level_id = contest_level.contest_level_id")
	session.Join("LEFT", "school", "teacher.school_id = school.school_id")
	session.Join("LEFT", "college", "teacher.college_id = college.college_id")
	session.Join("LEFT", "department", "teacher.department_id = department.department_id")
	session.Join("LEFT", "contest_entry", "contest.contest_entry_id = contest_entry.contest_entry_id")

	searchContestType, err := public.SearchContestTypeByName(contestType)
	if err != nil {
		logging.L.Error(err)
	}

	if contest != "" {
		session.Where("contest.contest like ?", "%"+contest+"%")
	}
	if contestLevel > 0 {
		session.Where("contest.contest_level_id = ?", contestLevel)
	}
	if contestEntry != "" {
		session.Where("contest_entry.contest_entry like ?", "%"+contestEntry+"%")
	}
	if startTime != "" && endTime != "" {
		session.Where("contest.create_time > ? and contest.create_time < ?", startTime, endTime)
	}
	if contestType != "" {
		session.Where("contest.contest_type_id = ?", searchContestType.ContestTypeID)
	}
	if state != -1 {
		session.Where("contest.state = ?", state)
	}
	if isGroup > 0 {
		session.Where("contest.is_group = ?", isGroup)
	}

	data := &[]models.ContestContestTypeTeacher{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.ContestReturn, len(*data))
	for i := 0; i < len(list); i++ {
		list[i].Username = (*data)[i].Username
		list[i].Name = (*data)[i].Name
		list[i].School = (*data)[i].School
		list[i].College = (*data)[i].College
		list[i].ID = (*data)[i].Contest.ID
		list[i].State = (*data)[i].Contest.State
		list[i].Contest = (*data)[i].Contest.Contest
		list[i].ContestType = (*data)[i].Contest.ContestType
		list[i].Describe = (*data)[i].Describe
		list[i].ContestLevel = (*data)[i].ContestLevel
		list[i].ContestEntry = (*data)[i].ContestEntry
		list[i].Department = (*data)[i].Department
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].Contest.CreateTime)
		list[i].StartTime = models.MysqlFormatString2String((*data)[i].Contest.StartTime)
		list[i].Deadline = models.MysqlFormatString2String((*data)[i].Contest.Deadline)
	}

	return &list, total, session.Commit()
}

func (self CmsContestLogic) InsertContest(username, contest, contestType, startTime, deadline string, state int) error {
	if username == "" || contest == "" || contestType == "" || startTime == "" || deadline == "" {
		return errors.New("竞赛信息不能为空")
	}

	StartTime := models.FormatString2OftenTime(startTime)
	DeadlineTime := models.FormatString2OftenTime(deadline)

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("CmsContestLogic InsertContest session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("CmsContestLogic InsertContest session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	account, err := public.SearchAccountByUsernameAndRole(username, 2)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	teacher, err := public.SearchTeacherByID(account.UserID)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	searchContestType, err := public.SearchContestTypeByName(contestType)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	has, err := session.Table("contest").Where("contest = ? and contest_type_id = ?", contest, searchContestType.ContestTypeID).Exist()
	if err != nil {
		fmt.Println("InsertContestInfo Exist error:", err)
		logging.L.Error(err)
		return err
	}
	if has {
		logging.L.Error("竞赛已存在")
		return err
	}

	NewContest := &models.NewContest{
		Contest:     contest,
		TeacherID:   account.UserID,
		SchoolID:    teacher.SchoolID,
		CollegeID:   teacher.CollegeID,
		ContestType: searchContestType.ContestTypeID,
		CreateTime:  models.NewOftenTime(),
		StartTime:   StartTime,
		Deadline:    DeadlineTime,
		State:       state,
	}

	_, err = session.Table("contest").Insert(NewContest)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(fail)
			return fail
		}
		fmt.Println("InsertContestInfo Insert error:", err)
		logging.L.Error(err)
		return err
	}
	return session.Commit()
}

func (self CmsContestLogic) UpdateContest(id int64, username, contest, contestType, startTime, deadline string, state int) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("CmsContestLogic UpdateContest session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("CmsContestLogic UpdateContest session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	has, err := session.Table("contest").Where("id = ?", id).Exist()
	if err != nil {
		fmt.Println("UpdateContestInfo Exist error:", err)
		logging.L.Error(err)
		return err
	}
	if !has {
		logging.L.Error("竞赛不存在")
		return errors.New("竞赛不存在")
	}

	account, err := public.SearchAccountByUsernameAndRole(username, 2)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	exist, err := session.Table("contest").Where("contest = ?", contest).Exist()
	if exist {
		logging.L.Error("已有同名竞赛")
		return errors.New("已有同名竞赛")
	}
	if err != nil {
		logging.L.Error(err)
		return err
	}

	searchType := &models.ContestType{}
	exist, err = session.Where("type = ?", contestType).Get(searchType)
	if err != nil {
		DPrintf("UpdateContest查询竞赛类型:", err)
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("竞赛类型不存在")
		return errors.New("竞赛类型不存在")
	}

	_, err = session.Where("id = ?", id).Update(&models.ContestInfo{
		TeacherID:   account.UserID,
		Contest:     contest,
		ContestType: searchType.ContestTypeID,
		StartTime:   models.FormatString2OftenTime(startTime),
		Deadline:    models.FormatString2OftenTime(deadline),
		State:       state,
	})
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(fail)
			return fail
		}
		logging.L.Error(err)
		return err
	}

	return session.Commit()
}

func (self CmsContestLogic) DeleteContest(ids *[]int64) (string, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("CmsContestLogic UpdateContest session.Begin() 发生错误:", err)
		return "", 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("CmsContestLogic UpdateContest session.Close() 发生错误:", err)
		}
	}()
	var count int64

	for _, id := range *ids {
		var contest models.ContestInfo
		if id < 1 {
			fmt.Println("非法id")
			continue
		}
		contest.ID = id
		affected, err := session.Delete(&contest)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return "", 0, fail
			}
			return "操作出错", 0, err
		}
		if affected > 0 {
			count += affected
		}
	}

	return "操作成功", 0, session.Commit()
}

func (self CmsContestLogic) GetContestCount() (int64, error) {
	session := MasterDB.NewSession()

	count, err := session.Table("contest").Count()
	if err != nil {
		DPrintf("GetContestCount Count 发生错误:", err)
		return count, err
	}
	return count, err
}

func (self CmsContestLogic) ProcessContest(id int64, state int) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("CmsContestLogic UpdateContest session.Begin() 发生错误:", err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("CmsContestLogic UpdateContest session.Close() 发生错误:", err)
		}
	}()

	exist, err := session.Table("contest").Where("id = ?", id).Exist()
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("竞赛不存在")
		return errors.New("竞赛不存在")
	}

	_, err = session.Where("id = ?", id).Update(&models.ContestInfo{State: state})
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			logging.L.Error(fail)
			return fail
		}
		logging.L.Error(err)
		return err
	}

	return session.Commit()
}
