package logic

import (
	"errors"
	. "server/database"
	"server/logic/public"
	"server/models"
	"server/utils/e"
	"server/utils/logging"
	. "server/utils/mydebug"
	"time"
)

type ContestLogic struct{}

var DefaultContestLogic = ContestLogic{}

func (self ContestLogic) DisplayContest(paginator *Paginator, contest, contestType string, contestID, userID int64, contestLevel, isGroup, year, role int) (*[]models.ContestReturn, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("DisplayContest session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("DisplayContest session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	startTime := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	endTime := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")

	if role == e.DepartmentRole {
		account, err := public.SearchDepartmentManagerByID(userID)
		if err != nil {
			return nil, 0, err
		}
		session.Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id")
		session.Where("contest.state = 1")
		session.Where("contest.school_id = ?", account.SchoolID)
		session.Where("contest.start_time > ? and contest.start_time < ?", startTime, endTime)
		session.Join("LEFT", "contest_level", "contest.contest_level_id = contest_level.contest_level_id")
	} else {
		account, err := public.SearchAccountByID(userID)
		if err != nil {
			return nil, 0, err
		}

		if account.Role == e.StudentRole {
			student, err := public.SearchStudentByID(account.UserID)
			if err != nil {
				return nil, 0, err
			}

			session.Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id")
			session.Where("contest.state = 1")
			session.Where("contest.school_id = ?", student.SchoolID)
			session.Where("contest.start_time > ? and contest.start_time < ?", startTime, endTime)
			session.Join("LEFT", "contest_level", "contest.contest_level_id = contest_level.contest_level_id")
		} else if account.Role == e.TeacherRole {
			teacher, err := public.SearchTeacherByID(account.UserID)
			if err != nil {
				return nil, 0, err
			}

			session.Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id")
			session.Where("contest.state = 1")
			session.Where("contest.school_id = ?", teacher.SchoolID)
			session.Where("contest.start_time > ? and contest.start_time < ?", startTime, endTime)
			session.Join("LEFT", "contest_level", "contest.contest_level_id = contest_level.contest_level_id")
		}
	}

	if contest != "" {
		session.Where("contest like ?", "%"+contest+"%")
	}
	if contestType != "" {
		searchContestType, err := public.SearchContestTypeByName(contestType)
		if err != nil {

		} else {
			session.Where("contest.contest_type_id = ?", searchContestType.ContestTypeID)
		}
	}
	if contestLevel > 0 {
		session.Where("contest.contest_level_id = ?", contestLevel)
	}
	if isGroup > 0 {
		session.Where("contest.is_group = ?", isGroup)
	}

	data := &[]models.ContestInfoType{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		DPrintf("InsertEnrollInformation 查询竞赛发生错误: ", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.ContestReturn, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].ID
		list[i].State = (*data)[i].State
		list[i].Contest = (*data)[i].Contest
		list[i].ContestType = (*data)[i].ContestType.ContestType
		list[i].CreateTime = (*data)[i].CreateTime.String()
		list[i].StartTime = (*data)[i].StartTime.String()
		list[i].Deadline = (*data)[i].Deadline.String()
		list[i].Describe = (*data)[i].Describe
		list[i].ContestLevel = (*data)[i].ContestLevel
		list[i].StartTime = (*data)[i].StartTime.String()
		list[i].EnrollTime = (*data)[i].EnrollTime.String()
		// 竞赛可报名条件，审核通过，在报名截至时间之前，且教师未关闭报名
		if (*data)[i].State == e.Pass && (*data)[i].ContestState == e.EnrollOpen && models.NewOftenTime().Before(&(*data)[i].Deadline) {
			list[i].ContestState = e.EnrollOpen
		} else {
			list[i].ContestState = e.EnrollClose
		}
	}

	return &list, total, session.Commit()
}

func (self ContestLogic) GetContestDetail(userID, contestID int64) (*[]models.ContestDetail, int64, error) {
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

	departmentAccount, err := public.SearchDepartmentManagerByID(userID)
	if err != nil {
		logging.L.Error()
		return nil, 0, err
	}

	session.Table("teacher").Where("teacher.school_id = ? and teacher.college_id = ? and teacher.department_id = ?", departmentAccount.SchoolID, departmentAccount.CollegeID, departmentAccount.DepartmentID)
	session.Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id")
	session.Where("contest.id = ?", contestID)
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("LEFT", "contest_level", "contest_level.contest_level_id = contest.contest_level_id")
	session.Join("LEFT", "department", "department.department_id = teacher.department_id")
	session.Join("LEFT", "contest_entry", "contest_entry.contest_entry_id = contest.contest_entry_id")

	data := &[]models.ContestDetail{}

	total, err := session.FindAndCount(data)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}

	for i := 0; i < len(*data); i++ {
		(*data)[i].StartTime = models.MysqlFormatString2String((*data)[i].StartTime)
		(*data)[i].EnrollTime = models.MysqlFormatString2String((*data)[i].EnrollTime)
		(*data)[i].Deadline = models.MysqlFormatString2String((*data)[i].Deadline)
	}

	return data, total, session.Commit()
}

