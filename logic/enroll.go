package logic

import (
	"errors"
	"fmt"
	. "server/database"
	"server/models"
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

func (self EnrollLogic) InsertEnrollInformation(username string, teamID string, contestName string, create_time string, school string, phone string, email string) error {
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
		return nil, 0, errors.New("分页器为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Search session.Begin() 发生错误:", err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Search session.Close() 发生错误:", err)
		}
	}()

	if userID > 0 {
		session.Table("enroll_information").Where("user_id = ?", userID)
	}
	if len(contest) > 0 {
		session.Table("enroll_information").Where("contest = ?", contest)
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		session.Table("enroll_information").Where("create_time >= ? AND create_time <= ?", startTime, endTime)
	}
	if state >= 0 {
		session.Table("enroll_information").Where("state = ?", state)
	}
	fmt.Println(state)

	temp := &[]models.EnrollInformation{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(temp)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return nil, 0, fail
		}
		DPrintf("Search 查找报名信息失败:", err)
		return nil, 0, err
	}

	list := make([]models.EnrollInformationReturn, len(*temp))
	for i := 0; i < len(*temp); i++ {
		list[i].ID = (*temp)[i].ID
		//list[i].Username = (*temp)[i].Username
		//list[i].UserID = (*temp)[i].UserID
		list[i].TeamID = (*temp)[i].TeamID
		//list[i].Contest = (*temp)[i].Contest
		list[i].CreateTime = (*temp)[i].CreateTime.String()
		list[i].School = (*temp)[i].School
		list[i].Phone = (*temp)[i].Phone
		list[i].Email = (*temp)[i].Email
		list[i].State = (*temp)[i].State
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
