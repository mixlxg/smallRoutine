接口文档



[toc]

---



#### 全局接口

---

##### 全局接口说明

​	前端传递的时间格式必须时间戳（s)

​    后端返回的时间格式时间戳（s)

​	任意接口返回{code: 600/401} ：

​		600：说明用户需要修改初始密码否则无法访问任何接口

​		401：用户未认证或者session到期或认证失败，需要重新登录

---



##### 修改密码接口

- **URL**

```
/pd/app/modifyPwd
```

- **method**

```
POST
```

- **传参说明**

```bash
# 包体传值，application/json
	username	string		用户名			 必传
	oldpwd		string		原始密码		必传
	newpwd		string		新密码			 必传
	wphone		int			手机号码		可选 		#如果用户不同意授权，也让用户修改密码
```

- **返回值**

```bash
// 返回json格式数据
{
    "code": "200/401/400/503",
    "errMsg": "xxx" 
}
```

解释：

```bash
# code 返回码
200： 修改密码成功
401：使用原始密码认证失败
400：非法传参
503：调用接口失败
```

---



##### 获取图形验证码接口

- ***URL地址：***

```bash
/pd/getCaptcha
```

- **method**

```bash
GET
```

- **传参说明**

```bash
# query param
captchaType		string		验证码类型		可选参数	# 默认是digit，可以为digit，string，math
```

 解释：

​	captchaType 值可以为digit，string，math。

- **返回值**

```bash
// 返回值格式json
{
    "code": 200,	// 200,503
    "captchaId": "xxxx", // 验证码唯一标识id用于图形验证码验证使用
    "img":"xxxx"	// base64的验证码加密数据	
    "errMsg":"xxx"	// 错误返回信息，正常当非200时返回
}
```

解释：

```bash
# code
200: 正确返回验证码数据
503：服务器生成验证码失败
#captchaId
验证图形验证码时和用户输入的数据一起传递到后端，用于验证码数据验证
# img
base64 加密的验证码图片数据
#errMsg
错误信息
```

---

#### 小程序接口

---



##### 登录接口

- **URL地址**

``` bash
/pd/login
```

- **method**

```bash
GET
```

- **传参说明**

```bash
# 采用parameter query 方式
	username 	string  	用户名					 必传 
	password	string		用户密码				必传
	wcode		string		小程序code				 必传
	captchaId	string		图形验证码后台验证的唯一id 必传
	captchaValue string		图形验证码输入框值		 必传
```

​	解释：无

- **返回值**

```bash
	// 返回格式json
	{
		"code":"200/600/401/503/400/601/602"，	 // 必返回
		"errMsg"："xxx"，					//错误信息， 返回503时候返回
		"username": "xxx"				// 当返回600时 返回用户名，这里username在系统是唯一不存在重复，会在用户导入时做											//校验
	}
```

​	解释：

```bash
# 返回code状态解释
200: 登录成功
401：用户名或者密码校验失败
400：传参不正确
503：接口调用失败
600：用户绑定成功但需要修改初始密码
601：captchaId，captchaValue 不能为空
602： 图形验证码校验失败
```

---



##### 登出接口

- **URL地址**

```bash
/pd/app/logout
```

- **method**

```bash
POST
```

- **传参说明**

```
无需传参
```

- **返回值**

```bash
# 返回格式json
{
    	"code":	200/503	，	//必返
		"errMsg": "xxx"		//返回码503时返回
}

```

​	解释：

```bash
# code码状态：
	200： 登出成功
	503： 接口调用失败
# errMsg
	当code返回503时，返回的一个错误信息
```

---

##### 小程序文件上传接口

- **URL**

  ```
  /pd/app/upload
  ```

  

- **method**

  ```
  POST
  ```

  

- **传参说明**

  ```
  #form-data
  FileKey  string 	上传文件类型	必传		// 这个filekey传值等于上传文件的key，用于后端获取文件使用
  ActivityName  string	活动名称	必传
  ```

- **返回参数说明**

  ```
  # 返回格式json
  {
      	"code":	200/503/400	，	//必返
  		"errMsg": "xxx"		//返回码503时返回
  		"url":"xxx"		//返回图片的url链接
  }
  ```

  解释

  ```
  # code码状态：
  	200： 登出成功
  	503： 接口调用失败
  	400: 传参失败
  # errMsg
  	当code返回503时，返回的一个错误信息
  # url
  	返回图片url链接
  # 案例
  #入参
  http请求：
  POST /pd/app/upload HTTP/1.1
  Host: 127.0.0.1:8080
  Cookie: sessionid=MTYzNjQ0NjkzOXxOd3dBTkVGUlZqWkZUMFpYVTFOTFRVZFlWVUpVV1ZNMlVsSXlVRXhCTWxaV1RsbFVVbGMyVWsxUVdEUk9VRU5QUmpSRFNGaERVMUU9fBhYY6PXrpYW21WS3ciaVvMG1XByYdg9YeouFORvIDoQ
  Content-Length: 418
  Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW
  
  ----WebKitFormBoundary7MA4YWxkTrZu0gW
  Content-Disposition: form-data; name="pic"; filename="/C:/Users/LvXiuGang/Desktop/吕乐作业/20210523/微信图片_20210524084646.png"
  Content-Type: image/png
  
  (data)
  ----WebKitFormBoundary7MA4YWxkTrZu0gW
  Content-Disposition: form-data; name="FileKey"
  
  pic
  ----WebKitFormBoundary7MA4YWxkTrZu0gW
  Content-Disposition: form-data; name="ActivityName"
  
  一个测试活动1
  ----WebKitFormBoundary7MA4YWxkTrZu0gW
  # 返回值参数
  {
      "code": 200,
      "url": "http://127.0.0.1:8080/pd/statics/一个测试活动1/98e0f013-77de-4a81-8f66-a2ce501bf6e8.png"
  }
  ```
  
  

---

