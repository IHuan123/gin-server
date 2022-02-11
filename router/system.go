package router

import (
	"github.com/gin-gonic/gin"
	"reactAdminServer/controllers/adminSystem"
	Middlewares "reactAdminServer/middlewares"
)

//系统数据
func InitSystemRouter(r *gin.Engine) {
	sCon := &adminSystem.SystemController{}
	sRouter := r.Group("/system")
	{
		//菜单
		sRouter.GET("/getAllMenus", Middlewares.TokenAuthMiddleware(), sCon.GetAllMenus)
		sRouter.GET("/getAllRoles", Middlewares.TokenAuthMiddleware(), sCon.GetAllRoles)
		sRouter.POST("/addMenu", Middlewares.TokenAuthMiddleware(), sCon.AddMenu)
		sRouter.PUT("/updateMenu", Middlewares.TokenAuthMiddleware(), sCon.UpdateMenu)
		sRouter.DELETE("/deleteMenu", Middlewares.TokenAuthMiddleware(), sCon.DeleteMenu)

		//角色
		sRouter.POST("/addRoles", Middlewares.TokenAuthMiddleware(), sCon.AddRole)
		sRouter.GET("/getRolesMenus", Middlewares.TokenAuthMiddleware(), sCon.GetRolesMenus)
		sRouter.PUT("/updateRoles", Middlewares.TokenAuthMiddleware(), sCon.UpdateRoles)
		sRouter.DELETE("/deleteRoles", Middlewares.TokenAuthMiddleware(), sCon.DelRoles)
		sRouter.PUT("/updateRolesMenus", Middlewares.TokenAuthMiddleware(), sCon.UpdateRolesMenus)

		//用户
		sRouter.GET("/s_users", Middlewares.TokenAuthMiddleware(), sCon.GetUser)
		sRouter.POST("/s_users", Middlewares.TokenAuthMiddleware(), sCon.AddUser)
		sRouter.PUT("/s_users", Middlewares.TokenAuthMiddleware(), sCon.EditUser)
		sRouter.DELETE("/s_users", Middlewares.TokenAuthMiddleware(), sCon.DelUser)
	}
}
