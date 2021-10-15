package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func SmallRoutineSessions(logger *logrus.Logger) gin.HandlerFunc  {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// 判断客户端是否带着有效的cookie过来的，如果未带cookie或者cookie失效了，则返回未认证码状态吗给客户端
		// 获取相关数据
		username := session.Get("username")
		password := session.Get("password")

		if username == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
				"code": http.StatusUnauthorized,
				"errorMg":"用户未认证或会话过时重新登录",
			})
			return
		}else {
			if password == nil{
				logger.Infof("用户%s为首次登录需要修改密码",username)
				c.AbortWithStatusJSON(http.StatusOK,gin.H{
					"code":600,
					"username":username,
				})
				return
			}else {
				c.Next()
			}
		}

	}
}