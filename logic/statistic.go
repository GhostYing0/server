package logic

import (
	"errors"
	. "server/database"
	"server/logic/public"
	"server/utils/e"
	"server/utils/logging"
	"time"
)

type StatisticLogic struct{}

var DefaultStatisticLogic = StatisticLogic{}

func (self StatisticLogic) StudentStatistic(userID int64) (map[string]int64, error) {
	account, err := public.SearchAccountByID(userID)
	data := make(map[string]int64)

	startTime := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
	endTime := time.Date(time.Now().Year()+1, time.January, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")

	contestCount, err := MasterDB.
		Table("contest").
		Where("contest.state = ? and contest.contest_state = ? and start_time > ? and start_time < ?", e.Pass, e.EnrollOpen, startTime, endTime).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	enrollPassCount, err := MasterDB.
		Table("enroll_information").
		Where("enroll_information.state = ? and student_id = ?", e.Pass, account.UserID).
		Count()

	enrollRejectCount, err := MasterDB.
		Table("enroll_information").
		Where("enroll_information.state = ? and student_id = ?", e.Reject, account.UserID).
		Count()

	enrollProcessCount, err := MasterDB.
		Table("enroll_information").
		Where("enroll_information.state = ? and student_id = ?", e.Processing, account.UserID).
		Count()

	gradeCount, err := MasterDB.
		Table("grade").
		Where("grade.state = ? and student_id = ?", e.Pass, account.UserID).
		Count()

	data["contest_count"] = contestCount
	data["enroll_pass_count"] = enrollPassCount
	data["enroll_reject_count"] = enrollRejectCount
	data["enroll_process_count"] = enrollProcessCount
	data["grade_count"] = gradeCount
	return data, err
}

func (self StatisticLogic) TeacherStatistic(userID int64) (map[string]int64, error) {
	account, err := public.SearchAccountByID(userID)
	data := make(map[string]int64)

	contestPassCount, err := MasterDB.
		Table("contest").
		Where("contest.state = ? and teacher_id = ?", e.Pass, account.UserID).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	contestRejectCount, err := MasterDB.
		Table("contest").
		Where("contest.state = ? and teacher_id = ?", e.Reject, account.UserID).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	contestProcessCount, err := MasterDB.
		Table("contest").
		Where("contest.state = ? and teacher_id = ?", e.Processing, account.UserID).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	enrollCount, err := MasterDB.
		Table("contest").
		Where("contest.teacher_id = ?", account.UserID).
		Join("LEFT", "enroll_information", "enroll_information.contest_id = contest.id").
		Where("enroll_information.state = ?", e.Pass).
		Count()

	gradePassCount, err := MasterDB.
		Table("contest").
		Where("contest.teacher_id = ?", account.UserID).
		Join("LEFT", "grade", "grade.contest_id = contest.id").
		Where("grade.state = ?", e.Pass).
		Count()

	gradeRejectCount, err := MasterDB.
		Table("contest").
		Where("contest.teacher_id = ?", account.UserID).
		Join("LEFT", "grade", "grade.contest_id = contest.id").
		Where("grade.state = ?", e.Pass).
		Count()

	gradeProcessCount, err := MasterDB.
		Table("contest").
		Where("contest.teacher_id = ?", account.UserID).
		Join("LEFT", "grade", "grade.contest_id = contest.id").
		Where("grade.state = ?", e.Pass).
		Count()

	data["contest_pass_count"] = contestPassCount
	data["contest_reject_count"] = contestRejectCount
	data["contest_process_count"] = contestProcessCount
	data["enroll_count"] = enrollCount
	data["grade_pass_count"] = gradePassCount
	data["grade_reject_count"] = gradeRejectCount
	data["grade_process_count"] = gradeProcessCount
	return data, err
}

func (self StatisticLogic) DepartmentStatistic(userID int64) (map[string]int64, error) {
	account, err := public.SearchDepartmentManagerByID(userID)
	data := make(map[string]int64)

	contestPassCount, err := MasterDB.
		Table("teacher").
		Where("teacher.department_id = ?", account.DepartmentID).
		Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id").
		Where("contest.state = ?", e.Pass).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	contestRejectCount, err := MasterDB.
		Table("teacher").
		Where("teacher.department_id = ?", account.DepartmentID).
		Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id").
		Where("contest.state = ?", e.Reject).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	contestProcessCount, err := MasterDB.
		Table("teacher").
		Where("teacher.department_id = ?", account.DepartmentID).
		Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id").
		Where("contest.state = ?", e.Processing).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	enrollPassCount, err := MasterDB.
		Table("teacher").
		Where("teacher.department_id = ?", account.DepartmentID).
		Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id").
		Join("LEFT", "enroll_information", "contest.id = enroll_information.contest_id").
		Where("enroll_information.state = ?", e.Pass).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	enrollProcessCount, err := MasterDB.
		Table("teacher").
		Where("teacher.department_id = ?", account.DepartmentID).
		Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id").
		Join("LEFT", "enroll_information", "contest.id = enroll_information.contest_id").
		Where("enroll_information.state = ?", e.Processing).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	enrollRejectCount, err := MasterDB.
		Table("teacher").
		Where("teacher.department_id = ?", account.DepartmentID).
		Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id").
		Join("LEFT", "enroll_information", "contest.id = enroll_information.contest_id").
		Where("enroll_information.state = ?", e.Reject).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	gradePassCount, err := MasterDB.
		Table("teacher").
		Where("teacher.department_id = ?", account.DepartmentID).
		Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id").
		Join("LEFT", "grade", "contest.id = grade.contest_id").
		Where("grade.state = ?", e.Pass).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	gradeProcessCount, err := MasterDB.
		Table("teacher").
		Where("teacher.department_id = ?", account.DepartmentID).
		Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id").
		Join("LEFT", "grade", "contest.id = grade.contest_id").
		Where("grade.state = ?", e.Processing).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	gradeRejectCount, err := MasterDB.
		Table("teacher").
		Where("teacher.department_id = ?", account.DepartmentID).
		Join("LEFT", "contest", "contest.teacher_id = teacher.teacher_id").
		Join("LEFT", "grade", "contest.id = grade.contest_id").
		Where("grade.state = ?", e.Reject).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	data["contest_pass_count"] = contestPassCount
	data["contest_reject_count"] = contestRejectCount
	data["contest_process_count"] = contestProcessCount
	data["enroll_pass_count"] = enrollPassCount
	data["enroll_process_count"] = enrollProcessCount
	data["enroll_reject_count"] = enrollRejectCount
	data["grade_pass_count"] = gradePassCount
	data["grade_reject_count"] = gradeRejectCount
	data["grade_process_count"] = gradeProcessCount
	return data, err
}

func (self StatisticLogic) ManagerStatistic(userID int64) (map[string]int64, error) {
	exist, err := MasterDB.Table("cms_account").Where("id = ?", userID).Exist()
	if !exist {
		return nil, errors.New("管理员不存在")
	}
	if err != nil {
		return nil, err
	}
	data := make(map[string]int64)

	contestPassCount, err := MasterDB.
		Table("contest").
		Where("contest.state = ?", e.Pass).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	contestCount, err := MasterDB.
		Table("contest").
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	enrollCount, err := MasterDB.
		Table("enroll_information").
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	enrollPassCount, err := MasterDB.
		Table("enroll_information").
		Where("enroll_information.state = ?", e.Pass).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	gradeCount, err := MasterDB.
		Table("grade").
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	gradePassCount, err := MasterDB.
		Table("grade").
		Where("grade.state = ?", e.Pass).
		Count()
	if err != nil {
		logging.L.Error(err)
	}

	data["contest_pass_count"] = contestPassCount
	data["contest_count"] = contestCount
	data["enroll_pass_count"] = enrollPassCount
	data["enroll_count"] = enrollCount
	data["grade_pass_count"] = gradePassCount
	data["grade_count"] = gradeCount
	return data, err
}
