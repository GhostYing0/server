package logic

import (
	"errors"
	"fmt"
	. "server/database"
	"server/logic/public"
	"server/models"
	. "server/utils/e"
	"server/utils/logging"
	. "server/utils/mydebug"
)

type EnrollLogic struct{}

var DefaultEnrollLogic = EnrollLogic{}

func (self EnrollLogic) InsertEnrollInformation(userID int64, name, teamID, contest string, school string, phone string, email string) error {
	if phone == "" || email == "" {
		logging.L.Error("请手机号和邮箱")
		return errors.New("请手机号和邮箱")
	}
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("InsertEnrollInformation session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("InsertEnrollInformation session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	searchContest, err := public.SearchContestByName(contest)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	searchSchool, err := public.SearchSchoolByName(school)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	exist, err := session.
		Table("enroll_information").
		Where("contest_id = ? AND student_id = ?", searchContest.ID, account.UserID).
		Exist()
	if err != nil {
		DPrintf("InsertEnrollInformation 查询重复报名发生错误: ", err)
		logging.L.Error(err)
		return err
	}
	if exist {
		DPrintf("请勿重复报名")
		logging.L.Error("请勿重复报名")
		return errors.New("请勿重复报名")
	}

	enroll := &models.NewEnroll{
		StudentID:  account.UserID,
		TeamID:     teamID,
		ContestID:  searchContest.ID,
		CreateTime: models.NewOftenTime(),
		SchoolID:   searchSchool.SchoolID,
		Phone:      phone,
		Email:      email,
		State:      3,
	}

	_, err = session.Insert(enroll)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		DPrintf("InsertEnrollInformation 添加报名信息发生错误:", err)
		return err
	}

	return session.Commit()
}

func (self EnrollLogic) Search(paginator *Paginator, userID int64, contest string, startTime string, endTime string, state int) (*[]models.EnrollInformationReturn, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		logging.L.Error("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Search session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Search session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	session.Join("LEFT", "contest", "contest.id = enroll_information.contest_id")
	session.Join("LEFT", "account", "account.user_id = enroll_information.student_id")
	if userID > 0 {
		session.Where("account.id = ?", userID)
	}
	if len(contest) > 0 {
		session.Where("contest.contest = ?", contest)
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		session.Where("enroll_information.create_time >= ? AND enroll_information.create_time <= ?", startTime, endTime)
	}
	if state >= 0 {
		session.Where("enroll_information.state = ?", state)
	}
	fmt.Println(state)

	data := &[]models.EnrollContestStudent{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		DPrintf("Search 查找报名信息失败:", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.EnrollInformationReturn, len(*data))
	for i := 0; i < len(*data); i++ {
		if err != nil {
			logging.L.Error(err)
		}
		list[i].ID = (*data)[i].ID
		//list[i].Username = (*temp)[i].Username
		//list[i].UserID = (*temp)[i].UserID
		list[i].TeamID = (*data)[i].TeamID
		list[i].Contest = (*data)[i].Contest
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].EnrollInformation.CreateTime)
		list[i].Phone = (*data)[i].Phone
		list[i].Email = (*data)[i].Email
		list[i].State = (*data)[i].EnrollInformation.State
		list[i].RejectReason = (*data)[i].RejectReason
		startTime := models.FormatString2OftenTime(models.MysqlFormatString2String((*data)[i].StartTime))
		list[i].StartTime = startTime.String()
		if models.NewOftenTime().After(&startTime) && list[i].State == Pass {
			fmt.Println("Start:", startTime)
			fmt.Println("timeNow:", models.NewOftenTime())
			list[i].DoUpload = true
		}
	}

	return &list, total, session.Rollback()
}