##### 小程序通用文件删除接口

**接口说明**

```
这个接口只会根据指定url删除对应的物理文件，不会修改数据库里面任何文件信息，如果需要删除url和数据库里面的信息，需要单独调用特定的接口
```



- **URL**

  ```
  /pd/app/delPic
  ```

  

- **method**

  ```
  POST
  ```

  

- **传参说明**

  ```
  #body json
  url		string		要删除的url链接		必传
  ```

  

- **返回参数说明**

  ```
  # 返回格式json
  {
      	"code":	200/503	，	//必返
  		"errMsg": "xxx"		//返回码503时返回
  }
  ```

  解释

  ```
  # code码状态：
  	200： 登出成功
  	503： 接口调用失败
  # errMsg
  	当code返回503时，返回的一个错误信息
  # 案例
  #入参
  {
      "url":"http://127.0.0.1:8080/pd/statics/一个测试活动1/98e0f013-77de-4a81-8f66-a2ce501bf6e8.png"
  }
  #出参
  {
      "code": 200
  }
  
  ```

  

---



##### 小程序用户获取当前用户信息

- **URL**

  ```
  /pd/app/getUserMess
  ```

  

- **method**

  ```
  GET
  ```

  

- **传参说明**

  ```
  无参数
  ```

  

- **返回参数说明**

  ```bash
  # 返回格式json
  {
      	"code":	200/503	，	//必返
  		"errMsg": "xxx"		//返回码503时返回
  }
  ```

  解释

  ```bash
  # code码状态：
  	200： 登出成功
  	503： 接口调用失败
  # errMsg
  	当code返回503时，返回的一个错误信息
  # 案例
  #入参
  http://127.0.0.1:8080/pd/app/getUserMess
  #出参
  {
      "code": 200,
      "data": {
          "Company": "攀登",
          "CreateTime": 1634193449,
          "Department": "",
          "Role": "admin",
          "UserName": "admin",
          "id": 3
      }
  }
  ```
  


---

##### 小程序用户获取活动-战队全量分页数据接口

- **URL**

  ```
  /pd/app/getPageActivityMess
  ```

  

- **method**

  ```
  GET
  ```

  

- **传参说明**

  ```
  # query param
  page  int 	当前页码	必传
  step  int	每页显示数据条数	必传
  ```

  

- **返回参数说明**

  ```
  # 返回格式json
  {
      	"code":	200/503	，	//必返
  		"errMsg": "xxx"		//返回码503时返回
  }
  ```

  解释：

  ```
  # code码状态：
  	200： 登出成功
  	503： 接口调用失败
  # errMsg
  	当code返回503时，返回的一个错误信息
  # 案例
  # admin角色用户查询
  # 入参
  http://127.0.0.1:8080/pd/app/getPageActivityMess?page=1&step=4
  #出参
  {
      "code": 200,
      "data": [
          {
              "ActivityName": "一个测试活动1",
              "Approver": [],
              "EndTime": 1638113130,
              "StartTime": 1635348330,
              "groups": [
                  {
                      "GroupLeader": "吕秀刚",
                      "GroupName": "西方战队",
                      "users": [
                          "admin",
                          "吕秀刚"
                      ]
                  }
              ],
              "id": 2
          },
          {
              "ActivityName": "这是我的update测试2",
              "Approver": [],
              "EndTime": 1635579983,
              "StartTime": 1635493583,
              "groups": [],
              "id": 3
          },
          {
              "ActivityName": "南京加油",
              "Approver": [],
              "EndTime": 1635579983,
              "StartTime": 1635493583,
              "groups": [],
              "id": 4
          }
      ],
      "total": 3
  }
  #样例2 普通用户查询自己活动信息
  # 入参
  http://127.0.0.1:8080/pd/app/getPageActivityMess?page=1&step=4
  # 出参
  {
      "code": 200,
      "data": [
          {
              "ActivityName": "一个测试活动1",
              "Approver": [],
              "EndTime": 1638113130,
              "StartTime": 1635348330,
              "groups": [
                  {
                      "GroupLeader": "吕秀刚",
                      "GroupName": "西方战队",
                      "users": [
                          "admin",
                          "吕秀刚"
                      ]
                  }
              ],
              "id": 2
          }
      ],
      "total": 1
  }
  ```

  

---

##### 小程序用户获取活动-战队搜索栏数据接口

- URL

  ```
  /pd/app/getSelectActivityList
  ```

- method

  ```
  GET
  ```

  

- 传参说明

  ```
  # query param
  ActivityId		int			活动id		 //可选
  ```

  

- 返回参数说明

  ```
  # 返回格式json
  {
      	"code":	200/503	，	//必返
  		"errMsg": "xxx"		//返回码503时返回
  }
  ```

  解释：

  ```
  # code码状态：
  	200： 登出成功
  	503： 接口调用失败
  # errMsg
  	当code返回503时，返回的一个错误信息
  # 案例
  #入参当用户为admin角色时
  # 入参
  http://127.0.0.1:8080/pd/app/getSelectActivityList
  # 出参
  {
      "code": 200,
      "data": [
          {
              "ActivityName": "一个测试活动1",
              "activityId": 2,
              "end_flag": false,
              "groups": [
                  {
                      "GroupLeader": "吕秀刚",
                      "GroupName": "西方战队",
                      "groupId": 2
                  }
              ]
          },
          {
              "ActivityName": "这是我的update测试2",
              "activityId": 3,
              "end_flag": true,
              "groups": null
          },
          {
              "ActivityName": "南京加油",
              "activityId": 4,
              "end_flag": true,
              "groups": null
          }
      ]
  }
  # 样例2 admin角色用户带参数ActivityId
  # 入参
  http://127.0.0.1:8080/pd/app/getSelectActivityList?ActivityId=2
  #出参
  {
      "code": 200,
      "data": [
          {
              "ActivityName": "一个测试活动1",
              "activityId": 2,
              "end_flag": false,
              "groups": [
                  {
                      "GroupLeader": "吕秀刚",
                      "GroupName": "西方战队",
                      "groupId": 2
                  }
              ]
          }
      ]
  }
  #样例3 普通用户请求
  #入参
  http://127.0.0.1:8080/pd/app/getSelectActivityList
  # 出参
  {
      "code": 200,
      "data": [
          {
              "ActivityName": "一个测试活动1",
              "activityId": 2,
              "end_flag": false,  
              "groups": [
                  {
                      "GroupLeader": "吕秀刚",
                      "GroupName": "西方战队",
                      "groupId": 2
                  }
              ]
          }
      ]
  }
  # 样例4 普通用户带ActivityId
  #入参
  http://127.0.0.1:8080/pd/app/getSelectActivityList?ActivityId=3
  #出参
  {
      "code": 200,
      "data": []  // 这个用户没有参加这个活动，所有返回空列表
  }
  ```




