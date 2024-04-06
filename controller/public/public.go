package public

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	logic "server/logic/public"
	"server/middleware/jwt"
	"server/utils/app"
	. "server/utils/mydebug"
	"strings"
)

type PublicController struct{}

func (self PublicController) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/get_info", jwt.JwtTokenCheck(), self.GetInfo) // 从Redis里读token
	g.GET("/logout", jwt.JwtTokenCheck(), self.Logout)
	g.POST("/v1/upload" /*jwt.JwtTokenCheck(),*/, self.Upload)
	g.StaticFS("/picture", http.Dir("D:/GDesign/picture/img"))
}

func (PublicController) GetInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	var err error
	data := make(map[string]interface{})
	token := c.Query("token")

	id, username, role, err := logic.DefaultPublic.GetInfo(token)
	if err != nil {
		DPrintf("User GetInfo 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	data["id"] = id
	data["username"] = username
	data["role"] = role

	appG.ResponseSucMsg(data)
}

// Logout
func (PublicController) Logout(c *gin.Context) {
	appG := app.Gin{C: c}

	token := c.Query("token")

	if len(token) <= 0 {
		DPrintf("token为空")
		appG.ResponseErr("token为空")
		return
	}

	err := logic.DefaultPublic.Logout(token)
	if err != nil {
		DPrintf("登出发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	appG.ResponseSucMsg("登出成功")
}

//func (PublicController) Upload(c *gin.Context) {
//	appG := app.Gin{C: c}
//
//	file, header, err := c.Request.FormFile("file")
//	if err != nil {
//		DPrintf("err:", err)
//		appG.ResponseErr(err.Error())
//		return
//	}
//	defer file.Close()
//
//	// 定义文件保存路径
//	filename := header.Filename
//	if filename == "" {
//		filename = "uploaded_file"
//	}
//	filePath := "D:/GDesign/picture/img/" + filename
//
//	// 创建保存文件的目录（如果不存在）
//	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
//		os.Mkdir("./uploads", 0755)
//	}
//
//	// 保存文件到磁盘
//	out, err := os.Create(filePath)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
//		return
//	}
//	defer out.Close()
//
//	// 将上传的文件拷贝到新的位置
//	_, err = io.Copy(out, file)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file"})
//		return
//	}
//
//	appG.ResponseSuc("图片上传成功")
//
//}

// Logout
func (PublicController) Upload(c *gin.Context) {
	appG := app.Gin{C: c}

	fmt.Println("Asdadadadasdasd")
	token := c.Query("token")
	fmt.Println("token:", token)

	//if len(token) <= 0 {
	//	DPrintf("token为空")
	//	appG.ResponseErr("token为空")
	//	return
	//}

	file, err := c.FormFile("file[raw]")

	if err != nil {
		DPrintf("FormFile err:", err)
		appG.ResponseErr(err.Error())
		return
	}
	if file == nil {
		DPrintf("logic.DefaultPublic.Upload file 为空")
		appG.ResponseErr("file 为空")
		return
	}
	saveDir, err := logic.DefaultPublic.UploadImg(file)
	if err != nil {
		DPrintf("logic.DefaultPublic.Upload 发生错误:", err)
		appG.ResponseErr(err.Error())
		return
	}

	err = c.SaveUploadedFile(file, saveDir)
	if err != nil {
		DPrintf("上传错误")
		// 返回值
		appG.ResponseErr("文件保存失败")
		return
	}

	imageurl := strings.Replace(saveDir, "D:/GDesign/picture/img", "http://127.0.0.1:9006/api/public/picture", -1)

	appG.ResponseSucMsg(gin.H{"imageurl": imageurl}, "上传成功")
}
