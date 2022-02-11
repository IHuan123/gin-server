package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"reactAdminServer/databases"
	_ "reactAdminServer/databases"
	"reactAdminServer/middlewares"
	"reactAdminServer/models"
	rASession "reactAdminServer/rASessions"
	"reactAdminServer/router"
)

func main() {
	//日志 -----------》
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色
	//gin.DisableConsoleColor()
	isExistsLog, _ := models.ExistsDir("/log")
	if !isExistsLog {
		err := models.CreateDir("/log")
		if err != nil {
			panic(err)
		}
	}
	// 记录日志到文件
	f, _ := os.Create("log/gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	// 记录错误日志到文件，同时输出到控制台
	fErr, _ := os.Create("log/gin_err.log")
	gin.DefaultErrorWriter = io.MultiWriter(fErr, os.Stdout)

	//连接数据库 -----------》
	databases.InitDB()
	defer databases.CloseMysql()
	fmt.Println("mysql connect success")
	r := gin.Default() //原理也是调用的gin.New()

	//设置中间件 -----------》
	r.Use(Middlewares.Cors(), rASession.Session("react_admin_session"), Middlewares.RegisterValidator())

	//设置静态资源目录 -----------》
	pathName, _ := os.Getwd()
	statusPath := pathName + "/static"
	r.StaticFS("/static", http.Dir(statusPath))
	//初始化router -----------》
	router.InitLoginRouter(r)
	router.InitCaptcha(r)
	router.InitSystemRouter(r)
	if err := r.Run(":9000"); err != nil {
		panic(err)
	}
	fmt.Println("Service started successfully")
}
