package adminSystem

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"reactAdminServer/databases"
)

type Menus struct {
	MenuId     int     `json:"menu_id" db:"menu_id"`
	Title      string  `json:"title" db:"title" binding:"required"`
	Path       string  `json:"path" db:"path" binding:"required"`
	Icon       string  `json:"icon" db:"icon" binding:"required"`
	RKey       string  `json:"key" db:"r_key" binding:"required"`
	Visible    *int     `json:"visible" db:"visible" binding:"required"`
	KeepAlive  *int     `json:"keep_alive" db:"keep_alive" binding:"required"`
	Weight     *int     `json:"weight" db:"weight" binding:"required"`
	ParentKey  string  `json:"parent_key" db:"parent_key"`
	Children   []Menus `json:"children" db:"children"`
	ParentName string  `json:"parent_name"`
}

func handleMenus(list []Menus) []Menus {
	var parent []Menus
	var children []Menus
	listLen := len(list)
	for i := 0; i < listLen; i++ {
		item := list[i]
		if item.ParentKey == "" {
			parent = append(parent, item)
		} else {
			children = append(children, item)
		}
	}
	for i, p := range parent {
		for _, c := range children {
			if p.RKey == c.ParentKey {
				c.ParentName = parent[i].Title
				parent[i].Children = append(parent[i].Children, c)
			}
		}
	}
	return parent
}

//获取全部菜单
func (con *SystemController) GetAllMenus(c *gin.Context) {
	var data []Menus
	sqlStr := `SELECT menu_id,title,r_path path,icon,r_key,visible,keep_alive,weight,parent_key FROM menus order by weight DESC`
	stmt, err := databases.DB.Prepare(sqlStr)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var menu Menus
		err = rows.Scan(&menu.MenuId, &menu.Title, &menu.Path, &menu.Icon, &menu.RKey, &menu.Visible, &menu.KeepAlive, &menu.Weight, &menu.ParentKey)
		if err != nil {
			con.Err(c, err.Error())
			return
		}
		data = append(data, menu)
	}
	res := handleMenus(data)
	con.Success(c, res)
}

//update
func (con *SystemController) UpdateMenu(c *gin.Context) {
	var params Menus
	if err := c.BindJSON(&params); err != nil {
		con.Err(c, err.Error())
		return
	}
	tx, err := databases.DB.Beginx() // 开启事务
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
	sqlStr := `UPDATE menus SET title=?, r_path=?, icon=?, r_key=?, visible=?, keep_alive=?, weight=?, parent_key=? WHERE menu_id=?`
	rs, err := tx.Exec(sqlStr, params.Title, params.Path, params.Icon, params.RKey, params.Visible, params.KeepAlive, params.Weight, params.ParentKey, params.MenuId)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	n, err := rs.RowsAffected() // 操作影响的行数
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	if n != 1 {
		con.Err(c, "数据更新失败！")
		return
	}
	con.Success(c, "update data success")
}

//add
func (con *SystemController) AddMenu(c *gin.Context) {
	var params Menus
	var sqlStr string
	err := c.BindJSON(&params)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	//adminSystem.Menus{MenuId:0, Title:"test", Path:"/test", Icon:"android", RKey:"test", Visible:(*int)(0xc00048a448), KeepAlive:(*int)(0xc00048a450), Weight:(*int)(0xc00048a458), ParentKey:"", Children:[]adminSystem.Menus(nil), ParentName:""}
	sqlStr = `INSERT INTO menus (title,r_path,icon,r_key,visible,keep_alive,weight,parent_key) values(:title,:path,:icon,:r_key,:visible,:keep_alive,:weight,:parent_key)`
	_,err = databases.DB.NamedExec(sqlStr,map[string]interface{}{
		"title": params.Title,
		"path": params.Path,
		"icon":params.Icon,
		"r_key":params.RKey,
		"visible":params.Visible,
		"keep_alive":params.KeepAlive,
		"weight":params.Weight,
		"parent_key":params.ParentKey,
	})
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	con.Success(c, "success")
}

//delete
func (con *SystemController) DeleteMenu(c *gin.Context) {
	var params struct {
		MenuIds []int `json:"menu_ids" db:"menu_ids" binding:"required"`
	}
	err := c.BindJSON(&params)
	if err != nil {
		con.Err(c, err.Error())
		return
	}

	tx, err := databases.DB.Beginx() // 开启事务
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
	sqlStr := `DELETE FROM menus where menu_id IN (?) `
	query, args, err := sqlx.In(sqlStr, params.MenuIds)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	_, err = tx.Exec(query, args...)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	roleMenuSql:=`DELETE FROM role_menu where menu_id IN (?) `
	rQuery, rArgs, err := sqlx.In(roleMenuSql, params.MenuIds)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	_, err = tx.Exec(rQuery, rArgs...)
	if err != nil {
		con.Err(c, err.Error())
		return
	}

	con.Success(c, "success")
}
