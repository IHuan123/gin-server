package adminSystem

import (
	"github.com/gin-gonic/gin"
	"reactAdminServer/controllers/base"
	"reactAdminServer/databases"
)

type SystemController struct {
	base.BaseController
}
type Syer interface {
	GetAllRoles(c *gin.Context)
}
type Roles struct {
	Id int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	DataScope string `json:"dataScope" db:"dataScope"`
}


func (con *SystemController) GetAllRoles(c *gin.Context){
	var data []Roles
	sqlStr := `SELECT * FROM roles`
	err := databases.DB.Select(&data,sqlStr)
	if err!=nil{
		con.Err(c,err.Error())
		return
	}
	con.Success(c,data)
}