package views

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"smallRoutine/model"
	"smallRoutine/utils"
	"time"
)

func CreateActivity(logger *logrus.Logger,gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var param model.CreateActivity
		if err:= c.ShouldBindJSON(&param);err!=nil{
			logger.Errorf("创建活动失败，获取前端传递参数失败,错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		// 判断endTime 是否大于开始时间
		if (param.EndTime- param.StartTime) <0{
			logrus.Errorf("开始时间：%v,结束时间:%v,结束时间早于开始时间",param.StartTime,param.EndTime)
			c.JSON(http.StatusOK,gin.H{
				"code": 606,
			})
			return
		}
		// 判断活动是否存在
		var activity model.Activity
		if err := gdb.Where(model.Activity{ActivityName: param.ActivityName}).First(&activity).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				// 活动不存在创建活动
				activity.ActivityName = param.ActivityName
				activity.ActivityContent = param.ActivityContent
				activity.ActivityType = param.ActivityType
				activity.StartTime = utils.StampToTime(param.StartTime)
				activity.EndTime = utils.StampToTime(param.EndTime)
				if err=gdb.Create(&activity).Error;err != nil{
					logger.Errorf("创建活动：%#v,失败错误信息：%s",activity,err.Error())
					c.JSON(http.StatusOK,gin.H{
						"code": http.StatusServiceUnavailable,
						"errMsg": err.Error(),
					})
					return
				}
				logger.Infof("活动：%v创建成功",activity)
				c.JSON(http.StatusOK,gin.H{
					"code":http.StatusOK,
				})

			}else {
				logger.Errorf("param:%v,查询活动信息失败，错误信息：%s",param,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
			}
		}else {
			logger.Errorf("%s>>>活动已存在，不能创建同名活动",param.ActivityName)
			c.JSON(http.StatusOK,gin.H{
				"code":605,
			})
			return
		}

	}
}

