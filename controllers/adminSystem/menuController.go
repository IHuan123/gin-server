package adminSystem

import (
	"github.com/gin-gonic/gin"
	"reactAdminServer/databases"
)

type Menus struct {
	MenuId    int    `json:"menu_id" db:"menu_id"`
	Title     string `json:"title" db:"title"`
	Path      string `json:"path" db:"path"`
	Icon      string `json:"icon" db:"icon"`
	RKey      string `json:"key" db:"r_key"`
	Visible   int    `json:"visible" db:"visible"`
	KeepAlive int    `json:"keep_alive" db:"keep_alive"`
	Weight    int    `json:"weight" db:"weight"`
	ParentKey string `json:"parent_key" db:"parent_key"`
	Children []Menus `json:"children" db:"children"`
	ParentName string `json:"parent_name"`
}
func handleMenus(list []Menus) []Menus {
	var parent []Menus
	var children []Menus
	listLen := len(list)
	for i := 0; i < listLen; i++ {
		item := list[i]
		if item.ParentKey == "" {
			parent = append(parent, item)
		}else{
			children = append(children,item)
		}
	}
	for i,p := range parent{
		for _,c := range children {
			if p.RKey == c.ParentKey{
				c.ParentName = parent[i].Title
				parent[i].Children = append(parent[i].Children,c)
			}
		}
	}
	return parent
}

//获取全部菜单
func (con *SystemController) GetAllMenus(c *gin.Context){
	var data []Menus
	sqlStr := `SELECT menu_id,title,r_path path,icon,r_key,visible,keep_alive,weight,parent_key FROM menus`
	stmt,err := databases.DB.Prepare(sqlStr)
	if err != nil {
		con.Err(c,err.Error())
		return
	}
	defer stmt.Close()
	rows,err := stmt.Query()
	if err != nil {
		con.Err(c,err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var menu Menus
		err = rows.Scan(&menu.MenuId,&menu.Title,&menu.Path,&menu.Icon,&menu.RKey,&menu.Visible,&menu.KeepAlive,&menu.Weight,&menu.ParentKey)
		if err != nil {
			con.Err(c,err.Error())
			return
		}
		data = append(data,menu)
	}
	res := handleMenus(data)
	con.Success(c,res)
}

//update
func (con *SystemController) UpdateMenu(c *gin.Context)  {
	con.Success(c,"success")
}
//add
func (con *SystemController) AddMenu(c *gin.Context)  {
	con.Success(c,"success")
}
//delete
func (con *SystemController) DeleteMenu(c *gin.Context)  {
	con.Success(c,"success")
}