func (self ContestLogic) StudentGetOneContest(paginator *Paginator, contestID, userID int64, isGroup int) (*[]models.ContestReturn, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("DisplayContest session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("DisplayContest session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		return nil, 0, err
	}

	if contestID <= 0 {
		return nil, 0, err
	}

	if account.Role == e.StudentRole {
		student, err := public.SearchStudentByID(account.UserID)
		if err != nil {
			return nil, 0, err
		}

		session.Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id")
		session.Where("contest.state = 1")
		session.Where("contest.id = ?", contestID)
		session.Where("contest.school_id = ?", student.SchoolID)
		session.Join("LEFT", "contest_level", "contest.contest_level_id = contest_level.contest_level_id")
	} else if account.Role == e.TeacherRole {
		teacher, err := public.SearchTeacherByID(account.UserID)
		if err != nil {
			return nil, 0, err
		}

		session.Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id")
		session.Where("contest.state = 1")
		session.Where("contest.school_id = ?", teacher.SchoolID)
	}
	//if isGroup > 0 {
	//	session.Where("contest.is_group = ?", isGroup)
	//}

	data := &[]models.ContestInfoType{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		DPrintf("InsertEnrollInformation 查询竞赛发生错误: ", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.ContestReturn, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].ID
		list[i].State = (*data)[i].State
		list[i].Contest = (*data)[i].Contest
		list[i].ContestType = (*data)[i].ContestType.ContestType
		list[i].CreateTime = (*data)[i].CreateTime.String()
		list[i].StartTime = (*data)[i].StartTime.String()
		list[i].Deadline = (*data)[i].Deadline.String()
		list[i].Describe = (*data)[i].Describe
		list[i].ContestLevel = (*data)[i].ContestLevel
		list[i].StartTime = (*data)[i].StartTime.String()
		list[i].EnrollTime = (*data)[i].EnrollTime.String()
		list[i].IsGroup = (*data)[i].IsGroup
		// 竞赛可报名条件，审核通过，在报名截至时间之前，且教师未关闭报名
		if (*data)[i].State == e.Pass && (*data)[i].ContestState == e.EnrollOpen && models.NewOftenTime().Before(&(*data)[i].Deadline) {
			list[i].ContestState = e.EnrollOpen
		} else {
			list[i].ContestState = e.EnrollClose
		}
	}

	return &list, total, session.Commit()
}

