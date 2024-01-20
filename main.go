package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	controlleCms "server/controller/cms"
	controlleStudent "server/controller/student"
	controlleTeacher "server/controller/teacher"
	"server/database"
	_ "server/database"
	"server/middleware/jwt"
	"time"
)

func main() {
	database.Setup()
	r := gin.New()

	r.Use(jwt.JwtTokenCheck())
	fmt.Println("Server start")

	studentG := r.Group("api/student")
	//studentG.Use(jwt.JwtTokenCheck())
	controlleStudent.RegisterRoutes(studentG)

	teacherG := r.Group("api/teacher")
	//teacherG.Use(jwt.JwtTokenCheck())
	controlleTeacher.RegisterRoutes(teacherG)

	backG := r.Group("api/cms")
	//backG.Use(jwt.JwtTokenCheck())
	controlleCms.RegisterRoutes(backG)

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
