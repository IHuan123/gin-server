package router

import (
	"github.com/gin-gonic/gin"
	"reactAdminServer/controllers/captcha"
)

func InitCaptcha(r *gin.Engine){
	//生成验证码
	cCon := captcha.CaptchaController{}
	captchaGroup:=r.Group("/captcha")
	{
		captchaGroup.GET("/code",cCon.GetCaptcha)
	}
}
