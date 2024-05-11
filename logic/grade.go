package logic

import (
	"errors"
	"fmt"
	. "server/database"
	"server/logic/public"
	"server/models"
	"server/utils/e"
	"server/utils/logging"
	. "server/utils/mydebug"
	"time"

	"github.com/polaris1119/times"
)

type GradeLogic struct{}

var DefaultGradeLogic = GradeLogic{}

//
//func (self GradeLogic) InsertGradeInformation(user_id int64, contest, grade, certificate, ps string) error {
//	if len(contest) <= 0 {
//		logging.L.Info("请填写竞赛名称")
//		return errors.New("请填写竞赛名称")
//	}
//	if len(grade) <= 0 {
//		logging.L.Info("请填写成绩")
//		return errors.New("请填写成绩")
//	}
//
//	session := MasterDB.NewSession()
//	if err := session.Begin(); err != nil {
//		DPrintf("InsertGradeInformation session.Begin() 发生错误:", err)
//		logging.L.Error(err)
//		return err
//	}
//	defer func() {
//		err := session.Close()
//		if err != nil {
//			DPrintf("InsertGradeInformation session.Close() 发生错误:", err)
//			logging.L.Error(err)
//
//		}
//	}()
//
//	account, err := public.SearchAccountByID(user_id)
//	if err != nil {
//		DPrintf("InsertEnrollInformation 查询用户失败:", err)
//		logging.L.Error(err)
//		return err
//	}
//
//	searchContest, err := public.SearchContestByName(contest)
//	if err != nil {
//		DPrintf("InsertEnrollInformation 查询竞赛失败:", err)
//		logging.L.Error(err)
//		return err
//	}
//
//	student, err := public.SearchStudentByID(account.UserID)
//	if err != nil {
//		DPrintf("InsertEnrollInformation 学生失败:", err)
//		logging.L.Error(err)
//		return err
//	}
//
//	exist, err := session.Table("enroll_information").Where("student_id = ? and contest_id = ?", student.StudentID, searchContest.ID).Exist()
//	if !exist {
//		logging.L.Error("未报名该竞赛，无法上传成绩")
//		return errors.New("未报名该竞赛，无法上传成绩")
//	}
//	if err != nil {
//		logging.L.Error(err)
//		return err
//	}
//
//	exist, err = session.Table("grade").Where("student_id = ? and contest_id = ?", student.StudentID, searchContest.ID).Exist()
//	if exist {
//		logging.L.Error("已上传成绩，不能重复上传")
//		return errors.New("已上传成绩，不能重复上传")
//	}
//	if err != nil {
//		logging.L.Error(err)
//		return err
//	}
//
//	gradeInformation := &models.GradeInformation{
//		StudentID:   account.UserID,
//		SchoolID:    student.SchoolID,
//		ContestID:   searchContest.ID,
//		CreateTime:  models.NewOftenTime().String(),
//		Certificate: certificate,
//		Grade:       grade,
//		UpdateTime:  models.NewOftenTime().String(),
//		PS:          ps,
//		State:       3,
//	}
//
//	_, err = session.Insert(gradeInformation)
//	if err != nil {
//		fail := session.Rollback()
//		if fail != nil {
//			DPrintf("回滚失败")
//			logging.L.Error(fail)
//			return fail
//		}
//		logging.L.Error(err)
//		DPrintf("InsertEnrollInformation 上传成绩信息发生错误:", err)
//		return err
//	}
//
//	return session.Commit()
//}

func (self GradeLogic) UpdateStudentGrade(form *models.UpdateGradeForm) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("UpdateGradeInformation session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("UpdateGradeInformation session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	has, err := session.Table("grade").Where("id = ?", form.ID).Exist()
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !has {
		logging.L.Error("成绩信息不存在")
		return errors.New("成绩信息不存在")
	}

	teacher, err := public.SearchTeacherByName(form.GuidanceTeacher)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	grade := &models.GradeInformation{
		GuidanceTeacher: teacher.TeacherID,
		UpdateTime:      models.NewOftenTime().String(),
		Grade:           form.Prize,
		Certificate:     form.Certificate,
		State:           e.Processing,
	}
	_, err = session.Where("id = ?", form.ID).Update(grade)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(fail)
			return fail
		}
		logging.L.Error(err)
		DPrintf("cms enroll Update 失败:", err)
		return err
	}

	return session.Commit()
}