---

##### 小程序用户创建订单接口

- **URL**

  ```
  /pd/app/createOrder
  ```

  

- **method**

  ```
  POST
  ```

  

- **传参说明**

  ```
  	Customer 				string			客户（个人/公司）					必传
  	CustomerPhone 			string			客户手机号码						 可选
  	CustomerContent 		string			签约产品						   必传
  	OrderTimeLimit			int			    签约合同/产品期限（单位年）			 可选
  	OrderMoney				float			金额（元）						  可选
  	OrderPicUrl				array			合同或者订购产品截图url			  必传
  	OrderCompleteTime 		uint64			合同/产品签约时间					必传
  	ActivityId 				uint64			活动id							必传
  ```

  

- **返回参数说明**

  ```
  # 返回格式json
  {
      	"code":	200/400/503/701/702/704/705，	//必返
  		"errMsg": "xxx"		//返回码503时返回
  }
  ```

  解释

  ```
  # code码状态：
  	200： 登出成功
  	400: 传参错误
  	503： 接口调用失败
  	701：用户还未在这个活动的战队或者分组中，不能参加这个活动
  	702：admin 角色用户不能创建订单
  	703： 活动已经结束不能在创建订单
  	704： 活动不存在
  	705： 活动还没有开始不能创建订单
  # errMsg
  	当code返回503时，返回的一个错误信息
  # 案例
  # 入参 // 这里用户的角色是admin 返回702
  {
      "Customer":"中国铁通有限公司",
      "CustomerPhone":"18260087527",
      "CustomerContent":"集团宽带业务",
      "OrderTimeLimit": 3,
      "OrderMoney": 89000.23,
      "OrderPicUrl": ["http://127.0.0.1:8080/pd/statics/一个测试活动1/98e0f013-77de-4a81-8f66-a2ce501bf6e8.png"],
      "OrderCompleteTime":1636620329,
      "ActivityId":2
  }
  # 出参
  {
      "code": 702
  }
  案例2
  # 入参 用户不再是admin角色
  {
      "Customer":"中国铁通有限公司",
      "CustomerPhone":"18260087527",
      "CustomerContent":"集团宽带业务",
      "OrderTimeLimit": 3,
      "OrderMoney": 89000.23,
      "OrderPicUrl": ["http://127.0.0.1:8080/pd/statics/一个测试活动1/98e0f013-77de-4a81-8f66-a2ce501bf6e8.png"],
      "OrderCompleteTime":1636620329,
      "ActivityId":2
  }
  #出参
  {
      "code": 200
  }
  ```

  

---

##### 小程序用户查询个人订单接口

**接口说明**

```
普通用户根据活动id查询个人的订单全部信息，admin根据活动id查询这个活动的所有订单
```

- **URL**

  ```
  /pd/app/getOrders
  ```

  

- **method**

  ```
  GET
  ```

  

- **传参说明**

  ```
  # param query
  page	int 	当前页码					必传
  step	int		步长，每页多少条数据			必传
  activityId int	活动id					 必传
  ```

  

- **返回参数说明**

  ```
  {
      // 返回格式json
      "code":200/400/503,
      "errMsg":"xxxx",   // 当返回503时返回错误信息
      "data": array,
      "total": int
  }
  ```

  解释
  
  ```
  # code
  200: 成功登录
  400：传参不正确
  503：接口不可用
  #errMsg
  503时返回错误信息 
  # data
  返回的分页数据
  # total
  一共多少条数据
  #样例1
  # 入参
  http://127.0.0.1:8080/pd/app/getOrders?page=1&step=1&activityId=2
  # 出参
  {
      "code": 200,
      "data": [
          {
              "AgreeName": "",						// 审批人
              "Customer": "中国铁通有限公司",				// 客户
              "CustomerContent": "集团宽带业务",		 // 签约业务	
              "CustomerPhone": "18260087527",			// 客户手机
              "IsAgree": false,						// 是否通过审批false 是未审批，true已经审批过了
              "OrderCompleteTime": 1636620329,		// 合同/产品签约时间
              "OrderMoney": 89000.23,					// 签约金额
              "OrderPicUrl": [						//附件图片
                  "http://127.0.0.1:8080/pd/statics/一个测试活动1/98e0f013-77de-4a81-8f66-a2ce501bf6e8.png"
              ],
              "OrderTimeLimit": 3,					//合同/产品签约期限
              "Reason": "",							// 不为空则说明审批被拒绝，拒绝原因
              "orderId": 1							// 订单id
          }
      ],
      "total": 1
  }
  ```
  
  ---

##### 小程序用户查询战队订单接口

**接口说明**

```
这个接口只能当前会活动的队长才能查询，普通用户查询不了
```

- **URL**

  ```
  /pd/app/getGroupOrders
  ```

  

- **method**

  ```
  GET
  ```

  

- **传参说明**

  ```
  # param query
  page	int 	当前页码					必传
  step	int		步长，每页多少条数据			必传
  activityId int	活动id					 必传
  ```

  

