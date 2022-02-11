package adminSystem

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"reactAdminServer/databases"
)

//查询
//获取所有用户
func (con *SystemController) GetUser(c *gin.Context) {
	var data []databases.User
	sqlStr := `SELECT uid,username,avatar,deptId,email,enabled,phone,sex,roles,createTime FROM t_user`
	err := databases.DB.Select(&data, sqlStr)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	con.Success(c, data)
}

//新增
type AddUserParams struct {
	databases.User
}

func (con *SystemController) AddUser(c *gin.Context) {
	var params EditUserParams
	if err := c.BindJSON(&params); err != nil {
		con.Err(c, err.Error())
		return
	}
	//sqlStr:=`UPDATE t_user SET username=:username,avatar=:avatar,phone=:phone,roles=:roles,enabled=:enabled,sex=:sex,email=:email WHERE uid = :uid`
	sqlStr := `INSERT INTO t_user(username,avatar,phone,roles,enabled,sex,email,password,deptId) VALUES (:username,:avatar,:phone,:roles,:enabled,:sex,:email,:password,:deptId)`
	_, err := databases.DB.NamedExec(sqlStr, map[string]interface{}{
		"avatar":   params.Avatar,
		"deptId":   0,
		"email":    params.Email,
		"enabled":  params.Enabled,
		"phone":    params.Phone,
		"roles":    params.Roles,
		"sex":      params.Sex,
		"uid":      params.Uid,
		"password": "123456",
		"username": params.UserName,
	})
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	con.Success(c, "success")
}

//编辑
type EditUserParams = databases.User

func (con *SystemController) EditUser(c *gin.Context) {
	var params EditUserParams
	if err := c.BindJSON(&params); err != nil {
		con.Err(c, err.Error())
		return
	}
	sqlStr := `UPDATE t_user SET username=:username,avatar=:avatar,phone=:phone,roles=:roles,enabled=:enabled,sex=:sex,email=:email WHERE uid = :uid`
	_, err := databases.DB.NamedExec(sqlStr, map[string]interface{}{
		"avatar":   params.Avatar,
		"deptId":   params.DeptId,
		"email":    params.Email,
		"enabled":  params.Enabled,
		"phone":    params.Phone,
		"roles":    params.Roles,
		"sex":      params.Sex,
		"uid":      params.Uid,
		"username": params.UserName,
	})
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	con.Success(c, "success")
}

//删除
func (con *SystemController) DelUser(c *gin.Context) {
	var params struct {
		UIds []int `json:"uIds" binding:"required"`
	}
	err := c.BindJSON(&params)
	fmt.Println(params)
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
	sqlStr := `DELETE FROM t_user where uid IN (?) `
	query, args, err := sqlx.In(sqlStr, params.UIds)
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
