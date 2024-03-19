package cms

import (
	"errors"
	"fmt"
	. "server/database"
	. "server/logic"
	"server/models"
	. "server/utils/mydebug"
)

type CmsRegistrationLogic struct{}

var DefaultRegistrationContest = CmsRegistrationLogic{}

func (self CmsRegistrationLogic) Display(paginator *Paginator) (*[]models.EnrollInformation, int64, error) {
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

	list := &[]models.EnrollInformation{}

	total, err := session.Table("registration").Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(list)
	if err != nil {
		DPrintf("Display 查询报名信息发生错误: ", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return nil, 0, fail
		}
		return nil, 0, err
	}

	return list, total, session.Commit()
}

func (self CmsRegistrationLogic) Add(username string, teamID int64, contestID int64, create_time string, school string, phone string, email string) error {
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

func (self CmsRegistrationLogic) Update(param *models.EnrollInformation) error {
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

	has, err := session.Where("id = ?", param.ID).Exist()
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

	_, err = session.Where("id = ?", param.ID).Update(param)
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

func (self CmsRegistrationLogic) Delete(ids *[]int64) (int64, error) {
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
