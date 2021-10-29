package model

type MUser struct {
	//修改密码接口，定义binding model
	Username	string	`json:"username" binding:"required"`
	OPassword	string	`json:"oldpwd" binding:"required"`
	NPassword	string	`json:"newpwd" binding:"required"`
	WPhone	string	`json:"wphone" binding:"-"`
}


type AdminUser struct {
	// admin login 接口binding model
	Username 	string		`form:"username" binding:"required"`
	Password	string		`form:"password" binding:"required"`
	CaptchaId	string		`form:"captchaId" binding:"required"`
	CaptchaValue string		`form:"captchaValue" binding:"required"`
}

type QueryUser struct {
	QueryType string	`json:"query_type" binding:"required"`
	Username string		`json:"username" binding:"-"`
	Role string			`json:"role" binding:"-"`
}

type CreateUser struct {
	UserName string		`json:"UserName" binding:"required"`
	Password string		`json:"Password" binding:"required"`
	Phone 	string		`json:"Phone" binding:"-"`
	Company	string		`json:"Company" binding:"required"`
	Department string	`json:"Department" binding:"-"`
	RoleName	string	`json:"Role" binding:"-"`
}
type UpdateUser struct {
	UserName string		`json:"UserName" binding:"required"`
	Password string		`json:"Password" binding:"-"`
	Phone 	string		`json:"Phone" binding:"-"`
	Company	string		`json:"Company" binding:"-"`
	Department string	`json:"Department" binding:"-"`
	RoleName	string	`json:"Role" binding:"-"`
}
type CreateActivity struct {
	ActivityName string			`json:"ActivityName" binding:"required"`
	ActivityContent	string		`json:"ActivityContent" binding:"required"`
	StartTime int64		`json:"StartTime" binding:"required"`
	EndTime	int64		`json:"EndTime" binding:"required"`
}

type UpdateActivity struct {
	ActivityName string			`json:"ActivityName" binding:"required"`
	ActivityContent	string		`json:"ActivityContent" binding:"-"`
	StartTime int64		`json:"StartTime" binding:"-"`
	EndTime	int64			`json:"EndTime" binding:"-"`
}
type ParamGroup struct {
	GroupName string		`json:"GroupName" binding:"required"`
	ActivityName string		`json:"ActivityName" binding:"required"`
}

type AddUsersToGroup struct {
	GroupName string		`json:"GroupName" binding:"required"`
	ActivityName string		`json:"ActivityName" binding:"required"`
	Users []string			`json:"Users" binding:"required"`
}
type DelUsersFromGroup struct {
	GroupName string		`json:"GroupName" binding:"required"`
	ActivityName string		`json:"ActivityName" binding:"required"`
	Users []string			`json:"Users" binding:"required"`
}
