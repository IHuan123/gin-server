// session中间件
package rASession

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)
// 中间件，处理session
func Session(keyPairs string) gin.HandlerFunc {
	store := SessionConfig()
	return sessions.Sessions(keyPairs, store)
}
func SessionConfig() sessions.Store {
	var store sessions.Store
	sessionMaxAge := 3600
	sessionSecret := "REACT-ADMIN-AKNBNHYN-SECRET"
	store = cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		MaxAge: sessionMaxAge, //seconds
		Path:   "/",
	})
	return store
}

func SetSession(c *gin.Context,sessionId string, value string) string {
	// 初始化session对象
	session := sessions.Default(c)
	session.Set(sessionId, value)
	session.Save()
	return sessionId
}

func GetSession(c *gin.Context, sessionId string) (value string) {
	// 初始化session对象
	session := sessions.Default(c)
	result := session.Get(sessionId)
	if value,ok := result.(string);ok{
		return value
	}
	return ""
}
//删除
func DelSession(c *gin.Context,sessionId string){
	// 初始化session对象
	session := sessions.Default(c)
	session.Delete(sessionId)
	session.Save()
}
//清空
func ClearSession(c *gin.Context){
	// 初始化session对象
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}
