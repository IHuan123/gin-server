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

func (con *SystemController) GetAllRoles(c *gin.Context) {
	var data []databases.Role
	sqlStr := `SELECT * FROM roles`
	err := databases.DB.Select(&data, sqlStr)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	con.Success(c, data)
}

type RoleMenu struct {
	MenuId int `json:"menu_id" db:"menu_id"`
	databases.RoleMenu
	databases.Menu
}

func (con *SystemController) GetRolesMenus(c *gin.Context) {
	rid := c.Query("rid")
	var menu []RoleMenu
	sqlStr := `SELECT * FROM role_menu r JOIN menus m USING(menu_id) WHERE r.role_id = ?`
	err := databases.DB.Select(&menu, sqlStr, rid)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	con.Success(c, menu)
}
