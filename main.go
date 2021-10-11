package main

import (
	"fmt"
	"smallRoutine/router"
)
func main() {
	err = router.NewRouter(conf,logger,gdb,basePath,store,wSdk)
	if err != nil{
		fmt.Println(err)
	}
}
