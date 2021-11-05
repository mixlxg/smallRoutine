package views

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"smallRoutine/model"
	"time"
)

func CreateGroup(logger *logrus.Logger,gdb *gorm.DB) gin.HandlerFunc  {
	return func(c *gin.Context) {
		var param model.ParamGroup
		if err := c.ShouldBindJSON(&param);err !=nil{
			logger.Errorf("创建战队接口，绑定前端传递参数失败，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		// 判断当前活动中是否已经存在这个战队
		// 判断当前活动是否已存在
		var activity = model.Activity{ActivityName: param.ActivityName}
		if err:= gdb.Where(&activity).First(&activity).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("活动：%s,活动不存在传参有问题",activity.ActivityName)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("获取活动：%s,信息失败，错误信息：%s",activity.ActivityName,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code":http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}
		// 判断当前活动是否过期
		if activity.EndTime.Unix() < time.Now().Unix(){
			logger.Errorf("param:%#v,活动：%s 活动已结束", param,param.ActivityName)
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		// 判断当前战队名是否存在
		var group model.Group
		if err := gdb.Model(&activity).Where(model.Group{GroupName: param.GroupName}).Association("Groups").Find(&group); err !=nil{
			logger.Errorf("查询战队:%s，信息失败，错误信息：%s",activity.ActivityName,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		// 判断是否存在
		if group.GroupName == param.GroupName{
			logger.Infof("活动：%s,中战队：%s，已经存在，不能在添加",param.ActivityName,param.GroupName)
			c.JSON(http.StatusOK,gin.H{
				"code":609,
			})
			return
		}
		// 创建战队
		group.ActivityID= activity.ID
		group.GroupName = param.GroupName
		if err:= gdb.Create(&group).Error;err !=nil{
			logger.Errorf("创建活动：%s 的战队：%s 失败，错误信息：%s",param.ActivityName,param.GroupName,err.Error())
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

func AddUsersToGroup(logger *logrus.Logger,gdb *gorm.DB) gin.HandlerFunc  {
	return func(c *gin.Context) {
		var param model.AddUsersToGroup
		if err := c.ShouldBindJSON(&param);err !=nil{
			logger.Errorf("前端传参失败，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusBadRequest,
			})
			return
		}
		// 校验活动战队是否存在
		var group model.Group
		if err:=gdb.Joins("inner join activities on activities.id = groups.activity_id and activities.activity_name =?",param.ActivityName).Where(&model.Group{GroupName: param.GroupName}).First(&group).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("活动：%s,战队：%s 不存在",param.ActivityName,param.GroupName)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("查询战队：%s 失败，错误信息：%s",param.GroupName,err.Error())
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}
		// 查询用户信息
		var users []*model.User
		if err:= gdb.Where("user_name in ?",param.Users).Find(&users).Error;err !=nil{
			logger.Errorf("查询用户:%#v 信息报错，错误信息:%s",param.Users,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		// 添加用户到战队中
		if err := gdb.Model(&group).Association("Users").Append(users);err !=nil{
			logger.Errorf("添加战队失败，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
		})
		return
	}
}

func DelUserFromGroup(logger *logrus.Logger, gdb *gorm.DB)gin.HandlerFunc  {
	return func(c *gin.Context) {
		var param model.DelUsersFromGroup
		if err := c.ShouldBindJSON(&param);err !=nil{
			logger.Errorf("绑定前端参数失败，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		// 判断group是否存在
		var group model.Group
		if err:=gdb.Joins("inner join activities on activities.id = groups.activity_id and activities.activity_name =?",param.ActivityName).Where(&model.Group{GroupName: param.GroupName}).First(&group).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("活动：%s,战队：%s 不存在",param.ActivityName,param.GroupName)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("查询战队：%s 失败，错误信息：%s",param.GroupName,err.Error())
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}
		// 查询用户信息
		var users []*model.User
		if err:= gdb.Where("user_name in ?",param.Users).Find(&users).Error;err !=nil{
			logger.Errorf("查询用户:%#v 信息报错，错误信息:%s",param.Users,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		// 删除战队和用户的关联关系
		if err :=gdb.Model(&group).Association("Users").Delete(users);err !=nil{
			logger.Errorf("param:%#v,从战队：%s 中删除用户报错，错误信息：%s",param,param.GroupName,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		// 判断队长是否在要删除的用户中，如果在将队长字段设置为“”
		for _,u := range param.Users{
			if u == group.GroupLeader{
				group.GroupLeader = ""
				if err := gdb.Select("GroupLeader").Updates(group).Error;err !=nil{
					logger.Errorf("修改战队队长属性失败，错误信息：%s", err.Error())
					c.JSON(http.StatusOK,gin.H{
						"code": http.StatusServiceUnavailable,
						"errMsg": err.Error(),
					})
					return
				}
			}
		}
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
		})
		return
	}
}

func DelGroup(logger *logrus.Logger, gdb *gorm.DB)gin.HandlerFunc  {
	return func(c *gin.Context) {
		var param model.ParamGroup
		// 获取前端传递过来的参数
		if err:= c.ShouldBindJSON(&param);err !=nil{
			logger.Errorf("前端传参失败，错误信息：%s",err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		//判断这个活动的战队是否存在
		var activity  = model.Activity{ActivityName: param.ActivityName}
		if err:= gdb.Where(activity).First(&activity).Error;err != nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("活动：%s 不存在",activity.ActivityName)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("查询活动：%s 的信息失败，错误信息：%s",activity.ActivityName,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}

		var group = model.Group{GroupName: param.GroupName,ActivityID: activity.ID}
		if err := gdb.Where(group).First(&group).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("活动：%s,战队：%s不错在",param.ActivityName,param.GroupName)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("查询战队：%s 的信息报错，错误信息：%s",param.GroupName,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}
		// 删除战队里面的用户关联关系，然后删除用户
		if err:= gdb.Model(&group).Association("Users").Clear();err !=nil{
			logger.Errorf("删除活动：%s,战队失败：%s,错误信息：%s",param.ActivityName,param.GroupName,err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
			return
		}
		// 删除用户
		if err:=gdb.Delete(&group).Error;err !=nil{
			logger.Errorf("删除活动：%s,战队失败：%s,错误信息：%s",param.ActivityName,param.GroupName,err.Error())
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

func GroupModify(logger *logrus.Logger, gdb *gorm.DB)gin.HandlerFunc  {
	return func(c *gin.Context) {
		var param model.ModifyGroup
		if err := c.ShouldBindJSON(&param);err !=nil{
			logger.Errorf("前端传过来的参数有问题，报错信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		// 判断现在的战队是否存在
		var group model.Group
		if err :=gdb.Model(&group).Joins("inner join activities on activities.id=groups.activity_id and activities.activity_name=?",param.ActivityName).Where(model.Group{GroupName: param.GroupName}).First(&group).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("param:%#v,活动：%s 或者战队:%s 不存在",param,param.ActivityName,param.GroupName)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("查询活动：%s,战队:%s 的信息失败，错误信息:%s",param.ActivityName,param.GroupName,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}
		// 判断要修改的属性
		var nactivity model.Activity
		if param.NGroupName != ""{
			if param.NActivityName == ""{
				group.GroupName = param.NGroupName
			}else {
				if err := gdb.Where(model.Activity{ActivityName: param.NActivityName}).First(&nactivity).Error;err !=nil{
					if errors.Is(err,gorm.ErrRecordNotFound){
						logger.Errorf("新活动：%s不存在",param.NActivityName)
						c.JSON(http.StatusOK,gin.H{
							"code": http.StatusBadRequest,
						})
						return
					}else {
						logger.Errorf("查询新活动：%s 报错，错误信息:%s",param.NActivityName,err.Error())
						c.JSON(http.StatusOK,gin.H{
							"code": http.StatusServiceUnavailable,
							"errMsg": err.Error(),
						})
						return
					}
				}
				group.GroupName = param.NGroupName
				group.ActivityID = nactivity.ID
			}
		}else {
			if param.NActivityName !=""{
				if err := gdb.Where(model.Activity{ActivityName: param.NActivityName}).First(&nactivity).Error;err !=nil{
					if errors.Is(err,gorm.ErrRecordNotFound){
						logger.Errorf("新活动：%s不存在",param.NActivityName)
						c.JSON(http.StatusOK,gin.H{
							"code": http.StatusBadRequest,
						})
						return
					}else {
						logger.Errorf("查询新活动：%s 报错，错误信息:%s",param.NActivityName,err.Error())
						c.JSON(http.StatusOK,gin.H{
							"code": http.StatusServiceUnavailable,
							"errMsg": err.Error(),
						})
						return
					}
				}
				group.ActivityID = nactivity.ID
			}else {
				logger.Errorf("param:%#v,新战队名称和新的活动名称不能同时为空",param)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}
		}
		// 判断完成，更新战队信息
		if err := gdb.Updates(&group).Error;err !=nil{
			logger.Errorf("param:%#v,更新战队信息失败，错误信息：%s",param,err.Error())
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

func SetGroupLeader(logger *logrus.Logger,gdb *gorm.DB)gin.HandlerFunc  {
	return func(c *gin.Context) {
		var param model.SetGroupLeader
		if err := c.ShouldBindJSON(&param); err!=nil{
			logger.Errorf("获取前端传参数失败，错误信息：%s",err.Error())
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusBadRequest,
			})
			return
		}
		// 查询战队信息
		var group model.Group
		if err :=gdb.Model(&group).Joins("inner join activities on activities.id = groups.activity_id and activities.activity_name = ?",param.ActivityName).Where(model.Group{GroupName: param.GroupName}).Preload("Users").First(&group).Error;err !=nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				logger.Errorf("param:%#v,活动：%s 或者战队:%s 不存在", param,param.ActivityName,param.GroupName)
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusBadRequest,
				})
				return
			}else {
				logger.Errorf("param:%#v,查询数据库报错，错误信息：%s",param,err.Error())
				c.JSON(http.StatusOK,gin.H{
					"code": http.StatusServiceUnavailable,
					"errMsg": err.Error(),
				})
				return
			}
		}
		// 查询到战队信息，判断要设置的队长是否在这个战队中
		var userIsExist = false
		for _,user := range group.Users{
			if user.UserName == param.LeaderName{
				userIsExist = true
				break
			}
		}
		if userIsExist == false{
			logger.Errorf("要设置的战队队长不在这个战队中")
			c.JSON(http.StatusOK,gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
		switch param.OpType {
		case "add":
			group.GroupLeader = param.LeaderName
		case "update":
			group.GroupLeader =param.LeaderName
		case "del":
			group.GroupLeader =""
		}
		// 更新数据库
		if err:=gdb.Select("GroupLeader").Updates(group).Error;err !=nil{
			logger.Errorf("param:%#v,更新战队队长信息失败，错误信息：%s",param,err.Error())
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