package logic

import (
	"errors"
	. "server/database"
	"server/logic/public"
	"server/models"
	"server/utils/logging"
	. "server/utils/mydebug"
)

type ContestLogic struct{}

var DefaultContestLogic = ContestLogic{}

func (self ContestLogic) DisplayContest(paginator *Paginator, contest, contestType string) (*[]models.ContestReturn, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("DisplayContest session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("DisplayContest session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	session.Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id")
	session.Where("contest.state = 1")

	if contest != "" {
		session.Where("contest = ?", contest)
	}
	if contestType != "" {
		searchContestType, err := public.SearchContestTypeByName(contestType)
		if err != nil {

		} else {
			session.Where("contest.contest_type_id = ?", searchContestType.ContestTypeID)
		}
	}

	data := &[]models.ContestInfoType{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		DPrintf("InsertEnrollInformation 查询竞赛发生错误: ", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.ContestReturn, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].ContestState = (*data)[i].ContestState
		list[i].ID = (*data)[i].ID
		list[i].State = (*data)[i].State
		list[i].Contest = (*data)[i].Contest
		list[i].ContestType = (*data)[i].ContestType.ContestType
		list[i].CreateTime = (*data)[i].CreateTime.String()
		list[i].StartTime = (*data)[i].StartTime.String()
		list[i].Deadline = (*data)[i].Deadline.String()
	}

	return &list, total, session.Commit()
}

func (self ContestLogic) ViewTeacherContest(paginator *Paginator, userID int64, contest, contestType string, state int) (*[]models.ContestReturn, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("DisplayContest session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("DisplayContest session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return nil, 0, err
	}

	session.Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id")
	session.Where("contest.teacher_id = ?", account.UserID)

	if state != -1 {
		session.Where("contest.state = ?", state)
	}
	if contest != "" {
		session.Where("contest = ?", contest)
	}
	if contestType != "" {
		searchContestType, err := public.SearchContestTypeByName(contestType)
		if err != nil {

		} else {
			session.Where("contest.contest_type_id = ?", searchContestType.ContestTypeID)
		}
	}

	data := &[]models.ContestInfoType{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		DPrintf("InsertEnrollInformation 查询竞赛发生错误: ", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.ContestReturn, len(*data))
	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].ID
		list[i].ContestState = (*data)[i].ContestState
		list[i].State = (*data)[i].State
		list[i].Contest = (*data)[i].Contest
		list[i].ContestType = (*data)[i].ContestType.ContestType
		list[i].CreateTime = (*data)[i].CreateTime.String()
		list[i].StartTime = (*data)[i].StartTime.String()
		list[i].Deadline = (*data)[i].Deadline.String()
	}

	return &list, total, session.Commit()
}

func (self ContestLogic) UpdateContest(id, userID int64, contest, contestType, startTime, deadline string, contestState, state int) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("ProcessEnroll session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			logging.L.Error(err)
			DPrintf("ProcessEnroll session.Close() 发生错误:", err)
		}
	}()

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	session.Table("contest")
	session.Where("teacher_id = ?", account.UserID)
	exist, err := session.Where("id = ?", id).Exist()
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("竞赛不存在")
		return errors.New("竞赛不存在")
	}

	searchContestType, err := public.SearchContestTypeByName(contestType)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	updateContest := &models.ContestInfo{
		Contest:      contest,
		ContestType:  searchContestType.ContestTypeID,
		ContestState: contestState,
		State:        state,
	}

	if startTime != "" {
		updateContest.StartTime = models.FormatString2OftenTime(startTime)
	}
	if deadline != "" {
		updateContest.Deadline = models.FormatString2OftenTime(deadline)
	}
	_, err = session.Where("id = ?", id).Update(updateContest)
	if err != nil {
		fail := session.Rollback()
		if err != nil {
			logging.L.Error(err)
			return fail
		}
		logging.L.Error(err)
		return err
	}
	return session.Commit()
}

func (self ContestLogic) UploadContest(userID int64, contest, contestType, startTime, deadline string, describe *string) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("UploadContest session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			logging.L.Error(err)
			DPrintf("UploadContest session.Close() 发生错误:", err)
		}
	}()

	account, err := public.SearchAccountByID(userID)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	searchContestType, err := public.SearchContestTypeByName(contestType)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	teacher, err := public.SearchTeacherByID(account.UserID)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	newContest := &models.ContestInfo{
		TeacherID:    account.UserID,
		Contest:      contest,
		ContestType:  searchContestType.ContestTypeID,
		SchoolID:     teacher.SchoolID,
		CollegeID:    teacher.CollegeID,
		ContestState: 2,
		CreateTime:   models.NewOftenTime(),
		Describe:     *describe,
		State:        3,
	}

	if startTime != "" {
		newContest.StartTime = models.FormatString2OftenTime(startTime)
	}
	if deadline != "" {
		newContest.Deadline = models.FormatString2OftenTime(deadline)
	}
	_, err = session.Insert(newContest)
	if err != nil {
		fail := session.Rollback()
		if err != nil {
			logging.L.Error(err)
			return fail
		}
		logging.L.Error(err)
		return err
	}
	return session.Commit()
}
