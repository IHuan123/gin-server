package adminSystem

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"reactAdminServer/databases"
	"strings"
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
//新增
func (con *SystemController) AddRole(c *gin.Context) {
	var params databases.Role
	var sqlStr string
	err := c.BindJSON(&params)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	sqlStr = `INSERT INTO roles (name,dataScope) values(:name,:dataScope)`
	_,err = databases.DB.NamedExec(sqlStr,map[string]interface{}{
		"name": params.Name,
		"dataScope": params.DataScope,
	})
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	con.Success(c, "success")
}

//编辑角色基本信息
type UpdateRolesParams struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	DataScope string `json:"dataScope"`
}
func (con *SystemController) UpdateRoles(c *gin.Context) {
	var params UpdateRolesParams
	if err := c.BindJSON(&params); err != nil {
		con.Err(c, err.Error())
		return
	}
	fmt.Printf("%#v\n", params)
	sqlStr := `UPDATE roles SET name = :name,dataScope = :dataScope WHERE id = :id`
	_, err := databases.DB.NamedExec(sqlStr, map[string]interface{}{
		"name":      params.Name,
		"dataScope": params.DataScope,
		"id":        params.Id,
	})
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	con.Success(c, "update success")
}

//删除
func (con *SystemController) DelRoles(c *gin.Context) {
	var params struct {
		Ids []int `json:"ids"`
	}
	if err := c.BindJSON(&params); err != nil {
		con.Err(c, err.Error())
		return
	}

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
	//删除角色
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
	//删除角色所有的访问权限
	delSqlStr := `DELETE FROM role_menu WHERE role_id IN (?)`
	dQuery, dArgs, err := sqlx.In(delSqlStr, params.Ids)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	_, err = tx.Exec(dQuery, dArgs...)
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
	rid := c.Query("id")
	var menu []int
	sqlStr := `SELECT m.menu_id FROM role_menu r JOIN menus m USING(menu_id) WHERE r.role_id = ?`
	err := databases.DB.Select(&menu, sqlStr, rid)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	con.Success(c, menu)
}

//更新角色菜单
type UpdateRoleMenu struct {
	MIds []int `json:"mIds"`
	Rid  int   `json:"rid"`
}

type RoleMenuInsert struct {
	RoleId int `json:"role_id" db:"role_id"`
	MenuId int `json:"menu_id" db:"menu_id"`
}
func (rm RoleMenuInsert) Value() (driver.Value, error) {
	return []interface{}{rm.RoleId, rm.MenuId}, nil
}

//处理参数
func (rmi UpdateRoleMenu) handleUpdateRoleParams() ([]RoleMenuInsert, error) {
	var res []RoleMenuInsert
	mIds := rmi.MIds
	for _, v := range mIds {
		item := RoleMenuInsert{
			RoleId: rmi.Rid,
			MenuId: v,
		}
		res = append(res, item)
	}
	return res, nil
}

//批量插入
func BatchInsertRoleMenu(tx *sql.Tx, roleMenu []interface{}) error {
	sqlStr := "INSERT INTO role_menu (role_id, menu_id) VALUES (?),(?)"
	query, args, err := sqlx.In(
		sqlStr,
		roleMenu..., // 如果arg实现了 driver.Valuer, sqlx.In 会通过调用 Value()来展开它
	)
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, args...)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

// BatchInsertUsers 自行构造批量插入的语句
func BatchInsertRoleMenu2(tx *sql.Tx, roleMenu []RoleMenuInsert) error {
	// 存放 (?, ?) 的slice
	valueStrings := make([]string, 0, len(roleMenu))
	// 存放values的slice
	valueArgs := make([]interface{}, 0, len(roleMenu) * 2)
	// 遍历users准备相关数据
	for _, u := range roleMenu {
		// 此处占位符要与插入值的个数对应
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, u.RoleId)
		valueArgs = append(valueArgs, u.MenuId)
	}
	// 自行拼接要执行的具体语句
	stmt := fmt.Sprintf("INSERT INTO role_menu (role_id, menu_id) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := tx.Exec(stmt, valueArgs...)
	return err
}
//更新角色菜单权限
func (con *SystemController) UpdateRolesMenus(c *gin.Context) {
	var params UpdateRoleMenu
	if err := c.BindJSON(&params); err != nil {
		con.Err(c, err.Error())
		return
	}
	res, _ := params.handleUpdateRoleParams()
	tx, err := databases.DB.Begin()
	if err != nil {
		con.Err(c, err.Error())
		return
	}
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
	//清除旧数据
	clearSql := "DELETE FROM role_menu WHERE role_id = ?"
	_, err = tx.Exec(clearSql, params.Rid)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	err = BatchInsertRoleMenu2(tx, res)
	if err != nil {
		con.Err(c, err.Error())
		return
	}

	con.Success(c, "success")
}
