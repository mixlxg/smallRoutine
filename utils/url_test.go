package utils

import (
	"fmt"
	url2 "net/url"
	"testing"
)

func TestUrl(t *testing.T)  {
	url,err := url2.Parse("https://api.weixin.qq.com/cgi-bin/token")
	if err !=nil{
		fmt.Println(err)
	}
	q := url.Query()
	q.Set("appid","test")

	url.RawQuery = q.Encode()
	fmt.Println(q.Encode())
	fmt.Println(url.String())
}
