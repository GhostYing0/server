package models

type EnrollForm struct {
	ContestID         int64  `json:"contest_id" xorm:"contest_id"`
	StudentName       string `json:"student_name"`
	TeamName          string `json:"team_name"`
	Handle            int64  `json:"handle_team"`
	CollegeID         int64  `json:"college"`
	MajorID           int64  `json:"major"`
	GuidanceTeacher   string `json:"guidance_teacher"`
	TeacherDepartment string `json:"teacher_department"`
	TeacherTitle      string `json:"teacher_title"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	State             int    `json:"state"`
}

type EnrollInformationForm struct {
	ID         int64  `json:"id"`
	UserName   string `json:"username"`
	Name       string `json:"name"`
	TeamID     string `json:"team_id"`
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
	TeamID    string `json:"team_id" xorm:"team_id"`
	//ContestID    string    `json:"contest_id" xorm:"contest_id"`
	ContestID    int64     `json:"contest_id" xorm:"contest_id"`
	CreateTime   string    `json:"create_time" xorm:"create_time"`
	SchoolID     int64     `json:"school_id" xorm:"school_id"`
	Phone        string    `json:"phone" xorm:"phone"`
	Email        string    `json:"email" xorm:"email"`
	State        int       `json:"state" xorm:"state"`
	RejectReason string    `json:"reject_reason" xorm:"reject_reason"`
	Deleted      OftenTime `json:"deleted" xorm:"deleted"`
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
	CreateTime        string    `json:"create_time" xorm:"create_time"`
	StartTime         string    `json:"start_time" xorm:"start_time"`
	Deadline          string    `json:"deadline" xorm:"deadline"`
	State             int       `json:"state" xorm:"state"`
	Describe          string    `json:"desc" xorm:"describe"`
	Deleted           OftenTime `json:"deleted" xorm:"deleted"`
	Student           `xorm:"extends"`
	School            string `xorm:"school"`
	College           string `xorm:"college"`
	Semester          string `xorm:"semester"`
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
	ID        int64  `json:"id" xorm:"e_id"`
	StudentID string `json:"student_id" xorm:"student_id"`
	TeamID    string `json:"team_id" xorm:"team_id"`
	//ContestID    string    `json:"contest_id" xorm:"contest_id"`
	ContestID int64 `json:"contest_id" xorm:"contest_id"`

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
	Student      `xorm:"extends"`
	School       string `xorm:"school"`
	College      string `xorm:"college"`
	Semester     string `xorm:"semester"`
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
	ID           int64  `json:"id" xorm:"id"`
	Username     string `json:"username" xorm:"username"`
	StudentID    string `json:"student_id" xorm:"student_id"`
	TeamID       string `json:"team_id" xorm:"team_id"`
	ContestType  string `json:"contest_type" xorm:"contest_type"`
	CreateTime   string `json:"create_time" xorm:"create_time"`
	StartTime    string `json:"start_time" xorm:"start_time"`
	School       string `json:"school" xorm:"school"`
	College      string `json:"college" xorm:"college"`
	Semester     string `json:"semester" xorm:"semester"`
	Class        string `json:"student_class" xorm:"student_class"`
	Phone        string `json:"phone" xorm:"phone"`
	Email        string `json:"email" xorm:"email"`
	State        int    `json:"state" xorm:"state"`
	Name         string `json:"name" xorm:"name"`
	Contest      string `json:"contest" xorm:"contest"`
	RejectReason string `json:"reject_reason" xorm:"reject_reason"`
	DoUpload     bool   `json:"do_upload"`
}

type TeacherGetOneEnrollInformationReturn struct {
	ID              int64  `json:"id" xorm:"id"`
	Username        string `json:"username" xorm:"username"`
	StudentID       string `json:"student_id" xorm:"student_id"`
	TeamID          string `json:"team_id" xorm:"team_id"`
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
