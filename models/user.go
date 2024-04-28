package models

type UserParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type OldUser struct {
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

type User struct {
}

type AccountStudent struct {
	Account `xorm:"extends"`
	Student `xorm:"extends"`
}

func (*AccountStudent) TableName() string {
	return "student"
}

func (*AccountTeacher) TableName() string {
	return "teacher"
}

type AccountTeacher struct {
	Account `xorm:"extends"`
	Teacher `xorm:"extends"`
}

type UpdateUserInfo OldUser

type UserDeleteId struct {
	ID []int64 `json:"ids"`
}

func (OldUser) TableName() string {
	return "account"
}

type Student struct {
	StudentID    string    `json:"student_id" xorm:"student_id"`
	Name         string    `json:"name" xorm:"name"`
	Gender       string    `json:"gender" xorm:"gender"`
	SchoolID     int64     `json:"school_id" xorm:"school_id"`
	SemesterID   int64     `json:"semester_id" xorm:"semester_id"`
	CollegeID    int64     `json:"college_id" xorm:"college_id"`
	DepartmentID int64     `json:"department_id" xorm:"department_id"`
	Class        string    `json:"class" xorm:"class"`
	Avatar       string    `json:"avatar" xorm:"avatar"`
	Deleted      OftenTime `xorm:"deleted"`
}

type StudentForm struct {
	ID       int64  `json:"id" xorm:"id"`
	Username string `json:"username" xorm:"username"`
	Password string `json:"password" xorm:"password"`
	Name     string `json:"name" xorm:"name"`
	Gender   string `json:"gender" xorm:"gender"`
	School   string `json:"school" xorm:"school"`
	Semester string `json:"semester" xorm:"semester"`
	College  string `json:"college"" xorm:"college"`
	Class    string `json:"class" xorm:"class"`
	Avatar   string `json:"avatar" xorm:"avatar"`
}

type StudentReturn struct {
	ID        int64  `json:"id" xorm:"id"`
	Username  string `json:"username" xorm:"username"`
	Password  string `json:"password" xorm:"password"`
	Role      int    `json:"role" xorm:"role"`
	StudentID string `json:"student_id" xorm:"student_id"`
	Name      string `json:"name" xorm:"name"`
	Gender    string `json:"gender" xorm:"gender"`
	School    string `json:"school" xorm:"school"`
	Semester  string `json:"semester" xorm:"semester"`
	College   string `json:"college" xorm:"college"`
	Class     string `json:"class" xorm:"class"`
	Phone     string `json:"phone" xorm:"phone"`
	Email     string `json:"email" xorm:"email"`
	Avatar    string `json:"avatar" xorm:"avatar"`
}

func (Student) TableName() string {
	return "student"
}

type Teacher struct {
	TeacherID    string    `json:"teacher_id" xorm:"teacher_id"`
	Name         string    `json:"name" xorm:"name"`
	Gender       string    `json:"gender" xorm:"gender"`
	SchoolID     int64     `json:"school_id" xorm:"school_id"`
	CollegeID    int64     `json:"college_id" xorm:"college_id"`
	DepartmentID int64     `json:"department_id" xorm:"department_id"`
	Title        string    `json:"title" xorm:"title"`
	Avatar       string    `json:"avatar" xorm:"avatar"`
	Deleted      OftenTime `xorm:"deleted"`
}

type TeacherForm struct {
	ID       int64  `json:"id" xorm:"id"`
	Username string `json:"username" xorm:"username"`
	Password string `json:"password" xorm:"password"`
	Name     string `json:"name" xorm:"name"`
	Gender   string `json:"gender" xorm:"gender"`
	School   string `json:"school" xorm:"school"`
	College  string `json:"college" xorm:"college"`
	Avatar   string `json:"avatar" xorm:"avatar"`
}

type TeacherReturn struct {
	ID        int64  `json:"id" xorm:"id"`
	Username  string `json:"username" xorm:"username"`
	Password  string `json:"password" xorm:"password"`
	Role      int    `json:"role" xorm:"role"`
	TeacherID string `json:"teacher_id" xorm:"teacher_id"`
	Name      string `json:"name" xorm:"name"`
	Gender    string `json:"gender" xorm:"gender"`
	School    string `json:"school" xorm:"school"`
	College   string `json:"college" xorm:"college"`
	Phone     string `json:"phone" xorm:"phone"`
	Email     string `json:"email" xorm:"email"`
	Avatar    string `json:"avatar" xorm:"avatar"`
}

func (Teacher) TableName() string {
	return "teacher"
}

type Semester struct {
	SemesterID int64  `xorm:"semester_id"`
	Semester   string `xorm:"semester"`
}

func (Semester) TableName() string {
	return "semester"
}

type School struct {
	SchoolID int64  `xorm:"school_id"`
	School   string `xorm:"school"`
}

func (School) TableName() string {
	return "school"
}

type College struct {
	CollegeID int64  `xorm:"college_id"`
	College   string `xorm:"college"`
}

func (College) TableName() string {
	return "college"
}
