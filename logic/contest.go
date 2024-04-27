package logic

import (
	"errors"
	. "server/database"
	"server/logic/public"
	"server/models"
	"server/utils/e"
	"server/utils/logging"
	. "server/utils/mydebug"
)

type ContestLogic struct{}

var DefaultContestLogic = ContestLogic{}

func (self ContestLogic) DisplayContest(paginator *Paginator, contest, contestType string, userID int64) (*[]models.ContestReturn, int64, error) {
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

	if account.Role == e.StudentRole {
		student, err := public.SearchStudentByID(account.UserID)
		if err != nil {
			return nil, 0, err
		}

		session.Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id")
		session.Where("contest.state = 1")
		session.Where("contest.school_id = ?", student.SchoolID)
	} else if account.Role == e.TeacherRole {
		teacher, err := public.SearchTeacherByID(account.UserID)
		if err != nil {
			return nil, 0, err
		}

		session.Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id")
		session.Where("contest.state = 1")
		session.Where("contest.school_id = ?", teacher.SchoolID)
	}

	if contest != "" {
		session.Where("contest = ?", contest)
	}
	if contestType != "" {
		searchContestType, err := public.SearchContestTypeByName(contestType)
		if err != nil {

		} else {
			session.Where("contest.contest_type_id = ?", searchContestType.ContestTypeID)
		}
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
		// 竞赛可报名条件，审核通过，在报名截至时间之前，且教师未关闭报名
		if (*data)[i].State == e.Pass && (*data)[i].ContestState == e.EnrollOpen && models.NewOftenTime().Before(&(*data)[i].Deadline) {
			list[i].ContestState = e.EnrollOpen
		} else {
			list[i].ContestState = e.EnrollClose
		}
	}

	return &list, total, session.Commit()
}

func (self ContestLogic) ViewTeacherContest(paginator *Paginator, userID int64, contest, contestType string, state int) (*[]models.ContestReturn, int64, error) {
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

	session.Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id")
	session.Where("contest.teacher_id = ?", account.UserID)

	if state != -1 {
		session.Where("contest.state = ?", state)
	}
	if contest != "" {
		session.Where("contest = ?", contest)
	}
	if contestType != "" {
		searchContestType, err := public.SearchContestTypeByName(contestType)
		if err != nil {

		} else {
			session.Where("contest.contest_type_id = ?", searchContestType.ContestTypeID)
		}
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
		// 竞赛可报名条件，审核通过，在报名截至时间之前，且教师未关闭报名
		if (*data)[i].State == e.Pass && (*data)[i].ContestState == e.EnrollOpen && models.NewOftenTime().Before(&(*data)[i].Deadline) {
			list[i].ContestState = e.EnrollOpen
		} else {
			list[i].ContestState = e.EnrollClose
		}
	}

	return &list, total, session.Commit()
}

func (self ContestLogic) UpdateContest(id, userID int64, contest, contestType, startTime, deadline string, contestState, state int) error {
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

	if contest != "" {
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

		searchContestType, err := public.SearchContestTypeByName(contestType)
		if err != nil {
			logging.L.Error(err)
			return err
		}

		exist, err = session.Table("contest").Where("contest = ? and contest_type_id = ?", contest, searchContestType.ContestTypeID).Exist()
		if exist {
			logging.L.Error("已有同名竞赛")
			return errors.New("已有同名竞赛")
		}
		if err != nil {
			logging.L.Error(err)
			return err
		}
	}
	newContestType := int64(0)
	if contestType != "" {
		searchContestType, err := public.SearchContestTypeByName(contestType)
		if err != nil {
			logging.L.Error(err)
			return err
		}
		newContestType = searchContestType.ContestTypeID
	}

	updateContest := &models.ContestInfo{
		Contest:      contest,
		ContestType:  newContestType,
		ContestState: contestState,
		State:        e.Processing,
	}

	if startTime != "" {
		updateContest.StartTime = models.FormatString2OftenTime(startTime)
	}
	if deadline != "" {
		updateContest.Deadline = models.FormatString2OftenTime(deadline)
	}
	_, err = session.Where("id = ?", id).Update(updateContest)
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

func (self ContestLogic) UploadContest(userID int64, contest, contestType, startTime, deadline string, describe *string) error {
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

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	searchContestType, err := public.SearchContestTypeByName(contestType)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	teacher, err := public.SearchTeacherByID(account.UserID)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	exist, err := session.Table("contest").Where("contest = ? and contest_type_id = ?", contest, searchContestType.ContestTypeID).Exist()
	if exist {
		logging.L.Error("已有同名竞赛")
		return errors.New("已有同名竞赛")
	}
	if err != nil {
		logging.L.Error(err)
		return err
	}

	newContest := &models.ContestInfo{
		TeacherID:    account.UserID,
		Contest:      contest,
		ContestType:  searchContestType.ContestTypeID,
		SchoolID:     teacher.SchoolID,
		CollegeID:    teacher.CollegeID,
		ContestState: 2,
		CreateTime:   models.NewOftenTime(),
		Describe:     *describe,
		State:        3,
	}

	if startTime != "" {
		newContest.StartTime = models.FormatString2OftenTime(startTime)
	}
	if deadline != "" {
		newContest.Deadline = models.FormatString2OftenTime(deadline)
	}
	_, err = session.Insert(newContest)
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

func (self ContestLogic) GetContestForTeacher(userID int64) (*[]models.ContestAndType, error) {
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

	_, err = MasterDB.
		Table("contest").
		Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id").
		Where("contest.teacher_id = ?", teacher.UserID).
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

func (self ContestLogic) DepartmentManagerGetContest(paginator *Paginator, contest, contestType string, contestLevel int, userID int64) (*[]models.DepartmentContestEnrollReturn, int64, error) {
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
		list[i].StartTime = models.MysqlFormatString2String((*data)[i].Contest.StartTime)
		list[i].Deadline = models.MysqlFormatString2String((*data)[i].Contest.Deadline)
		list[i].RejectedCount, _ = session.Table("enroll_information").Where("state = ? and contest_id = ?", e.Reject, (*data)[i].Contest.ID).Count()
		list[i].ProcessingCount, _ = session.Table("enroll_information").Where("state = ? and contest_id = ?", e.Processing, (*data)[i].Contest.ID).Count()
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

func (self ContestLogic) DepartmentManagerGetContestGrade(paginator *Paginator, contest, contestType string, contestLevel int, userID int64) (*[]models.DepartmentContestGradeReturn, int64, error) {
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
		list[i].RejectedCount, _ = session.Table("grade").Where("state = ? and contest_id = ?", e.Reject, (*data)[i].Contest.ID).Count()
		list[i].ProcessingCount, _ = session.Table("grade").Where("state = ? and contest_id = ?", e.Processing, (*data)[i].Contest.ID).Count()
	}

	return &list, total, session.Commit()
}