func (self GradeLogic) InsertGradeInformation(user_id, enrollID int64, grade int, rewardTime, certificate, teacherName, teahcerDepartment, teacherTitle string) error {
	fmt.Print("asdasd:", grade)
	if grade <= 0 {
		logging.L.Info("请填写成绩")
		return errors.New("请填写成绩")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("InsertGradeInformation session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("InsertGradeInformation session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	enroll := &models.EnrollInformation{}
	exist, err := session.Where("id = ?", enrollID).Get(enroll)
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("报名信息不存在")
		return errors.New("报名信息不存在")
	}

	account, err := public.SearchAccountByID(user_id)
	if err != nil {
		DPrintf("InsertEnrollInformation 查询用户失败:", err)
		logging.L.Error(err)
		return err
	}

	teacher, err := public.SearchTeacherByID(account.UserID)
	if err != nil {
		DPrintf("InsertEnrollInformation 失败:", err)
		logging.L.Error(err)
		return err
	}

	exist, err = session.Table("contest").Where("teacher_id = ?", teacher.TeacherID).Exist()
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("无该竞赛")
		return errors.New("无该竞赛")
	}

	student, err := public.SearchStudentByID(enroll.StudentID)
	if err != nil {
		DPrintf("InsertEnrollInformation 失败:", err)
		logging.L.Error(err)
		return err
	}

	department, err := public.SearchDepartmentByName(teahcerDepartment)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	teacherAccount := &models.Teacher{}
	exist, err = MasterDB.Where("name = ? and department_id = ?", teacherName, department.DepartmentID).Get(teacherAccount)
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("教师不存在")
		return errors.New("教师不存在")
	}

	prizeSearch, err := public.SearchPrizeByID(grade)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	exist, err = session.Table("grade").Where("enroll_id = ?", enrollID).Exist()
	//exist, err = session.Table("grade").Where("student_id = ? and contest_id = ? and state != ? and state != ?", enroll.StudentID, enroll.ContestID, e.Revoked, e.Reject).Exist()
	if exist {
		logging.L.Error("已上传成绩，不能重复上传")
		return errors.New("已上传成绩，不能重复上传")
	}
	if err != nil {
		logging.L.Error(err)
		return err
	}

	gradeInformation := &models.GradeInformation{
		StudentID:       enroll.StudentID,
		SchoolID:        student.SchoolID,
		ContestID:       enroll.ContestID,
		CreateTime:      models.NewOftenTime().String(),
		Certificate:     certificate,
		Grade:           prizeSearch.PrizeID,
		UpdateTime:      models.NewOftenTime().String(),
		GuidanceTeacher: teacherAccount.TeacherID,
		EnrollID:        enrollID,
		RewardTime:      rewardTime,
		//PS:          ps,
		State: 3,
	}

	_, err = session.Insert(gradeInformation)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(fail)
			return fail
		}
		logging.L.Error(err)
		DPrintf("InsertEnrollInformation 上传成绩信息发生错误:", err)
		return err
	}

	return session.Commit()
}

