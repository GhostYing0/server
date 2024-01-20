package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Test() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("test")
		//fmt.Println(utils.GenSnowflakeID())
		//fmt.Println(uuid.CreateUUIDByTimeAndMAC())
		c.Next()
	}
}
