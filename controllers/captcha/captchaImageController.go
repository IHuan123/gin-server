package captcha

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
)
var Store = base64Captcha.DefaultMemStore
type CaptchaImageController struct {}
//生成验证码base64
func (con *CaptchaImageController) GenerateCaptcha()(string, string){
	var driver base64Captcha.Driver
	driverString := base64Captcha.DriverString{
		Height:          80, //图片高
		Width:           230, //图片宽
		NoiseCount:      0, //噪点数
		ShowLineOptions: 2 | 4, //设置线条数量
		Length:          4, //验证码长度
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm", //数据源
		//BgColor:         &color.RGBA{R: 3, G: 102, B: 214, A: 125}, //背景颜色
		Fonts:           []string{"wqy-microhei.ttc"}, //字体
	}
	//ConvertFonts 按名称加载字体
	driver = driverString.ConvertFonts()
	//生成图片base64
	captcha := base64Captcha.NewCaptcha(driver,Store)
	id, b64s, err := captcha.Generate()
	if err != nil {
		fmt.Println("Register GetCaptchaPhoto get base64Captcha has err:", err)
		return "", ""
	}
	return id, b64s
}
//校验验证码
func (con *CaptchaImageController) Verify(id string,val string) bool {
	if id == "" || val == "" {
		return false
	}
	// id,b64s是空 也会返回true 需要在加判断
	// 同时在内存清理掉这个图片
	return Store.Verify(id, val, true)
}
