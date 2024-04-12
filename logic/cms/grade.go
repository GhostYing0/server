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
)

type CmsGradeLogic struct{}

var DefaultGradeContest = CmsGradeLogic{}

func (self CmsGradeLogic) Display(paginator *Paginator, username, name, contest, school, startTime, endTime, grade string, state int) (*[]models.ReturnGradeInformation, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("DisplayGrade session.Begin() 发生错误:", err)
		logging.L.Error()
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("DisplayGrade session.Close() 发生错误:", err)
			logging.L.Error()
		}
	}()

	session.Table("grade")
	session.Join("LEFT", "student", "student.student_id = grade.student_id")
	session.Join("LEFT", "school", "school.school_id = grade.school_id")
	session.Join("LEFT", "contest", "contest.id = grade.contest_id")
	session.Join("LEFT", "account", "account.user_id = student.student_id")

	if username != "" {
		session.Where("account.username = ?", username)
	}
	if name != "" {
		session.Where("student.name = ?", name)
	}
	if contest != "" {
		session.Where("contest.contest = ?", contest)
	}
	if startTime != "" && endTime != "" {
		session.Where("grade.create_time > ? and grade.create_time < ?", startTime, endTime)
	}
	if grade != "" {
		session.Where("grade.grade = ?", grade)
	}
	if school != "" {
		session.Where("school.school = ?", school)
	}
	if state != -1 {
		session.Where("grade.state = ?", state)
	}

	data := &[]models.GradeStudentSchoolContestAccount{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		DPrintf("Display 查询成绩信息发生错误: ", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.ReturnGradeInformation, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].GradeInformation.ID
		list[i].Username = (*data)[i].Username
		list[i].Name = (*data)[i].Name
		list[i].Contest = (*data)[i].Contest
		list[i].School = (*data)[i].School.School
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].GradeInformation.CreateTime)
		list[i].Certificate = (*data)[i].Certificate
		list[i].Grade = (*data)[i].Grade
		list[i].State = (*data)[i].GradeInformation.State
	}

	return &list, total, session.Commit()
}

func (self CmsGradeLogic) Add(username string, contest string, grade string, createTime string, certificate string, state int) error {
	if len(username) <= 0 {
		return errors.New("请填写姓名")
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

	searchContest, err := public.SearchContestByName(contest)
	if err != nil {
		DPrintf("InsertGradeInformation 查询竞赛发生错误: ", err)
		logging.L.Error(err)
		return err
	}

	account, err := public.SearchAccountByUsernameAndRole(username, 1)
	if err != nil {
		DPrintf("InsertGradeInformation 查询参赛者失败:", err)
		logging.L.Error(err)
		return err
	}

	student, err := public.SearchStudentByID(account.UserID)
	if err != nil {
		DPrintf("InsertGradeInformation 查询学生:", err)
		logging.L.Error(err)
		return err
	}

	enroll := &models.GradeInformation{
		StudentID:   account.UserID,
		ContestID:   searchContest.ID,
		SchoolID:    student.SchoolID,
		CreateTime:  models.MysqlFormatString2String(createTime),
		Grade:       grade,
		Certificate: certificate,
		State:       state,
	}

	_, err = session.Insert(enroll)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(fail)
			return fail
		}
		logging.L.Error(err)
		DPrintf("InsertGradeInformation 添加成绩信息发生错误:", err)
		return err
	}

	return session.Commit()
}

func (self CmsGradeLogic) Update(id int64, username string, contest string, grade string, createTime string, certificate string, state int) error {
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

	searchContest, err := public.SearchContestByName(contest)
	if err != nil {
		DPrintf("InsertGradeInformation 查询竞赛发生错误: ", err)
		logging.L.Error(err)
		return err
	}

	account, err := public.SearchAccountByUsernameAndRole(username, 1)
	if err != nil {
		DPrintf("InsertGradeInformation 查询参赛者失败:", err)
		logging.L.Error(err)
		return err
	}

	student, err := public.SearchStudentByID(account.UserID)
	if err != nil {
		DPrintf("InsertGradeInformation 查询学生:", err)
		logging.L.Error(err)
		return err
	}

	enroll := &models.GradeInformation{
		StudentID:   account.UserID,
		ContestID:   searchContest.ID,
		SchoolID:    student.SchoolID,
		CreateTime:  models.MysqlFormatString2String(createTime),
		Grade:       grade,
		Certificate: certificate,
		State:       state,
	}
	_, err = session.Where("id = ?", id).Update(enroll)
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

func (self CmsGradeLogic) Delete(ids *[]int64) (int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("DeleteGradeInformation session.Begin() 发生错误:", err)
		return 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("DeleteGradeInformation session.Close() 发生错误:", err)
		}
	}()

	var count int64

	for _, id := range *ids {
		var gradeInformation models.GradeInformation
		if id < 1 {
			fmt.Println("非法id")
			continue
		}
		gradeInformation.ID = id
		affected, err := session.Delete(&gradeInformation)
		if err != nil {
			DPrintf("CmsGradeLogic Delete 发生错误:", err)
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

func (self CmsGradeLogic) GetGradeCount() (int64, error) {
	session := MasterDB.NewSession()

	count, err := session.Table("grade").Count()
	if err != nil {
		DPrintf("GetGradeCount Count 发生错误:", err)
		return count, err
	}
	return count, err
}
