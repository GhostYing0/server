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

func (self GradeLogic) InsertGradeInformation(user_id, id int64, grade, certificate, ps string) error {
	if len(grade) <= 0 {
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
	exist, err := session.Where("id = ?", id).Get(enroll)
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

	exist, err = session.Table("grade").Where("student_id = ? and contest_id = ? and state != ? and state != ?", enroll.StudentID, enroll.ContestID, e.Revoked, e.Reject).Exist()
	if exist {
		logging.L.Error("已上传成绩，不能重复上传")
		return errors.New("已上传成绩，不能重复上传")
	}
	if err != nil {
		logging.L.Error(err)
		return err
	}

	gradeInformation := &models.GradeInformation{
		StudentID:   enroll.StudentID,
		SchoolID:    student.SchoolID,
		ContestID:   enroll.ContestID,
		CreateTime:  models.NewOftenTime().String(),
		Certificate: certificate,
		Grade:       grade,
		UpdateTime:  models.NewOftenTime().String(),
		PS:          ps,
		State:       3,
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

func (self GradeLogic) Search(paginator *Paginator, grade string, contest string, startTime string, endTime string, state int, user_id int64, role int) (*[]models.ReturnGradeInformation, int64, error) {
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

	session.Table("grade").Where("student_id = ?", account.UserID)
	session.Join("LEFT", "contest", "contest.id = grade.contest_id")

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
		session.Where("grade.createTime >= ? AND grade.createTime <= ?", start, end)
	}
	if state > 0 {
		session.Where("grade.state = ?", state)
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
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].GradeInformation.CreateTime)
		list[i].Certificate = (*data)[i].Certificate
		list[i].Grade = (*data)[i].Grade
		list[i].State = (*data)[i].GradeInformation.State
		list[i].PS = (*data)[i].PS
		list[i].RejectReason = (*data)[i].RejectReason
	}

	return &list, total, session.Rollback()
}

func (self GradeLogic) TeacherSearch(paginator *Paginator, grade string, contest string, startTime string, endTime string, state int, user_id int64, role int) (*[]models.ReturnGradeInformation, int64, error) {
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

	session.Table("contest").Where("teacher_id = ?", account.UserID)
	session.Join("LEFT", "grade", "grade.contest_id = contest.id")
	session.Join("LEFT", "school", "grade.school_id = school.school_id")
	session.Join("LEFT", "student", "grade.student_id = student.student_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")

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
		session.Where("grade.createTime >= ? AND grade.createTime <= ?", start, end)
	}
	if state > 0 {
		session.Where("grade.state = ?", state)
	}

	data := &[]models.CurStudentGrade{}

	total, err := session.Where("teacher_id = ?", account.UserID).Select("grade.id as id, grade.*, school.*,contest.*,student.*,contest_type.*").Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
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
		list[i].Certificate = (*data)[i].Certificate
		list[i].Grade = (*data)[i].Grade
		list[i].School = (*data)[i].School
		list[i].ContestType = (*data)[i].ContestType
		list[i].Name = (*data)[i].Name
		list[i].State = (*data)[i].GradeInformation.State
		list[i].RejectReason = (*data)[i].GradeInformation.RejectReason
		list[i].PS = (*data)[i].PS
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
			Grade:       grade,
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

func (self GradeLogic) DepartmentManagerSearchGrade(paginator *Paginator, grade string, contest string, startTime string, endTime string, state int, contestID, user_id int64, role int) (*[]models.ReturnGradeInformation, int64, error) {
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
	session.Join("RIGHT", "grade", "grade.student_id = student.student_id")
	session.Join("LEFT", "contest", "contest.id = grade.contest_id")
	session.Join("LEFT", "contest_type", "contest_type.id = contest.contest_type_id")

	if contestID > 0 {
		session.Where("grade.contest_id = ?", contestID)
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
	if len(startTime) > 0 && len(endTime) > 0 {
		start := times.StrToLocalTime(startTime)
		end := times.StrToLocalTime(endTime)
		session.Where("grade.createTime >= ? AND grade.createTime <= ?", start, end)
	}
	if state > 0 {
		session.Where("grade.state = ?", state)
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
		list[i].Certificate = (*data)[i].Certificate
		list[i].Grade = (*data)[i].Grade
		list[i].School = (*data)[i].School
		list[i].ContestType = (*data)[i].ContestType
		list[i].Name = (*data)[i].Name
		list[i].State = (*data)[i].GradeInformation.State
		list[i].RejectReason = (*data)[i].GradeInformation.RejectReason
		list[i].PS = (*data)[i].PS
	}

	return &list, total, session.Rollback()
}
