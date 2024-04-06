package cms

import (
	"fmt"
	. "server/database"
	. "server/logic"
	"server/models"
	. "server/utils/mydebug"
)

type CmsContestLogic struct{}

var DefaultCmsContest = CmsContestLogic{}

func (self CmsContestLogic) Display(paginator *Paginator) (*[]models.DisplayContestForm, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("CmsContestLogic Display session.Begin() 发生错误:", err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("CmsContestLogic Display session.Close() 发生错误:", err)
		}
	}()

	var total int64
	var err error

	var List []models.DisplayContestForm

	if paginator.PerPage() > 0 {
		total, err = session.Table("contest").Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(&List)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return nil, 0, fail
			}
			return nil, 0, err
		}
	} else {
		total, err = session.Table("contest").Limit(10, 10*(paginator.CurPage()-1)).FindAndCount(&List)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return nil, 0, fail
			}
			return nil, 0, err
		}
	}

	return &List, total, session.Commit()
}

func (self CmsContestLogic) InsertContest(Name string, Type string, StartDate string, Deadline string) (string, error) {
	if Name == "" || Type == "" || Deadline == "" || StartDate == "" {
		return "竞赛信息不能为空", nil
	}

	StartTime := models.FormatString2OftenTime(StartDate)
	DeadlineTime := models.FormatString2OftenTime(Deadline)

	NewContest := &models.NewContest{
		Name:      Name,
		Type:      Type,
		StartDate: StartTime,
		Deadline:  DeadlineTime,
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("CmsContestLogic InsertContest session.Begin() 发生错误:", err)
		return "", err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("CmsContestLogic InsertContest session.Close() 发生错误:", err)
		}
	}()

	has, err := session.Table("contest").Where("name = ? and type = ?", Name, Type).Exist()
	if err != nil {
		fmt.Println("InsertContestInfo Exist error:", err)
		return "操作错误", err
	}

	if has {
		return "竞赛已存在", err
	}

	_, err = session.Table("contest").Insert(NewContest)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return "", fail
		}
		fmt.Println("InsertContestInfo Insert error:", err)
		return "操作错误", err
	}
	return "操作成功", session.Commit()
}

func (self CmsContestLogic) UpdateContest(ID int64, Name string, Type string, StartDate string, Deadline string) (string, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("CmsContestLogic UpdateContest session.Begin() 发生错误:", err)
		return "", err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("CmsContestLogic UpdateContest session.Close() 发生错误:", err)
		}
	}()

	has, err := session.Table("contest").Where("id = ?", ID).Exist()
	if err != nil {
		fmt.Println("UpdateContestInfo Exist error:", err)
		return "操作错误", err
	}
	if !has {
		return "竞赛不存在", err
	}

	session.Table("contest").Where("id = ?", ID)

	var TimeStartDate models.OftenTime
	var TimeDeadline models.OftenTime
	if len(StartDate) > 0 {
		TimeStartDate = models.FormatString2OftenTime(StartDate)
		if err != nil {
			fmt.Println("UpdateContestInfo StartDate time.Parse error:", err)
			return "时间解析出错", err
		}
	}
	if len(Deadline) > 0 {
		TimeDeadline = models.FormatString2OftenTime(Deadline)
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
		_, err = session.Table("contest").Where("id = ?", ID).Cols("name").Update(param)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return "", fail
			}
			fmt.Println("UpdateContestInfo Update Name error:", err)
			return "操作错误", err
		}
	}
	if len(Type) > 0 {
		_, err = session.Table("contest").Where("id = ?", ID).Cols("type").Update(param)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return "", fail
			}
			fmt.Println("UpdateContestInfo Update Type error:", err)
			return "操作错误", err
		}
	}
	if len(StartDate) > 0 {
		_, err = session.Table("contest").Where("id = ?", ID).Cols("start_date").Update(param)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return "", fail
			}
			fmt.Println("UpdateContestInfo Update StartDate error:", err)
			return "操作错误", err
		}
	}
	if len(Deadline) > 0 {
		_, err = session.Table("contest").Where("id = ?", ID).Cols("deadline").Update(param)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return "", fail
			}
			fmt.Println("UpdateContestInfo Update Deadline error:", err)
			return "操作错误", err
		}
	}

	return "操作成功", session.Commit()
}

func (self CmsContestLogic) DeleteContest(ids *[]int64) (string, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("CmsContestLogic UpdateContest session.Begin() 发生错误:", err)
		return "", 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("CmsContestLogic UpdateContest session.Close() 发生错误:", err)
		}
	}()
	var count int64

	for _, id := range *ids {
		var contest models.ContestInfo
		if id < 1 {
			fmt.Println("非法id")
			continue
		}
		contest.ID = id
		affected, err := session.Table("contest").Delete(&contest)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return "", 0, fail
			}
			return "操作出错", 0, err
		}
		if affected > 0 {
			count += affected
		}
	}

	return "操作成功", 0, session.Commit()
}
