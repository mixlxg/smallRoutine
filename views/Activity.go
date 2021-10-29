package views

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"smallRoutine/model"
	"smallRoutine/utils"
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
		activityName := c.Query("activityName")
		var activity model.Activity
		if err:= gdb.Where("activity_name=?",activityName).Find(&activity).Error;err !=nil{
			logger.Errorf("查询活动失败，错误信息：%s",activityName)
			c.JSON(http.StatusOK,gin.H{
				"code":http.StatusServiceUnavailable,
				"errMsg": err.Error(),
			})
		}
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusOK,
			"data":activity,
		})
	}
}
