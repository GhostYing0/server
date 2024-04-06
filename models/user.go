package models

type UserParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int64     `xorm:"id"`
	Username string    `xorm:"username"`
	Password string    `xorm:"password"`
	Role     int       `xorm:"role"`
	Deleted  OftenTime `xorm:"deleted"`
}

type UserRedis struct {
	ID       int64  `json:"id"`
	Role     int    `json:"role"`
	Username string `json:"username"`
}

type UpdateUserInfo User

type UserDeleteId struct {
	ID []int64 `json:"ids"`
}

func (User) TableName() string {
	return "account"
}

type Student struct {
	StudentID  string    `json:"student_id" xorm:"student_id"`
	Name       string    `json:"name" xorm:"name"`
	Gender     string    `json:"gender" xorm:"gender"`
	SchoolID   int64     `json:"school_id" xorm:"school_id"`
	SemesterID int64     `json:"semester_id" xorm:"semester_id"`
	CollegeID  int64     `json:"college_id" xorm:"college_id"`
	Class      string    `json:"class" xorm:"class"`
	Avatar     string    `json:"avatar" xorm:"avatar"`
	Deleted    OftenTime `xorm:"deleted"`
}

func (Student) TableName() string {
	return "student"
}

type Teacher struct {
	TeacherID string    `json:"teacher_id" xorm:"teacher_id"`
	Name      string    `json:"name" xorm:"name"`
	Gender    string    `json:"gender" xorm:"gender"`
	SchoolID  int64     `json:"school_id" xorm:"school_id"`
	CollegeID int64     `json:"college_id" xorm:"college_id"`
	Deleted   OftenTime `xorm:"deleted"`
}

func (Teacher) TableName() string {
	return "teacher"
}

type Semester struct {
	SemesterID int64  `xorm:"semester_id"`
	Semester   string `xorm:"semester"`
}

type School struct {
	SchoolID int64  `xorm:"school_id"`
	School   string `xorm:"school"`
}

type College struct {
	CollegeID int64  `xorm:"college_id"`
	College   string `xorm:"college"`
}
