package models

type EnrollForm struct {
	ID                int64  `json:"id" xorm:""`
	ContestID         int64  `json:"contest_id" xorm:"contest_id"`
	StudentName       string `json:"student_name" xorm:""`
	TeamName          string `json:"team_name" xorm:""`
	Handle            int64  `json:"handle_team" xorm:""`
	College           string `json:"college" xorm:""`
	Teacher           string `json:"teacher_name"`
	Major             string `json:"major" xorm:""`
	GuidanceTeacher   string `json:"guidance_teacher" xorm:""`
	Class             string `json:"student_class" xorm:""`
	TeacherDepartment string `json:"teacher_department" xorm:"`
	Department        string `json:"department" xorm:""`
	TeacherTitle      string `json:"teacher_title" xorm:""`
	Phone             string `json:"phone" xorm:""`
	Email             string `json:"email" xorm:""`
	State             int    `json:"state" xorm:""`
}

func (EnrollForm) TableName() string {
	return "enroll_information"
}

type EnrollInformationForm struct {
	ID         int64  `json:"id"`
	UserName   string `json:"username"`
	Name       string `json:"name"`
	TeamID     int64  `json:"team_id"`
	ContestID  int64  `json:"contest_id"`
	StudentID  string `json:"student_id"`
	TeacherID  string `json:"teacher_id"`
	Contest    string `json:"contest"`
	CreateTime string `json:"create_time"`
	School     string `json:"school"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	State      int    `json:"state"`
}

type EnrollInformation struct {
	ID        int64  `json:"id" xorm:"id"`
	StudentID string `json:"student_id" xorm:"student_id"`
	TeamID    int64  `json:"team_id" xorm:"team_id"`
	//ContestID    string    `json:"contest_id" xorm:"contest_id"`
	ContestID       int64     `json:"contest_id" xorm:"contest_id"`
	CreateTime      string    `json:"create_time" xorm:"create_time"`
	SchoolID        int64     `json:"school_id" xorm:"school_id"`
	Phone           string    `json:"phone" xorm:"phone"`
	Email           string    `json:"email" xorm:"email"`
	State           int       `json:"state" xorm:"state"`
	RejectReason    string    `json:"reject_reason" xorm:"reject_reason"`
	GuidanceTeacher string    `json:"guidance_teacher" xorm:"guidance_teacher"`
	Deleted         OftenTime `json:"deleted" xorm:"deleted"`
}

type NewEnroll struct {
	StudentID       string    `json:"student_id" xorm:"student_id"`
	TeamID          int64     `json:"team_id" xorm:"team_id"`
	ContestID       int64     `json:"contest_id" xorm:"contest_id"`
	CreateTime      OftenTime `json:"create_time" xorm:"create_time"`
	SchoolID        int64     `json:"school_id" xorm:"school_id"`
	Phone           string    `json:"phone" xorm:"phone"`
	Email           string    `json:"email" xorm:"email"`
	State           int       `json:"state" xorm:"state"`
	GuidanceTeacher string    `json:"guidance_teacher" xorm:"guidance_teacher"`
}

func (NewEnroll) TableName() string {
	return "enroll_information"
}

type EnrollContestStudent struct {
	EnrollInformation `xorm:"extends"`
	Username          string    `json:"username" xorm:"username"`
	Contest           string    `json:"contest" xorm:"contest"`
	ContestState      int       `json:"contest_state" xorm:"contest_state"`
	ContestType       string    `json:"type" xorm:"type"`
	ContestLevel      string    `json:"contest_level" xorm:"contest_level"`
	StartTime         string    `json:"start_time" xorm:"start_time"`
	Deadline          string    `json:"deadline" xorm:"deadline"`
	State             int       `json:"state" xorm:"state"`
	Describe          string    `json:"desc" xorm:"describe"`
	Deleted           OftenTime `json:"deleted" xorm:"deleted"`
	Student           `xorm:"extends"`
	School            string `xorm:"school"`
	College           string `xorm:"college"`
	Team              string `json:"team_name" xorm:"team_name"`
	Semester          string `xorm:"semester"`
	IsGroup           int    `json:"is_group" xorm:"is_group"`
	TeacherName       string `json:"teacher_name" xorm:"t_name"`
	Department        string `json:"department" xorm:"t_department"`
	Title             string `json:"title" xorm:"title"`
	Major             string `json:"major" xorm:"major"`
	ContestEntry      string `json:"contest_entry" xorm:"contest_entry"`
}

type TeacherUploadGetEnroll struct {
	EnrollInformation `xorm:"extends"`
	Name              string    `json:"name" xorm:"name"`
	Contest           string    `json:"contest" xorm:"contest"`
	ContestState      int       `json:"contest_state" xorm:"contest_state"`
	ContestType       string    `json:"type" xorm:"type"`
	ContestLevel      string    `json:"contest_level" xorm:"contest_level"`
	CreateTime        string    `json:"create_time" xorm:"create_time"`
	StartTime         string    `json:"start_time" xorm:"start_time"`
	StudentSchoolID   string    `json:"student_school_id" xorm:"student_school_id"`
	Deadline          string    `json:"deadline" xorm:"deadline"`
	State             int       `json:"state" xorm:"state"`
	Describe          string    `json:"desc" xorm:"describe"`
	Deleted           OftenTime `json:"deleted" xorm:"deleted"`
	Major             string    `json:"major" xorm:"major"`
	Student           `xorm:"extends"`
	School            string `xorm:"school"`
	College           string `xorm:"college"`
	Semester          string `xorm:"semester"`
}

type EnrollContestStudent_e_id struct {
	ID           int64  `json:"id" xorm:"e_id"`
	StudentID    string `json:"student_id" xorm:"student_id"`
	TeamID       int64  `json:"team_id" xorm:"team_id"`
	Team         string `json:"team" xorm:"team_name"`
	ContestID    int64  `json:"contest_id" xorm:"contest_id"`
	Teacher      string `json:"teacher_name" xorm:"teacher_name"`
	Title        string `json:"title"`
	Department   string `json:"department"`
	ContestLevel string `json:"contest_level"`
	IsGroup      int    `json:"is_group" xorm:"is_group"`

	SchoolID int64  `json:"school_id" xorm:"school_id"`
	Phone    string `json:"phone" xorm:"phone"`
	Email    string `json:"email" xorm:"email"`

	RejectReason string `json:"reject_reason" xorm:"reject_reason"`

	Username     string    `json:"username" xorm:"username"`
	Contest      string    `json:"contest" xorm:"contest"`
	ContestState int       `json:"contest_state" xorm:"contest_state"`
	ContestType  string    `json:"type" xorm:"type"`
	CreateTime   string    `json:"create_time" xorm:"create_time"`
	StartTime    string    `json:"start_time" xorm:"start_time"`
	Deadline     string    `json:"deadline" xorm:"deadline"`
	State        int       `json:"state" xorm:"state"`
	Describe     string    `json:"desc" xorm:"describe"`
	Deleted      OftenTime `json:"deleted" xorm:"deleted"`
	Major        string    `json:"major"`
	Student      `xorm:"extends"`
	School       string `xorm:"school"`
	College      string `xorm:"college"`
	Semester     string `xorm:"semester"`
}

func (EnrollContestStudent_e_id) TableName() string {
	return "enroll_information"
}

type EnrollContest struct {
	EnrollInformation `xorm:"extends"`
	Contest           `xorm:"extends"`
	Account           `xorm:"extends"`
}

func (EnrollContest) TableName() string {
	return "enroll_information"
}

func (EnrollContestStudent) TableName() string {
	return "enroll_information"
}

type ContestInfoAccount struct {
	Account     `xorm:"extends"`
	ContestInfo `xorm:"extends"`
}

type EnrollInformationReturn struct {
	ID              int64  `json:"id" xorm:"id"`
	Username        string `json:"username" xorm:"username"`
	StudentID       string `json:"student_id" xorm:"student_id"`
	TeamID          int64  `json:"team_id" xorm:"team_id"`
	ContestType     string `json:"contest_type" xorm:"contest_type"`
	ContestLevel    string `json:"contest_level" xorm:"contest_level"`
	CreateTime      string `json:"create_time" xorm:"create_time"`
	StartTime       string `json:"start_time" xorm:"start_time"`
	School          string `json:"school" xorm:"school"`
	StudentSchoolID string `json:"student_school_id" xorm:"student_school_id"`
	College         string `json:"college" xorm:"college"`
	Semester        string `json:"semester" xorm:"semester"`
	Class           string `json:"student_class" xorm:"student_class"`
	Phone           string `json:"phone" xorm:"phone"`
	Email           string `json:"email" xorm:"email"`
	Team            string `json:"team_name" xorm:"team_name"`
	State           int    `json:"state" xorm:"state"`
	IsGroup         int    `json:"is_group"`
	Name            string `json:"name" xorm:"name"`
	Contest         string `json:"contest" xorm:"contest"`
	RejectReason    string `json:"reject_reason" xorm:"reject_reason"`
	DoUpload        bool   `json:"do_upload"`
	TeacherName     string `json:"teacher_name" xorm:"t_name"`
	Department      string `json:"department" xorm:"t_department"`
	Title           string `json:"title" xorm:"title"`
	Major           string `json:"major" xorm:"major"`
	ContestEntry    string `json:"contest_entry"`
}

type TeacherGetOneEnrollInformationReturn struct {
	ID              int64  `json:"id" xorm:"id"`
	Username        string `json:"username" xorm:"username"`
	StudentID       string `json:"student_id" xorm:"student_id"`
	TeamID          int64  `json:"team_id" xorm:"team_id"`
	ContestID       int64  `json:"contest_id"`
	ContestType     string `json:"contest_type" xorm:"contest_type"`
	CreateTime      string `json:"create_time" xorm:"create_time"`
	ContestLevel    string `json:"contest_level" xorm:"contest_level"`
	StudentSchoolID string `json:"student_school_id" xorm:"student_school_id"`
	StartTime       string `json:"start_time" xorm:"start_time"`
	School          string `json:"school" xorm:"school"`
	College         string `json:"college" xorm:"college"`
	Major           string `json:"major" xorm:"major"`
	Semester        string `json:"semester" xorm:"semester"`
	Class           string `json:"student_class" xorm:"student_class"`
	Phone           string `json:"phone" xorm:"phone"`
	Email           string `json:"email" xorm:"email"`
	Name            string `json:"name" xorm:"name"`
	Contest         string `json:"contest" xorm:"contest"`
}

func (EnrollInformationReturn) TableName() string {
	return "enroll_information"
}

func (EnrollInformation) TableName() string {
	return "enroll_information"
}

type EnrollDeleteId struct {
	ID []int64 `json:"ids"`
}

type EnrollIds EnrollDeleteId
type GradeIds EnrollDeleteId

type PassEnrollID struct {
	IDS []int64 `json:"ids"`
}

type PassGradeID struct {
	IDS []int64 `json:"ids"`
}

type PassContestID struct {
	IDS []int64 `json:"ids"`
}

type Team struct {
	TeamID    int64     `json:"team_id" xorm:"team_id pk autoincr" `
	TeamName  string    `json:"team_name" xorm:"team_name"`
	ContestID int64     `json:"contest_id" xorm:"contest_id"`
	Deleted   OftenTime `json:"deleted" xorm:"deleted"`
}

func (Team) TableName() string {
	return "team"
}
