package cms

import (
	"errors"
	"fmt"
	. "server/database"
	. "server/logic"
	"server/models"
	"server/utils/logging"
	. "server/utils/mydebug"
)

type CmsEnrollLogic struct{}

var DefaultEnrollContest = CmsEnrollLogic{}

func (self CmsEnrollLogic) Display(paginator *Paginator, name string, contest, startTime, endTime string, state int) (*[]models.EnrollInformationReturn, int64, error) {
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
	if name != "" {
		session.Where("name = ?", name)
	}
	if contest != "" {
		session.Where("contest = ?", contest)
	}
	if startTime != "" && endTime != "" {
		session.Table("enroll_information").Where("create_time > ? and create_time < ?", startTime, endTime)
	}
	if state != -1 {
		session.Table("enroll_information").Where("state = ?", state)
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
		list[i].Name = (*data)[i].Name
		list[i].TeamID = (*data)[i].TeamID
		list[i].Contest = (*data)[i].Contest.Contest
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].EnrollInformation.CreateTime.String())
		list[i].School = (*data)[i].School
		list[i].Phone = (*data)[i].Phone
		list[i].Email = (*data)[i].Email
		list[i].State = (*data)[i].EnrollInformation.State
	}

	return &list, total, session.Commit()
}

func (self CmsEnrollLogic) Add(username string, teamID string, contestName string, create_time string, school string, phone string, email string, state int) error {
	if len(username) <= 0 {
		return errors.New("请填写姓名")
	}
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

	contest := &models.ContestInfo{}
	exist, err := session.Where("name = ?", contestName).Get(contest)
	if err != nil {
		DPrintf("InsertEnrollInformation 查询竞赛发生错误: ", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}
	if !exist {
		DPrintf("InsertEnrollInformation 竞赛不存在")
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return errors.New("竞赛不存在")
	}

	user := &models.Account{}
	exist, err = session.Table("account").In("username", username).Get(user)
	if err != nil {
		DPrintf("InsertEnrollInformation 查询参赛者失败:", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}
	if !exist {
		DPrintf("InsertEnrollInformation 用户不存在")
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return errors.New("用户不存在")
	}

	exist, err = session.Table("enroll_information").Where("contest = ? AND username = ?", contest.Contest, user.Username).Exist()
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		DPrintf("InsertEnrollInformation 查询重复报名发生错误: ", err)
		return err
	}
	if exist {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		DPrintf("请勿重复报名")
		return errors.New("请勿重复报名")
	}

	enroll := &models.EnrollInformation{
		//Username:   user.Username,
		//UserID:     user.ID,
		//Contest:    contest.Contest,
		CreateTime: models.FormatString2OftenTime(create_time),
		School:     school,
		Phone:      phone,
		Email:      email,
		State:      state,
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

func (self CmsEnrollLogic) Update(id int64, username string, teamID string, contestName string, create_time string, school string, phone string, email string, state int) error {
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

	has, err := session.Table("enroll_information").Where("id = ?", id).Exist()
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}
	if !has {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return errors.New("报名信息不存在")
	}

	newEnroll := &models.EnrollInformation{
		ID: id,
		//Contest:    contestName,
		//Username:   username,
		TeamID:     teamID,
		CreateTime: models.FormatString2OftenTime(create_time),
		School:     school,
		Phone:      phone,
		Email:      email,
		State:      state,
	}
	_, err = session.Where("id = ?", id).Update(newEnroll)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
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
