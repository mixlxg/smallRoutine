package views

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UpLoad(logger *logrus.Logger)gin.HandlerFunc  {
	return func(c *gin.Context) {
		fmt.Printf("%#v",c.Request.Header)
		//form,_:=c.MultipartForm()

		//f,_:=os.Create("test.jpg")
		//defer func() {_=f.Close()}()
		//res,_:=ioutil.ReadAll(c.Request.Body)
		//f.Write(res)
	}
}