func (self GradeLogic) Search(paginator *Paginator, grade int, contest string, startTime string, endTime string, state int, user_id int64, role, contestLevel, isGroup int) (*[]models.ReturnGradeInformation, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		logging.L.Error("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	if user_id <= 0 {
		logging.L.Error("UserID Error")
		return nil, 0, errors.New("UserID Error")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		logging.L.Error(err)
		DPrintf("Search session.Begin() 发生错误:", err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Search session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	// 查看自身上传竞赛成绩
	account, err := public.SearchAccountByID(user_id)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}

	session.Table("grade").Where("student_id = ?", account.UserID)
	session.Join("LEFT", "contest", "contest.id = grade.contest_id")
	session.Join("LEFT", "prize", "prize.prize_id = grade.grade_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("LEFT", "contest_level", "contest_level.contest_level_id = contest.contest_level_id")

	if grade > 0 {
		session.Where("grade.grade_id = ?", grade)
	}
	if len(contest) > 0 {
		session.Where("contest.contest like ?", "%"+contest+"%")
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		start := times.StrToLocalTime(startTime)
		end := times.StrToLocalTime(endTime)
		session.Where("grade.createTime >= ? AND grade.createTime <= ?", start, end)
	}
	if state > 0 {
		session.Where("grade.state = ?", state)
	}
	if isGroup > 0 {
		session.Where("contest.is_group = ?", isGroup)
	}
	if contestLevel > 0 {
		session.Where("contest.contest_level_id = ?", contestLevel)
	}

	data := &[]models.CurStudentGrade{}

	total, err := session.Where("student_id = ?", account.UserID).Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	//total, err := session.Limit(paginator.PerPage(), paginator.Offset()).Select("g.id as id, g.*, account.*, student.*, contest.*, contest_type.*").FindAndCount(data)
	if err != nil {
		logging.L.Error(err)
		DPrintf("Search 查找成绩信息失败:", err)
		return nil, 0, err
	}

	list := make([]models.ReturnGradeInformation, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].ID
		list[i].Contest = (*data)[i].Contest
		list[i].RewardTime = models.MysqlFormatString2String((*data)[i].GradeInformation.RewardTime)
		list[i].Certificate = (*data)[i].Certificate
		list[i].Grade = (*data)[i].Grade
		list[i].State = (*data)[i].GradeInformation.State
		list[i].PS = (*data)[i].PS
		list[i].ContestLevel = (*data)[i].ContestLevel
		list[i].ContestType = (*data)[i].ContestType
		list[i].RejectReason = (*data)[i].RejectReason
	}

	return &list, total, session.Rollback()
}

func (self GradeLogic) TeacherSearch(paginator *Paginator, grade, contest, class, major, name, teamName string /* startTime string, endTime string,*/, state int, contestID, user_id int64, role, year int) (*[]models.ReturnGradeInformation, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		logging.L.Error("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	if user_id <= 0 {
		logging.L.Error("UserID Error")
		return nil, 0, errors.New("UserID Error")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		logging.L.Error(err)
		DPrintf("Search session.Begin() 发生错误:", err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Search session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	// 查看自身上传竞赛成绩
	account, err := public.SearchAccountByID(user_id)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}

	//session.Table("account").Where("user_id = ?", account.UserID)
	//session.Join("LEFT", "contest", "contest.teacher_id = account.user_id")
	//session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	//session.Join("LEFT", "grade as g", "g.contest_id = contest.id")
	//session.Join("RIGHT", "student", "student.student_id = g.student_id")

	//startTime := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	//endTime := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")

	session.Table("contest").Where("contest.teacher_id = ?", account.UserID)
	session.Join("LEFT", "grade", "grade.contest_id = contest.id")
	session.Join("LEFT", "prize", "grade.grade_id = prize.prize_id")
	session.Join("LEFT", "school", "grade.school_id = school.school_id")
	session.Join("LEFT", "student", "grade.student_id = student.student_id")
	session.Join("LEFT", "college", "student.college_id = college.college_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("LEFT", "contest_level", "contest_level.contest_level_id = contest.contest_level_id")
	session.Join("LEFT", "teacher", "grade.guidance_teacher = teacher.teacher_id")
	session.Join("LEFT", "department", "department.department_id = teacher.department_id")
	session.Join("LEFT", "major", "student.major_id = major.major_id")
	session.Join("LEFT", "enroll_information", "grade.enroll_id = enroll_information.id")
	session.Join("LEFT", "team", "enroll_information.team_id = team.team_id")
	session.Select("grade.id as id, teacher.name as t_name, grade.*, school.school,contest.*, college.college," +
		"student.*,contest_type.type, team.team_name, " +
		"department.department,teacher.title, prize.prize, major.major")
	if contestID > 0 {
		session.Where("contest.id = ?", contestID)
	}
	if len(grade) > 0 {
		session.Where("grade.grade = ?", grade)
	}
	if len(contest) > 0 {
		searchContest, err := public.SearchContestByName(contest)
		if err != nil {
			logging.L.Error(err)
		} else {
			session.Where("grade.contest_id = ?", searchContest.ID)
		}
	}
	if class != "" {
		session.Where("student.class like ?", "%"+class+"%")
	}
	if name != "" {
		session.Where("student.name like ?", "%"+name+"%")
	}
	if major != "" {
		session.Where("major.major like ?", "%"+major+"%")
	}
	if teamName != "" {
		session.Where("team.team_name like ?", "%"+teamName+"%")
	}
	//if len(startTime) > 0 && len(endTime) > 0 {
	//	start := times.StrToLocalTime(startTime)
	//	end := times.StrToLocalTime(endTime)
	//	session.Where("grade.create_time >= ? AND grade.create_time <= ?", start, end)
	//}
	if state > 0 {
		session.Where("grade.state = ?", state)
	}

	data := &[]models.CurStudentGrade{}

	//total, err := session.Where("teacher_id = ?", account.UserID).Select("grade.id as id, grade.*, school.*,contest.*,student.*,contest_type.*").Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	//total, err := session.Limit(paginator.PerPage(), paginator.Offset()).Select("g.id as id, g.*, account.*, student.*, contest.*, contest_type.*").FindAndCount(data)
	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		logging.L.Error(err)
		DPrintf("Search 查找成绩信息失败:", err)
		return nil, 0, err
	}

	list := make([]models.ReturnGradeInformation, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].ID
		list[i].Contest = (*data)[i].Contest
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].GradeInformation.CreateTime)
		list[i].RewardTime = models.MysqlFormatString2String((*data)[i].GradeInformation.RewardTime)
		list[i].Certificate = (*data)[i].Certificate
		list[i].Grade = (*data)[i].Grade
		list[i].School = (*data)[i].School
		list[i].Team = (*data)[i].Team
		list[i].ContestType = (*data)[i].ContestType
		list[i].Name = (*data)[i].Name
		list[i].State = (*data)[i].GradeInformation.State
		list[i].RejectReason = (*data)[i].GradeInformation.RejectReason
		list[i].PS = (*data)[i].PS
		list[i].ContestLevel = (*data)[i].ContestLevel
		list[i].GuidanceTeacher = (*data)[i].TeacherName
		list[i].Title = (*data)[i].Title
		list[i].Department = (*data)[i].Department
		list[i].Class = (*data)[i].Class
		list[i].Major = (*data)[i].Major
		list[i].StudentSchoolID = (*data)[i].StudentSchoolID
		list[i].College = (*data)[i].College
	}

	return &list, total, session.Rollback()
}