func (self ContestLogic) ViewTeacherContest(paginator *Paginator, contestID, contestLevel, userID int64, contest, contestType string, state, year, isGroup int) (*[]models.ContestReturn, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("DisplayContest session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("DisplayContest session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}

	startTime := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	endTime := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")

	session.Table("contest")
	session.Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id")
	session.Join("LEFT", "contest_level", "contest.contest_level_id = contest_level.contest_level_id")
	session.Where("contest.teacher_id = ?", account.UserID)
	session.Where("contest.start_time > ? and contest.start_time < ?", startTime, endTime)
	if isGroup > 0 {
		session.Where("contest.is_group = ?", isGroup)
	}
	if contestLevel > 0 {
		session.Where("contest.contest_level_id = ?", contestLevel)
	}

	if state != -1 {
		session.Where("contest.state = ?", state)
	}
	if contest != "" {
		session.Where("contest like ?", "%"+contest+"%")
	}
	if contestType != "" {
		searchContestType, err := public.SearchContestTypeByName(contestType)
		if err != nil {

		} else {
			session.Where("contest.contest_type_id = ?", searchContestType.ContestTypeID)
		}
	}
	if contestID > 0 {
		session.Where("contest.id = ?", contestID)
	}

	data := &[]models.ContestInfoType{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		DPrintf("InsertEnrollInformation 查询竞赛发生错误: ", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.ContestReturn, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].ID
		list[i].State = (*data)[i].State
		list[i].Contest = (*data)[i].Contest
		list[i].ContestType = (*data)[i].ContestType.ContestType
		list[i].CreateTime = (*data)[i].CreateTime.String()
		list[i].StartTime = (*data)[i].StartTime.String()
		list[i].Deadline = (*data)[i].Deadline.String()
		list[i].Describe = (*data)[i].Describe
		list[i].RejectReason = (*data)[i].RejectReason
		list[i].Prize1Count = (*data)[i].Prize1Count
		list[i].Prize2Count = (*data)[i].Prize2Count
		list[i].Prize3Count = (*data)[i].Prize3Count
		list[i].Prize4Count = (*data)[i].Prize4Count
		list[i].ContestLevelID = (*data)[i].ContestLevelID
		list[i].IsGroup = (*data)[i].IsGroup
		list[i].ContestLevel = (*data)[i].ContestLevel
		list[i].ContestEntry = (*data)[i].ContestEntry
		list[i].MaxGroupNumber = (*data)[i].MaxGroupNumber
		list[i].EnrollTime = (*data)[i].EnrollTime.String()
		// 竞赛可报名条件，审核通过，在报名截至时间之前，且教师未关闭报名
		if (*data)[i].State == e.Pass && (*data)[i].ContestState == e.EnrollOpen && models.NewOftenTime().Before(&(*data)[i].Deadline) {
			list[i].ContestState = e.EnrollOpen
		} else {
			list[i].ContestState = e.EnrollClose
		}
	}

	return &list, total, session.Commit()
}

func (self ContestLogic) ViewTeacherContestGrade(paginator *Paginator, userID int64, contest, contestType string, state, isGroup, contestLevel, year int) (*[]models.TeacherUploadGradeContestReturn, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("DisplayContest session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("DisplayContest session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}

	startTime := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	endTime := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")

	session.Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id")
	session.Join("LEFT", "contest_level", "contest.contest_level_id = contest_level.contest_level_id")
	session.Where("contest.teacher_id = ?", account.UserID)
	session.Where("contest.state = ?", e.Pass)
	session.Where("contest.start_time > ? and contest.start_time < ?", startTime, endTime)
	session.Where("contest.is_group = ?", isGroup)

	if contest != "" {
		session.Where("contest = ?", contest)
	}
	if contestType != "" {
		searchContestType, err := public.SearchContestTypeByName(contestType)
		if err != nil {
			logging.L.Error()
			return nil, 0, err
		} else {
			session.Where("contest.contest_type_id = ?", searchContestType.ContestTypeID)
		}
	}

	if contestLevel > 0 {
		session.Where("contest.contest_level_id = ?", contestLevel)
	}

	data := &[]models.ContestInfoType{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		DPrintf("InsertEnrollInformation 查询竞赛发生错误: ", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.TeacherUploadGradeContestReturn, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].ID
		list[i].State = (*data)[i].State
		list[i].Contest = (*data)[i].Contest
		list[i].ContestType = (*data)[i].ContestType.ContestType
		list[i].CreateTime = (*data)[i].CreateTime.String()
		list[i].StartTime = (*data)[i].StartTime.String()
		list[i].Deadline = (*data)[i].Deadline.String()
		list[i].Describe = (*data)[i].Describe
		list[i].RejectReason = (*data)[i].RejectReason
		list[i].ContestLevel = (*data)[i].ContestLevel
		list[i].Prize1Count = (*data)[i].Prize1Count
		list[i].Prize2Count = (*data)[i].Prize2Count
		list[i].Prize3Count = (*data)[i].Prize3Count
		list[i].Prize4Count = (*data)[i].Prize4Count
		list[i].RewardCount, _ = session.Table("grade").Where("state = ? and contest_id = ?", e.Pass, list[i].ID).Count()
		list[i].EnrollCount, _ = session.Table("enroll_information").Where("state = ? and contest_id = ?", e.Pass, list[i].ID).Count()
		list[i].ProcessCount, _ = session.Table("grade").Where("state = ? and contest_id = ?", e.Processing, list[i].ID).Count()
		list[i].RejectedCount, _ = session.Table("grade").Where("state = ? and contest_id = ?", e.Reject, list[i].ID).Count()
		// 竞赛可报名条件，审核通过，在报名截至时间之前，且教师未关闭报名
		if (*data)[i].State == e.Pass && (*data)[i].ContestState == e.EnrollOpen && models.NewOftenTime().Before(&(*data)[i].Deadline) {
			list[i].ContestState = e.EnrollOpen
		} else {
			list[i].ContestState = e.EnrollClose
		}
	}

	return &list, total, session.Commit()
}

