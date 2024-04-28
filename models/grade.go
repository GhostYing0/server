package models

type GradeForm struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Contest      string `json:"contest"`
	School       string `json:"school"`
	Grade        string `json:"grade"`
	Certificate  string `json:"certificate"`
	CreateTime   string `json:"create_time"`
	PS           string `json:"ps"`
	RejectReason string `json:"reject_reason"`
	State        int    `json:"state" xorm:"state"`
}

type UploadGradeForm struct {
	EnrollID          int64  `json:"enroll_id"`
	Grade             int    `json:"prize"`
	Certificate       string `json:"certificate"`
	PS                string `json:"ps"`
	RewardTime        string `json:"reward_time"`
	GuidanceTeacher   string `json:"guidance_teacher"`
	TeacherDepartment string `json:"teacher_department"`
	TeacherTitle      string `json:"teacher_title"`
	ContestLevel      string `json:"contest_level"`
}

type GradeInformation struct {
	ID              int64     `json:"id" xorm:"id"`
	StudentID       string    `json:"student_id" xorm:"student_id" `
	ContestID       int64     `json:"contest_id" xorm:"contest_id"`
	SchoolID        int64     `json:"school_id" xorm:"school_id"`
	Grade           int       `json:"grade_id" xorm:"grade_id"`
	Certificate     string    `json:"certificate" xorm:"certificate"`
	State           int       `json:"state" xorm:"state"`
	CreateTime      string    `json:"create_time" xorm:"create_time"`
	UpdateTime      string    `json:"update_time" xorm:"update_time"`
	PS              string    `json:"ps" xorm:"ps"`
	RejectReason    string    `json:"reject_reason" xorm:"reject_reason"`
	GuidanceTeacher string    `json:"guidance_teacher" xorm:"guidance_teacher"`
	EnrollID        int64     `json:"enroll_id" xorm:"enroll_id"`
	RewardTime      OftenTime `json:"reward_time" xorm:"reward_time"`
	Deleted         OftenTime `json:"deleted" xorm:"deleted"`
}

type CurStudentGrade struct {
	GradeInformation `xorm:"extends"`
	School           string `xorm:"school"`
	Name             string `xorm:"name"`
	Contest          string `xorm:"contest"`
	ContestType      string `xorm:"type"`
}

func (CurStudentGrade) TableName() string {
	return "grade"
}

type GradeStudentSchoolContestAccount struct {
	GradeInformation `xorm:"extends"`
	Student          `xorm:"extends"`
	School           `xorm:"extends"`
	Contest          string `xorm:"contest"`
	ContestType      string `xorm:"type"`
	Username         string `xorm:"username"`
}

func (GradeStudentSchoolContestAccount) TableName() string {
	return "grade"
}

type ReturnGradeInformation struct {
	ID           int64  `json:"id" xorm:"id"`
	Username     string `json:"username" xorm:"username"`
	Name         string `json:"name" xorm:"name" `
	Contest      string `json:"contest" xorm:"contest"`
	School       string `json:"school" xorm:"school"`
	Grade        string `json:"grade" xorm:"grade"`
	Certificate  string `json:"certificate" xorm:"certificate"`
	State        int    `json:"state" xorm:"state"`
	ContestType  string `json:"contest_type" xorm:"contest_type"`
	PS           string `json:"ps" xorm:"ps"`
	RejectReason string `json:"reject_reason" xorm:"reject_reason"`
	CreateTime   string `json:"create_time" xorm:"create_time"`
}

func (GradeInformation) TableName() string {
	return "grade"
}

type Prize struct {
	PrizeID int    `json:"prize_id" xorm:"prize_id"`
	Prize   string `json:"prize" xorm:"prize"`
}

func (Prize) TableName() string {
	return "prize"
}