func DelActivity(logger *logrus.Logger,gdb *gorm.DB) gin.HandlerFunc  {
	return func(c *gin.Context) {
		ActivityName := c.Query("ActivityName")
		if ActivityName == ""{
			logger.Errorf("删除活动时，ActivityName不能为空")
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusBadRequest,
			})
			return
		}
		// 判断活动是否存在
		var activity model.Activity
		if err := gdb.Where(model.Activity{ActivityName: ActivityName}).Preload("Groups").Preload("Orders").First(&activity).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("ActivityName:%s，活动不存在",ActivityName)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("ActivityName:%s，校验是否存在时，查询数据库报错，错误信息：%s",ActivityName,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),

				})
			}
		}
		// 活动存在，判断是否有订单存在
		if len(activity.Orders) > 0{
			logger.Errorf("ActivityName:%s，活动已经开始存在订单不能删除", activity.ActivityName)
			c.JSON(http.StatusOK,gin.H{
				"code":607,
			})
			return
		}
		if len(activity.Groups) >0{
			logger.Errorf("ActivityName:%s，活动已存在战队，需要先求改现存战队所属的活动在删除。",activity.ActivityName)
			c.JSON(http.StatusOK,gin.H{
				"code":608,
			})
			return
		}
		if err := gdb.Delete(&activity).Error;err !=nil{
			logger.Errorf("删除活动：%s,失败错误信息：%s",activity.ActivityName,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		logger.Infof("删除活动：%s,成功",activity.ActivityName)
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
		})
		return
	}
}
func UpdateActivity(logger *logrus.Logger, gdb *gorm.DB) gin.HandlerFunc  {
	return func(c *gin.Context) {
		var param model.UpdateActivity
		if err:= c.ShouldBindJSON(&param);err !=nil{
			logger.Errorf("更新活动信息绑定前端传参失败，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusBadRequest,
			})
			return
		}
		// 判断活动是否存在
		var activity model.Activity
		if err := gdb.Where(model.Activity{ActivityName: param.ActivityName}).First(&activity).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("活动：%s，不存在，传参有问题",param.ActivityName)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("param:%#v，更新活动：%s失败，错误信息：%s",param,param.ActivityName,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code":http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}
		// 活动存在更新活动
		activity.ActivityContent=param.ActivityContent
		activity.ActivityType = param.ActivityType
		// 判断活动结束时间是否早于开始时间
		if param.StartTime != 0{
			if param.EndTime != 0{
				// 判断结束时间是否晚于开始时间
				if (param.EndTime - param.StartTime) >0{
					activity.StartTime = utils.StampToTime(param.StartTime)
					activity.EndTime = utils.StampToTime(param.EndTime)
				}else {
					logger.Errorf("修改活动：%s 时间，结束时间不能早于开始时间",param.ActivityName)
					c.JSON(http.StatusOK,gin.H{
						"code":606,
					})
					return
				}
			}else {
				logger.Errorf("修改活动：%s，时间开始和结束时间必须传值不能为空",param.ActivityName)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}
		}else {
			if param.EndTime !=0{
				logger.Errorf("修改活动：%s，时间开始和结束时间必须传值不能为空",param.ActivityName)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}
		}
		// 开始修改活动
		if err:= gdb.Updates(&activity).Error;err !=nil{
			logger.Errorf("param:%#v，修改活动：%s,信息失败，错误信息：%s",param,param.ActivityName,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusServiceUnavailable,
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
func QueryActivity(logger *logrus.Logger,gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var param model.QueryActivity
		if err := c.ShouldBindJSON(&param);err !=nil{
			logger.Errorf("前端传参不正确，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		var activities [] *model.Activity
		if param.QueryType == "all"{
			if err := gdb.Preload("Groups.Users").Find(&activities).Error;err !=nil{
				logger.Errorf("param:%#v,错误信息：%s",param,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}else if param.QueryType == "one" {
			if err := gdb.Where(model.Activity{ActivityName: param.ActivityName}).Preload("Groups.Users").Find(&activities).Error;err !=nil{
				logger.Errorf("param:%#v,错误信息：%s",param,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}else {
			logger.Errorf("param:%#v,传参数不正确",param)
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		// 格式化返回的数据
		var data = make(map[string]interface{})
		for _,activity := range activities{
			// 判断当前活动是否已经结束，如果结束添加一个标志位
			var ac = make(map[string]interface{})
			if activity.EndTime.Unix() >time.Now().Unix(){
				ac["end_flag"]=false
			}else {
				ac["end_flag"]=true
			}
			//判断groups是否为空
			if activity.Groups == nil{
				ac["groups"]=make([]interface{},0)
			}
			var groups [] interface{}
			for _,group := range activity.Groups{
				g:=make(map[string]interface{})
				g["group_name"] = group.GroupName
				g["group_leader"] = group.GroupLeader
				// 判断group是否有用户
				if group.Users ==nil{
					g["users"]=make([]interface{},0)
				}else {
					var users []interface{}
					for _,user := range group.Users{
						users = append(users,map[string]interface{}{
							"id":user.ID,
							"Openid":user.Openid,
							"UserName": user.UserName,
							"WPhone": user.WPhone,
							"Phone": user.Phone,
							"WxName": user.WxName,
							"Company": user.Company,
							"Department":user.Department,
						})
					}
					g["users"]=users
				}
				groups = append(groups, g)
			}
			ac["groups"] = groups
			// 判断活动是否为空
			if activity.Approver !=""{
				var u []string
				_=json.Unmarshal([]byte(activity.Approver),&u)
				ac["approver"] = u
			}else {
				ac["approver"]=[]string{}
			}
			data[activity.ActivityName] = ac
		}

		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusOK,
			"data":data,
		})
	}
}
func QueryActivityMess(logger *logrus.Logger,gdb *gorm.DB)gin.HandlerFunc  {
	return func(c *gin.Context) {
		ActivityName := c.Query("ActivityName")
		var activities [] *model.Activity
		if ActivityName ==""{
			if err := gdb.Model(model.Activity{}).Preload("Groups.Users").Preload(clause.Associations).Find(&activities).Error;err !=nil{
				logger.Errorf("查询活动信息失败，错误信息：%s", err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code":http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}else {
			if err := gdb.Where(model.Activity{ActivityName: ActivityName}).Preload("Groups.Users").Preload(clause.Associations).Find(&activities).Error;err !=nil{
				logger.Errorf("查询活动:%s 信息失败，错误信息：%s",ActivityName,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code":http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}
		// 格式化数据返回
		data := make([]map[string]interface{},0)
		for _,activity := range activities{
			ac := make(map[string]interface{})
			ac["ActivityName"] = activity.ActivityName
			ac["ActivityContent"] = activity.ActivityContent
			ac["ActivityType"]=activity.ActivityType
			ac["StartTime"]=utils.TimeToStamp(activity.StartTime)
			ac["EndTime"]=utils.TimeToStamp(activity.EndTime)
			if activity.Approver == ""{
				ac["Approver"] = make([]string,0)
			}else {
				u:=make([]string,0)
				err := json.Unmarshal([]byte(activity.Approver),&u)
				if err !=nil{
					logger.Errorf("反序列化：%s数据失败，错误信息：%s",activity.Approver,err.Error())
					ac["Approver"]=[]string{}
				}else {
					ac["Approver"]=u
				}
			}
			groups := make([]interface{},0)
			for _,g := range activity.Groups{
				group := map[string]interface{}{}
				group["GroupName"] = g.GroupName
				group["GroupLeader"] = g.GroupLeader
				users := make([]map[string]interface{},0)
				for _,u := range g.Users{
					users = append(users, map[string]interface{}{
						"UserName":u.UserName,
						"Openid":u.Openid,
						"WPhone":u.WPhone,
						"Phone": u.Phone,
						"WxName":u.WxName,
						"Company":u.Company,
						"Department":u.Department,
					})
				}
				group["users"]=users
				groups = append(groups,group)
			}
			ac["groups"] = groups
			data = append(data,ac)
		}
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusOK,
			"data": data,
		})
		return
	}
}
func MdApprover(logger *logrus.Logger,gdb *gorm.DB)gin.HandlerFunc  {
	return func(c *gin.Context) {
		var param model.MdApprover
		if err := c.ShouldBindJSON(&param);err !=nil{
			logger.Errorf("获取前端参数失败，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		// 查询活动是否存在
		var activity model.Activity
		if err := gdb.Where(model.Activity{ActivityName: param.ActivityName}).First(&activity).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("前端传过来的活动：%s 系统不存在",param.ActivityName)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("param:%#v,查询活动：%s 失败，错误信息：%s",param,param.ActivityName,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}
		// 判断 审批人是否为空
		if activity.Approver == ""{
			activity.Approver="[]"
		}
		//活动存在，判断操作类型
		var user []string
		//反序列化出数据库中，用户列表
		if err:=json.Unmarshal([]byte(activity.Approver),&user);err !=nil{
			logger.Errorf("反序列化数据库已存储的审批人信息报错，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		switch param.OpType {
		case "add":
			if len(user) == 0{
				user = param.Users
			}else {
				// 判断新增用户是否已存在
				for _,u := range param.Users{
					for _,un := range user{
						if u != un{
							user = append(user,u)
						}
					}
				}
			}

		case "del":
			for _,u := range param.Users{
				for i,un := range user{
					if u == un{
						user = append(user[:i],user[i+1:]...)
					}
				}
			}
		case "update":
			user = param.Users
		}
		// 更新数据库操作
		// 序列化审批人
		bu,err:=json.Marshal(user)
		if err !=nil{
			logger.Errorf("序列化：%#v失败，错误信息：%s",user,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		activity.Approver = string(bu)
		if err := gdb.Updates(activity).Error;err !=nil{
			logger.Errorf("更新activity信息：%#v,失败，错误信息：%s",activity,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusOK,
		})
	}
}

