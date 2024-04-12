package logic

import (
	"errors"
	"fmt"
	"github.com/polaris1119/times"
	. "server/database"
	"server/logic/public"
	"server/models"
	"server/utils/logging"
	. "server/utils/mydebug"
)

type GradeLogic struct{}

var DefaultGradeLogic = GradeLogic{}

func (self GradeLogic) InsertGradeInformation(user_id int64, contest, grade, certificate string) error {
	if len(contest) <= 0 {
		logging.L.Info("请填写竞赛名称")
		return errors.New("请填写竞赛名称")
	}
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

	account, err := public.SearchAccountByID(user_id)
	if err != nil {
		DPrintf("InsertEnrollInformation 查询用户失败:", err)
		logging.L.Error(err)
		return err
	}

	searchContest, err := public.SearchContestByName(contest)
	if err != nil {
		DPrintf("InsertEnrollInformation 查询竞赛失败:", err)
		logging.L.Error(err)
		return err
	}

	student, err := public.SearchStudentByID(account.UserID)
	if err != nil {
		DPrintf("InsertEnrollInformation 学生失败:", err)
		logging.L.Error(err)
		return err
	}

	gradeInformation := &models.GradeInformation{
		StudentID:   account.UserID,
		SchoolID:    student.SchoolID,
		ContestID:   searchContest.ID,
		CreateTime:  models.NewOftenTime().String(),
		Certificate: certificate,
		Grade:       grade,
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

	session.Table("account").Where("user_id = ?", account.UserID)
	session.Join("LEFT", "contest", "contest.teacher_id = account.user_id")
	session.Join("LEFT", "grade as g", "g.contest_id = contest.id")
	session.Join("RIGHT", "student", "student.student_id = g.student_id")

	if len(grade) > 0 {
		session.Where("g.grade = ?", grade)
	}
	if len(contest) > 0 {
		searchContest, err := public.SearchContestByName(contest)
		if err != nil {
			logging.L.Error(err)
		} else {
			session.Where("g.contest_id = ?", searchContest.ID)
		}
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		start := times.StrToLocalTime(startTime)
		end := times.StrToLocalTime(endTime)
		session.Where("g.createTime >= ? AND g.createTime <= ?", start, end)
	}
	if state > 0 {
		session.Where("g.state = ?", state)
	}

	data := &[]models.GradeStudentSchoolContestAccount{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).Select("g.id as id, g.*, account.*, student.*, contest.*").FindAndCount(data)
	if err != nil {
		logging.L.Error(err)
		DPrintf("Search 查找成绩信息失败:", err)
		return nil, 0, err
	}

	list := make([]models.ReturnGradeInformation, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].School = (*data)[i].School.School
		list[i].ID = (*data)[i].GradeInformation.ID
		list[i].Contest = (*data)[i].Contest
		list[i].Username = (*data)[i].Username
		list[i].Name = (*data)[i].Name
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].GradeInformation.CreateTime)
		list[i].Certificate = (*data)[i].Certificate
		list[i].Grade = (*data)[i].Grade
		list[i].State = (*data)[i].GradeInformation.State
	}

	return &list, total, session.Rollback()
}

func (self GradeLogic) ProcessGrade(id int64, state int) error {
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
	_, err = session.Where("id = ?", id).Update(models.GradeInformation{State: state})
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
