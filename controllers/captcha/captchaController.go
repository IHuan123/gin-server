package captcha

import (
	"bytes"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	rASession "reactAdminServer/rASessions"
	"time"
)

type CaptchaController struct {
	CaptchaImageController
}



func (con *CaptchaController) GetCaptcha(c *gin.Context){
	id,b64s := con.GenerateCaptcha()
	//captchaIdKey := "captchaId" + "_c_t_" + strconv.Itoa(int(time.Now().Unix()))
	//fmt.Println(captchaIdKey)
	rASession.SetSession(c,"captchaId",id)
	//go_redis.RedisClient.SetValue("captchaId", id,"5000")
	con.B64sServe(c.Writer,c.Request,struct {
		download      bool
		captchaId     string
		captchaBase64 string
		ext           string
	}{download: false, captchaId: id, captchaBase64: b64s, ext:".png" })
}
type B64ServeParams struct {
	download bool
	captchaId string
	captchaBase64 string
	ext string
}
func (con *CaptchaController) B64sServe(w http.ResponseWriter,r *http.Request,options B64ServeParams){
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	b64sStart := len("data:image/png;base64,")
	b64s := options.captchaBase64[b64sStart:]
	b64sByte,_ := base64.StdEncoding.DecodeString(b64s)
	content := bytes.NewBuffer(b64sByte)
	if options.download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, options.captchaId+options.ext, time.Time{}, bytes.NewReader(content.Bytes()))
}

