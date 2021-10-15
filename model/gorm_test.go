package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"smallRoutine/config"
	"smallRoutine/utils"
	"testing"
	"time"
)

func TestGorm(t *testing.T)  {
	current,err := os.Getwd()
	if err != nil{
		t.Errorf("current path失败错误信息：%v", err)
	}
	basePath := filepath.Dir(current)
	configPath := filepath.Join(basePath,"application.yml")
	t.Logf(configPath)
	config,err := config.NewConfig(configPath)
	if err != nil{
		t.Errorf("初始化解析配置文件失败，错误信息：%v", err)
	}
	gdb,err:=utils.NewGorm(config)
	if err !=nil{
		t.Errorf("初始化数据库%v报错，错误信息:%v",config.Mysql.MysqlUrl,err)
	}
	t.Logf("初始化数据库%v成功",gdb)
	user := User{
		UserName:   "test1",
		Password:   "E9C98C7556481D33ABDEB8BC39907402",
	}
	err = gdb.Where(&user).First(&user).Error
	if errors.Is(err,gorm.ErrRecordNotFound){
		t.Log("user不存在")
	}
	if user.LoginTime == nil{
		fmt.Printf("用户为首次登录")
		currentTime := time.Now()
		user.LoginTime=&currentTime
		gdb.Select("LoginTime").Updates(user)
	}
	//创建一个测试用户
	//user := User{
	//	UserName:   "test2",
	//	Password:   "E9C98C7556481D33ABDEB8BC39907402",
	//	Company:    "江苏联通",
	//}
	//err = gdb.Create(&user).Error
	//if err !=nil{
	//	t.Logf(err.Error())
	//}
	//fmt.Printf("user结果：%v",user)
	//gdb.Create(&Role{
	//	RoleName: "generalUser",
	//})
	//user := User{
	//	UserName:   "小吕",
	//	Password:   "gangzi2010",
	//	Email:      "1183566623@qq.com",
	//	PhoneNum:   "15651602025",
	//}
	//tx:=gdb.Create(&user)
	//if tx.Error != nil{
	//	t.Logf("入库失败，错误信息：%s", err.Error())
	//}else {
	//	fmt.Printf("影响业务条数：%d",tx.RowsAffected)
	//	fmt.Printf("user:%#v",user)
	//}

	//var p model.Permission=model.Permission{UrlPath:"/goApp/app/auth"}
	//tx:=gdb.Debug().Find(&p,"url_path=?",p.UrlPath)
	//fmt.Println(tx.RowsAffected)
	//fmt.Println(p)
	//var roles []*model.Role
	//gdb.Debug().Model(&model.User{ID:1}).Association("Role").Find(&roles)
	//var u model.User
	//gdb.Debug().Model(&u).Preload("Role").Preload("Role.Permission").Find(&u,"user_name=?","吕秀刚")
	//for _,role := range u.Role{
	//	for _,vp := range role.Permission{
	//		fmt.Printf("角色名:%s,权限:%s",role.RoleName,vp.UrlPath)
	//	}
	//}
	//for k,v:=range roles{
	//	fmt.Printf("index:%d  value:%#v\n",k,v)
	//}
	//for i:=0;i<len(roles) ;i++  {
	//	fmt.Printf("%#v",roles[i])
	//}
	//count:=gdb.Debug().Model(&model.User{ID:1}).Association("Role").Count()
	//fmt.Println(count)
	//var user model.User
	//var res map[string]interface{}= map[string]interface{}{}
	//gdb.Debug().Find(&user)
	//err =gdb.Model(&user).Association("Role").Find(res)
	//if err != nil{
	//	t.Log(err)
	//}
	//fmt.Println(res)
	//gdb.First()
	//result := map[string]interface{}{}
	//gdb.Table("users").Take(result)
	//fmt.Println(result)
	//result := map[string]interface{}{}
	//gdb.Model(&model.User{}).First(result)
	//fmt.Println(result)
	//var user *model.User = new(model.User)
	//gdb.First(user)
	//fmt.Println(user.UserName)
	//var user model.User=model.User{
	//	UserName:   "吕秀刚",
	//	Password:   MyMd5("gangzi2010"),
	//	Email:      "1183566623@qq.com",
	//	PhoneNum:   "15651602025",
	//	IsActive:   false,
	//	CreateTime: time.Time{},
	//	LoginTime:  time.Time{},
	//	Role:       []*model.Role{&model.Role{RoleName:"admin"},},
	//}

	//Permission := []*model.Permission{
	//		&model.Permission{
	//			UrlPath: "/goApp/app/auth",
	//		},
	//	}
	//var role model.Role
	//gdb.Model(&role).Where("role_name=?","admin").Find(&role)
	//err = gdb.Debug().Model(&role).Association("Permission").Append(Permission)
	//fmt.Println(err)
	//var user model.User
	//tx :=gdb.Debug().Where("user_name=?","吕秀刚").Preload("Role.Permission").Preload(clause.Associations).Find(&user)
	//fmt.Println(tx.RowsAffected)
	//for _,role := range user.Role{
	//	fmt.Printf("%#v",role.Permission)
	//}
	//var Roles []*model.Role
	//err = gdb.Debug().Table("user").Where("user_name=?","吕秀刚").Association("Role").Find(&Roles)
	//if err !=nil{
	//	fmt.Println(err)
	//}
	//fmt.Printf("%#v",Roles)
	//gdb.Debug().Create(&user)
	//tx :=gdb.Debug().Where("user_name=?","吕秀刚").Where("password=?",MyMd5("gangzi2010")).Find(&user)
	//fmt.Println(user)
	//fmt.Println(tx.RowsAffected)
	//fmt.Println(tx.Error)
	//var role model.Role
	//var user model.User = model.User{ID:4}
	//var roles []model.Role
	//gdb.Debug().Model(&roles).Preload("User").Find(&roles)
	//fmt.Printf("%#v",roles)
	//gdb.Debug().Where("name=?","小强")
	//fmt.Printf("%#v",user)
	//gdb.Debug().Where("user_name=?","吕秀刚").Assign("password","gangzi").FirstOrInit(&user)
	//fmt.Println(user)
	//gdb.Debug().Where("role_name=?","dev").Find(&role)
	//gdb.Debug().Where("user_name=?","小吕").Find(&user)
	//user.Role=append(user.Role,role)
	//user.ID=7
	//gdb.Debug().Model(&user).Omit("Role").Association("Role").Append(&role)
	//tx := gdb.Session(&gorm.Session{
	//	DryRun:true,
	//})
	// 查询role对象
	//var roles []model.Role
	//gdb.Debug().Find(&roles)
	//var users []*model.User
	//for i:=1001;i<=10000;i++{
	//	var user=model.User{
	//		UserName:   fmt.Sprintf("小吕%d",i),
	//		Password:   "gangzi2010",
	//		Email:      "1183566623@qq.com",
	//		PhoneNum:   "15651602025",
	//		CreateTime: time.Now(),
	//		LoginTime:  time.Now(),
	//		Role:roles,
	//	}
	//	gdb.Debug().Model(&user).Create(&user)
	//}
	//tx.Debug().CreateInBatches(&users,10)
	//fmt.Printf("%v",user.ID)
	//if err=gdb.AutoMigrate(&model.User{},&model.Role{});err !=nil{
	//	t.Fatalf("创建表失败错误信息：%v",err)
	//}
	//t.Logf("创建表User表成功")

	//插入测试数据
	//role :=model.Role{ID:2}
	//gdb.Debug().Model(&model.User{}).Where("user_name=?","仲新玲").Association("Role").Find(&role)
	//gdb.Debug().Model(&model.User{UserName:"仲新玲"}).Association("Role").Append(&model.Role{RoleName:"test"})
	//user:=model.User{}
	//role:=model.Role{}
	//gdb.Debug().Model(&user).Where("user_name=?","仲新玲").First(&user)
	//gdb.Debug().Model(&role).Where("role_name=?","test").First(&role)
	//gdb.Debug().Model(&user).Association("Role").Delete(&role)
	//gdb.Debug().Model(&model.User{ID:4}).Association("Role").Clear()
	//var users []model.User
	//var roles []model.Role
	//gdb.Debug().Model(&model.User{}).Find(&users)
	//gdb.Debug().Model(&model.Role{}).Find(&roles)
	//for _,user:= range users{
	//	fmt.Println(user)
	//	err = gdb.Debug().Model(&user).Association("Role").Append(&roles)
	//	fmt.Println(err)
	//}
	//var roles []*model.Role
	//err=gdb.Debug().Table("users").Where("user_name=?","吕秀刚").Association("ole").Find(&roles)
	//fmt.Println(roles)
	//fmt.Println(err)
	//var users []map[string]interface{}
	//gdb.Debug().Table("users").Select("users.id,users.user_name,roles.id,roles.role_name").Joins("inner join users_roles on users_roles.user_id=users.id").Joins("inner join roles on roles.id=users_roles.role_id").Find(&users)
	//fmt.Printf("%#v\n",users)
	//gdb.Debug().Find(&users)
	//fmt.Printf("%#v\n",users)
	//for _,user := range users {
	//	fmt.Print("id",user.ID,"username",user.UserName,"phonemun",user.PhoneNum,"password",user.Password,"email",user.Email,"rolename")
	//	for _,role := range user.Role{
	//		fmt.Print("roleName",role.RoleName)
	//	}
	//	fmt.Println()
	//}
}