package main

import (
	"fmt"
	"smallRoutine/router"
)
func main() {
	err = router.NewRouter(conf,logger,gdb,store,cstore,wSdk)
	if err != nil{
		fmt.Println(err)
	}
}
