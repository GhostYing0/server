package public

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	. "server/database"
	"server/models"
	. "server/utils/e"
	"server/utils/gredis"
	"server/utils/logging"
	. "server/utils/mydebug"
	"server/utils/util"
	"strconv"
	"time"
)

type PublicLogic struct{}

var DefaultPublic = PublicLogic{}

func (self PublicLogic) GetInfo(token string) (int64, string, int, error) {
	tokenHasExpired, err := gredis.Get(token)
	if err != nil {
		fmt.Println("")
		return 0, "", -1, err
	}
	if tokenHasExpired != "1" {
		return 0, "", -1, errors.New("token已过期, 请重新登录")
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		fmt.Println("")
		return 0, "", -1, err
	}

	id, err := strconv.Atoi(claims.ID)
	if err != nil {
		fmt.Println("")
		return 0, "", -1, err
	}

	return int64(id), claims.Username, claims.Role, err
}

func (self PublicLogic) Logout(token string) error {
	tokenHasExpired, err := gredis.Get(token)
	if err != nil {
		fmt.Println("")
		return err
	}
	if tokenHasExpired != "1" {
		return errors.New("已登出")
	}

	err = gredis.Del(token)
	if err != nil {
		DPrintf("Logout Del token 失败:", err)
	}
	return err
}

func (self PublicLogic) GetContestType() (*[]models.ContestType, error) {
	list := &[]models.ContestType{}
	err := MasterDB.Find(list)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}
	return list, err
}

func (self PublicLogic) GetSchool() (*[]models.School, error) {
	list := &[]models.School{}
	err := MasterDB.Find(list)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}
	return list, err
}

func (self PublicLogic) GetCollege() (*[]models.College, error) {
	list := &[]models.College{}
	err := MasterDB.Find(list)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}
	return list, err
}

func (self PublicLogic) GetSemester() (*[]models.Semester, error) {
	list := &[]models.Semester{}
	err := MasterDB.Find(list)
	if err != nil {
		logging.L.Error(err)
		return nil, err
	}
	return list, err
}

func (self PublicLogic) UploadImg(file *multipart.FileHeader) (string, error) {
	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		// 返回值
		DPrintf("文件类型不合法")
		return "", errors.New("文件类型不合法")
	}

	//currentTime := time.Now().Format("20060102")
	// 生成目录文件夹，并错误判断
	//if err := os.MkdirAll("D:/GDesign/picture/img"+currentTime, 0755); err != nil {
	//	DPrintf("上传错误")
	//	appG.ResponseErr("MkdirAll失败")
	//	return
	//}
	if err := os.MkdirAll("D:/GDesign/picture/img", 0755); err != nil {
		DPrintf("上传错误")
		return "", errors.New("MkdirAll失败")
	}

	fileUnixName := strconv.FormatInt(time.Now().UnixNano(), 10)

	saveDir := path.Join("D:/GDesign/picture/img", fileUnixName+extName)

	return saveDir, nil
}

func SearchSchoolByName(name string) (*models.School, error) {
	school := &models.School{}
	exist, err := MasterDB.Where("school = ?", name).Get(school)
	if err != nil {
		DPrintf("查询学校失败")
		logging.L.Error(err)
		return school, err
	}
	if !exist {
		return school, errors.New("学校不存在")
	}
	return school, err
}

func SearchSemesterByName(name string) (*models.Semester, error) {
	semester := &models.Semester{}
	exist, err := MasterDB.Where("semester = ?", name).Get(semester)
	if err != nil {
		DPrintf("查询学期失败")
		logging.L.Error(err)
		return semester, err
	}
	if !exist {
		return semester, errors.New("学年不存在")
	}
	return semester, err
}

func SearchCollegeByName(name string) (*models.College, error) {
	college := &models.College{}
	exist, err := MasterDB.Where("college = ?", name).Get(college)
	if err != nil {
		DPrintf("查询学院失败")
		logging.L.Error(err)
		return college, err
	}
	if !exist {
		return college, errors.New("学院不存在")
	}
	return college, err
}

