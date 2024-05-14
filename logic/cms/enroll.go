package cms

import (
	"errors"
	"fmt"
	. "server/database"
	. "server/logic"
	"server/logic/public"
	"server/models"
	"server/utils/e"
	"server/utils/logging"
	. "server/utils/mydebug"
)

type CmsEnrollLogic struct{}

var DefaultEnrollContest = CmsEnrollLogic{}

func (self CmsEnrollLogic) Display(paginator *Paginator, name string, contest, startTime, endTime, school, major string, state int) (*[]models.EnrollInformationReturn, int64, error) {
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

	session.Join("LEFT", "student", "student.student_id = enroll_information.student_id")
	session.Join("LEFT", "contest", "contest.id = enroll_information.contest_id")
	session.Join("LEFT", "account", "account.user_id = student.student_id")
	session.Join("LEFT", "school", "school.school_id = student.school_id")
	session.Join("LEFT", "major", "student.major_id = major.major_id")
	if name != "" {
		session.Where("student.name like ?", "%"+name+"%")
	}
	if contest != "" {
		session.Where("contest.contest like ?", "%"+contest+"%")
	}
	if startTime != "" && endTime != "" {
		session.Where("enroll_information.create_time > ? and enroll_information.create_time < ?", startTime, endTime)
	}
	if state != -1 {
		session.Where("enroll_information.state = ?", state)
	}
	if major != "" {
		session.Where("major.major like ?", "%"+major+"%")
	}
	if school != "" {
		searchSchool, err := public.SearchSchoolByName(school)
		if err != nil {
			logging.L.Error()
		} else {
			session.Where("enroll_information.school_id = ?", searchSchool.SchoolID)
		}
	}

	data := &[]models.EnrollContestStudent{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)

	if err != nil {
		DPrintf("Display 查询报名信息发生错误: ", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(fail)
			return nil, 0, fail
		}
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.EnrollInformationReturn, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].EnrollInformation.ID
		list[i].Username = (*data)[i].Username
		list[i].StudentID = (*data)[i].EnrollInformation.StudentID
		list[i].Name = (*data)[i].Name
		list[i].TeamID = (*data)[i].TeamID
		list[i].Contest = (*data)[i].Contest
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].EnrollInformation.CreateTime)
		list[i].School = (*data)[i].School
		list[i].Phone = (*data)[i].Phone
		list[i].Major = (*data)[i].Major
		list[i].Class = (*data)[i].Class
		list[i].StudentSchoolID = (*data)[i].StudentSchoolID
		list[i].Email = (*data)[i].Email
		list[i].State = (*data)[i].EnrollInformation.State
		list[i].RejectReason = (*data)[i].EnrollInformation.RejectReason
	}

	return &list, total, session.Commit()
}

func (self CmsEnrollLogic) Add(username string, name string, contest string, createTime string, school string, state int) error {
	if len(username) <= 0 {
		return errors.New("请填写姓名")
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

	contestInfo, err := public.SearchContestByName(contest)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	account, err := public.SearchAccountByUsernameAndRole(username, e.StudentRole)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	student, err := public.SearchStudentByName(name)
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
		Where("student_id = ? AND contest_id = ?", student.StudentID, contestInfo.ID).
		Exist()
	if err != nil {
		logging.L.Error(err)
		DPrintf("InsertEnrollInformation 查询重复报名发生错误: ", err)
		return err
	}
	if exist {
		logging.L.Error(err)
		DPrintf("请勿重复报名")
		return errors.New("请勿重复报名")
	}

	enroll := &models.NewEnroll{
		StudentID:  student.StudentID,
		ContestID:  contestInfo.ID,
		CreateTime: models.FormatString2OftenTime(createTime),
		SchoolID:   searchSchool.SchoolID,
		Phone:      account.Phone,
		Email:      account.Email,
		State:      state,
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
		DPrintf("InsertEnrollInformation 添加报名信息发生错误:", err)
		return err
	}

	return session.Commit()
}

func (self CmsEnrollLogic) Update(id int64, username string, name string, contest string, createTime string, school string, phone string, email string, state int) error {
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

	contestInfo, err := public.SearchContestByName(contest)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	_, err = public.SearchAccountByUsernameAndRole(username, 1)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	student, err := public.SearchStudentByName(name)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	has, err := session.Table("enroll_information").Where("id = ?", id).Exist()
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !has {
		logging.L.Error("报名信息不存在")
		return errors.New("报名信息不存在")
	}

	exist, err := session.Table("enroll_information").Where("student_id = ? && contest_id = ? && id != ?", student.StudentID, contestInfo.ID, id).Exist()
	if exist {
		logging.L.Error("已有相同报名信息")
		return errors.New("已有相同报名信息")
	}
	if err != nil {
		logging.L.Error(err)
		return err
	}

	_, err = session.Where("id = ?", id).Update(&models.NewEnroll{
		StudentID:  student.StudentID,
		ContestID:  contestInfo.ID,
		CreateTime: models.FormatString2OftenTime(createTime),
		Phone:      phone,
		Email:      email,
		State:      state,
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

func (self CmsEnrollLogic) Delete(ids *[]int64) (int64, error) {
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
		enrollInformation.ID = id
		affected, err := session.Delete(&enrollInformation)
		if err != nil {
			DPrintf("CmsRegistrationLogic Delete 发生错误:", err)
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

func (self CmsEnrollLogic) GetEnrollCount() (int64, error) {
	session := MasterDB.NewSession()

	count, err := session.Table("enroll").Count()
	if err != nil {
		DPrintf("GetEnrollCount Count 发生错误:", err)
		return count, err
	}
	return count, err
}
