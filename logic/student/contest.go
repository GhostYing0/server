package student

import (
	. "server/database"
	. "server/logic"
	"server/models"
	"time"
)

type StudentContestLogic struct{}

var DefaultStudentContest = StudentContestLogic{}

func (self StudentContestLogic) Display(paginator *Paginator) (*[]models.DisplayContestForm, int64, error) {
	tx := MasterDB.NewSession()
	var total int64
	var err error

	var List []models.DisplayContestForm

	if paginator.PerPage() > 0 {
		total, err = tx.Table("contest").Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(&List)
		if err != nil {
			tx.Rollback()
			return nil, 0, err
		}
	} else {
		total, err = tx.Table("contest").Limit(10, 10*(paginator.CurPage()-1)).FindAndCount(&List)
		if err != nil {
			tx.Rollback()
			return nil, 0, err
		}
	}

	tx.Commit()
	return &List, total, err
}

func (self StudentContestLogic) FindGrade(user_id int, paginator *Paginator) (*[]models.ContestGrade, int64, error) {
	tx := MasterDB.NewSession()
	var total int64
	var err error

	var List []models.ContestGrade

	if paginator.PerPage() > 0 {
		total, err = tx.Table("registration").
			Where("contestant_id = ? and awards is not NULL", user_id).
			//Where("contestant_id = ?", user_id).
			Limit(paginator.PerPage(), paginator.Offset()).
			FindAndCount(&List)
		if err != nil {
			tx.Rollback()
			return nil, 0, err
		}
	} else {
		total, err = tx.Table("registration").
			Where("contestant_id = ? and awards is not NULL", user_id).
			//Where("contestant_id = ?", user_id).
			Limit(10, 10*(paginator.CurPage()-1)).
			FindAndCount(&List)
		if err != nil {
			tx.Rollback()
			return nil, 0, err
		}
	}

	tx.Commit()
	return &List, total, err
}

func (self StudentContestLogic) RegisterContest(user_id int, contest_int int) (string, error) {
	tx := MasterDB.NewSession()

	has, err := tx.Table("registration").
		Where("contestant_id = ?", user_id).
		And("contest_id = ?", contest_int).
		Exist()

	if err != nil {
		tx.Rollback()
		return "操作错误", err
	}
	if has {
		return "您已经报名该比赛", err
	}

	param := &models.ContestantInfo{
		ContestantID: user_id,
		ContestID:    contest_int,
		EntryTime:    time.Now(),
	}

	_, err = tx.Table("registration").Insert(param)
	if err != nil {
		tx.Rollback()
		return "操作出错", err
	}

	tx.Commit()
	return "报名成功", err
}
