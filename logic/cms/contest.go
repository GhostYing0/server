package cms

import (
	"fmt"
	. "server/database"
	. "server/logic"
	"server/models"
	. "server/utils/mydebug"
	"time"
)

type CmsContestLogic struct{}

var DefaultCmsContest = CmsContestLogic{}

func (self CmsContestLogic) Display(paginator *Paginator) (*[]models.ContestReturn, int64, error) {
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

	var List []models.ContestReturn

	total, err = session.Table("contest").Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(&List)
	if err != nil {
		return nil, 0, err
	}
	for _, l := range List {
		fmt.Println(time.ParseInLocation("2006-01-02 15:04:05", l.CreateTime, time.Local))
	}

	return &List, total, session.Commit()
}

func (self CmsContestLogic) InsertContest(contest string, contestType string, stateTime string, deadline string) (string, error) {
	if contest == "" || contestType == "" || stateTime == "" || deadline == "" {
		return "竞赛信息不能为空", nil
	}

	StartTime := models.FormatString2OftenTime(stateTime)
	DeadlineTime := models.FormatString2OftenTime(deadline)

	NewContest := &models.NewContest{
		Contest:     contest,
		ContestType: contestType,
		StartTime:   StartTime,
		Deadline:    DeadlineTime,
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

	has, err := session.Table("contest").Where("name = ? and type = ?", contest, contestType).Exist()
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
		Contest:     Name,
		ContestType: Type,
		StartTime:   TimeStartDate,
		Deadline:    TimeDeadline,
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
