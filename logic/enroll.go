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

func (self EnrollLogic) InsertEnrollInformation(userID, contestID, handle int64,
	teamName, guidanceTeacher, teacherTitle, teacherDepartment, phone, email, college, major string) error {
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

	EnrollInformation := &models.EnrollInformation{}
	exist, err := session.
		Table("enroll_information").
		Where("contest_id = ? AND student_id = ?", contestID, account.UserID).
		Get(EnrollInformation)
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

	searchContest, err := public.SearchContestByID(contestID)
	if err != nil {
		logging.L.Error()
		return err
	}
	if searchContest.State != Pass || searchContest.ContestState != EnrollOpen {
		return errors.New("竞赛不可报名")
	}
	searchStudent, err := public.SearchStudentByID(account.UserID)
	if err != nil {
		logging.L.Error()
		return err
	}
	searchSchool, err := public.SearchSchoolByID(searchStudent.SchoolID)
	if err != nil {
		logging.L.Error()
		return err
	}
	_, err = public.SearchCollegeByName(college)
	if err != nil {
		logging.L.Error()
		return err
	}
	_, err = public.SearchMajorByName(major)
	if err != nil {
		logging.L.Error()
		return err
	}
	department, err := public.SearchDepartmentByName(teacherDepartment)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	teacherAccount := &models.Teacher{}
	exist, err = MasterDB.Where("name = ? and department_id = ?", guidanceTeacher, department.DepartmentID).Get(teacherAccount)
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("教师不存在")
		return errors.New("教师不存在")
	}

	teamID := int64(0)
	if searchContest.IsGroup == 1 {
		if handle == 2 {
			//加入队伍
			searchTeam, err := public.SearchTeamByNameAndContest(teamName, contestID)
			if err != nil && err.Error() != "队伍已存在" {
				logging.L.Error(err)
				return err
			}
			//看满没满队

			count, err := MasterDB.
				Table("enroll_information").
				Where("contest_id = ? and team_id = ? and (state = ? or state = ?)", contestID, searchTeam.TeamID, Pass, Processing).
				Count()
			if int(count) >= searchContest.MaxGroupNumber {
				logging.L.Error("队伍已满")
				return errors.New("队伍已满")
			}
			teamID = searchTeam.TeamID
		} else {
			//创建队伍
			//看队伍存不存在
			_, err := public.SearchTeamByNameAndContest(teamName, contestID)
			if err != nil {
				logging.L.Error(err)
				return err
			}

			newTeam := &models.Team{TeamName: teamName, ContestID: contestID}

			_, err = session.Insert(newTeam)
			if err != nil {
				fail := session.Rollback()
				if fail != nil {
					logging.L.Error(fail)
					return err
				}
				logging.L.Error(err)
				return err
			}
			teamID = newTeam.TeamID
		}
	}

	enroll := &models.NewEnroll{
		StudentID:       account.UserID,
		TeamID:          teamID,
		ContestID:       contestID,
		CreateTime:      models.NewOftenTime(),
		SchoolID:        searchSchool.SchoolID,
		Phone:           phone,
		GuidanceTeacher: guidanceTeacher,
		Email:           email,
		State:           3,
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

func (self EnrollLogic) Search(paginator *Paginator, userID, EnrollID, contestLevel int64, contest string, startTime string, endTime string, contestType string, state, isGroup, role int) (*[]models.EnrollInformationReturn, int64, error) {
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

	if isGroup > 0 {
		session.Where("contest.is_group = ?", isGroup)
	}
	if EnrollID > 0 {
		session.Where("enroll_information.id = ?", EnrollID)
	}

	session.Join("LEFT", "teacher", "teacher.teacher_id = enroll_information.guidance_teacher")
	session.Join("LEFT", "department", "teacher.department_id = department.department_id")
	session.Join("LEFT", "student", "student.student_id = enroll_information.student_id")
	session.Join("LEFT", "major", "student.major_id = major.major_id")
	session.Join("LEFT", "contest_level", "contest_level.contest_level_id = contest.contest_level_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("LEFT", "contest_entry", "contest_entry.contest_entry_id = contest.contest_entry_id")
	session.Join("LEFT", "college", "college.college_id = student.college_id")
	session.Join("LEFT", "team", "team.team_id = enroll_information.team_id")
	session.Select("enroll_information.*, contest.contest, contest.is_group, contest_level.contest_level, contest_type.type, contest.start_time," +
		"major.major, teacher.name as t_name, teacher.title, department.department t_department, student.*, contest_entry.contest_entry," +
		"college.college, team.team_name")
	if contestLevel > 0 {
		session.Where("contest.contest_level_id = ?", contestLevel)
	}

	if role == StudentRole {
		session.Join("LEFT", "account", "account.user_id = enroll_information.student_id")
		if userID > 0 {
			session.Where("account.id = ?", userID)
		}
	}
	if len(contest) > 0 {
		session.Where("contest.contest like ?", "%"+contest+"%")
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
		list[i].Name = (*data)[i].Name
		list[i].StudentSchoolID = (*data)[i].StudentSchoolID
		list[i].Class = (*data)[i].Class
		list[i].ContestEntry = (*data)[i].ContestEntry
		list[i].College = (*data)[i].College
		list[i].TeamID = (*data)[i].TeamID
		list[i].IsGroup = (*data)[i].IsGroup
		list[i].Contest = (*data)[i].Contest
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].EnrollInformation.CreateTime)
		list[i].Phone = (*data)[i].Phone
		list[i].Email = (*data)[i].Email
		list[i].ContestLevel = (*data)[i].ContestLevel
		list[i].ContestType = (*data)[i].ContestType
		list[i].State = (*data)[i].EnrollInformation.State
		list[i].RejectReason = (*data)[i].RejectReason
		list[i].Major = (*data)[i].Major
		list[i].TeacherName = (*data)[i].TeacherName
		list[i].Department = (*data)[i].Department
		list[i].Title = (*data)[i].Title
		list[i].Team = (*data)[i].Team
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].CreateTime)
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
	session.Join("LEFT", "major", "major.major_id = student.major_id")
	session.Where("contest.school_id = ?", teacher.SchoolID)
	session.Where("enroll_information.id = ?", enrollID)
	session.Select("account.user_id, enroll_information.*, contest.contest, contest.id as contest_id, contest_type.type, student.*, contest_level.contest_level, school.school, college.college, semester.semester, major.major")
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
		list[i].ContestID = (*data)[i].ContestID
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

