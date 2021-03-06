package views

import (
	"encoding/json"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/v3"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"net/http"
	"smallRoutine/config"
	"smallRoutine/model"
	"smallRoutine/utils"
	"strconv"
	"time"
)

func Login(logger *logrus.Logger, config *config.Config, wsdk *weapp.Client,gdb *gorm.DB, store base64Captcha.Store) gin.HandlerFunc  {
	return func(c *gin.Context) {
		captchaId := c.Query("captchaId")
		captchaValue := c.Query("captchaValue")
		if captchaValue =="" || captchaId ==""{
			logger.Errorf("captchaId,captchaValue不能为空")
			c.JSON(http.StatusOK,gin.H{
				"code": 601,
			})
			return
		}
		// 校验图形验证码
		if !store.Verify(captchaId,captchaValue,true){
			c.JSON(http.StatusOK,gin.H{
				"code":602,
			})
			return
		}
		code := c.Query("wcode")
		username := c.Query("username")
		password := c.Query("password")
		// username,password,code是否为空校验
		if username == "" || password == "" || code == "" {
			logger.Errorf("username,password,code不能存在为空情况")
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusBadRequest,
			})
			return
		}
		// 校验正好密码是否正确
		mpwd := utils.MyMd5(password)
		var user = model.User{
			UserName:   username,
			Password:   mpwd,
		}
		err := gdb.Where(user).First(&user).Error
		if errors.Is(err,gorm.ErrRecordNotFound){
			logger.Errorf("username:%s,password:%s用户名密码不正确",username,password)
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusUnauthorized,
			})
			return
		}
		if err !=nil{
			logger.Errorf("在数据库校验username:%s,password:%s是报错，错误信息：%s",username,password,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusServiceUnavailable,
				"errMsg":err.Error(),
			})
			return
		}

		// 用户密码是正确的，判断是否为首次登录，如果用户首次登录返回600和username让用户重置密码
		//如果用户首次登录通过code查询用户的openid
		if user.LoginTime == nil{
			// 用户是首次登录
			rep, err:=wsdk.Login(code)
			if err != nil{
				logger.Errorf("访问微信小程序登录接口失败，username:%s,password:%s,code:%s，错误信息：%s",username,password,code,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code":http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
			// 做返回值状态码判断
			utils.WeixinCodeJudge(c,rep,logger)
			openid := rep.OpenID
			user.Openid=openid
			// 将openid和用户信息关联并记录到数据库中
			err = gdb.Model(&user).Select("Openid").Updates(user).Error
			if err != nil{
				logger.Errorf("更新用户%s的openid:%s失败，错误信息：%s",username,openid,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
			// 写session 因为是首次登录所以只写入username，当修改完成密码之后写入也写入password到session中
			session := sessions.Default(c)
			// 设置session的有效期
			option := sessions.Options{MaxAge:config.Session.MaxAge}
			session.Options(option)
			session.Set("username",username)
			if err = session.Save();err !=nil{
				logger.Errorf("保存用户：%s 的session失败，错误信息：%s",username,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code":http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
			// session写入成功，返回状态码前端告知用户需要修改密码
			c.JSON(http.StatusOK,gin.H{
				"code":600,
				"username":username,
			})
			return
		}else {
			// 用户为非首次登录,更新用户登录时间
			currentTime := time.Now()
			user.LoginTime = &currentTime
			err = gdb.Model(&user).Select("LoginTime").Updates(user).Error
			if err !=nil{
				logger.Errorf("用户%s登录时跟新LoginTime失败，错误信息：%s",username,err.Error())
			}
			session := sessions.Default(c)
			// 设置session的有效期
			option := sessions.Options{MaxAge:config.Session.MaxAge}
			session.Options(option)
			// 设置session值
			session.Set("username",username)
			session.Set("password",password)
			if err = session.Save();err !=nil{
				logger.Errorf("保存用户%s 的 session失败，错误信息：%s",username,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code":http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}else {
				logger.Debug("保存用户%s的session成功",username)
				c.JSON(http.StatusOK,gin.H{
					"code":http.StatusOK,
				})
				return
			}
		}
	}
}

func Logout(logger *logrus.Logger) gin.HandlerFunc  {
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

func ModifyPwd(logger *logrus.Logger,gdb *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		// 定义binding struct
		var mUser model.MUser
		if err := c.ShouldBindJSON(&mUser);err !=nil{
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
				"errMsg": err.Error(),
			})
			return
		}
		// 判断session信息是否为session中用户修改自己的用户账号密码，如果不是说明可能有当前用户修改别人密码情况，这种情况时不允许的
		session := sessions.Default(c)
		suser := session.Get("username")
		if suser != mUser.Username{
			logger.Errorf("用户:%s 非法修改用户：%s的密码",suser,mUser.Username)
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": "非法修改密码",
			})
			return
		}

		// 校验老的用户名密码是否正确，如果正确开始更新密码并更新手机号码
		var user model.User
		err := gdb.Where("user_name=? and password=?",mUser.Username,utils.MyMd5(mUser.OPassword)).First(&user).Error
		if errors.Is(err,gorm.ErrRecordNotFound){
			logger.Errorf("用户名：%s 或者密码不正确，请重新输入用户密码",mUser.Username)
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusUnauthorized,
				"errMsg": "用户或密码不正确",
			})
			return
		}
		// 修改用户密码
		newPwd := utils.MyMd5(mUser.NPassword)
		user.Password = newPwd
		if mUser.WPhone !="" {
			user.WPhone = mUser.WPhone
		}
		currentTime := time.Now()
		user.LoginTime = &currentTime
		if user.LoginTime == nil{
			// 用户首次登录修改密码
			err = gdb.Model(&user).Updates(user).Error
			if err != nil{
				logger.Errorf("修改用户%s密码失败，错误信息：%s",user.UserName,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}else {
				// 获取该用户的session
				session.Set("password",utils.MyMd5(mUser.NPassword))
				if err = session.Save(); err !=nil{
					logger.Errorf("保存新用户：%s, session的password信息失败，错误信息：%s",user.UserName,err.Error())
					c.JSON(http.StatusOK,gin.H{
						"code": http.StatusOK,
					})
					return
				}else {
					logger.Infof("新用户%s修改密码成功", user.UserName)
					c.JSON(http.StatusOK,gin.H{
						"code": http.StatusOK,
					})
					return
				}
			}

		}else {
			err = gdb.Model(&user).Updates(user).Error
			if err != nil{
				logger.Errorf("修改用户%s密码失败，错误信息：%s",user.UserName,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}else {
				logger.Infof("修改用户:%s 密码成功",user.UserName)
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
				})
				return
			}
		}
	}
}

