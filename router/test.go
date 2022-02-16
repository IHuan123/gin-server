package router

import (
	"github.com/gin-gonic/gin"
	"reactAdminServer/controllers/testController"
)


func InitTestRouter(r *gin.Engine){
	tCon := testController.TestController{}
	rGroup := r.Group("/test")
	{
		rGroup.GET("/v1",tCon.V1)
	}
}
