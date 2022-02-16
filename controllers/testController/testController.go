package testController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reactAdminServer/controllers/base"
	"reactAdminServer/go_redis"
)

type TestController struct {
	base.BaseController
}

func (con *TestController) V1(c *gin.Context){
	value, err := go_redis.RedisClient.GetValue("test")
	if err != nil {
		con.Err(c,err.Error())
		return
	}
	fmt.Println(value)
	con.Success(c,value)
}



