package logic

import (
	"errors"
	"fmt"
	. "server/database"
	"server/logic/public"
	"server/models"
	"server/utils/logging"
	. "server/utils/mydebug"
)

type EnrollLogic struct{}

var DefaultEnrollLogic = EnrollLogic{}

func (self EnrollLogic) DisplayContest(paginator *Paginator) (*[]models.DisplayContestForm, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("DisplayContest session.Begin() 发生错误:", err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("DisplayContest session.Close() 发生错误:", err)
		}
	}()

	List := &[]models.DisplayContestForm{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(List)
	if err != nil {
		DPrintf("InsertEnrollInformation 查询竞赛发生错误: ", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return nil, 0, fail
		}
		return nil, 0, err
	}

	return List, total, session.Commit()
}

func (self EnrollLogic) InsertEnrollInformation(username, name, teamID, contest string, school string, phone string, email string) error {
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

	account, err := public.SearchAccountByUsername(username)
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
		State:      0,
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
		list[i].ID = (*data)[i].EnrollInformation.ID
		//list[i].Username = (*temp)[i].Username
		//list[i].UserID = (*temp)[i].UserID
		list[i].TeamID = (*data)[i].TeamID
		list[i].Contest = (*data)[i].Contest.Contest
		list[i].CreateTime = models.MysqlFormatString2String((*data)[i].EnrollInformation.CreateTime)
		list[i].Phone = (*data)[i].Phone
		list[i].Email = (*data)[i].Email
		list[i].State = (*data)[i].EnrollInformation.State

	}

	return &list, total, session.Rollback()
}

func (self EnrollLogic) ProcessEnroll(ids *[]int64, state int) (int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("ProcessEnroll session.Begin() 发生错误:", err)
		return 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("ProcessEnroll session.Close() 发生错误:", err)
		}
	}()

	var count int64

	for _, id := range *ids {
		if id < 1 {
			fmt.Println("非法id")
			continue
		}
		exist, err := session.Table("enroll_information").Where("id = ?", id).Exist()
		if !exist {

		}
		if err != nil {
			DPrintf("ProcessEnroll 查询竞赛信息发生错误:", err)
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return 0, fail
			}
			return count, err
		}
		affected, err := session.Where("id = ?", id).Update(models.EnrollInformation{State: state})
		if err != nil {
			DPrintf("EnrollLogic Update 发生错误:", err)
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