- **返回参数说明**

  ```
  {
      // 返回格式json
      "code":200/400/503/706,
      "errMsg":"xxxx",   // 当返回503时返回错误信息
      "data": array,
      "total": int
  }
  ```

  解释

  ```
  # code
  200: 成功登录
  400：传参不正确
  503：接口不可用
  706: 用户不是战队队长，无权查看整个战队订单信息
  #errMsg
  503时返回错误信息 
  # data
  返回的分页数据
  # total
  一共多少条数据
  #样例1
  # 入参
  http://127.0.0.1:8080/pd/app/getGroupOrders?page=1&step=1&activityId=2
  # 出参
  {
      "code": 200,
      "data": [
          {
              "AgreeName": "",
              "Customer": "中国铁通有限公司",
              "CustomerContent": "集团宽带业务",
              "CustomerPhone": "18260087527",
              "IsAgree": false,
              "OrderCompleteTime": 1636620329,
              "OrderMoney": 89000.23,
              "OrderPicUrl": [
                  "http://127.0.0.1:8080/pd/statics/一个测试活动1/98e0f013-77de-4a81-8f66-a2ce501bf6e8.png"
              ],
              "OrderTimeLimit": 3,
              "Reason": "",
              "orderId": 1
          }
      ],
      "total": 1
  }
  ```

---

##### 小程序待审批订单查询接口

**接口说明**

```
配置的审批用户有权限查询特定的活动，admin校色用户能查询所有未审批订单
```

- **URL**

  ```
  /pd/app/getApproveOrders
  ```

  

- **method**

  ```
  GET
  ```

  

- **传参说明**

  无

- **返回参数说明**

  ```
  {
      // 返回格式json
      "code":200/503,
      "errMsg":"xxxx",   // 当返回503时返回错误信息
      "data": array
  }
  ```

  解释：

  ```
  # code
  200: 成功登录
  503：接口不可用
  #errMsg
  错误信息
  # data
  待审批的订单
  # 样例1
  http://127.0.0.1:8080/pd/app/getApproveOrders
  # 返回参数
  {
      "code": 200,
      "data": [
          {
              "ActivityEndTime": 1638113130,
              "ActivityName": "一个测试活动1",
              "ActivityStartTime": 1635348330,
              "ActivityType": "",
              "Customer": "中国铁通有限公司",
              "CustomerContent": "集团宽带业务",
              "CustomerPhone": "18260087527",
              "OrderCompleteTime": 1636620329,
              "OrderMoney": 89000.23,
              "OrderPicUrl": [
                  "http://127.0.0.1:8080/pd/statics/一个测试活动1/98e0f013-77de-4a81-8f66-a2ce501bf6e8.png"
              ],
              "OrderTimeLimit": 3,
              "orderId": 1,
              "userCompany": "攀登",
              "userPhone": "0",
              "username": "admin"
          },
          {
              "ActivityEndTime": 1638113130,
              "ActivityName": "一个测试活动1",
              "ActivityStartTime": 1635348330,
              "ActivityType": "",
              "Customer": "中国铁通有限公司2",
              "CustomerContent": "集团宽带业务222",
              "CustomerPhone": "18260087527",
              "OrderCompleteTime": 1636620329,
              "OrderMoney": 890200.25,
              "OrderPicUrl": [
                  "http://127.0.0.1:8080/pd/statics/一个测试活动1/98e0f013-77de-4a81-8f66-a2ce501bf6e8.png"
              ],
              "OrderTimeLimit": 2,
              "orderId": 2,
              "userCompany": "江苏联通",
              "userPhone": "18260087527",
              "username": "吕秀刚"
          }
      ]
  }
  ```

  

---

#### 管理后台接口

---



##### 管理后台登录接口

- **URL地址**

```bash
/pd/admin/login
```

- **method**

```bash
GET
```

- **传参说明**

```bash
# query param
username	string		用户名							必传
password	string		密码							 必传
captchaId	string		图形验证码后台验证的唯一id 		  必传
captchaValue string		图形验证码输入框值		 		必传
```

- **返回值**

```bash
{
    // 返回格式json
    "code":200/400/503/403/602/401,
    "errMsg":"xxxx",   // 当返回503时返回错误信息
}
```

解释：

```bash
# code
200: 成功登录
400：传参不正确
503：接口不可用
602：图形验证码验证失败
403：权限不足，普通用户无权限登录，只有admin角色用户可以登录
401：账号或者密码不正确
```

---



##### 公司查询列表接口

- **URL**

  ```
  /pd/admin/getCompany
  ```

  

- **method**

  ```
  GET
  ```

  

- **入参说明**

  ```
  # param query
  Company 	string		公司名称		可选
  ```

  

- **返回参数说明**

  ```
  {
      // 返回格式json
      "code":200/503,
      "errMsg":"xxxx",   // 当返回503时返回错误信息
  }
  ```

  解释

  ```
  #code:
  200：正常返回值
  503：服务不可用
  #errMsg
  503时返回错误信息 
  #样例1
  #入参
  http://127.0.0.1:8080/pd/admin/getCompany
  #出参
  {
      "code": 200,
      "data": [
          {
              "Company": "攀登"
          },
          {
              "Company": "江苏联通"
          }
      ]
  }
  # 样例2
  #入参
  http://127.0.0.1:8080/pd/admin/getCompany?Company=江
  #出参
  {
      "code": 200,
      "data": [
          {
              "Company": "江苏联通"
          }
      ]
  }
  ```

  

---



##### 管理后台用户登出接口

- **URL**

```
/pd/admin/logout
```

- **method**

```
POST
```

- **传参说明**

无

- **返回值说明**

```bash
{
    // 返回格式json
    "code":200/503,
    "errMsg":"xxxx",   // 当返回503时返回错误信息
}
```

---



##### 用户账号信息查询接口

