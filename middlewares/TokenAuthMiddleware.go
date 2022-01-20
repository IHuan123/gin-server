package Middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"reactAdminServer/models"
)

//获取要验证token接口的数据
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var msg string
		tokenString := ctx.GetHeader("Authorization")
		//fmt.Println(tokenString)
		//vcalidate token formate
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token不存在"})
			ctx.Abort()
			return
		}
		token, _, err := models.ParseToken(tokenString)
		if err != nil || !token.Valid {
			fmt.Println(token.Valid)
		}
		if ve, ok := err.(*jwt.ValidationError); ok { //官方写法招抄就行
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				fmt.Println("错误的token")
				msg = "错误的token"
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				fmt.Println("token过期或未启用")
				msg = "token过期或未启用"
			} else {
				fmt.Println("无法处理这个token", err)
				msg = "无法处理这个token"
			}
			ctx.JSON(http.StatusOK, gin.H{"code": 401, "msg": msg})
			ctx.Abort()
			return
		}
	}
}
