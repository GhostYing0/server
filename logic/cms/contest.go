package cms

import (
	"fmt"
	. "server/database"
	. "server/logic"
	"server/models"
)

type CmsContestLogic struct{}

var DefaultCmsContest = CmsContestLogic{}

func (self CmsContestLogic) Display(paginator *Paginator) (*[]models.DisplayContestForm, int64, error) {
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

func (self CmsContestLogic) InsertContest(Name string, Type string, StartDate string, Deadline string) (string, error) {
	if Name == "" || Type == "" || Deadline == "" || StartDate == "" {
		return "竞赛信息不能为空", nil
	}

	StartTime := models.FormatString2OfenTime(StartDate)
	DeadlineTime := models.FormatString2OfenTime(Deadline)

	NewContest := &models.NewContest{
		Name:      Name,
		Type:      Type,
		StartDate: StartTime,
		Deadline:  DeadlineTime,
	}

	tx := MasterDB.NewSession()

	has, err := tx.Table("contest").Where("name = ? and type = ?", Name, Type).Exist()
	if err != nil {
		fmt.Println("InsertContestInfo Exist error:", err)
		return "操作错误", err
	}

	if has {
		return "竞赛已存在", err
	}

	_, err = tx.Table("contest").Insert(NewContest)
	if err != nil {
		tx.Rollback()
		fmt.Println("InsertContestInfo Insert error:", err)
		return "操作错误", err
	}
	tx.Commit()
	return "操作成功", err
}

func (self CmsContestLogic) UpdateContest(ID int64, Name string, Type string, StartDate string, Deadline string) (string, error) {
	tx := MasterDB.NewSession()

	has, err := tx.Table("contest").Where("id = ?", ID).Exist()
	if err != nil {
		fmt.Println("UpdateContestInfo Exist error:", err)
		return "操作错误", err
	}
	if !has {
		return "竞赛不存在", err
	}

	tx.Table("contest").Where("id = ?", ID)

	var TimeStartDate models.OftenTime
	var TimeDeadline models.OftenTime
	if len(StartDate) > 0 {
		TimeStartDate = models.FormatString2OfenTime(StartDate)
		if err != nil {
			fmt.Println("UpdateContestInfo StartDate time.Parse error:", err)
			return "时间解析出错", err
		}
	}
	if len(Deadline) > 0 {
		TimeDeadline = models.FormatString2OfenTime(Deadline)
		if err != nil {
			fmt.Println("UpdateContestInfo Deadline time.Parse error:", err)
			return "时间解析出错", err
		}
	}

	param := &models.ContestInfo{
		Name:      Name,
		Type:      Type,
		StartDate: TimeStartDate,
		Deadline:  TimeDeadline,
	}

	if len(Name) > 0 {
		_, err = tx.Table("contest").Where("id = ?", ID).Cols("name").Update(param)
		if err != nil {
			tx.Rollback()
			fmt.Println("UpdateContestInfo Update Name error:", err)
			return "操作错误", err
		}
	}
	if len(Type) > 0 {
		_, err = tx.Table("contest").Where("id = ?", ID).Cols("type").Update(param)
		if err != nil {
			tx.Rollback()
			fmt.Println("UpdateContestInfo Update Type error:", err)
			return "操作错误", err
		}
	}
	if len(StartDate) > 0 {
		_, err = tx.Table("contest").Where("id = ?", ID).Cols("start_date").Update(param)
		if err != nil {
			tx.Rollback()
			fmt.Println("UpdateContestInfo Update StartDate error:", err)
			return "操作错误", err
		}
	}
	if len(Deadline) > 0 {
		_, err = tx.Table("contest").Where("id = ?", ID).Cols("deadline").Update(param)
		if err != nil {
			tx.Rollback()
			fmt.Println("UpdateContestInfo Update Deadline error:", err)
			return "操作错误", err
		}
	}

	tx.Commit()
	return "操作成功", err
}

func (self CmsContestLogic) DeleteContest(ids *[]int64) (string, error, int64) {
	tx := MasterDB.NewSession()
	var count int64

	for _, id := range *ids {
		var contest models.ContestInfo
		if id < 1 {
			fmt.Println("非法id")
			continue
		}
		contest.ID = id
		affected, err := tx.Table("contest").Delete(&contest)
		if err != nil {
			tx.Rollback()
			return "操作出错", err, 0
		}
		if affected > 0 {
			count += affected
		}
	}

	tx.Commit()
	return "操作成功", nil, count
}
