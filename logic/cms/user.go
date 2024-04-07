package cms

import (
	"errors"
	"fmt"
	. "server/database"
	. "server/logic"
	"server/models"
	. "server/utils/mydebug"
)

type CmsUserLogic struct{}

var DefaultCmsUser = CmsUserLogic{}

func (self CmsUserLogic) DisplayStudent(paginator *Paginator, username, gender, school, semester, college, class, name string) (*[]models.StudentReturn, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser Display session.Begin() 发生错误:", err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser Display session.Close() 发生错误:", err)
		}
	}()

	if username != "" {
		session.Table("account").Where("username = ?", username)
	}
	//查询学生
	session.Table("account").Where("role = ?", 1)
	if gender != "" {
		session.Table("student").Where("gender = ?", gender)
	}
	if school != "" {
		session.Table("student").Where("school = ?", school)
	}
	if semester != "" {
		session.Table("student").Where("name = ?", name)
	}
	if college != "" {
		session.Table("student").Where("college = ?", college)
	}
	if name != "" {
		session.Table("student").Where("name = ?", name)
	}
	if class != "" {
		session.Table("student").Where("class = ?", class)
	}

	data := &[]models.AccountStudent{}

	total, err := session.
		Join("LEFT", "student", "account.user_id = student.student_id").
		Limit(paginator.PerPage(), paginator.Offset()).
		FindAndCount(data)
	if err != nil {
		DPrintf("Search 查找成绩信息失败:", err)
		return nil, 0, err
	}

	list := make([]models.StudentReturn, total)

	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].ID
		list[i].Username = (*data)[i].Username
		list[i].Password = (*data)[i].Password
		list[i].Role = (*data)[i].Role
		list[i].StudentID = (*data)[i].StudentID
		list[i].Name = (*data)[i].Name
		list[i].Gender = (*data)[i].Gender
		list[i].Class = (*data)[i].Class
		list[i].Avatar = (*data)[i].Avatar

		searchSchool := &models.School{}
		searchSemester := &models.Semester{}
		searchCollege := &models.College{}
		_, err := session.Where("school_id = ?", (*data)[i].SchoolID).Get(searchSchool)
		if err != nil {
			list[i].School = "查询出错"
		}
		list[i].School = searchSchool.School
		_, err = session.Where("semester_id = ?", (*data)[i].SemesterID).Get(searchSemester)
		if err != nil {
			list[i].Semester = "查询出错"
		}
		list[i].Semester = searchSemester.Semester
		_, err = session.Where("college_id = ?", (*data)[i].CollegeID).Get(searchCollege)
		if err != nil {
			list[i].College = "查询出错"
		}
		list[i].College = searchCollege.College
	}

	return &list, total, session.Rollback()
}

func (self CmsUserLogic) DisplayTeacher(paginator *Paginator, username, gender, school, semester, college, class, name string) (*[]models.TeacherReturn, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser Display session.Begin() 发生错误:", err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser Display session.Close() 发生错误:", err)
		}
	}()

	if username != "" {
		session.Table("account").Where("username = ?", username)
	}
	//查询学生
	session.Table("account").Where("role = ?", 2)
	if gender != "" {
		session.Table("teacher").Where("gender = ?", gender)
	}
	if school != "" {
		session.Table("teacher").Where("school = ?", school)
	}
	if college != "" {
		session.Table("teacher").Where("college = ?", college)
	}
	if name != "" {
		session.Table("teacher").Where("name = ?", name)
	}

	data := &[]models.AccountTeacher{}

	total, err := session.
		Join("LEFT", "teacher", "account.user_id = teacher.teacher_id").
		Limit(paginator.PerPage(), paginator.Offset()).
		FindAndCount(data)
	if err != nil {
		DPrintf("Search 查找成绩信息失败:", err)
		return nil, 0, err
	}

	list := make([]models.TeacherReturn, total)

	for i := 0; i < len(*data); i++ {
		list[i].ID = (*data)[i].ID
		list[i].Username = (*data)[i].Username
		list[i].Password = (*data)[i].Password
		list[i].Role = (*data)[i].Role
		list[i].TeacherID = (*data)[i].TeacherID
		list[i].Name = (*data)[i].Name
		list[i].Gender = (*data)[i].Gender

		searchSchool := &models.School{}
		searchCollege := &models.College{}
		_, err := session.Where("school_id = ?", (*data)[i].SchoolID).Get(searchSchool)
		if err != nil {
			list[i].School = "查询出错"
		}
		list[i].School = searchSchool.School
		_, err = session.Where("college_id = ?", (*data)[i].CollegeID).Get(searchCollege)
		if err != nil {
			list[i].College = "查询出错"
		}
		list[i].College = searchCollege.College
	}

	return &list, total, session.Rollback()
}
func (self CmsUserLogic) AddUser(username string, password string, role int) (string, error) {
	if len(username) == 0 || len(password) == 0 {
		return "账号和密码不能为空 ", nil
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser Display session.Begin() 发生错误:", err)
		return err.Error(), err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser Display session.Close() 发生错误:", err)
		}
	}()

	exist, err := session.Table("account").Where("username = ? and role = ?", username, role).Exist()
	if exist {
		return "用户已存在", err
	}

	param := &models.OldUser{
		Username: username,
		Password: password,
		Role:     role,
	}

	_, err = session.Insert(param)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return err.Error(), fail
		}
		return "操作出错", err
	}

	return "操作成功", session.Commit()
}

func (self CmsUserLogic) UpdateUser(ID int64, newUsername string, newPassword string, role int) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser UpdateUser session.Begin() 发生错误:", err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser UpdateUser session.Close() 发生错误:", err)
		}
	}()

	exist, err := session.Table("account").Where("id = ?", ID).Exist()
	if !exist {
		return errors.New("用户不存在")
	}
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		DPrintf("UpdateUser 查询用户失败:", err)
		return err
	}

	param := &models.OldUser{
		Username: newUsername,
		Password: newPassword,
		Role:     role,
	}

	_, err = session.Where("id = ?", ID).Update(param)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		DPrintf("UpdateUser Update 更新用户失败:", err)
		return err
	}

	return session.Commit()
}

func (self CmsUserLogic) DeleteUser(ids *[]int64) (int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser DeleteUser session.Begin() 发生错误:", err)
		return 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser DeleteUser session.Close() 发生错误:", err)
		}
	}()

	var count int64

	for _, id := range *ids {
		var contest models.ContestInfo
		if id < 1 {
			fmt.Println("非法id")
			continue
		}
		contest.ID = int64(id)
		affected, err := session.Table("account").Delete(&contest)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return 0, fail
			}
			return 0, err
		}
		if affected > 0 {
			count += affected
		}
	}

	return count, session.Commit()
}

func (self CmsUserLogic) GetUserCount() (int64, error) {
	session := MasterDB.NewSession()

	count, err := session.Table("account").Count()
	if err != nil {
		DPrintf("GetUserCount Count 发生错误:", err)
		return count, err
	}
	return count, err
}
