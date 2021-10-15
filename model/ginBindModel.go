package model
// 修改密码接口，定义binding model
type MUser struct {
	Username	string	`json:"username" binding:"required"`
	OPassword	string	`json:"oldpwd" binding:"required"`
	NPassword	string	`json:"newpwd" binding:"required"`
	WPhone	string	`json:"wphone" binding:"-"`
}
