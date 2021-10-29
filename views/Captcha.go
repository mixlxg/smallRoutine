package views

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetCaptcha(logger *logrus.Logger,store base64Captcha.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		captchaType := c.DefaultQuery("captchaType","digit")
		// 定义一个driver变量
		var driver base64Captcha.Driver
		// 根据传递的参数获取验证码类型的driver
		switch captchaType {
		case "string":
			driver = (&base64Captcha.DriverString{
				Height:          60,
				Width:           240,
				NoiseCount:      0,
				ShowLineOptions: 4,
				Length:          4,
				Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",

			}).ConvertFonts()
		case "math":
			driver = (&base64Captcha.DriverMath{
				Height:          60,
				Width:           240,
				NoiseCount:      0,
				ShowLineOptions: 4,
			}).ConvertFonts()
		default:
			driver = &base64Captcha.DriverDigit{
				Height:   60,
				Width:    240,
				Length:   4,
				MaxSkew:  0.9,
				DotCount: 90,
			}
		}
		// 实例化一个图形验证码对象
		cap:= base64Captcha.NewCaptcha(driver,store)
		id, b64s,err := cap.Generate()
		if err != nil{
			logger.Errorf("生成图形验证码失败，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg":err.Error(),
			})
			return
		}
		// 返回正常的数据给前端
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
			"captchaId":id,
			"img":b64s,
		})
		return
	}
}