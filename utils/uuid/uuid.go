package uuid

import (
	. "github.com/satori/go.uuid"
	"strconv"
	"time"
)

func CreateUUIDByTimeAndMAC() UUID {
	return NewV1()
}

func CreateUUIDByPOSIX() UUID {
	return NewV2(0)
}

func CreateUUIDByNameSpace(username, password, name, gender, semester, college, school, class string, role int, time time.Time) UUID {
	uuidName := username + password + name + gender + semester + college + school + class + strconv.Itoa(role) + time.String()
	return NewV3(CreateUUIDByPOSIX(), uuidName)
}

func CreateUUIDByNameSpaceTeacher() UUID {
	return NewV3(CreateUUIDByPOSIX(), "teacher")
}

func CreateUUIDByRandom() UUID {
	return NewV4()
}

func CreateUUIDBySHA_1() UUID {
	return NewV5(CreateUUIDByPOSIX(), "user")
}

func CreateUUIDByNewV5() UUID {
	return Must(CreateUUIDBySHA_1(), nil)
}
