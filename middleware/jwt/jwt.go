package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"server/utils/e"
	"server/utils/util"
	"strings"
)

func JwtTokenCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.SUCCESS
		data := make(map[string]interface{})
		userIdInt := int64(0)
		userRole := -1
		token := c.Request.Header.Get("BackServer-token")
		message := "解析成功"

		if !strings.Contains(c.Request.RequestURI, "/login") &&
			!strings.Contains(c.Request.RequestURI, "/register") {
			if token == "" {
				code = e.ERROR_AUTH_CHECK_TOKEN_EMPTY
				message = "Token为空"
			} else {
				claims, err := util.ParseToken(token)
				if err != nil {
					switch err.(*jwt.ValidationError).Errors {
					case jwt.ValidationErrorExpired:
						code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
						message = "Token已过期"
					default:
						code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
						fmt.Println(err)
						message = "解析Token失败"
					}
				} else if claims == nil {
					code = e.ERROR_AUTH_TOKEN_PARSE
					message = "Token错误"
				} else if strings.Contains(c.Request.RequestURI, "/cms") && claims.Role != 0 {
					code = e.ERROR_AUTH_TOKEN_DIFF
					message = "Token权限不足"
				} else {
					userIdInt, err = com.StrTo(claims.ID).Int64()
					if err != nil {
						code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
						message = "解析Token失败"
					}
					userRole = claims.Role
				}
			}
			if code == e.SUCCESS {
				c.Set("user_id", userIdInt)
				c.Set("role", userRole)
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code":    code,
					"message": message,
					"data":    data,
				})

				c.Abort()
				return
			}
		}
		c.Next()
	}
}
