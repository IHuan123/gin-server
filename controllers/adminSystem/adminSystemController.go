package adminSystem

import (
	"github.com/gin-gonic/gin"
	"reactAdminServer/controllers/base"
)

type SystemController struct {
	base.BaseController
}
type Syer interface {
	GetAllRoles(c *gin.Context)
}



