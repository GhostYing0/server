package cms

import (
	"errors"
	"fmt"
	. "server/database"
	. "server/logic"
	"server/models"
	. "server/utils/mydebug"
)

type CmsGradeLogic struct{}

var DefaultGradeContest = CmsGradeLogic{}

func (self CmsGradeLogic) Display(paginator *Paginator, username string, contest, startTime, endTime, grade string, state int) (*[]models.ReturnGradeInformation, int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("DisplayGrade session.Begin() 发生错误:", err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("DisplayGrade session.Close() 发生错误:", err)
		}
	}()

	fmt.Println(startTime)
	fmt.Println(endTime)
	if username != "" {
		session.Table("grade").Where("username = ?", username)
	}
	if contest != "" {
		session.Table("grade").Where("contest = ?", contest)
	}
	if startTime != "" && endTime != "" {
		session.Table("grade").Where("create_time > ? and create_time < ?", startTime, endTime)
	}
	if grade != "" {
		session.Table("grade").Where("grade = ?", grade)
	}
	if state != -1 {
		session.Table("grade").Where("state = ?", state)
	}

	temp := &[]models.GradeInformation{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(temp)
	//total, err = session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(list)
	if err != nil {
		DPrintf("Display 查询成绩信息发生错误: ", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return nil, 0, fail
		}
		return nil, 0, err
	}

	list := make([]models.ReturnGradeInformation, len(*temp))
	for i := 0; i < len(*temp); i++ {
		list[i].ID = (*temp)[i].ID
		list[i].Username = (*temp)[i].Username
		list[i].Contest = (*temp)[i].Contest
		list[i].CreateTime = (*temp)[i].CreateTime.String()
		list[i].Certificate = (*temp)[i].Certificate
		list[i].Grade = (*temp)[i].Grade
		list[i].State = (*temp)[i].State
	}

	return &list, total, session.Commit()
}

func (self CmsGradeLogic) Add(username string, contestName string, grade string, create_time string, certificate string, state int) error {
	if len(username) <= 0 {
		return errors.New("请填写姓名")
	}
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("InsertGradeInformation session.Begin() 发生错误:", err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("InsertGradeInformation session.Close() 发生错误:", err)
		}
	}()

	contest := &models.ContestInfo{}
	exist, err := session.Where("name = ?", contestName).Get(contest)
	if err != nil {
		DPrintf("InsertGradeInformation 查询竞赛发生错误: ", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}
	if !exist {
		DPrintf("InsertGradeInformation 竞赛不存在")
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
		DPrintf("InsertGradeInformation 查询参赛者失败:", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}
	if !exist {
		DPrintf("InsertGradeInformation 用户不存在")
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return errors.New("用户不存在")
	}

	enroll := &models.GradeInformation{
		Username:    user.Username,
		Contest:     contest.Contest,
		CreateTime:  models.FormatString2OftenTime(create_time),
		Grade:       grade,
		Certificate: certificate,
		State:       state,
	}

	_, err = session.Insert(enroll)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		DPrintf("InsertGradeInformation 添加成绩信息发生错误:", err)
		return err
	}

	return session.Commit()
}

func (self CmsGradeLogic) Update(id int64, username string, contestName string, grade string, create_time string, certificate string, state int) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("UpdateGradeInformation session.Begin() 发生错误:", err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("UpdateGradeInformation session.Close() 发生错误:", err)
		}
	}()

	has, err := session.Table("grade").Where("id = ?", id).Exist()
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
		return errors.New("成绩信息不存在")
	}

	exist, err := session.Table("account").In("username", username).Exist()
	if err != nil {
		DPrintf("InsertGradeInformation 查询参赛者失败:", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}
	if !exist {
		DPrintf("InsertGradeInformation 用户不存在")
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return errors.New("用户不存在")
	}

	newGrade := &models.GradeInformation{
		ID:          id,
		Contest:     contestName,
		Username:    username,
		CreateTime:  models.FormatString2OftenTime(create_time),
		Grade:       grade,
		Certificate: certificate,
		State:       state,
	}
	_, err = session.Where("id = ?", id).Update(newGrade)
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

func (self CmsGradeLogic) Delete(ids *[]int64) (int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("DeleteGradeInformation session.Begin() 发生错误:", err)
		return 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("DeleteGradeInformation session.Close() 发生错误:", err)
		}
	}()

	var count int64

	for _, id := range *ids {
		var gradeInformation models.GradeInformation
		if id < 1 {
			fmt.Println("非法id")
			continue
		}
		gradeInformation.ID = id
		affected, err := session.Delete(&gradeInformation)
		if err != nil {
			DPrintf("CmsGradeLogic Delete 发生错误:", err)
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

func (self CmsGradeLogic) GetGradeCount() (int64, error) {
	session := MasterDB.NewSession()

	count, err := session.Table("grade").Count()
	if err != nil {
		DPrintf("GetGradeCount Count 发生错误:", err)
		return count, err
	}
	return count, err
}
