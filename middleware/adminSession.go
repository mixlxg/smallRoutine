package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminSessions() gin.HandlerFunc  {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// 判断客户端是否带着有效的cookie过来的，如果未带cookie或者cookie失效了，则返回未认证码状态吗给客户端
		// 获取相关数据
		username := session.Get("username")
		password := session.Get("password")
		role := session.Get("role")

		if username == nil || password == nil || role == nil{
			c.AbortWithStatusJSON(http.StatusOK,gin.H{
				"code": http.StatusUnauthorized,
				"errorMg":"用户未认证或会话过时重新登录",
			})
			return
		} else if role !="admin" {
			c.AbortWithStatusJSON(http.StatusOK,gin.H{
				"code":http.StatusForbidden,
				"errorMg":"当前用户无权限访问此接口",
			})
			return
		}else {
			c.Next()
		}
	}
}
