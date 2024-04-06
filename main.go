package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	controllerUser "server/controller"
	controllerCms "server/controller/cms"
	controllerPublic "server/controller/public"
	"server/database"
	_ "server/database"
	"server/middleware/jwt"
	"server/utils/gredis"
	"time"
)

func init() {
	database.Setup()
	gredis.Setup()
}

func main() {
	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080", "http://10.1.31.107:8080"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"BackServer-token", "Content-Type"}

	r.Use(cors.New(config))
	fmt.Println("Server start")

	//r.StaticFS("/picture", http.Dir("D:/GDesign/code/server/picture/img"))
	public := r.Group("api/public")
	controllerPublic.RegisterRoutes(public)

	r.Use(jwt.JwtTokenCheck())
	backGround := r.Group("api")
	controllerUser.RegisterRoutes(backGround)

	backCmsGround := r.Group("api/cms")
	controllerCms.RegisterRoutes(backCmsGround)

	readTimeout := time.Second * 10
	writeTimeout := time.Second * 10
	maxHeaderBytes := 1 << 25

	server := &http.Server{
		Addr:           "127.0.0.1:9006",
		Handler:        r,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	server.ListenAndServe()
}