- **URL**

```bash
/pd/admin/queryUser
```

- **method**

```bash
POST
```

- **传参说明**

```bash
# 包体传值，application/json
query_type		string		查询用户信息的类型		必传		# 可以传role，user，all，			                                                                            # detail_user_by_role,detail_role_by_user
															# detail_user_by_company
username	string			用户名					可选		
role		string			角色名					可选
company		string			公司名称				可选
```

​	解释：

```bash
# queryTrype
传参role时，返回所有角色名
传参user时,返回所有用户名
传参all时，返回所有的角色和用户，以及他们的对应关系
传参detail_user_by_role 根据参数role查询用户信息
传参detail_role_by_user 根据username查询用户所属的角色
传参detail_user_by_company 根据company查询用户所属的角色，如果company为空自查询所有公司对应用户信息
# username
# role
```

- **返回参数**

```bash
// 返回json数据
{
    "code": "400/200/503",
    "errMsg": "xxxx",
    "data": array/map 		
}
```

解释：

```bash
#code:
400: 传参不正确
200：正常返回值
503：服务不可用
#errMsg
503时返回错误信息 
#data
返回查询数据

#例子1 查询系统所有角色
# 入参
{
    "query_type":"role"
}
# 出参：
{
    "code": 200,
    "data": [
        "admin",
        "generalUser"
    ]
}

# 例子2 查询系统所有用户
# 入参
{
    "query_type":"user"
}
# 出参：
{
    "code": 200,
    "data": [
        {
            "UserName": "admin"
        },
        {
            "UserName": "test1"
        },
        {
            "UserName": "test2"
        }
    ]
}
# 例子3，查询所有用户和角色

# 入参
{
    "query_type":"all"
}

#出参
{
    "code": 200,
    "data": {
        "admin": [
            {
                "Username": "admin"
            }
        ],
        "generalUser": [
            {
                "Username": "test1"
            },
            {
                "Username": "test2"
            }
        ]
    }
}
# 例子4 根据传过来的角色名称查询所有的用户信息
# 入参
{
    "query_type":"detail_user_by_role",
    "role": "generalUser"
}
# 出参
{
    "code": 200,
    "data": [
        {
            "UserName": "test1"
        },
        {
            "UserName": "test2"
        }
    ]
}
# 例子5
# 入参
{
    "query_type":"detail_role_by_user",
    "username": "admin"
}
# 出参
{
    "code": 200,
    "data": "admin"
}
# 例子6 查用公司和用户对应关系
# 入参：
{
    "query_type":"detail_user_by_company"
}
#出参
{
    "code": 200,
    "data": {
        "攀登": [
            "admin"
        ],
        "江苏联通": [
            "test2",
            "吕秀刚"
        ]
    }
}
# 入参
{
    "query_type":"detail_user_by_company",
    "company":"攀登"
}
# 出参
{
    "code": 200,
    "data": {
        "攀登": [
            "admin"
        ]
    }
}

```

---

##### 分页用户账号信息查询接口

- **URL**

  ```bash
  /pd/admin/pageQueryUsers
  ```

  

- **method**

  ```
  POST
  ```

  

- **传参说明**

  ```bash
  # body json
  	page	int		当前页					 必须   
  	step 	int		每页数据条数			   必须
      UserName string	 用户名称				可选		
  	Role string		角色名称				可选   // 不为空会根据role名称查询用户返回
  	Company string	公司名称				可选   // 不为空会根据公司名称查询用户返回
  ```

  

- **返回参数说明**

```bash
{
    "code": "400/200/503",
    "errMsg": "xxxx",
    "data": array/map 		
}
```

解释

```
#code:
400: 传参不正确
200：正常返回值
503：服务不可用
#errMsg
503时返回错误信息 
# 样例  查询所有用户
#入参
{
    "page":2,
    "step":3
}
#出参
{
    "code": 200,
    "data": [
        {
            "Company": "江苏联通",
            "CreateTime": 1634193449,
            "Department": "产户",
            "LoginTime": 1634281460,
            "Openid": "",
            "Phone": "0",
            "Role": "generalUser",
            "UserName": "test2",
            "WPhone": "0",
            "WxName": "",
            "id": 2
        },
        {
            "Company": "攀登",
            "CreateTime": 1634193449,
            "Department": "",
            "LoginTime": 1635127102,
            "Openid": "",
            "Phone": "0",
            "Role": "admin",
            "UserName": "admin",
            "WPhone": "0",
            "WxName": "",
            "id": 3
        },
        {
            "Company": "江苏联通",
            "CreateTime": 1635301913,
            "Department": "产户",
            "LoginTime": null,
            "Openid": "",
            "Phone": "18260087527",
            "Role": "admin",
            "UserName": "吕秀刚",
            "WPhone": "",
            "WxName": "",
            "id": 4
        }
    ],
    "total": 1
}
```



---



##### 用户创建接口

- **URL**

```bash
/pd/admin/createUser
```

- **method**

```bash
POST
```

- **传参说明**

```bash
# body体传值 json
UserName 		string		用户名			必传
Password		string		密码			必传
Phone			string		手机号码		可选
Company			string		公司名称		必传
Department		string		部门名称		可选
Role			string		角色名称		必传
```

- **返回参数说明**

```bash
# 返回数据为json格式
{
	"code": 200/400/503,
	"errMsg": "xxxx",
}
```

解释：

```bash
# code
200: 正常返回值
400：传参错误
503： 服务端处理失败
700: 存在重名用户
#errMsg
503是返回错误信息，用于排错

# 调用案例
#入参：
{
    "UserName":"吕秀刚",
    "Password":"gangzi2010",
    "Phone":"18260087527",	//可选
    "Company":"江苏联通",	
    "Department":"产户",		//可选
    "Role":"admin"
}
#出参：
{
    "code": 200
}
```

---

##### 获取当前用户信息

