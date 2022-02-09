package adminSystem

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"reactAdminServer/databases"
)

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


//编辑角色基本信息
type UpdateRolesParams struct {
	Id int `json:"id"`
	Name string `json:"name"`
	DataScope string `json:"dataScope"`
}
func (con *SystemController) UpdateRoles(c *gin.Context){
	var params UpdateRolesParams
	if err:=c.BindJSON(&params);err!=nil{
		con.Err(c,err.Error())
		return
	}
	fmt.Printf("%#v\n",params)
	sqlStr:=`UPDATE roles SET name = :name,dataScope = :dataScope WHERE id = :id`
	_, err := databases.DB.NamedExec(sqlStr, map[string]interface{}{
		"name":params.Name,
		"dataScope":params.DataScope,
		"id":params.Id,
	})
	if err != nil {
		con.Err(c,err.Error())
		return
	}
	con.Success(c,"update success")
}
//删除
func (con *SystemController) DelRoles(c *gin.Context)  {
	var params struct{
		Ids []int `json:"ids"`
	}
	if err:=c.BindJSON(&params);err!=nil{
		con.Err(c,err.Error())
		return
	}
	fmt.Printf("%#v\n",params)

	tx, err := databases.DB.Beginx() // 开启事务
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	//处理事务回滚和提交
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			fmt.Println("rollback")
			_ = tx.Rollback() // err is non-nil; don't change it
		} else {
			_ = tx.Commit() // err is nil; if Commit returns error update err
			fmt.Println("commit")
		}
	}()

	sqlStr := `DELETE FROM roles WHERE id IN (?)`
	query, args, err := sqlx.In(sqlStr, params.Ids)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	_, err = tx.Exec(query, args...)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	con.Success(c, "success")
}

type RoleMenu struct {
	MenuId int `json:"menu_id" db:"menu_id"`
	databases.RoleMenu
	databases.Menu
}
//通过role_id获取菜单
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
