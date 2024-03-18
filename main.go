package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	controllerUser "server/controller"
	controllerCms "server/controller/cms"
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
	r.Use(jwt.JwtTokenCheck())
	fmt.Println("Server start")

	backGround := r.Group("api/")
	controllerUser.RegisterRoutes(backGround)

	backCmsGround := r.Group("api/cms")
	controllerCms.RegisterRoutes(backCmsGround)

	readTimeout := time.Second * 10
	writeTimeout := time.Second * 10
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           "127.0.0.1:9006",
		Handler:        r,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	server.ListenAndServe()
}