func (self EnrollLogic) TeacherSearch(paginator *Paginator, contestID, userID int64, contest, class, major, name, teamName string /* startTime, endTime string*/, state int, contestType string, year int) (*[]models.EnrollInformationReturn, int64, error) {
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

	//startTime := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	//endTime := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")

	//session.Table("account").Where("user_id = ?", account.UserID)
	//session.Join("LEFT", "contest", "contest.teacher_id = account.user_id")
	//session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	//session.Join("RIGHT", "enroll_information", "contest.id = enroll_information.contest_id")
	//session.Where("contest.id = ?", contestID)
	//session.Where("contest.school_id = ? and enroll_information.state = ?", teacher.SchoolID, Pass)
	//session.Join("LEFT", "student", "student.student_id = enroll_information.student_id")
	//session.Join("LEFT", "major", "student.major_id = major.major_id")
	//session.Join("LEFT", "school", "student.school_id = school.school_id")
	//session.Join("LEFT", "teacher", "teacher.teacher_id = enroll_information.guidance_teacher")
	//session.Join("LEFT", "department", "teacher.department_id = department.department_id")
	//session.Join("LEFT", "team", "enroll_information.team_id = team.team_id")
	//session.Select("account.user_id, enroll_information.*, contest.contest, contest_type.type, student.*," +
	//	"major.*, school.school, college.college, semester.semester, team.team_name, major.major," +
	//	"teacher.name as t_name, teacher.title, department.department as t_department")

	session.Table("teacher").Where("teacher.teacher_id = ?", teacher.TeacherID)
	session.Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("LEFT", "contest_level", "contest_level.contest_level_id = contest.contest_level_id")
	session.Join("RIGHT", "enroll_information", "contest.id = enroll_information.contest_id")
	session.Where("contest.id = ?", contestID)
	session.Where("enroll_information.state = ?", Pass)
	session.Select("contest.contest, contest_type.type contest_level.contest_level, enroll_information.*")
	session.Join("LEFT", "student", "student.student_id = enroll_information.student_id")
	session.Join("LEFT", "major", "student.major_id = major.major_id")
	session.Join("LEFT", "college", "college.college_id = student.college_id")
	session.Join("LEFT", "teacher as guidance", "guidance.teacher_id = enroll_information.guidance_teacher")
	session.Join("LEFT", "department", "teacher.department_id = department.department_id")
	session.Join("LEFT", "semester", "student.semester_id = semester.semester_id")
	session.Join("LEFT", "team", "enroll_information.team_id = team.team_id")
	session.Select("contest.contest, contest_type.type, contest_level.contest_level, enroll_information.*," +
		"major.major, student.*, department.department as t_department, guidance.name as t_name, guidance.title, college.college")

	if len(contest) > 0 {
		session.Where("contest.contest like ?", "%"+contest+"%")
	}
	//if len(startTime) > 0 && len(endTime) > 0 {
	//	session.Where("enroll_information.create_time >= ? AND enroll_information.create_time <= ?", startTime, endTime)
	//}
	if contestType != "" {
		session.Where("contest_type.type = ?", contestType)
	}
	if name != "" {
		session.Where("student.name like ?", "%"+name+"%")
	}
	if major != "" {
		session.Where("major.major like ?", "%"+major+"%")
	}
	if class != "" {
		session.Where("student.class like ?", "%"+class+"%")
	}
	if teamName != "" {
		session.Where("team.team_name like ?", "%"+teamName+"%")
	}

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
		list[i].StudentSchoolID = (*data)[i].StudentSchoolID
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].EnrollInformation.CreateTime)
		list[i].Phone = (*data)[i].Phone
		list[i].Email = (*data)[i].Email
		list[i].Major = (*data)[i].Major
		list[i].Team = (*data)[i].Team
		list[i].Title = (*data)[i].Title
		list[i].TeacherName = (*data)[i].TeacherName
		list[i].Department = (*data)[i].Department
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

