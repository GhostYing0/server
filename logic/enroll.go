package logic

import (
	"encoding/json"
	"errors"
	. "server/database"
	"server/models"
	. "server/utils/mydebug"
)

type EnrollLogic struct{}

var DefaultEnrollLogic = EnrollLogic{}

func (self EnrollLogic) InsertEnrollInformation(members []byte, contest string, create_time string, school string, phone string, email string) error {
	session := MasterDB.NewSession()
	if err := session.Begin(); err != nil {
		DPrintf("InsertEnrollInformation session.Begin() 发生错误:", err)
		return err
	}
	defer func() {
		err := session.Close()
		if err != nil {
			DPrintf("InsertEnrollInformation session.Close() 发生错误:", err)
		}
	}()

	exist, err := session.Where("name = ?", contest).Exist()
	if err != nil {
		DPrintf("InsertEnrollInformation 查询竞赛发生错误:", err)
		return err
	}
	if !exist {
		DPrintf("InsertEnrollInformation 竞赛不存在")
		return errors.New("竞赛不存在")
	}

	list := new([]string)
	err = json.Unmarshal(members, list)
	if err != nil {
		DPrintf("InsertEnrollInformation Unmarshal 发生错误:", err)
		return err
	}

	_ = &models.EnrollInformation{
		Contest: contest,
		School:  school,
		Phone:   phone,
		Email:   email,
	}

	return session.Commit()
}
