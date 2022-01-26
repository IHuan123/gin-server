package databases

// menus
type Menu struct {
	MenuId    int    `json:"menu_id" db:"menu_id"`
	Title     string `json:"title" db:"title"`
	RPath     string `json:"r_path" db:"r_path"`
	Icon      string `json:"icon" db:"icon"`
	RKey      string `json:"r_key" db:"r_key"`
	Visible   *int   `json:"visible" db:"visible"`
	KeepAlive *int   `json:"keep_alive" db:"keep_alive"`
	Weight    *int   `json:"weight" db:"weight"`
	ParentKey string `json:"parent_key" db:"parent_key"`
}

//roles
type Role struct {
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	DataScope string `json:"dataScope" db:"dataScope"`
}

//role_menu
type RoleMenu struct {
	Id     int `json:"id" db:"id"`
	RoleId int `json:"role_id" db:"role_id"`
	MenuId int `json:"menu_id" db:"menu_id"`
}

//t_uesr
type User struct {
	Uid        int    `json:"uid" db:"uid"`
	UserName   string `json:"username" db:"username"`
	Avatar     string `json:"avatar" db:"avatar"`
	DeptId     int    `json:"deptId" db:"deptId"`
	Email      string `json:"email" db:"email"`
	Enabled    int    `json:"enabled" db:"enabled"`
	Phone      string `json:"phone" db:"phone"`
	Sex        string `json:"sex" db:"sex"`
	Roles      string `json:"roles" db:"roles"`
	CreateTime string `json:"createTime" db:"createTime"`
}
