package public

import (
	"errors"
	"fmt"
	"server/utils/gredis"
	. "server/utils/mydebug"
	"server/utils/util"
	"strconv"
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