func GetUserMess(logger *logrus.Logger, gdb *gorm.DB)gin.HandlerFunc  {
	return func(c *gin.Context) {
		// 获取session
		session := sessions.Default(c)
		username := session.Get("username").(string)
		var user = model.User{UserName: username}
		if err:=gdb.Where(user).Preload("Role").First(&user).Error;err !=nil{
			logger.Errorf("用户：%s 查询战队活动信息失败，错误信息：%s",username,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
			"data": map[string]interface{}{
				"id": user.ID,
				"UserName":user.UserName,
				"Company":user.Company,
				"Department":user.Department,
				"CreateTime": utils.TimeToStamp(user.CreateTime),
				"Role": user.Role.RoleName,
			},
		})
		return
	}
}

func GetSelectActivityList(logger *logrus.Logger,gdb *gorm.DB)gin.HandlerFunc  {
	return func(c *gin.Context) {
		// 获取session
		session := sessions.Default(c)
		username := session.Get("username").(string)
		var err error
		var user = model.User{UserName: username}
		if err:=gdb.Where(user).Preload("Role").First(&user).Error;err !=nil{
			logger.Errorf("用户：%s 查询战队活动信息失败，错误信息：%s",username,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		activityId := c.Query("ActivityId")
		var acId uint64
		if activityId == ""{
			acId = 0
		}else {
			id,err := strconv.Atoi(activityId)
			if err !=nil{
				logger.Errorf("activitid:%s是一个不正确的传参，错误信息：%s",activityId,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}
			acId = uint64(id)
		}
		// 判断用户权限

		var data = make([]interface{},0)
		if user.Role.RoleName == "admin"{
			var activities [] model.Activity
			err = gdb.Where(model.Activity{ID: acId}).Order("end_time desc").Preload("Groups").Find(&activities).Error
			if err !=nil{
				logger.Errorf("用户：%s,查询数据库失败，错误信息：%s",username,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg":err.Error(),
				})
				return
			}
			//格式化数据
			for _,activity := range activities{
				ac := make(map[string]interface{})
				ac["activityId"]=activity.ID
				if activity.ActivityContent ==""{
					ac["ActivityName"] = activity.ActivityName
				}else {
					ac["ActivityName"]=activity.ActivityContent
				}
				if activity.EndTime.Unix() > time.Now().Unix(){
					ac["end_flag"] = false
				}else {
					ac["end_flag"] = true
				}
				var groups []map[string]interface{}
				for _,group := range activity.Groups{
					groups = append(groups, map[string]interface{}{
						"groupId": group.ID,
						"GroupName":group.GroupName,
						"GroupLeader":group.GroupLeader,
					})
				}
				ac["groups"] = groups
				data = append(data,ac)
			}
		}else {
			err = gdb.Model(&user).Preload("Groups.Activity").Preload(clause.Associations).Find(&user).Error
			if err !=nil{
				logger.Errorf("用户：%s,查询数据报错，错信息：%s",username,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg":err.Error(),
				})
				return
			}
			// 格式化数据
			if acId == 0{
				// 返回所有的活动和战队数据
				for _,group := range user.Groups{
					ac := make(map[string]interface{})
					ac["activityId"]=group.Activity.ID
					if group.Activity.ActivityContent ==""{
						ac["ActivityName"] = group.Activity.ActivityName
					}else {
						ac["ActivityName"]=group.Activity.ActivityContent
					}
					if group.Activity.EndTime.Unix() > time.Now().Unix(){
						ac["end_flag"] = false
					}else {
						ac["end_flag"] = true
					}
					var groups []map[string]interface{}
					groups = append(groups, map[string]interface{}{
						"groupId": group.ID,
						"GroupName":group.GroupName,
						"GroupLeader":group.GroupLeader,
					})
					ac["groups"] = groups
					data = append(data,ac)
				}
			}else {
				for _,group := range user.Groups{
					if group.Activity.ID == acId{
						ac := make(map[string]interface{})
						ac["activityId"]=group.Activity.ID
						if group.Activity.ActivityContent ==""{
							ac["ActivityName"] = group.Activity.ActivityName
						}else {
							ac["ActivityName"]=group.Activity.ActivityContent
						}
						if group.Activity.EndTime.Unix() > time.Now().Unix(){
							ac["end_flag"] = false
						}else {
							ac["end_flag"] = true
						}
						var groups []map[string]interface{}
						groups = append(groups, map[string]interface{}{
							"groupId": group.ID,
							"GroupName":group.GroupName,
							"GroupLeader":group.GroupLeader,
						})
						ac["groups"] = groups
						data = append(data,ac)
					}
				}
			}
		}
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
			"data": data,
		})
	}
}

