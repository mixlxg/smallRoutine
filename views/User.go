package views

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/v3"
	"github.com/sirupsen/logrus"
	"net/http"
	"smallRoutine/config"
	"smallRoutine/utils"
)

func Login(logger *logrus.Logger, config *config.Config, wsdk *weapp.Client) gin.HandlerFunc  {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == ""{
			logger.Error("code 不能为空字符串")
			c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{
				"code":400,
				"errMsg":"code不能为空值",
			})
		}
		// 请求微信登录接口获取appid 和session_key,unionid
		res, err := wsdk.Login(code)
		if err != nil{
			logger.Errorf("访问微信小程序登录接口失败，code:%s，错误信息：%s",code,err.Error())
			c.AbortWithStatusJSON(http.StatusServiceUnavailable,gin.H{
				"code":http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
		}
		// 做返回值状态码判断
		utils.WeixinCodeJudge(c,res,logger)
		openid := res.OpenID
		session_key := res.SessionKey
		unionid := res.UnionID
		session := sessions.Default(c)
		// 设置session的有效期
		option := sessions.Options{MaxAge:config.Session.MaxAge}
		session.Options(option)
		// 设置session值
		session.Set("openid",openid)
		session.Set("session_key",session_key)
		session.Set("unionid",unionid)
		if err :=session.Save();err !=nil{
			logger.Errorf("保存session失败，错误信息：%s",err.Error())
			c.JSON(http.StatusUnauthorized,gin.H{
				"code":http.StatusUnauthorized,
				"errMsg": err.Error(),
			})
		}else {
			logger.Debug("保存session成功")
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusOK,
			})
		}
	}
}

func Logout(logger *logrus.Logger) gin.HandlerFunc  {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		if err:=session.Save();err !=nil{
			logger.Errorf("用户登出失败，错信息：%s", err.Error())
			c.JSON(http.StatusServiceUnavailable,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
		}else {
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusOK,
			})
		}

	}
}

func SessionKey(logger *logrus.Logger)gin.HandlerFunc  {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		openid := session.Get("openid")
		session_key := session.Get("session_key")
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusOK,
			"openid": openid,
			"session_key": session_key,
		})
	}
}
