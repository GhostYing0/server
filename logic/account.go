package logic

import (
	"errors"
	"github.com/unknwon/com"
	. "server/database"
	"server/logic/public"
	"server/models"
	. "server/utils/e"
	"server/utils/gredis"
	"server/utils/logging"
	. "server/utils/mydebug"
	"server/utils/util"
	"server/utils/uuid"
	"time"
)

type UserAccountLogic struct{}

var DefaultUserAccount = UserAccountLogic{}

func (self UserAccountLogic) Login(username string, password string, role int) (string, error) {
	if username == "" || password == "" {
		DPrintf("用户名和密码不能为空")
		return "", errors.New("用户名和密码不能为空")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("InsertEnrollInformation session.Begin() 发生错误:", err)
		return "", err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("InsertEnrollInformation session.Close() 发生错误:", err)
		}
	}()

	account := &models.Account{}
	has, err := session.Where("username = ?", username).And("role = ?", role).Get(account)
	if err != nil {
		DPrintf("Login 查询用户失败:", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return "", fail
		}
		return "", err
	}
	if !has {
		DPrintf("Login 用户不存在:")
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return "", fail
		}
		return "", errors.New("用户不存在")
	}

	if username != account.Username || util.EncodeMD5(password) != account.Password {
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return "", fail
		}
		return "", errors.New("用户名或密码错误")
	}

	token, err := util.GenerateToken(com.ToStr(account.ID), account.Username, account.Role)
	if err != nil {
		DPrintf("生成token失败")
		return "", err
	}

	err = gredis.Set(token, 1, time.Minute*60)
	if err != nil {
		DPrintf("login Set 失败:", err)
		return "", err
	}
	_, err = session.Where("id = ?", account.ID).Update(models.Account{LastLoginTime: models.NewOftenTime()})
	if err != nil {
		logging.L.Error(err)
		return "", err
	}

	return token, session.Commit()
}

func (self UserAccountLogic) Register(username string, password string, confirmPassword string, role int, name, gender, semester, college, school, class, phone, email string) error {
	if username == "" || password == "" {
		return errors.New("用户名或密码不能为空")
	}
	if password != confirmPassword {
		return errors.New("密码两次输入不一致")
	}
	if role != 1 && role != 2 {
		return errors.New("用户不是学生或老师")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Register session.Begin() 发生错误:", err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Register session.Close() 发生错误:", err)
		}
	}()

	has, err := session.Table("account").Where("username = ?", username).And("role = ?", role).Exist()
	if err != nil {
		DPrintf("Register 查询用户发生错误:", err)
		return err
	}
	if has {
		return errors.New("用户已存在")
	}

	searchSchoolResult := &models.School{}
	has, err = session.Where("school = ?", school).Get(searchSchoolResult)
	if err != nil {
		DPrintf("Register 查询学校发生错误:", err)
		return err
	}
	if !has {
		return errors.New("学校不存在")
	}

	searchCollegeResult := &models.College{}
	has, err = session.Where("college = ?", college).Get(searchCollegeResult)
	if err != nil {
		DPrintf("Register 查询学院发生错误:", err)
		return err
	}
	if !has {
		return errors.New("学院不存在")
	}

	newAccount := models.Account{
		Username:   username,
		Password:   util.EncodeMD5(password),
		Role:       role,
		Phone:      phone,
		Email:      email,
		CreateTime: models.NewOftenTime(),
		UserID:     uuid.CreateUUIDByNameSpace(username, password, name, gender, semester, college, school, class, role, time.Now()).String(),
	}

	_, err = session.Insert(newAccount)
	if err != nil {
		DPrintf("Login 添加新用户失败:", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}

	if role == 1 {
		//学生注册
		searchSemesterResult := &models.Semester{}
		has, err = session.Where("semester = ?", semester).Get(searchSemesterResult)
		if err != nil {
			DPrintf("Register 查询学期发生错误:", err)
			return err
		}
		if !has {
			return errors.New("学期不合法")
		}

		student := models.Student{
			StudentID:  newAccount.UserID,
			Name:       name,
			Gender:     gender,
			Class:      class,
			SchoolID:   searchSchoolResult.SchoolID,
			CollegeID:  searchCollegeResult.CollegeID,
			SemesterID: searchSemesterResult.SemesterID,
		}
		_, err = session.Insert(student)
		if err != nil {
			DPrintf("Register 添加学生失败:", err)
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return fail
			}
			return err
		}

	} else if role == 2 {
		//教师注册
		teacher := models.Teacher{
			TeacherID: newAccount.UserID,
			Name:      name,
			Gender:    gender,
			SchoolID:  searchSchoolResult.SchoolID,
			CollegeID: searchCollegeResult.CollegeID,
		}
		_, err = session.Insert(teacher)
		if err != nil {
			DPrintf("Register 添加教师失败:", err)
			fail := session.Rollback()
			if fail != nil {
				DPrintf("回滚失败")
				return fail
			}
			return err
		}
	}

	return session.Commit()
}

