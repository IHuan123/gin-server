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
		sRouter.GET("/getAllMenus", Middlewares.TokenAuthMiddleware(), sCon.GetAllMenus)
		sRouter.GET("/getAllRoles", Middlewares.TokenAuthMiddleware(), sCon.GetAllRoles)
		sRouter.POST("/addMenu", Middlewares.TokenAuthMiddleware(), sCon.AddMenu)
		sRouter.PUT("/updateMenu", Middlewares.TokenAuthMiddleware(), sCon.UpdateMenu)
		sRouter.DELETE("/deleteMenu", Middlewares.TokenAuthMiddleware(), sCon.DeleteMenu)
		sRouter.GET("/getRolesMenus", Middlewares.TokenAuthMiddleware(), sCon.GetRolesMenus)
		sRouter.PUT("/updateRoles",Middlewares.TokenAuthMiddleware(),sCon.UpdateRoles)
		sRouter.DELETE("/deleteRoles",Middlewares.TokenAuthMiddleware(),sCon.DelRoles)
	}
}
