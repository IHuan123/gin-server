package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"reactAdminServer/databases"
	_ "reactAdminServer/databases"
	"reactAdminServer/middlewares"
	rASession "reactAdminServer/rASessions"
	"reactAdminServer/router"
)

func main(){
	//连接数据库
	databases.InitDB()
	defer databases.CloseMysql()
	fmt.Println("mysql connect success")
	r := gin.Default() //原理也是调用的gin.New()
	//设置中间件
	r.Use(Middlewares.Cors(),rASession.Session("react_admin_session"),Middlewares.RegisterValidator())
	//设置静态资源目录
	pathName, _ := os.Getwd()
	statusPath := pathName + "/static"
	r.StaticFS("/status", http.Dir(statusPath))
	router.InitLoginRouter(r)
	router.InitCaptcha(r)
	router.InitSystemRouter(r)
	if err := r.Run(":9000");err != nil {
		panic(err)
	}
}
