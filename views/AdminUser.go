package views

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"smallRoutine/config"
	"smallRoutine/model"
	"smallRoutine/utils"
	"time"
)

func AdminLogin(logger *logrus.Logger,gdb *gorm.DB, store base64Captcha.Store,config *config.Config) gin.HandlerFunc  {
	return func(c *gin.Context) {
		// 绑定请求参数
		var aUser model.AdminUser
		if err := c.ShouldBindQuery(&aUser);err !=nil{
			logger.Errorf("绑定接口参数失败，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		// 校验图形验证码
		if !store.Verify(aUser.CaptchaId,aUser.CaptchaValue,true){
			logger.Errorf("用户：%s登录图形验证码验证失败",aUser.Username)
			c.JSON(http.StatusOK,gin.H{
				"code":602,
			})
			return
		}
		// 图形验证码验证成功，判断用户账号密码是否正确
		var user model.User
		err := gdb.Where("user_name=? and password=?",aUser.Username, utils.MyMd5(aUser.Password)).Preload("Role").First(&user).Error
		if errors.Is(err,gorm.ErrRecordNotFound){
			logger.Errorf("用户：%s登陆管理后台账号或者密码不正确",aUser.Username)
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusUnauthorized,
			})
			return
		}else if err !=nil {
			logger.Errorf("用户：%s登录后台时查询数据库失败，错误信息：%s",aUser.Username,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusServiceUnavailable,
				"errMsg":err.Error(),
			})
			return
		}
		// 账号密码验证通过，来判断一下角色是否时admin role的用户
		if user.Role.RoleName != "admin"{
			logger.Errorf("用户：%s非admin用户不允许登录管理后台",user.UserName)
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusForbidden,
			})
			return
		}
		// 全部判断完成，更新登录事件写session
		//更新登录时间
		currentTime:=time.Now()
		user.LoginTime=&currentTime
		err = gdb.Select("LoginTime").Updates(user).Error
		if err !=nil{
			logger.Errorf("用户%s登录时跟新LoginTime失败，错误信息：%s",user.UserName,err.Error())
		}
		session := sessions.Default(c)
		// 设置session的有效期
		option := sessions.Options{MaxAge:config.Session.MaxAge}
		session.Options(option)
		session.Set("username",aUser.Username)
		session.Set("password",user.Password)
		session.Set("role",user.Role.RoleName)
		if err = session.Save();err !=nil{
			logger.Errorf("保存用户：%s 的session失败，错误信息：%s",aUser.Username,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusOK,
		})
	}
}

func AdminLogout(logger *logrus.Logger)gin.HandlerFunc  {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		session.Clear()
		if err:=session.Save();err !=nil{
			logger.Errorf("用户%s登出失败，错信息：%s", username,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
		}else {
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusOK,
			})
		}

	}
}

func QueryUser(logger *logrus.Logger,gdb *gorm.DB) gin.HandlerFunc  {
	return func(c *gin.Context) {
		var param model.QueryUser
		if err:=c.ShouldBindJSON(&param);err !=nil{
			logger.Errorf("请求接口QueryUser传参不正确，错误信息：%s",err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		// 成功获取到用户传参，开始判断请求的类型
		if param.QueryType == "role"{
			var roles []*model.Role
			if err:= gdb.Find(&roles).Error;err!=nil{
				logger.Errorf("parma：%v,查询数据库报错，错误信息：%s",param,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code":http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
			// 格式化返回数据
			var data [] string
			for _,role:= range roles{
				data = append(data,role.RoleName)
			}
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusOK,
				"data":data,
			})
			return
		}
		//  查询所有用户名
		if param.QueryType == "user"{
			var data [] *struct{
				UserName string
			}
			if err :=gdb.Model(&model.User{}).Find(&data).Error;err !=nil{
				logger.Errorf("param:%v,获取user信息失败，错误信息：%s",param,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusOK,
				"data": data,
			})
			return
		}
		// 查看所有角色和用户关系
		if param.QueryType == "all"{
			var roles [] *model.Role
			if err:= gdb.Preload("User").Find(&roles).Error;err !=nil{
				logger.Errorf("parma：%v,查询数据库报错，错误信息：%s",param,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code":http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
			var data map[string][]interface{}=make(map[string][]interface{})
			for _,role := range roles{
				for _,user := range role.User{
					data[role.RoleName] = append(data[role.RoleName], struct{
						Username string
					}{Username: user.UserName})
				}
			}
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusOK,
				"data":data,
			})
			return
		}
		// 更具角色名称查询所有的用户信息
		if param.QueryType == "detail_user_by_role"{
			if param.Role == ""{
				logger.Errorf("param:%v,根据传入角色查询用户信息时role不能为空",param)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}
			var role *model.Role
			if err:=gdb.Where(model.Role{RoleName: param.Role}).Preload("User").Find(&role).Error;err !=nil{
				logger.Errorf("根据role:%s查询用户信息失败，错误信息：%s",param.Role,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
			var data [] interface{}
			for _,user := range role.User{
				data = append(data, struct {
					UserName string
				}{UserName: user.UserName})
			}
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusOK,
				"data": data,
			})
			return
		}
		// 根据用户名查询校色信息
		if param.QueryType == "detail_role_by_user"{
			if param.Username == ""{
				logger.Errorf("param:%v,根据用户查询角色信息时，用户名不能为空",param)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}
			var user = model.User{UserName: param.Username}
			if err:= gdb.Where(user).Preload("Role").Find(&user).Error;err !=nil{
				logger.Errorf("param:%v，根据用户查询角色是报错，错误信息：%s",param,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusOK,
				"data":user.Role.RoleName,
			})
			return
		}
		// 如果不是这五种类型查询返回错误
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusBadRequest,
		})
		return
	}
}

