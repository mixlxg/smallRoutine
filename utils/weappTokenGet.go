package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"smallRoutine/config"
)
type WeixinAccessTokenResponse struct {
	AccessToken string		`json:"access_token"`
	ExpiresIn int			`json:"expires_in"`
	ErrCode  int			`json:"errcode"`
	ErrMsg	string			`json:"errmsg"`
}

func GetToken(config *config.Config,logger *logrus.Logger) func()(token string, expireIn uint) {
	return func() (token string, expireIn uint) {
		baseurl := "https://api.weixin.qq.com/cgi-bin/token"
		urlobj,err := url.Parse(baseurl)
		if err != nil{
			logger.Errorf("解析url：%s失败，错误信息：%s",baseurl,err.Error())
			return "",0
		}
		q := urlobj.Query()
		q.Set("appid",config.WeixinConfig.AppId)
		q.Set("secret",config.WeixinConfig.Secret)
		q.Set("grant_type","client_credential")
		urlobj.RawQuery=q.Encode()
		apiurl := urlobj.String()
		rep,err:=http.Get(apiurl)
		if err !=nil{
			logger.Errorf("请求地址：%s失败，错误信息：%s",apiurl,err.Error())
			return "",0
		}
		// 关闭respose 请求体
		defer func() {_=rep.Body.Close()}()
		// 反序列化json响应包体
		var result *WeixinAccessTokenResponse = new(WeixinAccessTokenResponse)
		err = json.NewDecoder(rep.Body).Decode(result)
		if err !=nil{
			logger.Errorf("从微信小程序获取access_token失败，错误信息：%s", err.Error())
			return "",0
		}
		if result.ErrCode != 0 {
			logger.Errorf("从微信小程序获取access_token失败错误，错误码为：%v,错误信息：%s",result.ErrCode,result.ErrMsg)
			return "",0
		}
		return result.AccessToken,config.WeixinConfig.ExpiresIn
	}
}