func SearchStudentByName(name string) (*models.Student, error) {
	student := &models.Student{}
	exist, err := MasterDB.Where("name = ?", name).Get(student)
	if err != nil {
		logging.L.Error(err)
		return student, err
	}
	if !exist {
		return student, errors.New("学生不存在")
	}
	return student, err
}

func SearchTeacherByName(name string) (*models.Teacher, error) {
	teacher := &models.Teacher{}
	exist, err := MasterDB.Where("name = ?", name).Get(teacher)
	if err != nil {
		logging.L.Error(err)
		return teacher, err
	}
	if !exist {
		return teacher, errors.New("教师不存在")
	}
	return teacher, err
}

func SearchTeacherByID(id string) (*models.Teacher, error) {
	teacher := &models.Teacher{}
	exist, err := MasterDB.Where("teacher_id = ?", id).Get(teacher)
	if err != nil {
		logging.L.Error(err)
		return teacher, err
	}
	if !exist {
		return teacher, errors.New("教师不存在")
	}
	return teacher, err
}

func SearchContestByName(name string) (*models.ContestInfo, error) {
	contest := &models.ContestInfo{}
	exist, err := MasterDB.Where("contest = ?", name).Get(contest)
	if err != nil {
		logging.L.Error(err)
		return contest, err
	}
	if !exist {
		return contest, errors.New("竞赛不存在")
	}
	return contest, err
}

func SearchAccountByUsernameAndRole(username string, role int) (*models.Account, error) {
	account := &models.Account{}
	exist, err := MasterDB.Where("username = ? and role = ?", username, role).Get(account)
	if err != nil {
		logging.L.Error(err)
		return account, err
	}
	if !exist {
		return account, errors.New("账户不存在")
	}
	return account, err
}

func SearchSchoolByID(id int64) (*models.School, error) {
	school := &models.School{}
	exist, err := MasterDB.Where("school_id = ?", id).Get(school)
	if err != nil {
		logging.L.Error(err)
		return school, err
	}
	if !exist {
		return school, errors.New("学校不存在")
	}
	return school, err
}

func SearchAccountByID(id int64) (*models.Account, error) {
	account := &models.Account{}
	exist, err := MasterDB.Where("id = ?", id).Get(account)
	if err != nil {
		logging.L.Error(err)
		return account, err
	}
	if !exist {
		return account, errors.New("账户不存在")
	}
	return account, err
}

func SearchStudentByID(id string) (*models.Student, error) {
	student := &models.Student{}
	exist, err := MasterDB.Where("student_id = ?", id).Get(student)
	if err != nil {
		logging.L.Error(err)
		return student, err
	}
	if !exist {
		return student, errors.New("学生不存在")
	}
	return student, err
}

func SearchCollegeByID(id int64) (*models.College, error) {
	college := &models.College{}
	exist, err := MasterDB.Where("college_id = ?", id).Get(college)
	if err != nil {
		logging.L.Error(err)
		return college, err
	}
	if !exist {
		return college, errors.New("学院不存在")
	}
	return college, err
}

func SearchSemesterByID(id int64) (*models.Semester, error) {
	semester := &models.Semester{}
	exist, err := MasterDB.Where("semester_id = ?", id).Get(semester)
	if err != nil {
		logging.L.Error(err)
		return semester, err
	}
	if !exist {
		return semester, errors.New("学期不存在")
	}
	return semester, err
}

func SearchContestTypeByName(name string) (*models.ContestType, error) {
	contestType := &models.ContestType{}
	exist, err := MasterDB.Where("type = ?", name).Get(contestType)
	if err != nil {
		logging.L.Error(err)
		return contestType, err
	}
	if !exist {
		return contestType, errors.New("竞赛类型不存在")
	}
	return contestType, err
}

func (self PublicLogic) GetContest() (*[]models.ContestAndType, error) {
	contest := &[]models.ContestAndType{}

	_, err := MasterDB.
		Table("contest").
		Join("LEFT", "contest_type", "contest.contest_type_id = contest_type.id").
		Where("contest.state = ?", Pass).
		FindAndCount(contest)
	if err != nil {
		return nil, err
	}

	return contest, err
}
