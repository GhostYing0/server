package cms

import (
	. "server/database"
	. "server/logic"
	"server/models"
	"time"
)

type CmsRegistrationLogic struct{}

var DefaultRegistrationContest = CmsRegistrationLogic{}

func (self CmsRegistrationLogic) Display(paginator *Paginator) (*[]models.RegistrationInfo, int64, error) {
	tx := MasterDB.NewSession()
	var total int64
	var err error

	var List []models.RegistrationInfo

	if paginator.PerPage() > 0 {
		total, err = tx.Table("registration").Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(&List)
		if err != nil {
			tx.Rollback()
			return nil, 0, err
		}
	} else {
		total, err = tx.Table("registration").Limit(10, 10*(paginator.CurPage()-1)).FindAndCount(&List)
		if err != nil {
			tx.Rollback()
			return nil, 0, err
		}
	}

	tx.Commit()
	return &List, total, err
}

func (self CmsRegistrationLogic) AddRegistration(contestant string, contest string) (string, error) {
	tx := MasterDB.NewSession()

	has, err := tx.Table("account").Where("username = ?", contestant).Exist()
	if err != nil {
		return "操作错误", err
	}
	if !has {
		return "用户不存在", err
	}

	has, err = tx.Table("contest").Where("name = ?", contest).Exist()
	if err != nil {
		return "操作错误", err
	}
	if !has {
		return "竞赛不存在", err
	}

	var contestant_id []int
	err = tx.Table("account").
		Where("username = ?", contestant).
		Cols("id").
		Find(&contestant_id)
	if err != nil {
		return "操作出错", err
	}

	var contest_id []int
	err = tx.Table("contest").
		Where("name = ?", contest).
		Cols("id").
		Find(&contest_id)
	if err != nil {
		return "操作出错", err
	}

	param := &models.ContestantInfo{
		ContestantID: contestant_id[0],
		ContestID:    contest_id[0],
		EntryTime:    time.Now(),
	}

	has, err = tx.Table("registration").
		Where("contestant_id = ? and contest_id = ?", contestant_id[0], contest_id[0]).
		Exist()
	if has {
		return "用户已报名该比赛,不能重复报名", err
	}

	_, err = tx.Table("registration").Insert(param)
	if err != nil {
		return "操作出错", err
	}

	return "添加成功", err
}
