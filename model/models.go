package model

import (
	"time"
)



type Activity struct {
	/*
	id 主键自增
	ActivityName:活动名称
	ActivityContent:活动内容描述
	CreateTime:活动创建时间
	EndTime: 活动开始时间
	EndTime:活动结束时间
	*/
	ID uint64
	ActivityName string			`gorm:"type:varchar(256) not null"`
	ActivityContent string		`gorm:"type:varchar(256)"`
	CreateTime *time.Time		`gorm:"autoCreateTime"`
	StartTime *time.Time
	EndTime	*time.Time
	Groups	[]*Group
	Orders	[]*Order
}

type Group struct {
	/*
		ID 主键，自增id
	    GroupName: 组名称，战队名称
	    CreateTime: 组创建时间
	    ActivityID： 活动id,外键
	*/
	ID uint64
	GroupName string			`gorm:"type:varchar(256) not null"`
	CreateTime *time.Time		`gorm:"autoCreateTime"`
	ActivityID	uint64
	Activity	Activity
	Users []*User				`gorm:"many2many:users_groups"`
}
type Role struct {
	/*
		id: 自增id
	    RoleName: 角色名称，id=1 是admin用户id=2 是普通用户，默认用户都为普通用户
	*/
	ID uint64
	RoleName	string	`gorm:"type:varchar(20);not null"`
	User []*User
}
type User struct {
	/*
	ID id号 自增
	Openid: 小程序的openid 可以为空，绑定小程序后update
	UserName: 用户名，不能为空，用户登录的唯一标识，在数据库里面存储唯一不重复
	Password： 用户密码 md5加密存储
	Wphone： 微信绑定手机号码，可以为空
	Phone： 用户自己提供的手机号码
	WxName: 微信名称
	Company：公司名称
	Department：部门名称
	CreateTime: 创建时间
	LoginTime: 登录时间
	RoleID： 角色id,默认id=2 即普通用户
	*/
	ID uint64
	Openid string			`gorm:"type:varchar(100)"`
	UserName string 		`gorm:"type:varchar(100) not null;unique"`
	Password string			`gorm:"type:varchar(50) not null"`
	WPhone string			`gorm:"type:varchar(12)"`
	Phone  string			`gorm:"type:varchar(12)"`
	WxName string 			`gorm:"type:varchar(50)"`
	Company string			`gorm:"type:varchar(100) not null"`
	Department string		`gorm:"type:varchar(100)"`
	CreateTime *time.Time	`gorm:"autoCreateTime"`
	LoginTime	*time.Time
	Groups [] *Group		`gorm:"many2many:users_groups"`
	Orders [] *Order
	RoleID uint64			`gorm:"type:int(10);default:2"`
	Role Role
}
type Order struct {
	/*
		id: 主键，自增id
		Customer： 客户，可以是个人客户或者公司
	    CustomerPhone: 客户的手机号码
		CustomerContent: 达成业务内容
	    OrderTimeLimit： 合同周期年为单位
	    OrderMoney： 合同金额，可以为空
	    OrderPicUrl： 用于存放上传合同副本的url
		IsAgree: 是否同意订单
		AgreeName:审批人
	    Reason: 审批拒绝原因
	    OrderCompleteTime： 签约时间
	    CreateTime: 创建订单时间
	    UserID： 用户id
		ActivityID:活动id
	*/
	ID uint64
	Customer string			`gorm:"type:varchar(256) not null"`
	CustomerPhone string	`gorm:"type:varchar(12)"`
	CustomerContent string	`gorm:"type:varchar(256) not null"`
	OrderTimeLimit uint64	`gorm:"type:int(3)"`
	OrderMoney uint64		`gorm:"type:float(21,2)"`
	OrderPicUrl string		`gorm:"type:varchar(512)"`
	IsAgree	bool			`gorm:"type:bool;default:false"`
	AgreeName string 		`gorm:"type:varchar(50) not null"`
	Reason string			`gorm:"type:varchar(512)"`
	OrderCompleteTime time.Time
	CreateTime *time.Time	`gorm:"autoCreateTime"`
	UserID uint64
	User User
	ActivityID uint64
	Activity Activity
}