func (self ContestLogic) UpdateContest(userID int64, form *models.UpdateContestForm) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("ProcessEnroll session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			logging.L.Error(err)
			DPrintf("ProcessEnroll session.Close() 发生错误:", err)
		}
	}()

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	session.Table("contest")
	exist, err := session.Where("teacher_id = ? and id = ?", account.UserID, form.ID).Exist()
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("竞赛不存在")
		return errors.New("竞赛不存在")
	}

	searchContestType, err := public.SearchContestTypeByName(form.ContestType)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	exist, err = session.Table("contest").Where("id != ? and contest = ? and contest_type_id = ? and contest_entry_id = ?", form.ID, form.Contest, searchContestType.ContestTypeID, form.ContestEntry).Exist()
	if exist {
		logging.L.Error("已有同名竞赛")
		return errors.New("已有同名竞赛")
	}
	if err != nil {
		logging.L.Error(err)
		return err
	}

	updateContest := &models.ContestInfo{
		Contest:        form.Contest,
		ContestType:    searchContestType.ContestTypeID,
		ContestState:   form.ContestState,
		ContestEntry:   form.ContestEntry,
		IsGroup:        form.IsGroup,
		MaxGroupNumber: form.MaxGroupNumber,
		ContestLevelID: form.ContestLevelID,
		Prize1Count:    int64(form.Prize1Count),
		Prize2Count:    int64(form.Prize2Count),
		Prize3Count:    int64(form.Prize3Count),
		Prize4Count:    int64(form.Prize4Count),
		State:          e.Processing,
		Describe:       form.Describe,
	}

	if form.StartTime != "" {
		updateContest.StartTime = models.FormatString2OftenTime(form.StartTime)
	}
	if form.Deadline != "" {
		updateContest.Deadline = models.FormatString2OftenTime(form.Deadline)
	}
	if form.EnrollTime != "" {
		updateContest.EnrollTime = models.FormatString2OftenTime(form.EnrollTime)
	}
	_, err = session.Where("id = ?", form.ID).Update(updateContest)
	if err != nil {
		fail := session.Rollback()
		if err != nil {
			logging.L.Error(err)
			return fail
		}
		logging.L.Error(err)
		return err
	}
	return session.Commit()
}

