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
		UserID:      user.ID,
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

func (self GradeLogic) Search(paginator *Paginator, username string, userID int64, contest string, startTime string, endTime string, state int, user_id int64, role int) (*[]models.GradeInformation, int64, error) {
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
	if role != 0 {
		session.Table("grade").Where("user_id = ?", user_id)
	}
	if len(username) > 0 {
		session.Table("grade").Where("username = ?", username)
	}
	if userID > 0 {
		session.Table("grade").Where("user_id = ?", userID)
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

	data := &[]models.GradeInformation{}

	total, err := session.Limit(paginator.PerPage(), paginator.Offset()).FindAndCount(data)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return data, 0, fail
		}
		DPrintf("Search 查找成绩信息失败:", err)
		return data, 0, err
	}

	return data, total, session.Rollback()
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