func (self UserAccountLogic) UpdatePassword(userID int64, role int, password string, confirmPassword string) error {
	if password == "" || confirmPassword == "" {
		return errors.New("用户名或密码不能为空")
	}
	if password != confirmPassword {
		return errors.New("密码两次输入不一致")
	}

	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Register session.Begin() 发生错误:", err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Register session.Close() 发生错误:", err)
		}
	}()

	has, err := session.Table("account").Where("id = ?", userID).And("role = ?", role).Exist()
	if err != nil {
		DPrintf("UpdatePassword 查找用户发生错误:", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}
	if !has {
		DPrintf("用户不存在")
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return errors.New("用户不存在")
	}

	newPassword := util.EncodeMD5(password)
	_, err = session.Table("account").Where("id = ? and role = ?", userID, role).Cols("password").Update(&models.Account{Password: newPassword})
	if err != nil {
		DPrintf("UpdatePassword 修改密码发生错误:", err)
		fail := session.Rollback()
		if fail != nil {
			DPrintf("回滚失败")
			return fail
		}
		return err
	}

	return session.Commit()
}

func (self UserAccountLogic) GetProfileStudent(userID int64) (*models.StudentReturn, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Register session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Register session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	data := &models.AccountStudent{}

	session.Join("LEFT", "account", "account.user_id = student.student_id")
	exist, err := session.Where("account.id = ? and role = 1", userID).Get(data)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}
	if !exist {
		logging.L.Error("用户不存在")
		return nil, errors.New("用户不存在")
	}

	school, err := public.SearchSchoolByID(data.SchoolID)
	if err != nil {
		logging.L.Error()
		school.School = "查询错误"
	}
	college, err := public.SearchCollegeByID(data.CollegeID)
	if err != nil {
		logging.L.Error()
		college.College = "查询错误"
	}
	semester, err := public.SearchSemesterByID(data.SemesterID)
	if err != nil {
		logging.L.Error()
		semester.Semester = "查询错误"
	}

	studentReturn := &models.StudentReturn{
		Name:     data.Name,
		Gender:   data.Gender,
		School:   school.School,
		Semester: semester.Semester,
		College:  college.College,
		Class:    data.Class,
		Phone:    data.Phone,
		Email:    data.Email,
		Avatar:   data.Avatar,
	}
	return studentReturn, err
}

func (self UserAccountLogic) GetProfileTeacher(userID int64) (*models.TeacherReturn, error) {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Register session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return nil, err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Register session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	data := &models.AccountTeacher{}

	session.Join("LEFT", "account", "account.user_id = teacher.teacher_id")
	exist, err := session.Where("account.id = ? and role = 2", userID).Get(data)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}
	if !exist {
		logging.L.Error("用户不存在")
		return nil, errors.New("用户不存在")
	}

	school, err := public.SearchSchoolByID(data.SchoolID)
	if err != nil {
		logging.L.Error()
		school.School = "查询错误"
	}
	college, err := public.SearchCollegeByID(data.CollegeID)
	if err != nil {
		logging.L.Error()
		college.College = "查询错误"
	}

	teacherReturn := &models.TeacherReturn{
		Name:    data.Name,
		Gender:  data.Gender,
		School:  school.School,
		College: college.College,
		Phone:   data.Phone,
		Email:   data.Email,
		Avatar:  data.Avatar,
	}
	return teacherReturn, err
}

func (self UserAccountLogic) UpdateProfile(userID int64, role int, phone, email string) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Register session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Register session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	data := &models.Account{}

	exist, err := session.Table("account").Where("id = ? and role = ?", userID, role).Get(data)
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("用户不存在")
		return errors.New("用户不存在")
	}

	_, err = session.Table("account").Where("id = ? and role = ?", userID, role).Update(&models.Account{Phone: phone, Email: email})
	if err != nil {
		fail := session.Rollback()
		if fail != nil {
			logging.L.Error(err)
			return fail
		}
		logging.L.Error(err)
		return err
	}
	return session.Commit()
}

func (self UserAccountLogic) UpdateAvatar(userID int64, role int, avatar string) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("Register session.Begin() 发生错误:", err)
		logging.L.Error(err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("Register session.Close() 发生错误:", err)
			logging.L.Error(err)
		}
	}()

	data := &models.Account{}

	exist, err := session.Table("account").Where("id = ? and role = ?", userID, role).Get(data)
	if err != nil {
		logging.L.Error(err)
		return err
	}
	if !exist {
		logging.L.Error("用户不存在")
		return errors.New("用户不存在")
	}

	if role == StudentRole {
		student, err := public.SearchStudentByID(data.UserID)
		if err != nil {
			logging.L.Error(err)
			return err
		}
		_, err = session.Table("student").Where("student_id = ?", student.StudentID).Update(&models.Student{Avatar: avatar})
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				logging.L.Error(err)
				return fail
			}
			logging.L.Error(err)
			return err
		}
	} else if role == TeacherRole {
		teacher, err := public.SearchTeacherByID(data.UserID)
		if err != nil {
			logging.L.Error(err)
			return err
		}
		_, err = session.Table("teacher").Where("teacher_id = ?", teacher.TeacherID).Update(&models.Teacher{Avatar: avatar})
		if err != nil {
			fail := session.Rollback()
			if fail != nil {
				logging.L.Error(err)
				return fail
			}
			logging.L.Error(err)
			return err
		}
	}
	return session.Commit()
}
