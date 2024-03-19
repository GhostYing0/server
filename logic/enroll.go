package logic

import (
	"errors"
	"github.com/polaris1119/times"
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

func (self EnrollLogic) InsertEnrollInformation(username string, teamID int64, contestID int64, create_time string, school string, phone string, email string) error {
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
	exist, err := session.Where("id = ?", contestID).Get(contest)
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

	exist, err = session.Table("enroll_information").Where("contest = ? AND username = ?", contest.Name, user.Username).Exist()
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
		Username:   user.Username,
		UserID:     user.ID,
		Contest:    contest.Name,
		CreateTime: models.FormatString2OfenTime(create_time),
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

func (self EnrollLogic) Search(paginator *Paginator, username string, userID int64, contest string, startTime string, endTime string, school string, phone string, email string, state int, user_id int64, role int) (*[]models.EnrollInformation, int64, error) {
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

	// 不是管理员只能查看自己的信息
	if role != 0 {
		session.Table("enroll_information").Where("user_id = ?", user_id)
	}
	if len(username) > 0 {
		session.Table("enroll_information").Where("username = ?", username)
	}
	if userID > 0 {
		session.Table("enroll_information").Where("user_id = ?", userID)
	}
	if len(contest) > 0 {
		session.Table("enroll_information").Where("contest = ?", contest)
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		start := times.StrToLocalTime(startTime)
		end := times.StrToLocalTime(endTime)
		session.Table("enroll_information").Where("createTime >= ? AND createTime <= ?", start, end)
	}
	if len(school) > 0 {
		session.Table("enroll_information").Where("school = ?", school)
	}
	if len(phone) > 0 {
		session.Table("enroll_information").Where("phone = ?", phone)
	}
	if len(email) > 0 {
		session.Table("enroll_information").Where("email = ?", email)
	}
	if state >= 0 {
		session.Table("enroll_information").Where("state = ?", state)
	}

	data := &[]models.EnrollInformation{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return data, 0, fail
		}
		DPrintf("Search 查找报名信息失败:", err)
		return data, 0, err
	}

	return data, total, session.Rollback()
}
