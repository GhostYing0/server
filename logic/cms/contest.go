package cms

import (
	"errors"
	"fmt"
	. "server/database"
	. "server/logic"
	"server/models"
	. "server/utils/mydebug"
)

type CmsContestLogic struct{}

var DefaultCmsContest = CmsContestLogic{}

func (self CmsContestLogic) Display(paginator *Paginator, contest, contestType string, state int) (*[]models.ContestReturn, int64, error) {
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

	searchType := &models.ContestType{}

	if contest != "" {
		session.Table("contest").Where("contest = ?", contest)
	}
	if contestType != "" {
		exist, err := session.Where("type = ?", contestType).Get(searchType)
		if err != nil {
			DPrintf("Display查询竞赛类型:", err)
			return nil, 0, err
		}
		if !exist {
			return nil, 0, errors.New("竞赛类型不存在")
		}
		session.Table("contest").Where("contest_type_id = ?", searchType.ContestTypeID)
	}
	if state != -1 {
		session.Table("contest").Where("state = ?", state)
	}

	var total int64
	var err error

	var res []models.ContestContestType

	total, err = session.
		Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id").
		Limit(paginator.PerPage(), paginator.Offset()).
		FindAndCount(&res)
	if err != nil {
		return nil, 0, err
	}

	list := make([]models.ContestReturn, total)
	for i := 0; i < len(list); i++ {
		list[i].ID = res[i].Contest.ID
		list[i].State = res[i].Contest.State
		list[i].Contest = res[i].Contest.Contest
		list[i].ContestType = res[i].ContestType
		list[i].CreateTime = models.MysqlFormatString2String(res[i].Contest.CreateTime)
		list[i].StartTime = models.MysqlFormatString2String(res[i].Contest.StartTime)
		list[i].Deadline = models.MysqlFormatString2String(res[i].Contest.Deadline)
	}

	return &list, total, session.Commit()
}

func (self CmsContestLogic) InsertContest(contest, contestType, startTime, deadline string, state int) (string, error) {
	if contest == "" || contestType == "" || startTime == "" || deadline == "" {
		return "竞赛信息不能为空", nil
	}

	StartTime := models.FormatString2OftenTime(startTime)
	DeadlineTime := models.FormatString2OftenTime(deadline)

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

	searchType := &models.ContestType{}
	exist, err := session.Where("type = ?", contestType).Get(searchType)
	if err != nil {
		DPrintf("InsertContest查询竞赛类型:", err)
		return "", err
	}
	if !exist {
		return "", errors.New("竞赛类型不存在")
	}

	has, err := session.Table("contest").Where("contest = ? and contest_type_id = ?", contest, searchType.ContestTypeID).Exist()
	if err != nil {
		fmt.Println("InsertContestInfo Exist error:", err)
		return "操作错误", err
	}

	if has {
		return "竞赛已存在", err
	}

	NewContest := &models.NewContest{
		Contest:     contest,
		ContestType: searchType.ContestTypeID,
		CreateTime:  models.NewOftenTime(),
		StartTime:   StartTime,
		Deadline:    DeadlineTime,
		State:       state,
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

func (self CmsContestLogic) UpdateContest(id int64, contest, contestType, startTime, deadline string, state int) (string, error) {
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

	has, err := session.Table("contest").Where("id = ?", id).Exist()
	if err != nil {
		fmt.Println("UpdateContestInfo Exist error:", err)
		return "操作错误", err
	}
	if !has {
		return "竞赛不存在", err
	}

	searchType := &models.ContestType{}
	exist, err := session.Where("type = ?", contestType).Get(searchType)
	if err != nil {
		DPrintf("UpdateContest查询竞赛类型:", err)
		return "", err
	}
	if !exist {
		return "", errors.New("竞赛类型不存在")
	}

	_, err = session.Where("id = ?", id).Update(&models.ContestInfo{
		Contest:     contest,
		ContestType: searchType.ContestTypeID,
		StartTime:   models.FormatString2OftenTime(startTime),
		Deadline:    models.FormatString2OftenTime(deadline),
		State:       state,
	})
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return "", fail
		}
		return "更改失败", err
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
		affected, err := session.Delete(&contest)
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

func (self CmsContestLogic) GetContestCount() (int64, error) {
	session := MasterDB.NewSession()

	count, err := session.Table("contest").Count()
	if err != nil {
		DPrintf("GetContestCount Count 发生错误:", err)
		return count, err
	}
	return count, err
}
