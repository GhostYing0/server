package models

import "time"

type StudentEntryParam struct {
	ContestantID int `json:"contestant_id" xorm:"contestant_id"`
	ContestID    int `json:"contest_id" xorm:"contest_id"`
}

type ContestantInfo struct {
	ID           int       `xorm:"id"`
	ContestantID int       `xorm:"contestant_id"`
	ContestID    int       `xorm:"contest_id"`
	EntryTime    time.Time `xorm:"entry_time"`
	Deleted      time.Time `xorm:"deleted"`
}

type ContestGrade struct {
	ID           int       `xorm:"id"`
	ContestantID int       `xorm:"contestant_id"`
	ContestID    int       `xorm:"contest_id"`
	EntryTime    time.Time `xorm:"entry_time"`
	Awards       string    `xorm:"awards"`
	Deleted      time.Time `xorm:"deleted"`
}

type RegistrationInfo ContestGrade

type EntryContestParam struct {
	Contestant string `json:"contestant"`
	Contest    string `json:"contest"`
}