func (self ContestLogic) UploadContest(userId int64, form *models.TeacherUploadContestForm) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("UploadContest session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			logging.L.Error(err)
			DPrintf("UploadContest session.Close() 发生错误:", err)
		}
	}()

	account, err := public.SearchAccountByID(userId)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	searchContestType, err := public.SearchContestTypeByName(form.ContestType)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	teacher, err := public.SearchTeacherByID(account.UserID)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	contestEntry, err := public.SearchContestEntryByID(form.ContestEntry)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	contestStartTime, err := time.Parse("2006-01-02 15:04:05", form.StartTime)
	startTime := time.Now()
	endTime := time.Now()
	switch contestEntry.Cycle {
	case 1:
		//每年举办一次
		startTime = time.Date(contestStartTime.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
		endTime = time.Date(contestStartTime.Year()+1, time.January, 1, 0, 0, 0, 0, time.Local)
		exist, err := MasterDB.Table("contest").
			Where("contest_entry_id = ? and contest_type_id = ? and start_time > ? and start_time < ?", form.ContestEntry, searchContestType.ContestTypeID, startTime, endTime).
			Exist()
		if err != nil {
			logging.L.Error(err)
			return err
		}
		if exist {
			logging.L.Error("已有同名竞赛")
			return errors.New("已有同名竞赛")
		}
	case 2:
		//每半年举办一次
		startMouth := models.FormatString2OftenTime(form.StartTime).OftenTimeMonth()
		num := (startMouth / 7) * 6
		startTime = time.Date(contestStartTime.Year(), time.January+num, 1, 0, 0, 0, 0, time.Local)
		if num == 0 {
			endTime = time.Date(contestStartTime.Year(), time.July, 1, 0, 0, 0, 0, time.Local)
		} else {
			endTime = time.Date(contestStartTime.Year()+1, time.January, 1, 0, 0, 0, 0, time.Local)
		}
		exist, err := MasterDB.Table("contest").
			Where("contest_entry_id = ? and contest_type_id = ? and start_time > ? and start_time < ?", form.ContestEntry, searchContestType.ContestTypeID, startTime, endTime).
			Exist()
		if err != nil {
			logging.L.Error(err)
			return err
		}
		if exist {
			logging.L.Error("已有同名竞赛")
			return errors.New("已有同名竞赛")
		}
	case 3:
		//每季度举办一次
		startMouth := models.FormatString2OftenTime(form.StartTime).OftenTimeMonth()
		num := (startMouth / 4) * 3
		startTime = time.Date(contestStartTime.Year(), time.January+num, 1, 0, 0, 0, 0, time.Local)
		if num == 0 {
			endTime = time.Date(contestStartTime.Year(), 4, 1, 0, 0, 0, 0, time.Local)
		} else if num == 1 {
			endTime = time.Date(contestStartTime.Year(), 7, 1, 0, 0, 0, 0, time.Local)
		} else if num == 2 {
			endTime = time.Date(contestStartTime.Year(), 10, 1, 0, 0, 0, 0, time.Local)
		} else if num == 3 {
			endTime = time.Date(contestStartTime.Year()+1, time.January, 1, 0, 0, 0, 0, time.Local)
		}
		exist, err := MasterDB.Table("contest").
			Where("contest_entry_id = ? and contest_type_id = ? and start_time > ? and start_time < ?", form.ContestEntry, searchContestType.ContestTypeID, startTime, endTime).
			Exist()
		if err != nil {
			logging.L.Error(err)
			return err
		}
		if exist {
			logging.L.Error("已有同名竞赛")
			return errors.New("已有同名竞赛")
		}
	}

	newContest := &models.ContestInfo{
		TeacherID:      account.UserID,
		Contest:        form.Contest,
		ContestType:    searchContestType.ContestTypeID,
		SchoolID:       teacher.SchoolID,
		CollegeID:      teacher.CollegeID,
		ContestState:   2,
		CreateTime:     models.NewOftenTime(),
		Describe:       form.Desc,
		ContestEntry:   form.ContestEntry,
		Prize1Count:    form.Prize1,
		Prize2Count:    form.Prize2,
		Prize3Count:    form.Prize3,
		Prize4Count:    form.Prize4,
		ContestLevelID: form.ContestLevel,
		StartTime:      models.FormatString2OftenTime(form.StartTime),
		EnrollTime:     models.FormatString2OftenTime(form.EnrolTime),
		Deadline:       models.FormatString2OftenTime(form.Deadline),
		MaxGroupNumber: form.MaxGroupNumber,
		IsGroup:        form.IsGroup,
		Ps:             form.Ps,
		State:          3,
	}
	_, err = session.Insert(newContest)
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

func (self ContestLogic) GetContestForTeacher(userID int64, year int) (*[]models.ContestAndType, error) {
	contest := &[]models.ContestAndType{}

	teacher, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	if teacher.Role != e.TeacherRole {
		logging.L.Error("无权限")
		return nil, errors.New("无权限")
	}

	startTime := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	endTime := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")

	_, err = MasterDB.
		Table("contest").
		Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id").
		Where("contest.teacher_id = ? and contest.start_time > ? and contest.start_time < ?", teacher.UserID, startTime, endTime).
		Where("contest.state = ?", e.Pass).
		FindAndCount(contest)
	if err != nil {
		return nil, err
	}

	return contest, err
}

func (self ContestLogic) TransformState(userID, id int64, contestState int) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("ProcessEnroll session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			logging.L.Error(err)
			DPrintf("ProcessEnroll session.Close() 发生错误:", err)
		}
	}()

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	session.Table("contest")
	session.Where("teacher_id = ?", account.UserID)
	exist, err := session.Where("id = ?", id).Exist()
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("竞赛不存在")
		return errors.New("竞赛不存在")
	}

	_, err = session.Where("id = ?", id).Update(models.Contest{ContestState: contestState})
	if err != nil {
		fail := session.Rollback()
		if err != nil {
			logging.L.Error(err)
			return fail
		}
		logging.L.Error(err)
		return err
	}
	return session.Commit()
}

