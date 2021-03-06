package user

//循环引入会导致 import cycle not allowed
import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"reactAdminServer/controllers/adminSystem"
	"reactAdminServer/controllers/base"
	"reactAdminServer/controllers/captcha"
	"reactAdminServer/databases"
	"reactAdminServer/go_redis"
	"reactAdminServer/models"
	rASession "reactAdminServer/rASessions"
	"strconv"
	"strings"
)

type UserController struct {
	base.BaseController
	captcha.CaptchaImageController
}

type LoginParams struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type Menus = adminSystem.Menus

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
				parent[i].Children = append(parent[i].Children, c)
			}
		}
	}
	return parent
}
//处理token
func storeToken(uid int,token string) error {
	//存储现在的token
	err := go_redis.RedisClient.SetValue(token, uid)
	if err != nil {
		return err
	}
	err = go_redis.RedisClient.SetValue(strconv.Itoa(uid), token)
	if err != nil {
		return err
	}
	return nil
}
func delStoreToken(uid int) error{
	//删除之前的token
	preToken, err := redis.String(go_redis.RedisClient.GetValue(strconv.Itoa(uid)))
	if err != nil {
		return err
	}
	err = go_redis.RedisClient.DelValue(strconv.Itoa(uid))
	if err != nil {
		return err
	}
	err = go_redis.RedisClient.DelValue(preToken)
	if err != nil {
		return err
	}
	return nil
}
//登陆
func (con *UserController) Login(c *gin.Context) {
	var params LoginParams
	var res databases.User
	err := c.Bind(&params)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	code := params.Code
	captchaId := rASession.GetSession(c, "captchaId")
	//captchaId ,err := redis.String(go_redis.RedisClient.GetValue("captchaId"))
	//if err!=nil{
	//	con.Err(c, err.Error())
	//	return
	//}
	isVerify := con.Verify(captchaId, code)
	if !isVerify {
		con.Err(c, "验证码不正确")
		return
	}
	querySql := `SELECT uid,username,avatar,deptId,email,enabled,phone,sex,roles,createTime FROM t_user WHERE username = ? AND password = ?`
	err = databases.DB.Get(&res, querySql, params.UserName, params.Password)
	if err != nil {
		con.Err(c, "账号或密码错误")
		return
	}
	if res.Enabled == 0 {
		con.Unauthorized(c, "该账号已禁用")
		return
	}
	//uid := res.Uid
	////删除旧的token
	//err = delStoreToken(uid)
	//if err != nil {
	//	con.Err(c, "del token:"+err.Error())
	//	return
	//}

	//生成token
	token, err := models.GenerateToken(res.Uid, res.UserName, params.Password)
	if err != nil {
		con.Err(c, err.Error())
		return
	}


	////储存新的token
	//err = storeToken(uid, token)
	//if err != nil {
	//	con.Err(c, "set token:"+err.Error())
	//	return
	//}


	data := map[string]interface{}{
		"token":    "Bearer " + token,
		"userInfo": res,
	}
	con.Success(c, data)
}

//添加用户

//登出

//通过token获取参数
func (con *UserController) GetUser(c *gin.Context) {
	var res databases.User
	tokenString := c.GetHeader("Authorization")
	_, userInfo, err := models.ParseToken(tokenString)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	uid := userInfo.Uid
	querySql := `SELECT uid,username,avatar,deptId,email,enabled,phone,sex,roles,createTime FROM t_user WHERE uid = ?`
	err = databases.DB.Get(&res, querySql, uid)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	if res.Enabled == 0 {
		con.Unauthorized(c, "该账号已禁用")
		return
	}
	data := map[string]interface{}{
		"token":    "Bearer " + tokenString,
		"userInfo": res,
	}
	con.Success(c, data)
}

//通过uid获取菜单
func (con *UserController) GetMenus(c *gin.Context) {
	roleStr := c.Query("roles")
	roles := strings.Split(roleStr, ",")
	var data []Menus
	//(select max(id) from table group by [去除重复的字段名列表,....])
	sqlStr := `SELECT distinct 
       r.menu_id,m.title,m.r_path path,m.icon,m.r_key,m.visible,m.keep_alive,m.weight,parent_key 
FROM 
     role_menu r join menus m using(menu_id) 
WHERE 
      r.role_id in (?)
ORDER BY m.weight DESC`
	//err := databases.DB.Select(&data,sqlStr,roles)

	query, args, err := sqlx.In(sqlStr, roles)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = databases.DB.Rebind(query)
	err = databases.DB.Select(&data, query, args...)
	if err != nil {
		con.Err(c, err.Error())
		return
	}
	res := handleMenus(data)
	con.Success(c, res)
}
