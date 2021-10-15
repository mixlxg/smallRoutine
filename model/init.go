package model

import (
	"fmt"
	"gorm.io/gorm"
)

func Init(gdb *gorm.DB)  {
	// 开始初始定义的model
	// 创建用户相关用户权限相关表
	err:=gdb.AutoMigrate(&User{},&Role{},&Group{},&Order{},&Activity{})
	if err !=nil{
		panic(fmt.Sprintf("初始化用户，权限相关表失败，错误信息:%#v",err))
	}
}