func CreateUser(logger *logrus.Logger,gdb *gorm.DB) gin.HandlerFunc  {
	return func(c *gin.Context) {
		var param model.CreateUser
		if err := c.ShouldBindJSON(&param);err !=nil{
			logger.Errorf("获取前端传值失败，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		// 查询role获取到role id
		var role model.Role= model.Role{RoleName: param.RoleName}
		if err := gdb.Where(&role).First(&role).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("param:%v,角色：%s不存在", param,param.RoleName)
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("param：%v创建用户失败，错误信息；%s",param,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}
		// 创建用户，构造一个用户结构体用于创建用户
		var user model.User = model.User{
			UserName: param.UserName,
			Password: utils.MyMd5(param.Password),
			Phone: param.Phone,
			Company: param.Company,
			Department: param.Department,
			RoleID: role.ID,
		}
		if err:= gdb.Create(&user).Error;err !=nil{
			logger.Errorf("param:%v,创建用户失败，错误信息：%s",param,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusOK,
		})
		return
	}
}

func UpdateUser(logger *logrus.Logger, gdb *gorm.DB)gin.HandlerFunc  {
	return func(c *gin.Context) {
		var param model.UpdateUser
		if err := c.ShouldBindJSON(&param);err !=nil{
			logger.Errorf("获取前端传值失败，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		// 查询用户是否存在
		var user model.User
		if err:= gdb.Where(model.User{UserName: param.UserName}).First(&user).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("param:%v,用户不存在",param)
				c.JSON(http.StatusOK,gin.H{
					"code":http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("param:%v,查询数据库失败错误信息：%s",param,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}
		// 判断是否修改角色
		if param.RoleName != ""{
			// 查询当前传递过来角色信息
			var role model.Role
			if err:= gdb.Where(model.Role{RoleName: param.RoleName}).First(&role).Error;err !=nil{
				if errors.Is(err,gorm.ErrRecordNotFound){
					logger.Errorf("角色:%s不存在",param.RoleName)
					c.JSON(http.StatusOK,gin.H{
						"code": http.StatusBadRequest,
					})
					return
				}else {
					logger.Errorf("param：%v，查询角色信息报错，错误信息：%s", param,err.Error())
					c.JSON(http.StatusOK,gin.H{
						"code": http.StatusServiceUnavailable,
						"errMsg": err.Error(),
					})
					return
				}
			}
			user.RoleID = role.ID
		}
		// 判断是否修改密码
		if param.Password != "" {
			user.Password = utils.MyMd5(param.Password)
		}
		// 剩余选项
		user.Phone = param.Phone
		user.Company = param.Company
		user.Department = param.Department
		// 更新用户信息
		if err:= gdb.Model(&user).Updates(&user).Error;err !=nil{
			logger.Errorf("param:%v,更新用户信息失败，错误信息：%s",param,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusOK,
		})
		return
	}
}

func DelUser(logger *logrus.Logger,gdb *gorm.DB)gin.HandlerFunc  {
	return func(c *gin.Context) {
		username := c.Query("UserName")
		if username == ""{
			logger.Errorf("传参不正确，请传递用户名")
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		// 查找用户是否存在
		var user model.User
		if err:= gdb.Where(model.User{UserName: username}).Preload("Orders").First(&user).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("要删除的用户：%s不存在",username)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("删除用户信息时，查询用户:%s报错，错误信息：%s",username,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}
		// 判断用户是否已经存在订单信息
		if len(user.Orders) != 0{
			logger.Infof("用户：%s,存在订单信息，不能删除",username)
			c.JSON(http.StatusOK,gin.H{
				"code":604,
			})
			return
		}
		// 删除用户的group关联关系
		if err:= gdb.Model(&user).Association("Groups").Clear();err !=nil{
			logger.Errorf("删除用户：%s 的战队级联关系失败，错误信息：%s",username,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		// 删除用户
		if err:= gdb.Delete(&user).Error;err !=nil{
			logger.Errorf("删除用户：%s 失败，错误信息：%s",username,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusOK,
		})
		return
	}
}