- **URL**

  ```
  /pd/admin/getCurrentUserMess
  ```

  

- **method**

  ```
  GET
  ```

  

- **传参说明**

  ```
  无
  ```

  

- **返回参数说明**

  ```
  # 返回数据为json格式
  {
  	"code": 200/503,
  	"errMsg": "xxxx",
  }
  ```

  解释

  ```
  # code
  200: 正常返回值
  503： 服务端处理失败
  #errMsg
  503是返回错误信息，用于排错
  # 样例
  # 入参
  http://127.0.0.1:8080/pd/admin/getCurrentUserMess
  #出参
  {
      "code": 200,
      "data": {
          "Company": "攀登",	   //公司
          "CreateTime": 1634193449, //活动创建时间
          "Department": "",			//部门
          "LoginTime": 1635127102,	//用户登录时间
          "Openid": "",				//openid
          "Phone": "0",				// 用户导入时手机号码
          "UserName": "admin",		//用户名
          "WPhone": "0",				//微信注册号码
          "WxName": "",				// 微信名
          "id": 3,					//用户唯一id
          "role": "admin"				// 角色名
      }
  }
  ```

  

---



##### 用户更新接口

- **URL**

```bash
/pd/admin/updateUser
```

- **method**

```bash
POST
```

- **传参说明**

```bash
# body传参，json格式
UserName 		string		用户名			必传
Password		string		密码			 可选
Phone			string		手机号码		可选
Company			string		公司名称		可选
Department		string		部门名称		可选
Role			string		角色名称		可选

```

- **返回参数说明**

```bash
# 返回数据为json格式
{
	"code": 200/400/503,
	"errMsg": "xxxx",
}
```

解释：

```bash
# code
200: 正常返回值
400：传参错误
503： 服务端处理失败
#errMsg
503是返回错误信息，用于排错

# 案例
#入参
{
    "UserName":"test2",
    "Password":"gangzi2010",	//修改了密码
    "Department":"产户",			// 修改了部门
    "Role":"admin"				// 修改了角色	
    							// 其他没有传的参数都没有修改
}
#出参
{
    "code": 200
}
```

---



##### 删除用户接口

- **URL**

```bash
/pd/admin/delUser
```

- **method**

```bash
GET
```

- **传参说明**

```bash
# query param
UserName	string	用户名		必传
```

- **返回参数说明**

```bash
{
	"code": 200/400/503/604,
	"errMsg": "xxxx",
}
```

解释：

```bash
# code
200: 正常返回值
400：传参错误
503： 服务端处理失败
604: 用户已存在订单信息不能删除
#errMsg
503是返回错误信息，用于排错
```

---



##### 创建活动接口

- **URL**

```bash
/pd/admin/createActivity
```

- **method**

```bash
POST
```

- **传参说明**

```bash
# body传值 json
ActivityName		string		活动名称	必传		// 主要于后台区分唯一活动的标识
ActivityContent		string		活动描述	必选	// 用于展示给用户看，可能存在活动相同的情况，所以这个字段类似别名
ActivityType		string		活动类型	必须	// 做成一个select框可选则为（B2C/B2B） B2C代表面向用户，B2B面向企业
StartTime			int		开始时间戳	必传		
EndTime				int		结束时间戳	必传		
```

- **返回参数说明**

```bash
{
	"code": 200/400/503/605,
	"errMsg": "xxxx",
}
```

解释：

```bash
# code
200: 正常返回值
400：传参错误
503： 服务端处理失败
605: 活动已存在
606: 活动结束时间早于开始时间
#errMsg
503是返回错误信息，用于排错
# 案例
# 入参
{
    "ActivityName":"修改时间戳类型测试活动",
    "ActivityContent":"南京加油",
    "StartTime": 1635493583,
    "EndTime": 1635579983
}
#出参
{
    "code": 200
}
```

---

##### 添加，删除，修改活动审批人接口

- **URL**

  ```
  /pd/admin/MdApprover
  ```

  

- **method**

  ```
  POST
  ```

  

- **传参说明**

  ```
  # body  json
  OpType		string			操作类型		必传 // 支持参数  add del update
  ActivityName	string		活动名称		必传
  Users		array			审批人列表		必传
  ```

  

- **返回参数说明**

  ```
  {
  	"code": 200/400/503,
  	"errMsg": "xxxx",
  }
  ```

  解释

  ```
  # code
  200: 正常返回值
  400：传参错误
  503： 服务端处理失败
  #errMsg
  503是返回错误信息，用于排错
  # 样例1 添加活动审批人
  # 入参
  {
      "OpType":"add",
      "ActivityName":"一个测试活动1",
      "Users": ["test2","吕秀刚"]
  }
  #出参
  {
      "code": 200
  }
  # 样例2 删除审批人
  # 入参
  {
      "OpType":"del",
      "ActivityName":"一个测试活动1",
      "Users": ["test2"]
  }
  # 出参
  {
      "code": 200
  }
  # 样例3 修改该更新审批人
  # 入参
  {
      "OpType":"update",
      "ActivityName":"一个测试活动1",
      "Users": ["test2"]
  }
  #出参
  {
      "code": 200
  }
  ```

  

---



##### 删除活动接口

- **URL**

```bash
/pd/admin/delActivity
```

- **method**

```bash
GET
```

- **传参说明**

```bash
# param query
ActivityName		string		活动名称		#必传
```

- **返回参数说明**

```bash
{
	"code": 200/400/503/607/608,
	"errMsg": "xxxx",
}
```

解释：

```bash
# code
200: 正常返回值
400：传参错误
503： 服务端处理失败
607: 活动已经开始存在订单，不能删除
608: 活动已经关联了战队，需要先新建活动，将战队修改到新建的活动中，在删除老的活动
#errMsg
503是返回错误信息，用于排错
# 样例
#入参
http://127.0.0.1:8080/pd/admin/delActivity?ActivityName=一个测试活动
#出参
{
    "code": 200
}
```