func (self GradeLogic) GetUserGrade(paginator *Paginator, grade string, contest string /* startTime string, endTime string,*/, state int, gradeID, user_id int64, role, year int) (*[]models.ReturnGradeInformation, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		logging.L.Error("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	if user_id <= 0 {
		logging.L.Error("UserID Error")
		return nil, 0, errors.New("UserID Error")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		logging.L.Error(err)
		DPrintf("Search session.Begin() 发生错误:", err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Search session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	// 查看自身上传竞赛成绩
	_, err := public.SearchAccountByID(user_id)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}

	//session.Table("account").Where("user_id = ?", account.UserID)
	//session.Join("LEFT", "contest", "contest.teacher_id = account.user_id")
	//session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	//session.Join("LEFT", "grade as g", "g.contest_id = contest.id")
	//session.Join("RIGHT", "student", "student.student_id = g.student_id")

	startTime := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	endTime := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")

	session.Table("grade").Where("grade.id = ?", gradeID)
	session.Join("LEFT", "contest", "contest.id = grade.contest_id")
	session.Join("LEFT", "prize", "grade.grade_id = prize.prize_id")
	session.Join("LEFT", "school", "grade.school_id = school.school_id")
	session.Join("LEFT", "student", "grade.student_id = student.student_id")
	session.Join("LEFT", "college", "student.college_id = college.college_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("LEFT", "contest_level", "contest_level.contest_level_id = contest.contest_level_id")
	session.Join("LEFT", "teacher", "grade.guidance_teacher = teacher.teacher_id")
	session.Join("LEFT", "department", "department.department_id = teacher.department_id")
	session.Join("LEFT", "major", "major.major_id = student.major_id")
	session.Join("LEFT", "contest_entry", "contest_entry.contest_entry_id = contest.contest_entry_id")
	session.Select("teacher.name as t_name, grade.*, school.school, contest.contest, college.college," +
		"student.*,contest_type.type, " +
		"department.department,teacher.title, prize.*, major.major, contest_entry.contest_entry")
	if len(grade) > 0 {
		session.Where("grade.grade = ?", grade)
	}
	if len(contest) > 0 {
		searchContest, err := public.SearchContestByName(contest)
		if err != nil {
			logging.L.Error(err)
		} else {
			session.Where("grade.contest_id = ?", searchContest.ID)
		}
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		start := times.StrToLocalTime(startTime)
		end := times.StrToLocalTime(endTime)
		session.Where("grade.create_time >= ? AND grade.create_time <= ?", start, end)
	}
	if state > 0 {
		session.Where("grade.state = ?", state)
	}

	data := &[]models.CurStudentGrade{}

	//total, err := session.Where("teacher_id = ?", account.UserID).Select("grade.id as id, grade.*, school.*,contest.*,student.*,contest_type.*").Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	//total, err := session.Limit(paginator.PerPage(), paginator.Offset()).Select("g.id as id, g.*, account.*, student.*, contest.*, contest_type.*").FindAndCount(data)
	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		logging.L.Error(err)
		DPrintf("Search 查找成绩信息失败:", err)
		return nil, 0, err
	}

	list := make([]models.ReturnGradeInformation, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].ID
		list[i].Contest = (*data)[i].Contest
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].GradeInformation.CreateTime)
		list[i].RewardTime = models.MysqlFormatString2String((*data)[i].GradeInformation.RewardTime)
		list[i].Certificate = (*data)[i].Certificate
		list[i].Grade = (*data)[i].Grade
		list[i].School = (*data)[i].School
		list[i].ContestType = (*data)[i].ContestType
		list[i].Name = (*data)[i].Name
		list[i].Major = (*data)[i].Major
		list[i].ContestEntry = (*data)[i].ContestEntry
		list[i].Prize = (*data)[i].Prize
		list[i].State = (*data)[i].GradeInformation.State
		list[i].RejectReason = (*data)[i].GradeInformation.RejectReason
		list[i].PS = (*data)[i].PS
		list[i].ContestLevel = (*data)[i].ContestLevel
		list[i].GuidanceTeacher = (*data)[i].TeacherName
		list[i].Title = (*data)[i].Title
		list[i].Department = (*data)[i].Department
		list[i].Class = (*data)[i].Class
		list[i].StudentSchoolID = (*data)[i].StudentSchoolID
		list[i].College = (*data)[i].College
	}

	return &list, total, session.Rollback()
}