func GetPageActivityMess(logger *logrus.Logger, gdb *gorm.DB)gin.HandlerFunc  {
	return func(c *gin.Context) {
		// 获取session
		session := sessions.Default(c)
		username := session.Get("username").(string)
		var err error
		var user = model.User{UserName: username}
		if err:=gdb.Where(user).Preload("Role").Preload("Groups.Activity").Preload("Groups.Users").Preload(clause.Associations).First(&user).Error;err !=nil{
			logger.Errorf("用户：%s 查询战队活动信息失败，错误信息：%s",username,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		var param model.GetPageActivityMess
		if err:=c.ShouldBindQuery(&param);err !=nil{
			logger.Errorf("获取前端传递参数报错，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		var activities []*model.Activity
		var total int64
		// 获取用户一共有多少数据
		if user.Role.RoleName == "admin"{
			err = gdb.Model(&model.Activity{}).Count(&total).Error
			if err !=nil{
				logger.Errorf("param:%#v,获取页码总数失败，错误信息：%s",param,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}else {
			total = int64(len(user.Groups))
		}
		// 总页码
		var totalPage int64
		var data = make([]interface{},0)
		switch total{
		case 0:
			totalPage = 1
		default:
			totalPage = int64(math.Ceil(float64(total)/float64(param.Size)))
		}
		if param.CurrentPage > totalPage{
			param.CurrentPage = totalPage
		}
		if param.CurrentPage <= 1{
			param.CurrentPage=1
		}
		if user.Role.RoleName == "admin"{
			err = gdb.Model(&model.Activity{}).Offset(int((param.CurrentPage-1)*param.Size)).Limit(int(param.Size)).Order("end_time desc").Preload("Groups.Users").Preload(clause.Associations).Find(&activities).Error
			if err !=nil{
				logger.Errorf("查询活动信息报错，错误信息：%s",err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
			// 格式化返回的数据
			for _,activity := range activities{
				ac := make(map[string]interface{})
				ac["id"] = activity.ID
				if activity.ActivityContent == ""{
					ac["ActivityName"] = activity.ActivityName
				}else {
					ac["ActivityName"] = activity.ActivityContent
				}
				if activity.Approver == ""{
					ac["Approver"] = []string{}
				}else {
					// 反序列化出来审批人信息
					err = json.Unmarshal([]byte(activity.Approver),ac["Approver"])
					if err !=nil{
						ac["Approver"] = []string{}
					}
				}
				ac["StartTime"] = utils.TimeToStamp(activity.StartTime)
				ac["EndTime"] = utils.TimeToStamp(activity.EndTime)
				groups := make([]map[string]interface{},0)
				for _,g :=range activity.Groups{
					var us []string
					for _,u := range g.Users{
						us= append(us,u.UserName)
					}
					groups = append(groups, map[string]interface{}{
						"GroupName":g.GroupName,
						"GroupLeader":g.GroupLeader,
						"users":us,
					})
				}
				ac["groups"] = groups
				data = append(data,ac)
			}
		}else {
			var tmp = make([]interface{},0)
			// 格式化返回的数据
			for _,gr := range user.Groups{
				ac := make(map[string]interface{})
				ac["id"] = gr.Activity.ID
				if gr.Activity.ActivityContent == ""{
					ac["ActivityName"] = gr.Activity.ActivityName
				}else {
					ac["ActivityName"] = gr.Activity.ActivityContent
				}
				if gr.Activity.Approver == ""{
					ac["Approver"] = []string{}
				}else {
					// 反序列化出来审批人信息
					err = json.Unmarshal([]byte(gr.Activity.Approver),ac["Approver"])
					if err !=nil{
						ac["Approver"] = []string{}
					}
				}
				ac["StartTime"] = utils.TimeToStamp(gr.Activity.StartTime)
				ac["EndTime"] = utils.TimeToStamp(gr.Activity.EndTime)
				groups := make([]map[string]interface{},0)
				var us []string
				for _,u := range gr.Users{
					us= append(us,u.UserName)
				}
				groups = append(groups, map[string]interface{}{
					"GroupName":gr.GroupName,
					"GroupLeader":gr.GroupLeader,
					"users":us,
				})
				ac["groups"] = groups
				tmp = append(tmp,ac)
				//开始分页
				offset := (param.CurrentPage -1)*param.Size
				if param.CurrentPage == totalPage{
					data = tmp[offset:]
				}else {
					data = tmp[offset:param.Size]
				}

			}

		}
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusOK,
			"data": data,
			"total": total,
		})
		return
	}
}