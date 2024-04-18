package cms

import (
	"errors"
	"fmt"
	. "server/database"
	. "server/logic"
	"server/logic/public"
	"server/models"
	"server/utils/logging"
	. "server/utils/mydebug"
	"server/utils/util"
	"server/utils/uuid"
	"time"
)

type CmsStudentLogic struct{}

var DefaultCmsStudent = CmsStudentLogic{}

func (self CmsStudentLogic) DisplayStudent(paginator *Paginator, username, gender, school, semester, college, class, name string) (*[]models.StudentReturn, int64, error) {
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

	//查询学生
	session.Table("account").Where("role = ?", 1)

	if username != "" {
		session.Table("account").Where("username = ?", username)
	}
	session.Join("LEFT", "student", "account.user_id = student.student_id")
	if gender != "" {
		session.Where("student.gender = ?", gender)
	}
	if school != "" {
		searchSchool, err := public.SearchSchoolByName(school)
		if err != nil {
			logging.L.Error(err)
		} else {
			session.Where("student.school_id = ?", searchSchool.SchoolID)
		}
	}
	if semester != "" {
		searchSemester, err := public.SearchSemesterByName(semester)
		if err != nil {
			logging.L.Error(err)
		} else {
			session.Where("student.semester_id = ?", searchSemester.SemesterID)
		}
	}
	if college != "" {
		searchCollege, err := public.SearchCollegeByName(college)
		if err != nil {
			logging.L.Error(err)
		} else {
			session.Where("student.college_id = ?", searchCollege.CollegeID)
		}
	}
	if name != "" {
		session.Where("student.name = ?", name)
	}
	if class != "" {
		session.Where("student.class = ?", class)
	}

	data := &[]models.AccountStudent{}

	total, err := session.
		//Join("LEFT", "student", "account.user_id = student.student_id").
		Limit(paginator.PerPage(), paginator.Offset()).
		FindAndCount(data)
	if err != nil {
		DPrintf("Search 查找成绩信息失败:", err)
		logging.L.Error(err)
		return nil, 0, err
	}

	list := make([]models.StudentReturn, len(*data))

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
			logging.L.Error(err)
		}
		list[i].School = searchSchool.School
		_, err = session.Where("semester_id = ?", (*data)[i].SemesterID).Get(searchSemester)
		if err != nil {
			list[i].Semester = "查询出错"
			logging.L.Error(err)
		}
		list[i].Semester = searchSemester.Semester
		_, err = session.Where("college_id = ?", (*data)[i].CollegeID).Get(searchCollege)
		if err != nil {
			list[i].College = "查询出错"
			logging.L.Error(err)
		}
		list[i].College = searchCollege.College
	}

	return &list, total, session.Rollback()
}

func (self CmsStudentLogic) AddStudent(username, password, name, gender, school, college, class, semester, avatar string) error {
	if len(username) == 0 || len(password) == 0 {
		logging.L.Error("账号和密码不能为空")
		return nil
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

	exist, err := session.Table("account").Where("username = ? and role = 1", username).Exist()
	if exist {
		logging.L.Error("用户已存在")
		return errors.New("用户已存在")
	}

	searchSchool, err := public.SearchSchoolByName(school)
	if err != nil {
		logging.L.Error(err)
		return err
	}
	searchSemester, err := public.SearchSemesterByName(semester)
	if err != nil {
		logging.L.Error(err)
		return err
	}
	searchCollege, err := public.SearchCollegeByName(college)
	if err != nil {
		logging.L.Error(err)
		return err
	}

	studentID := uuid.CreateUUIDByNameSpace(username, password, name, gender, semester, college, school, class, 1, time.Now()).String()
	account := &models.Account{
		Username: username,
		Password: util.EncodeMD5(password),
		Role:     1,
		UserID:   studentID,
	}

	student := &models.Student{
		StudentID:  studentID,
		Name:       name,
		Gender:     gender,
		SchoolID:   searchSchool.SchoolID,
		SemesterID: searchSemester.SemesterID,
		CollegeID:  searchCollege.CollegeID,
		Class:      class,
		Avatar:     avatar,
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

	_, err = session.Insert(student)
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

func (self CmsStudentLogic) UpdateStudent(id int64, username, password, name, gender, school, college, class, semester, avatar string) error {
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

	searchSchool, err := public.SearchSchoolByName(school)
	if err != nil {
		logging.L.Error(err)
		return err
	}
	searchSemester, err := public.SearchSemesterByName(semester)
	if err != nil {
		logging.L.Error(err)
		return err
	}
	searchCollege, err := public.SearchCollegeByName(college)
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

	_, err = session.Where("student_id = ?", searchAccount.UserID).
		Update(&models.Student{
			Name:       name,
			Gender:     gender,
			SchoolID:   searchSchool.SchoolID,
			CollegeID:  searchCollege.CollegeID,
			Class:      class,
			SemesterID: searchSemester.SemesterID,
			Avatar:     avatar,
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

func (self CmsStudentLogic) DeleteStudent(ids *[]int64) (int64, error) {
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

		student := &models.Student{StudentID: account.UserID}

		affected, err := session.Delete(student)
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

func (self CmsStudentLogic) GetStudentCount() (int64, error) {
	session := MasterDB.NewSession()

	count, err := session.Table("student").Count()
	if err != nil {
		DPrintf("GetUserCount Count 发生错误:", err)
		logging.L.Error(err)
		return count, err
	}
	return count, err
}
