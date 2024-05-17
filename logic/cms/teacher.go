package cms

import (
	"errors"
	"fmt"
	. "server/database"
	. "server/logic"
	"server/logic/public"
	"server/models"
	"server/utils/e"
	"server/utils/logging"
	. "server/utils/mydebug"
	"server/utils/util"
	"server/utils/uuid"
	"time"
)

type CmsTeacherLogic struct{}

var DefaultCmsTeacher = CmsTeacherLogic{}

func (self CmsTeacherLogic) DisplayTeacher(paginator *Paginator, username, gender, school, college, name, department, title string) (*[]models.TeacherReturn, int64, error) {
	if paginator == nil {
		DPrintf("Search 分页器为空")
		logging.L.Error("Search 分页器为空")
		return nil, 0, errors.New("分页器为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser Display session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser Display session.Close() 发生错误:", err)
		}
	}()

	//查询教师
	session.Table("account").Where("role = ?", e.TeacherRole)

	if username != "" {
		session.Table("account").Where("username like ?", "%"+username+"%")
	}
	session.Join("LEFT", "teacher", "account.user_id = teacher.teacher_id")
	session.Join("LEFT", "department", "department.department_id = teacher.department_id")
	session.Join("LEFT", "college", "college.college_id = teacher.college_id")
	if gender != "" {
		session.Where("teacher.gender = ?", gender)
	}
	if school != "" {
		searchSchool, err := public.SearchSchoolByName(school)
		if err != nil {
			logging.L.Error(err)
		} else {
			session.Where("teacher.school_id = ?", searchSchool.SchoolID)
		}
	}
	if college != "" {
		session.Where("college.college like ?", "%"+college+"%")
	}
	if title != "" {
		session.Where("teacher.title like ?", "%"+title+"%")
	}
	if department != "" {
		session.Where("department.department like ?", "%"+department+"%")
	}
	if name != "" {
		session.Where("teacher.name like ?", "%"+name+"%")
	}

	data := &[]models.AccountTeacher{}

	total, err := session.
		//Join("LEFT", "teacher", "account.user_id = teacher.teacher_id").
		Limit(paginator.PerPage(), paginator.Offset()).
		FindAndCount(data)
	if err != nil {
		DPrintf("Search 查找成绩信息失败:", err)
		logging.L.Error(err)
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
		list[i].Title = (*data)[i].Title
		list[i].Department = (*data)[i].Department
		list[i].Gender = (*data)[i].Gender
		list[i].College = (*data)[i].College
	}

	return &list, total, session.Rollback()
}

func (self CmsTeacherLogic) AddTeacher(username, password, name, gender, school, college, avatar, department, title string) error {
	if len(username) == 0 || len(password) == 0 {
		logging.L.Error("账号和密码不能为空")
		return errors.New("账号和密码不能为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser Display session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser Display session.Close() 发生错误:", err)
		}
	}()

	exist, err := session.Table("account").Where("username = ? and role = 2", username).Exist()
	if exist {
		logging.L.Error("用户已存在")
		return errors.New("用户已存在")
	}

	//searchSchool, err := public.SearchSchoolByName(school)
	//if err != nil {
	//	logging.L.Error(err)
	//	return err
	//}
	searchCollege, err := public.SearchCollegeByName(college)
	if err != nil {
		logging.L.Error(err)
		return err
	}
	searchDepartment, err := public.SearchDepartmentByName(department)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	TeacherID := uuid.CreateUUIDByNameSpace(username, password, name, gender, "", college, school, "", 1, time.Now()).String()
	account := &models.Account{
		Username: username,
		Password: util.EncodeMD5(password),
		Role:     2,
		UserID:   TeacherID,
	}

	teacher := &models.Teacher{
		TeacherID:    TeacherID,
		Name:         name,
		Gender:       gender,
		DepartmentID: searchDepartment.DepartmentID,
		//SchoolID:  searchSchool.SchoolID,
		CollegeID: searchCollege.CollegeID,
		Title:     title,
		Avatar:    avatar,
	}

	_, err = session.Insert(account)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(err)
			return fail
		}
		logging.L.Error(err)
		return err
	}

	_, err = session.Insert(teacher)
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(err)
			return fail
		}
		logging.L.Error(err)
		return err
	}

	return session.Commit()
}

func (self CmsTeacherLogic) UpdateTeacher(id int64, username, password, name, gender, school, college, avatar, department, title string) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser UpdateUser session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser UpdateUser session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	searchAccount := &models.Account{}
	exist, err := session.Table("account").Where("id = ?", id).Get(searchAccount)
	if !exist {
		logging.L.Error("用户不存在")
		return errors.New("用户不存在")
	}
	if err != nil {
		logging.L.Error(err)
		DPrintf("UpdateUser 查询用户失败:", err)
		return err
	}

	exist, err = session.Table("account").Where("username = ? and role = ? and id != ?", username, e.TeacherRole, id).Exist()
	if exist {
		logging.L.Error("已有同名用户")
		return errors.New("已有同名用户")
	}
	if err != nil {
		logging.L.Error(err)
		return err
	}

	//searchSchool, err := public.SearchSchoolByName(school)
	//if err != nil {
	//	logging.L.Error(err)
	//	return err
	//}
	searchCollege, err := public.SearchCollegeByName(college)
	if err != nil {
		logging.L.Error(err)
		return err
	}
	searchDepartment, err := public.SearchDepartmentByName(department)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	_, err = session.Where("id = ?", id).
		Update(&models.Account{
			Username: username,
			Password: util.EncodeMD5(password),
		})
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(fail)
			return fail
		}
		logging.L.Error(err)
		DPrintf("UpdateUser Update 更新用户失败:", err)
		return err
	}

	_, err = session.Where("teacher_id = ?", searchAccount.UserID).
		Update(&models.Teacher{
			Name:   name,
			Gender: gender,
			//SchoolID:  searchSchool.SchoolID,
			DepartmentID: searchDepartment.DepartmentID,
			Title:        title,
			CollegeID:    searchCollege.CollegeID,
			Avatar:       avatar,
		})
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			logging.L.Error(fail)
			return fail
		}
		logging.L.Error(err)
		DPrintf("UpdateUser Update 更新用户失败:", err)
		return err
	}

	return session.Commit()
}

func (self CmsTeacherLogic) DeleteTeacher(ids *[]int64) (int64, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("cmsUser DeleteUser session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return 0, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("cmsUser DeleteUser session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	var count int64

	for _, id := range *ids {
		account := &models.Account{}
		if id < 1 {
			fmt.Println("非法id")
			continue
		}

		_, err := session.Where("id = ?", id).Get(account)
		if err != nil {
			logging.L.Error(err)
			return 0, err
		}

		_, err = session.Delete(account)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				logging.L.Error(fail)
				return 0, fail
			}
			logging.L.Error(err)
			return 0, err
		}

		teacher := &models.Teacher{TeacherID: account.UserID}

		affected, err := session.Delete(teacher)
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				logging.L.Error(fail)
				return 0, fail
			}
			logging.L.Error(err)
			return 0, err
		}
		if affected > 0 {
			count += affected
		}
	}

	return count, session.Commit()
}

func (self CmsTeacherLogic) GetTeacherCount() (int64, error) {
	session := MasterDB.NewSession()

	count, err := session.Table("teacher").Count()
	if err != nil {
		DPrintf("GetUserCount Count 发生错误:", err)
		logging.L.Error(err)
		return count, err
	}
	return count, err
}
