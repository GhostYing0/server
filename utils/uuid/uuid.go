package uuid

import (
	. "github.com/satori/go.uuid"
)

func CreateUUIDByTimeAndMAC() UUID {
	return NewV1()
}

func CreateUUIDByPOSIX() UUID {
	return NewV2(0)
}

func CreateUUIDByNameSpace() UUID {
	return NewV3(CreateUUIDByPOSIX(), "user")
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
