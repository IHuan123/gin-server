package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct{}


func (con *BaseController) Success(c *gin.Context,data interface{}){
	c.JSON(http.StatusOK,gin.H{
		"code":http.StatusOK,
		"data":data,
	})
}


func (con *BaseController) Err(c *gin.Context,msg string){
	c.JSON(http.StatusOK,gin.H{
		"code":http.StatusBadRequest,
		"msg":msg,
		"data":nil,
	})
}

func (con *BaseController) Unauthorized(c *gin.Context,msg string){
	c.JSON(http.StatusOK,gin.H{
		"code":http.StatusUnauthorized,
		"msg":msg,
		"data":nil,
	})
}