---

##### 活动信息查询接口

- **URL**

  ```bash
  /pd/admin/queryActivity
  ```

  

- **method**

  ```bash
  GET
  ```

  

- **传参说明**

  ```bash
  #query param
  ActivityName		string		活动名称	可选		// 主要于后台区分唯一活动的标识
  ```

  解释：

  ```
  #ActivityName
  当ActivityName 不传时返回所有活动信息，如果ActivityName有值返回要查询的活动信息
  ```

  

- **返回参数说明**

  ```bash
  {
  	"code": 200/503,
  	"errMsg": "xxxx",
  }
  ```

  解释：

  ```bash
  # code
  200: 正常返回值
  503： 服务端处理失败
  #errMsg
  503是返回错误信息，用于排错
  # 样例1 查询全部活动信息
  # 入参
  http://127.0.0.1:8080/pd/admin/queryActivity
  #出参
  {
      "code": 200,
      "data": [
          {
              "ActivityContent": "",    //活动名称别名，用于小程序用户显示
              "ActivityName": "一个测试活动1", // 活动唯一名称
              "ActivityType": "",   // 活动类型，B2B,B2C用于创建不同类型订单
              "EndTime": 1638113130, //活动结束时间
              "StartTime": 1635348330 //后动开始时间
          },
          {
              "ActivityContent": "这是我的update测试2",
              "ActivityName": "一个测试活动2",
              "ActivityType": "",
              "EndTime": 1635579983,
              "StartTime": 1635493583
          },
          {
              "ActivityContent": "南京加油",
              "ActivityName": "修改时间戳类型测试活动",
              "ActivityType": "",
              "EndTime": 1635579983,
              "StartTime": 1635493583
          }
      ]
  }
  # 样例2 更加活动名称查询活动信息接口
  # 入参
  http://127.0.0.1:8080/pd/admin/queryActivity?ActivityName=一个测试活动1
  # 出参
  {
      "code": 200,
      "data": [
          {
              "ActivityContent": "",
              "ActivityName": "一个测试活动1",
              "ActivityType": "",
              "EndTime": 1638113130,
              "StartTime": 1635348330
          }
      ]
  }
  ```

  

---



##### 活动信息修改接口

- **URL**

```
/pd/admin/updateActivity
```

- **method**

```
POST
```

- **传参说明**

```bash
# body传值 json
ActivityName		string		活动名称	必传		// 主要于后台区分唯一活动的标识
ActivityContent		string		活动描述	可选	// 用于展示给用户看，可能存在活动相同的情况，所以这个字段类似别名
ActivityType		string		活动类型	可选	// 做成一个select框可选则为（B2C/B2B） B2C代表面向用户，B2B面向企业
StartTime			int		开始时间戳	可选		
EndTime				int		结束时间戳	可选		
```

- **返回参数说明**

```bash
{
	"code": 200/400/503/606,
	"errMsg": "xxxx",
}
```

解释：

```bash
# code
200: 正常返回值
400：传参错误
503： 服务端处理失败
606: 活动结束时间早于开始时间
#errMsg
503是返回错误信息，用于排错
# 案例1，修改活动别名
# 入参
{
    "ActivityName":"一个测试活动2",
    "ActivityContent":"这是我的update测试"
}
#出参
{
    "code": 200
}
# 案例2 修改活动时间
# 入参
{
    "ActivityName":"一个测试活动2",
    "StartTime":1635579983,
    "EndTime":1635493583
}
#出参
{
    "code": 606
}
# 案例3 同时修改别名和活动时间
# 入参
{
    "ActivityName":"一个测试活动2",
    "StartTime":1635493583,
    "EndTime":1635579983
}
#出参
{
    "code": 200
}

```

---



##### 活动战队用户信息查询接口

- **URL**

  ```bash
  /pd/admin/queryActivityGroupsUsers
  ```

  

- **method**

  ```
  POST
  ```

  

- **传参说明**

  ```bash
  # body json
  QueryType	string	查询类型	必选
  ActivityName	string	活动名称	可选
  ```

  解释

  ```bash
  #QueryType 参数可传递的值
  all: 查询所有活动
  one:查询一个活动的关联信息
  #ActivityName
  当QueryType=one时这个参数不能为空
  ```

  

- **返回参数说明**

  ```bash
  {
  	"code": 200/400/503,
  	"errMsg": "xxxx",
  }
  ```

  解释

  ```bash
  # code
  200: 正常返回值
  400：传参错误
  503： 服务端处理失败
  #errMsg
  503是返回错误信息，用于排错
  # 样例1
  #入参
  {
      "QueryType":"all"
  }
  # 出参
  {
      "code": 200,
      "data": {
          "一个测试活动1": {
              "approver": [             // 活动审批人列表
                  "test2"
              ],
              "end_flag": false,			// 活动进行中
              "groups": [
                  {
                      "group_leader": "",   // 队长
                      "group_name": "西方战队",
                      "users": [
                          {
                              "Company": "江苏联通",
                              "Department": "产户",
                              "Openid": "",
                              "Phone": "18260087527",
                              "UserName": "吕秀刚",
                              "WPhone": "",
                              "WxName": "",
                              "id": 4
                          }
                      ]
                  }
              ]
          },
          "一个测试活动2": {
              "approver": [],
              "end_flag": true,
              "groups": null
          },
          "修改时间戳类型测试活动": {
              "approver": [],
              "end_flag": true,
              "groups": null
          }
      }
  }
  
  # 样例2
  # 入参
  {
      "QueryType":"one",
      "ActivityName":"修改时间戳类型测试活动"
  }
  # 出参
  {
      "code": 200,
      "data": {
          "修改时间戳类型测试活动": {
              "approver": [],
              "end_flag": true,
              "groups": null
          }
      }
  }
  
  ```
  
  



---



##### 创建战队接口

- **URL**