func (self EnrollLogic) DepartmentManagerSearchEnroll(paginator *Paginator, contestID, userID int64, contest string, startTime string, endTime string, state int, contestType, name, major, class, studentSchoolID string) (*[]models.EnrollInformationReturn, int64, error) {
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

	session.Table("student")
	session.Join("RIGHT", "enroll_information", "student.student_id = enroll_information.student_id")
	session.Join("LEFT", "contest", "contest.id = enroll_information.contest_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("LEFT", "major", "major.major_id = student.major_id")
	session.Join("LEFT", "college", "college.college_id = student.college_id")
	session.Join("LEFT", "contest_level", "contest.contest_level_id = contest_level.contest_level_id")
	session.Join("LEFT", "teacher", "teacher.teacher_id = enroll_information.guidance_teacher")
	session.Join("LEFT", "department", "teacher.department_id = department.department_id")
	session.Join("LEFT", "team", "team.team_id = enroll_information.team_id")
	session.Select("enroll_information.id as e_id, enroll_information.*," +
		"student.*, contest_type.type, major.major, contest.*," +
		"contest_level.contest_level, college.college," +
		"teacher.name as teacher_name, teacher.title, department.department, team.team_name")
	session.Where("enroll_information.contest_id = ?", contestID)

	if contestID <= 0 {
		return nil, 0, err
	}
	if len(name) > 0 {
		session.Where("student.name like ?", "%"+name+"%")
	}
	if len(class) > 0 {
		session.Where("student.class like ?", "%"+class+"%")
	}
	if len(major) > 0 {
		session.Where("major.major like ?", "%"+major+"%")
	}
	if len(studentSchoolID) > 0 {
		session.Where("student.student_school_id = ?", studentSchoolID)
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
		list[i].ContestLevel = (*data)[i].ContestLevel
		list[i].TeacherName = (*data)[i].Teacher
		list[i].Major = (*data)[i].Major
		list[i].Team = (*data)[i].Team
		list[i].Department = (*data)[i].Department
		list[i].Title = (*data)[i].Title
		list[i].ContestType = (*data)[i].ContestType
		list[i].StudentSchoolID = (*data)[i].StudentSchoolID
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
		list[i].IsGroup = (*data)[i].IsGroup
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

func (self EnrollLogic) UpdateEnrollInformation(userID int64, form models.EnrollForm, role int) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("InsertEnrollInformation session.Begin() 发生错误:", err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("InsertEnrollInformation session.Close() 发生错误:", err)
		}
	}()

	if role == StudentRole {
		_, err := public.SearchAccountByID(userID)
		if err != nil {
			logging.L.Error(err)
			return err
		}
	} else if role == CmsManagerRole {

	} else {
		return errors.New("无权限")
	}

	enroll := models.EnrollInformation{}
	exist, err := MasterDB.Where("id = ?", form.ID).Get(&enroll)
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("un exist")
		return errors.New("报名信息不存在")
	}

	searchContest, err := public.SearchContestByID(enroll.ContestID)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	searchTeam := &models.Team{}
	if searchContest.IsGroup == 1 {
		searchTeam, err = public.SearchTeamByNameAndContest(form.TeamName, enroll.ContestID)
		if err != nil {
			logging.L.Error(err)
			return err
		}
	}

	searchTeacher, err := public.SearchTeacherByName(form.Teacher)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	_, err = public.SearchDepartmentByName(form.Department)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	newEnroll := &models.EnrollInformation{
		TeamID:          searchTeam.TeamID,
		Phone:           form.Phone,
		Email:           form.Email,
		GuidanceTeacher: searchTeacher.TeacherID,
		State:           form.State,
	}
	_, err = session.Where("id = ?", form.ID).Update(newEnroll)
	if err != nil {
		session.Rollback()
		logging.L.Error(err)
		return err
	}

	return session.Commit()
}