func (self ContestLogic) CancelContest(id, userID int64) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("ProcessEnroll session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			logging.L.Error(err)
			DPrintf("ProcessEnroll session.Close() 发生错误:", err)
		}
	}()

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	session.Table("contest")
	session.Where("teacher_id = ?", account.UserID)
	exist, err := session.Where("id = ?", id).Exist()
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("竞赛不存在")
		return errors.New("竞赛不存在")
	}

	_, err = session.Where("id = ?", id).Update(models.Contest{State: e.Revoked})
	if err != nil {
		fail := session.Rollback()
		if err != nil {
			logging.L.Error(err)
			return fail
		}
		logging.L.Error(err)
		return err
	}
	return session.Commit()
}

func (self ContestLogic) DepartmentManagerGetContest(paginator *Paginator, contest, contestType string, contestLevel int, userID int64, isGroup, year int) (*[]models.DepartmentContestEnrollReturn, int64, error) {
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

	departmentAccount, err := public.SearchDepartmentManagerByID(userID)
	if err != nil {
		logging.L.Error()
		return nil, 0, err
	}

	startTime := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	endTime := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")

	session.Table("teacher").Where("teacher.school_id = ? and teacher.college_id = ? and teacher.department_id = ?", departmentAccount.SchoolID, departmentAccount.CollegeID, departmentAccount.DepartmentID)
	session.Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("LEFT", "contest_level", "contest_level.contest_level_id = contest.contest_level_id")
	session.Where("contest.state = ?", e.Pass)
	session.Where("contest.start_time > ? and contest.start_time < ?", startTime, endTime)

	if isGroup > 0 {
		session.Where("contest.is_group = ?", isGroup)
	}
	if err != nil {
		logging.L.Error(err)
	}
	searchContestType := &models.ContestType{}
	if contest != "" {
		session.Where("contest.contest like ?", "%"+contest+"%")
	}
	if contestType != "" {
		searchContestType, err = public.SearchContestTypeByName(contestType)
		if err != nil {
			logging.L.Error(err)
		}
		session.Where("contest.contest_type_id = ?", searchContestType.ContestTypeID)
	}
	if contestLevel != -1 {
		session.Where("contest.contest_level_id = ?", contestLevel)
	}

	data := &[]models.ContestContestTypeTeacher{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}

	session.Join("LEFT", "enroll_information", "contest.id = enroll_information.contest_id")

	list := make([]models.DepartmentContestEnrollReturn, len(*data))
	for i := 0; i < len(list); i++ {
		list[i].Username = (*data)[i].Username
		list[i].Name = (*data)[i].Name
		list[i].School = (*data)[i].School
		list[i].College = (*data)[i].College
		list[i].ID = (*data)[i].Contest.ID
		list[i].State = (*data)[i].Contest.State
		list[i].Contest = (*data)[i].Contest.Contest
		list[i].ContestType = (*data)[i].Contest.ContestType
		list[i].ContestLevel = (*data)[i].Contest.ContestLevel
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].Contest.CreateTime)
		list[i].EnrollTime = models.MysqlFormatString2String((*data)[i].Contest.EnrollTime)
		list[i].StartTime = models.MysqlFormatString2String((*data)[i].Contest.StartTime)
		list[i].Deadline = models.MysqlFormatString2String((*data)[i].Contest.Deadline)
		list[i].PassCount, _ = MasterDB.Table("enroll_information").Where("state = ? and contest_id = ?", e.Pass, (*data)[i].Contest.ID).Count()
		list[i].RejectedCount, _ = MasterDB.Table("enroll_information").Where("state = ? and contest_id = ?", e.Reject, (*data)[i].Contest.ID).Count()
		list[i].ProcessingCount, _ = MasterDB.Table("enroll_information").Where("state = ? and contest_id = ?", e.Processing, (*data)[i].Contest.ID).Count()
	}

	return &list, total, session.Commit()
}

