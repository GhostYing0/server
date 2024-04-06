package public

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"server/utils/gredis"
	. "server/utils/mydebug"
	"server/utils/util"
	"strconv"
	"time"
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

func (self PublicLogic) UploadImg(file *multipart.FileHeader) (string, error) {
	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		// 返回值
		DPrintf("文件类型不合法")
		return "", errors.New("文件类型不合法")
	}

	//currentTime := time.Now().Format("20060102")
	// 生成目录文件夹，并错误判断
	//if err := os.MkdirAll("D:/GDesign/picture/img"+currentTime, 0755); err != nil {
	//	DPrintf("上传错误")
	//	appG.ResponseErr("MkdirAll失败")
	//	return
	//}
	if err := os.MkdirAll("D:/GDesign/picture/img", 0755); err != nil {
		DPrintf("上传错误")
		return "", errors.New("MkdirAll失败")
	}

	fileUnixName := strconv.FormatInt(time.Now().UnixNano(), 10)

	saveDir := path.Join("D:/GDesign/picture/img", fileUnixName+extName)

	return saveDir, nil
}