func (self GradeLogic) ProcessGrade(id int64, state int, rejectReason string) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("ProcessEnroll session.Begin() 发生错误:", err)
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
		return errors.New("非法id")
	}
	exist, err := session.Table("grade").Where("id = ?", id).Exist()
	if !exist {
		logging.L.Error("不存在")
		return errors.New("不存在")
	}
	if err != nil {
		DPrintf("ProcessEnroll 查询成绩信息发生错误:", err)
		logging.L.Error(err)
		return err
	}

	newInfo := &models.GradeInformation{State: state}
	if state == e.Reject {
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

func (self GradeLogic) PassGrade(ids *[]int64, state int) (int64, error) {
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
		var gradeInformation models.GradeInformation
		if id < 1 {
			fmt.Println("非法id")
			continue
		}
		gradeInformation.State = e.Pass
		affected, err := session.Where("id = ?", id).Update(&gradeInformation)
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

func (self GradeLogic) Update(id int64, grade string, certificate string) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("UpdateGradeInformation session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("UpdateGradeInformation session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	has, err := session.Table("grade").Where("id = ?", id).Exist()
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !has {
		logging.L.Error("成绩信息不存在")
		return errors.New("成绩信息不存在")
	}

	_, err = session.
		Table("grade").
		Where("id = ?", id).
		Update(&models.GradeInformation{
			//Grade:       grade,
			Certificate: certificate,
			UpdateTime:  models.NewOftenTime().String(),
		})
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(fail)
			return fail
		}
		logging.L.Error(err)
		DPrintf("cms enroll Update 失败:", err)
		return err
	}

	return session.Commit()
}

