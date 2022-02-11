package router

import (
	"github.com/gin-gonic/gin"
	"reactAdminServer/controllers/user"
	"reactAdminServer/middlewares"
)

func InitLoginRouter(r *gin.Engine) {
	uCon := user.UserController{}
	r.POST("/login", uCon.Login)
	userGroup := r.Group("/user")
	//使用token中间件
	userGroup.Use(Middlewares.TokenAuthMiddleware())

	{
		userGroup.GET("/info", Middlewares.TokenAuthMiddleware(), uCon.GetUser)
		userGroup.GET("/getMenus", Middlewares.TokenAuthMiddleware(), uCon.GetMenus)
	}
}