func (self ContestLogic) ViewTeacherContestEnroll(paginator *Paginator, contest, contestType string, contestLevel int, userID int64, year, isGroup int) (*[]models.TeacherContestEnrollReturn, int64, error) {
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

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error()
		return nil, 0, err
	}

	startTime := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	endTime := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")

	session.Table("teacher").Where("teacher.teacher_id = ?", account.UserID)
	session.Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id")
	session.Where("contest.is_group = ?", isGroup)
	session.Where("contest.start_time > ? and contest.start_time < ?", startTime, endTime)
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("LEFT", "contest_level", "contest_level.contest_level_id = contest.contest_level_id")
	session.Where("contest.state = ?", e.Pass)

	if err != nil {
		logging.L.Error(err)
	}
	searchContest := &models.ContestInfo{}
	searchContestType := &models.ContestType{}
	if contest != "" {
		searchContest, err = public.SearchContestByName(contest)
		if err != nil {
			logging.L.Error(err)
		}
		session.Where("contest.id = ?", searchContest.ID)
	}
	if contestType != "" {
		searchContestType, err = public.SearchContestTypeByName(contest)
		if err != nil {
			logging.L.Error(err)
		}
		session.Where("contest.contest_type_id = ?", searchContestType.ContestTypeID)
	}
	if contestLevel != -1 {
		session.Where("contest.contest_level_id = ?", contestLevel)
	}

	data := &[]models.ContestContestTypeTeacher{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}

	session.Join("LEFT", "enroll_information", "contest.id = enroll_information.contest_id")

	list := make([]models.TeacherContestEnrollReturn, len(*data))
	for i := 0; i < len(list); i++ {
		list[i].Username = (*data)[i].Username
		list[i].Name = (*data)[i].Name
		list[i].School = (*data)[i].School
		list[i].College = (*data)[i].College
		list[i].ID = (*data)[i].Contest.ID
		list[i].State = (*data)[i].Contest.State
		list[i].Contest = (*data)[i].Contest.Contest
		list[i].ContestType = (*data)[i].Contest.ContestType
		list[i].ContestLevel = (*data)[i].Contest.ContestLevel
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].Contest.CreateTime)
		list[i].EnrollTime = models.MysqlFormatString2String((*data)[i].Contest.EnrollTime)
		list[i].StartTime = models.MysqlFormatString2String((*data)[i].Contest.StartTime)
		list[i].Deadline = models.MysqlFormatString2String((*data)[i].Contest.Deadline)
		list[i].EnrollCount, err = MasterDB.Table("enroll_information").Where("state = ? and contest_id = ?", e.Pass, (*data)[i].Contest.ID).Count()
	}

	return &list, total, session.Commit()
}

func (self ContestLogic) ProcessContest(id int64, state int, userID int64, rejectReason string) error {
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

	departmentAccount, err := public.SearchDepartmentManagerByID(userID)
	if err != nil {
		logging.L.Error()
		return err
	}

	session.Table("contest").Where("id = ?", id)
	session.Join("LEFT", "teacher", "contest.teacher_id = teacher.teacher_id")

	exist, err := session.Where("teacher.school_id = ? and teacher.college_id = ? and teacher.department_id = ?", departmentAccount.SchoolID, departmentAccount.CollegeID, departmentAccount.DepartmentID).Exist()
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("竞赛不存在")
		return errors.New("竞赛不存在")
	}

	newContest := models.ContestInfo{State: state}
	if state == e.Reject {
		newContest.RejectReason = rejectReason
	}

	_, err = session.Where("id = ?", id).Update(&newContest)
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

