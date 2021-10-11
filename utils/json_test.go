package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJson(t *testing.T)  {
	type WeixinAccessTokenResponse struct {
		AccessToken string		`json:"access_token"`
		ExpiresIn int			`json:"expires_in"`
		ErrCode  int			`json:"errcode"`
		ErrMsg	string			`json:"errmsg"`
	}
	jsonstr := `{"access_token":"ACCESS_TOKEN","expires_in":7200}`
	w:=new(WeixinAccessTokenResponse)
	err:=json.Unmarshal([]byte(jsonstr),w)
	if err !=nil{
		fmt.Println(err.Error())
	}
	fmt.Printf("%#v", w)
}