func (self GradeLogic) DepartmentManagerSearchGrade(paginator *Paginator, grade int, contest string, startTime string, endTime string, state int, contestID, user_id int64, role int, name, major, class string) (*[]models.ReturnGradeInformation, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		logging.L.Error("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	if user_id <= 0 {
		logging.L.Error("UserID Error")
		return nil, 0, errors.New("UserID Error")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		logging.L.Error(err)
		DPrintf("Search session.Begin() 发生错误:", err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Search session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	// 查看本校同系部上传竞赛成绩
	account, err := public.SearchDepartmentManagerByID(user_id)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}

	//session.Table("student").Where("student.school_id = ? and student.college_id = ? and student.department_id = ?", account.SchoolID, account.CollegeID, account.DepartmentID)
	session.Table("student").Where("student.school_id = ?", account.SchoolID)
	session.Join("RIGHT", "grade", "grade.student_id = student.student_id")
	session.Join("LEFT", "contest", "contest.id = grade.contest_id")
	session.Join("LEFT", "contest_level", "contest.contest_level_id = contest_level.contest_level_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")
	session.Join("LEFT", "major", "major.major_id = student.major_id")
	session.Join("LEFT", "prize", "prize.prize_id = grade.grade_id")
	session.Join("LEFT", "college", "college.college_id = student.college_id")
	session.Join("LEFT", "teacher", "teacher.teacher_id = grade.guidance_teacher")
	session.Join("LEFT", "department", "department.department_id = teacher.department_id")
	session.Join("LEFT", "enroll_information", "enroll_information.id = grade.enroll_id")
	session.Select("teacher.name as t_name, teacher.title, contest.contest, department.department, grade.*, contest_level.contest_level, contest_type.type, major.major," +
		"college.college, student.*, prize.prize")

	if contestID > 0 {
		session.Where("grade.contest_id = ?", contestID)
	}
	if grade > 0 {
		session.Where("grade.grade_id = ?", grade)
	}
	if len(contest) > 0 {
		searchContest, err := public.SearchContestByName(contest)
		if err != nil {
			logging.L.Error(err)
		} else {
			session.Where("grade.contest_id = ?", searchContest.ID)
		}
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		start := times.StrToLocalTime(startTime)
		end := times.StrToLocalTime(endTime)
		session.Where("grade.createTime >= ? AND grade.createTime <= ?", start, end)
	}
	if state > 0 {
		session.Where("grade.state = ?", state)
	}
	if major != "" {
		session.Where("major.major like ?", "%"+major+"%")
	}
	if name != "" {
		session.Where("student.name like ?", "%"+name+"%")
	}
	if class != "" {
		session.Where("student.class like ?", "%"+class+"%")
	}

	data := &[]models.CurStudentGrade{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	//total, err := session.Limit(paginator.PerPage(), paginator.Offset()).Select("g.id as id, g.*, account.*, student.*, contest.*, contest_type.*").FindAndCount(data)
	if err != nil {
		logging.L.Error(err)
		DPrintf("Search 查找成绩信息失败:", err)
		return nil, 0, err
	}

	list := make([]models.ReturnGradeInformation, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].ID
		list[i].Contest = (*data)[i].Contest
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].GradeInformation.CreateTime)
		list[i].RewardTime = models.MysqlFormatString2String((*data)[i].GradeInformation.RewardTime)
		list[i].Certificate = (*data)[i].Certificate
		list[i].Grade = (*data)[i].Grade
		list[i].School = (*data)[i].School
		list[i].Major = (*data)[i].Major
		//list[i].Team = (*data)[i].Team
		list[i].ContestType = (*data)[i].ContestType
		list[i].College = (*data)[i].College
		list[i].StudentSchoolID = (*data)[i].StudentSchoolID
		list[i].Class = (*data)[i].Class
		list[i].Name = (*data)[i].Name
		list[i].State = (*data)[i].GradeInformation.State
		list[i].RejectReason = (*data)[i].GradeInformation.RejectReason
		list[i].PS = (*data)[i].PS
		list[i].GuidanceTeacher = (*data)[i].TeacherName
		list[i].Title = (*data)[i].Title
		list[i].Department = (*data)[i].Department
		list[i].ContestLevel = (*data)[i].ContestLevel

	}

	return &list, total, session.Rollback()
}
