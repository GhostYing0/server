package logic

import (
	"errors"
	"fmt"
	"github.com/polaris1119/times"
	. "server/database"
	"server/models"
	. "server/utils/mydebug"
)

type GradeLogic struct{}

var DefaultGradeLogic = GradeLogic{}

func (self GradeLogic) InsertGradeInformation(username string, contest string, grade string, certificate string, createTime string) error {
	if len(username) <= 0 {
		return errors.New("请填写姓名")
	}
	if len(contest) <= 0 {
		return errors.New("请填写竞赛名称")
	}
	if len(grade) <= 0 {
		return errors.New("请填写成绩")
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

	user := &models.Account{}
	exist, err := session.Table("account").In("username", username).Get(user)
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

	gradeInformation := &models.GradeInformation{
		Username:    user.Username,
		Contest:     contest,
		CreateTime:  models.FormatString2OftenTime(createTime),
		Certificate: certificate,
		Grade:       grade,
		State:       0,
	}

	_, err = session.Insert(gradeInformation)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		DPrintf("InsertEnrollInformation 上传成绩信息发生错误:", err)
		return err
	}

	return session.Commit()
}

func (self GradeLogic) Search(paginator *Paginator, grade string, contest string, startTime string, endTime string, state int, user_id int64, role int) (*[]models.ReturnGradeInformation, int64, error) {
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
	user := &models.Account{}
	exist, err := session.Where("id = ? and role = ?", user_id, role).Get(user)
	if err != nil {
		return nil, 0, err
	}
	if !exist {
		return nil, 0, errors.New("用户不存在")
	}
	if len(user.Username) > 0 {
		session.Table("grade").Where("username = ?", user.Username)
	}
	if len(grade) > 0 {
		session.Table("grade").Where("grade = ?", grade)
	}
	if len(contest) > 0 {
		session.Table("grade").Where("contest = ?", contest)
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		start := times.StrToLocalTime(startTime)
		end := times.StrToLocalTime(endTime)
		session.Table("grade").Where("createTime >= ? AND createTime <= ?", start, end)
	}
	if state > 0 {
		session.Table("grade").Where("state = ?", state)
	}

	temp := &[]models.GradeInformation{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(temp)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return nil, 0, fail
		}
		DPrintf("Search 查找成绩信息失败:", err)
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

	return &list, total, session.Rollback()
}

func (self GradeLogic) ProcessGrade(ids *[]int64, state int) (int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("ProcessEnroll session.Begin() 发生错误:", err)
		return 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("ProcessEnroll session.Close() 发生错误:", err)
		}
	}()

	var count int64

	for _, id := range *ids {
		if id < 1 {
			fmt.Println("非法id")
			continue
		}
		exist, err := session.Table("grade").Where("id = ?", id).Exist()
		if !exist {

		}
		if err != nil {
			DPrintf("ProcessEnroll 查询成绩信息发生错误:", err)
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return 0, fail
			}
			return count, err
		}
		affected, err := session.Where("id = ?", id).Update(models.GradeInformation{State: state})
		if err != nil {
			DPrintf("EnrollLogic Update 发生错误:", err)
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
