package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/v3"
	"github.com/sirupsen/logrus"
	"net/http"
)

func WeixinCodeJudge(c *gin.Context, res *weapp.LoginResponse, logger *logrus.Logger)  {
	if res.ErrCode !=0 {
		logger.Errorf("微信返回errorcode：%d,errmsg:%s",res.ErrCode,res.ErrMSG)
		c.AbortWithStatusJSON(http.StatusServiceUnavailable,gin.H{
			"code":http.StatusServiceUnavailable,
			"errCode": res.ErrCode,
			"errMsg": res.ErrMSG,
		})
	}
}