```bash
/pd/admin/createGroup
```

- **method**

```bash
POST
```

- **传参说明**

```bash
# body  json
GroupName 		string		战队名称	必传
ActivityName		string		活动名称	必传
```

- **返回参数说明**

```bash
{
	"code": 200/400/503/609,
	"errMsg": "xxxx",
}
```

解释：

```bash
# code
200: 正常返回值
400：传参错误
503： 服务端处理失败
609: 当前活动中这个战队已存在
#errMsg
503是返回错误信息，用于排错
# 样例1
#入参
{
    "GroupName":"东方战神",   // 战队不存在
    "ActivityName":"一个测试活动1"
}
#出参
{
    "code": 200
}

#样例2
# 入参
{
    "GroupName":"东方战神",  // 战队已经存在
    "ActivityName":"一个测试活动1"
}
# 出参
{
    "code": 609
}
```

---



##### 添加用户到战队接口

- **URL**

```bash
/pd/admin/addUsersToGroup
```

- **method**

```bash
POST
```

- **传参说明**

```bash
# body json
GroupName 		string		战队名称	必传
ActivityName		string		活动名称	必传
Users			array			用户名列表	必传
```

- **返回参数说明**

```bash
{
	"code": 200/400/503,
	"errMsg": "xxxx",
}
```

解释：

```bash
# code
200: 正常返回值
400：传参错误
503： 服务端处理失败
#errMsg
503是返回错误信息，用于排错
# 样例1
#入参
{
    "GroupName":"东方战神",				//存在
    "ActivityName":"一个测试活动1",		//存在
    "Users":["test2","吕秀刚"]			
}
#出参
{
    "code": 200
}
#样例2
# 入参
{
    "GroupName":"东方战神1",	// 不存在
    "ActivityName":"一个测试活动1",
    "Users":["test2","吕秀刚"]
}
#出参
{
    "code": 400
}
```

---

##### 设置，修改，删除战队队长接口

- **URL**

  ```
  /pd/admin/setGroupLeader
  ```

  

- **method**

  ```
  POST
  ```

  

- **传参说明**

  ```
  # body json
  	OpType	string		操作类型	必传		//可传参数 add,del,update
  	GroupName	string		战队名称	必传
  	ActivityName	string	活动名称	必传
  	LeaderName string		队长名称	必传 //必须是当前战队人员
  ```

  

- **返回参数说明**

  ```
  {
  	"code": 200/400/503,
  	"errMsg": "xxxx",
  }
  ```

  解释

  ```
  # code
  200: 正常返回值
  400：传参错误
  503： 服务端处理失败
  #errMsg
  503是返回错误信息，用于排错
  # 样例1 添加战队队长
  # 入参
  {
      "OpType":"add",
      "ActivityName":"一个测试活动1",
      "GroupName":"西方战队",
      "LeaderName":"吕秀刚"
  }
  #出参
  {
      "code": 200
  }
  # 样例2 删除队长
  {
      "OpType":"del",
      "ActivityName":"一个测试活动1",
      "GroupName":"西方战队",
      "LeaderName":"吕秀刚"
  }
  #出参
  {
      "code": 200
  }
  # 样例3 更换队长
  #入参
  {
      "OpType":"update",
      "ActivityName":"一个测试活动1",
      "GroupName":"西方战队",
      "LeaderName":"吕秀刚"
  }
  #出参
  {
      "code": 200
  }
  ```

  

---



##### 从战队中删除用户接口

- **URL**

```
/pd/admin/delUserFromGroup
```

- **method**

```
POST
```

- **传参说明**

```bash
# body json
GroupName 		string		战队名称	必传
ActivityName		string		活动名称	必传
Users			array			用户名列表	必传
```

- **返回参数说明**

```bash
{
	"code": 200/400/503,
	"errMsg": "xxxx",
}
```

解释：

```bash
# code
200: 正常返回值
400：传参错误
503： 服务端处理失败
#errMsg
503是返回错误信息，用于排错
# 样例1
# 入参
{
    "GroupName":"东方战神",
    "ActivityName":"一个测试活动1",
    "Users":["test2"]
}
#出参
{
    "code": 200
}
```



---



##### 删除战队接口

- **URL**

```bash
/pd/admin/delGroup
```

- **method**

```bash
POST
```

- **传参说明**

```bash
# body json
GroupName 		string		战队名称	必传
ActivityName		string		活动名称	必传
```

- **返回参数说明**

```bash
{
	"code": 200/400/503,
	"errMsg": "xxxx",
}
```

解释：

```bash
# code
200: 正常返回值
400：传参错误
503： 服务端处理失败
#errMsg
503是返回错误信息，用于排错
#样例
#入参
{
    "GroupName":"东方战神",
    "ActivityName":"一个测试活动1"
}
#出参
{
    "code": 200
}
```

---

##### 战队修改接口

- **URL**

  ```bash
  /pd/admin/modifyGroup
  ```

  

- **method**

  ```bash
  POST
  ```

  

- **传参说明**

  ```bash
  # body json
  GroupName 		string		战队名称	必传
  ActivityName		string		活动名称	必传
  NGroupName		string			新战队名称	可选
  NActivityName	string			新活动名称	可选
  ```

  解释：

  ```bash
  战队可以修改的属性是：战队名称和对应的活动归属
  ```

  

- **返回参数说明**

  ```bash
  {
  	"code": 200/400/503,
  	"errMsg": "xxxx",
  }
  ```

  解释：

  ```bash
  # code
  200: 正常返回值
  400：传参错误
  503： 服务端处理失败
  #errMsg
  503是返回错误信息，用于排错
  #样例 修改战队名称
  #入参
  {
      "GroupName":"东方战神",
      "ActivityName":"一个测试活动2",
      "NGroupName": "西方战队",
      "NActivityName": "一个测试活动1"
  }
  # 出参
  {
      "code": 200
  }
  ```

  



​			