func (self EnrollLogic) ProcessEnroll(id int64, state int, rejectReason string) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("ProcessEnroll session.Begin() 发生错误:", err)
		logging.L.Error()
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("ProcessEnroll session.Close() 发生错误:", err)
		}
	}()

	if id < 1 {
		fmt.Println("非法id")
		logging.L.Error("非法id")
		return errors.New("非法id")
	}

	exist, err := session.Table("enroll_information").Where("id = ?", id).Exist()
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("报名信息不存在")
		return errors.New("报名信息不存在")
	}

	newInfo := &models.EnrollInformation{State: state}
	if state == Reject {
		newInfo.RejectReason = rejectReason
	}

	_, err = session.Where("id = ?", id).Update(newInfo)
	if err != nil {
		DPrintf("EnrollLogic Update 发生错误:", err)
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

func (self EnrollLogic) TeacherGetOneEnroll(paginator *Paginator, enrollID, userID int64, contest string, startTime string, endTime string, state int, contestType string) (*[]models.TeacherGetOneEnrollInformationReturn, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		logging.L.Error("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Search session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Search session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}
	if account.Role != TeacherRole {
		logging.L.Error("权限错误")
		return nil, 0, errors.New("权限错误")
	}

	teacher, err := public.SearchTeacherByID(account.UserID)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}
	if enrollID <= 0 {
		return nil, 0, err
	}

	session.Table("account").Where("user_id = ?", account.UserID)
	session.Join("LEFT", "contest", "contest.teacher_id = account.user_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("RIGHT", "enroll_information", "contest.id = enroll_information.contest_id")
	session.Join("LEFT", "student", "student.student_id = enroll_information.student_id")
	session.Join("LEFT", "contest_level", "contest_level.contest_level_id = contest.contest_level_id")
	session.Join("LEFT", "major", "major.major_id = student.major")
	session.Where("contest.school_id = ?", teacher.SchoolID)
	session.Where("enroll_information.id = ?", enrollID)
	session.Select("account.user_id, enroll_information.*, contest.contest, contest_type.type, student.*, contest_level.contest_level, school.school, college.college, semester.semester, major.major")
	if len(contest) > 0 {
		session.Where("contest.contest = ?", contest)
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		session.Where("enroll_information.create_time >= ? AND enroll_information.create_time <= ?", startTime, endTime)
	}
	if state >= 0 {
		session.Where("enroll_information.state = ?", state)
	}
	if contestType != "" {
		session.Where("contest_type.type = ?", contestType)
	}
	session.Join("LEFT", "school", "student.school_id = school.school_id")
	session.Join("LEFT", "college", "student.college_id = college.college_id")
	session.Join("LEFT", "semester", "student.semester_id = semester.semester_id")

	data := &[]models.TeacherUploadGetEnroll{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).
		//Select("ei1.id as id, ei1.*, where account.*,contest.*,contest_type.*,student.*,school.*,college.*,semester.*").
		FindAndCount(data)
	if err != nil {
		DPrintf("Search 查找报名信息失败:", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.TeacherGetOneEnrollInformationReturn, len(*data))
	for i := 0; i < len(*data); i++ {
		if err != nil {
			logging.L.Error(err)
		}
		list[i].ID = (*data)[i].EnrollInformation.ID
		list[i].Name = (*data)[i].Name
		list[i].ContestType = (*data)[i].ContestType
		list[i].School = (*data)[i].School
		list[i].College = (*data)[i].College
		list[i].ContestLevel = (*data)[i].ContestLevel
		list[i].Semester = (*data)[i].Semester
		list[i].StudentSchoolID = (*data)[i].StudentSchoolID
		list[i].Class = (*data)[i].Class
		list[i].TeamID = (*data)[i].TeamID
		list[i].Contest = (*data)[i].Contest
		list[i].Major = (*data)[i].Major
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].EnrollInformation.CreateTime)
		list[i].Phone = (*data)[i].Phone
		list[i].Email = (*data)[i].Email
	}

	return &list, total, session.Rollback()
}

func (self EnrollLogic) TeacherSearch(paginator *Paginator, contestID, userID int64, contest string, startTime string, endTime string, state int, contestType string) (*[]models.EnrollInformationReturn, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		logging.L.Error("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Search session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Search session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}
	if account.Role != TeacherRole {
		logging.L.Error("权限错误")
		return nil, 0, errors.New("权限错误")
	}

	teacher, err := public.SearchTeacherByID(account.UserID)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}
	if contestID <= 0 {
		return nil, 0, err
	}

	session.Table("account").Where("user_id = ?", account.UserID)
	session.Join("LEFT", "contest", "contest.teacher_id = account.user_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("RIGHT", "enroll_information", "contest.id = enroll_information.contest_id")
	session.Where("contest.school_id = ?", teacher.SchoolID)
	session.Where("contest.id = ?", contestID)
	session.Select("account.user_id, enroll_information.*, contest.contest, contest_type.type, student.*, school.school, college.college, semester.semester")
	if len(contest) > 0 {
		session.Where("contest.contest = ?", contest)
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		session.Where("enroll_information.create_time >= ? AND enroll_information.create_time <= ?", startTime, endTime)
	}
	if state >= 0 {
		session.Where("enroll_information.state = ?", state)
	}
	if contestType != "" {
		session.Where("contest_type.type = ?", contestType)
	}
	session.Join("LEFT", "student", "student.student_id = enroll_information.student_id")
	session.Join("LEFT", "school", "student.school_id = school.school_id")
	session.Join("LEFT", "college", "student.college_id = college.college_id")
	session.Join("LEFT", "semester", "student.semester_id = semester.semester_id")

	data := &[]models.EnrollContestStudent{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).
		//Select("ei1.id as id, ei1.*, where account.*,contest.*,contest_type.*,student.*,school.*,college.*,semester.*").
		FindAndCount(data)
	if err != nil {
		DPrintf("Search 查找报名信息失败:", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.EnrollInformationReturn, len(*data))
	for i := 0; i < len(*data); i++ {
		if err != nil {
			logging.L.Error(err)
		}
		list[i].ID = (*data)[i].EnrollInformation.ID
		list[i].Username = (*data)[i].Username
		list[i].Name = (*data)[i].Name
		list[i].ContestType = (*data)[i].ContestType
		list[i].School = (*data)[i].School
		list[i].College = (*data)[i].College
		list[i].Semester = (*data)[i].Semester
		list[i].Class = (*data)[i].Class
		list[i].TeamID = (*data)[i].TeamID
		list[i].Contest = (*data)[i].Contest
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].EnrollInformation.CreateTime)
		list[i].Phone = (*data)[i].Phone
		list[i].Email = (*data)[i].Email
		list[i].RejectReason = (*data)[i].EnrollInformation.RejectReason
		list[i].State = (*data)[i].EnrollInformation.State
		startTime := models.FormatString2OftenTime(models.MysqlFormatString2String((*data)[i].StartTime))
		//list[i].StartTime = startTime.String()
		if models.NewOftenTime().After(&startTime) && list[i].State == Pass {
			fmt.Println("Start:", startTime)
			fmt.Println("timeNow:", models.NewOftenTime())
			list[i].DoUpload = true
		}
	}

	return &list, total, session.Rollback()
}

func (self EnrollLogic) DepartmentManagerSearchEnroll(paginator *Paginator, contestID, userID int64, contest string, startTime string, endTime string, state int, contestType string) (*[]models.EnrollInformationReturn, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		logging.L.Error("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Search session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Search session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	account, err := public.SearchDepartmentManagerByID(userID)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}
	if account.Role != DepartmentRole {
		logging.L.Error("权限错误")
		return nil, 0, errors.New("权限错误")
	}

	//sql := `SELECT enroll_information.id as e_id, ` +
	//	`student.name, student.student_id, student.school_id, contest.contest, contest_type.type FROM student ` +
	//	`RIGHT JOIN enroll_information ON student.student_id = enroll_information.student_id ` +
	//	`LEFT JOIN contest ON contest.id = enroll_information.contest_id ` +
	//	`LEFT JOIN contest_type ON contest_type.id = contest.contest_type_id `
	session.Table("student").Where("student.school_id = ?", account.SchoolID)
	session.Join("RIGHT", "enroll_information", "student.student_id = enroll_information.student_id")
	session.Join("LEFT", "contest", "contest.id = enroll_information.contest_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Select("enroll_information.id as e_id, enroll_information.*, student.*, contest_type.")
	session.Where("enroll_information.contest_id = ?", contestID)
	//session.Table("student").Where("student.school_id = ?", account.SchoolID)
	//session.Join("RIGHT", "enroll_information", "student.student_id = enroll_information.student_id")
	//session.Join("LEFT", "contest", "contest.id = enroll_information.contest_id")
	//session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	//session.Where("enroll_information.contest_id = ?", contestID)
	if contestID <= 0 {
		return nil, 0, err
	}
	if len(contest) > 0 {
		session.Where("contest.contest = ?", contest)
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		session.Where("enroll_information.create_time >= ? AND enroll_information.create_time <= ?", startTime, endTime)
	}
	if state >= 0 {
		session.Where("enroll_information.state = ?", state)
	}
	if contestType != "" {
		session.Where("contest_type.type = ?", contestType)
	}
	//session.Join("LEFT", "school", "student.school_id = school.school_id")
	//session.Join("LEFT", "college", "student.college_id = college.college_id")
	//session.Join("LEFT", "semester", "student.semester_id = semester.semester_id")

	data := &[]models.EnrollContestStudent_e_id{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).
		FindAndCount(data)
	if err != nil {
		DPrintf("Search 查找报名信息失败:", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.EnrollInformationReturn, len(*data))
	for i := 0; i < len(*data); i++ {
		if err != nil {
			logging.L.Error(err)
		}
		list[i].ID = (*data)[i].ID
		list[i].Username = (*data)[i].Username
		list[i].Name = (*data)[i].Name
		list[i].ContestType = (*data)[i].ContestType
		list[i].School = (*data)[i].School
		list[i].College = (*data)[i].College
		list[i].Semester = (*data)[i].Semester
		list[i].Class = (*data)[i].Class
		list[i].TeamID = (*data)[i].TeamID
		list[i].Contest = (*data)[i].Contest
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].CreateTime)
		list[i].Phone = (*data)[i].Phone
		list[i].Email = (*data)[i].Email
		list[i].RejectReason = (*data)[i].RejectReason
		list[i].State = (*data)[i].State
		startTime := models.FormatString2OftenTime(models.MysqlFormatString2String((*data)[i].StartTime))
		//list[i].StartTime = startTime.String()
		if models.NewOftenTime().After(&startTime) && list[i].State == Pass {
			fmt.Println("Start:", startTime)
			fmt.Println("timeNow:", models.NewOftenTime())
			list[i].DoUpload = true
		}
	}

	return &list, total, session.Rollback()
}

func (self EnrollLogic) PassEnroll(ids *[]int64, state int) (int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("InsertEnrollInformation session.Begin() 发生错误:", err)
		return 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("InsertEnrollInformation session.Close() 发生错误:", err)
		}
	}()

	var count int64

	for _, id := range *ids {
		var enrollInformation models.EnrollInformation
		if id < 1 {
			fmt.Println("非法id")
			continue
		}
		enrollInformation.State = Pass
		affected, err := session.Where("id = ?", id).Update(&enrollInformation)
		if err != nil {
			DPrintf("CmsRegistrationLogic Update 发生错误:", err)
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return 0, fail
			}
			return count, err
		}
		if affected > 0 {
			count += affected
		}
	}

	return count, session.Commit()
}
