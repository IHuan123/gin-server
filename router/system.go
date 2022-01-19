package router

import (
	"github.com/gin-gonic/gin"
	"reactAdminServer/controllers/adminSystem"
)

//系统数据
func InitSystemRouter(r *gin.Engine){
	sCon := &adminSystem.SystemController{}
	sRouter := r.Group("/system")
	{
		sRouter.GET("/getAllMenus",sCon.GetAllMenus)
		sRouter.GET("/getAllRoles",sCon.GetAllRoles)
		sRouter.PUT("/updateMenu",sCon.UpdateMenu)
		sRouter.DELETE("/deleteMenu",sCon.DeleteMenu)
	}
}
