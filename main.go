package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"reactAdminServer/databases"
	_ "reactAdminServer/databases"
	"reactAdminServer/go_redis"
	"reactAdminServer/middlewares"
	"reactAdminServer/models"
	rASession "reactAdminServer/rASessions"
	"reactAdminServer/router"
)
const (
	serverHost = "127.0.0.1:6379"
	password   = "123456"
)
func main() {
	err := go_redis.InitRedis(serverHost, password, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 关闭redis客户端链接
	defer func() {
		err := go_redis.RedisClose()
		if err != nil {
			panic(err)
		}
	}()
	go_redis.RedisClient.SetValue("test",123456)
	fmt.Println(go_redis.RedisClient)
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
	defer databases.MysqlClose()
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
	router.InitTestRouter(r)
	if err := r.Run(":9000"); err != nil {
		panic(err)
	}
	fmt.Println("Service started successfully")
}