func (self ContestLogic) OnlyGetDepartmentContest(userID int64) (*[]models.ContestAndType, error) {
	contest := &[]models.ContestAndType{}

	departmentAccount, err := public.SearchDepartmentManagerByID(userID)
	if err != nil {
		logging.L.Error()
		return nil, err
	}

	//MasterDB.Table("teacher").Where("teacher.school_id = ? and teacher.college_id = ? and teacher.department_id = ?", departmentAccount.SchoolID, departmentAccount.CollegeID, departmentAccount.DepartmentID)
	//MasterDB.Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id")
	//MasterDB.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	//MasterDB.Join("LEFT", "contest_level", "contest_level.contest_level_id = contest.contest_level_id")

	_, err = MasterDB.
		Table("teacher").
		Where("teacher.school_id = ? and teacher.college_id = ? and teacher.department_id = ?", departmentAccount.SchoolID, departmentAccount.CollegeID, departmentAccount.DepartmentID).
		Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id").
		Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id").
		Join("LEFT", "contest_level", "contest_level.contest_level_id = contest.contest_level_id").
		Where("contest.state = ?", e.Pass).
		FindAndCount(contest)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}

	return contest, err
}

func (self ContestLogic) DepartmentManagerGetContestGrade(paginator *Paginator, contest, contestType string, contestLevel int, userID int64, isGroup, year int) (*[]models.DepartmentContestGradeReturn, int64, error) {
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

	departmentAccount, err := public.SearchDepartmentManagerByID(userID)
	if err != nil {
		logging.L.Error()
		return nil, 0, err
	}

	startTime := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	endTime := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")

	session.Table("teacher").Where("teacher.school_id = ? and teacher.college_id = ? and teacher.department_id = ?", departmentAccount.SchoolID, departmentAccount.CollegeID, departmentAccount.DepartmentID)
	session.Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("LEFT", "contest_level", "contest_level.contest_level_id = contest.contest_level_id")
	session.Where("contest.state = ?", e.Pass)
	session.Where("contest.start_time > ? and contest.start_time < ?", startTime, endTime)

	if isGroup > 0 {
		session.Where("contest.is_group = ?", isGroup)
	}

	if err != nil {
		logging.L.Error(err)
	}
	searchContest := &models.ContestInfo{}
	searchContestType := &models.ContestType{}
	if contest != "" {
		searchContest, err = public.SearchContestByName(contest)
		if err != nil {
			logging.L.Error(err)
		}
		session.Where("contest.id = ?", searchContest.ID)
	}
	if contestType != "" {
		searchContestType, err = public.SearchContestTypeByName(contest)
		if err != nil {
			logging.L.Error(err)
		}
		session.Where("contest.contest_type_id = ?", searchContestType.ContestTypeID)
	}
	if contestLevel != -1 {
		session.Where("contest.contest_level_id = ?", contestLevel)
	}

	data := &[]models.ContestContestTypeTeacherGrade{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}

	session.Join("LEFT", "grade", "contest.id = grade.contest_id")

	list := make([]models.DepartmentContestGradeReturn, len(*data))
	for i := 0; i < len(list); i++ {
		list[i].Username = (*data)[i].Username
		list[i].School = (*data)[i].School
		list[i].College = (*data)[i].College
		list[i].ID = (*data)[i].Contest.ID
		list[i].State = (*data)[i].Contest.State
		list[i].Contest = (*data)[i].Contest.Contest
		list[i].ContestType = (*data)[i].Contest.ContestType
		list[i].ContestLevel = (*data)[i].Contest.ContestLevel
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].Contest.CreateTime)
		list[i].RejectReason = (*data)[i].RejectReason
		list[i].Prize1Count = (*data)[i].Prize1Count
		list[i].Prize2Count = (*data)[i].Prize2Count
		list[i].Prize3Count = (*data)[i].Prize3Count
		list[i].Prize4Count = (*data)[i].Prize4Count
		list[i].RewardCount, _ = session.Table("grade").Where("state = ? and contest_id = ?", e.Pass, (*data)[i].Contest.ID).Count()
		list[i].RejectedCount, _ = session.Table("grade").Where("state = ? and contest_id = ?", e.Reject, (*data)[i].Contest.ID).Count()
		list[i].ProcessingCount, _ = session.Table("grade").Where("state = ? and contest_id = ?", e.Processing, (*data)[i].Contest.ID).Count()
	}

	return &list, total, session.Commit()
}
