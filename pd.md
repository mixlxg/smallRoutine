#### 接口文档



[toc]

---

##### 全局接口说明

​	前端传递的时间格式必须是RFC3339 格式串

​    后端返回的时间格式也是RFC3339格式

​	任意接口返回{code: 600/401} ：

​		600：说明用户需要修改初始密码否则无法访问任何接口

​		401：用户未认证或者session到期或认证失败，需要重新登录

##### 获取图形验证码接口

###### URL地址：

```bash
/pd//getCaptcha
```

###### method

```bash
GET
```

###### 传参说明

```bash
# query param
captchaType		string		验证码类型		可选参数	# 默认是digit，可以为digit，string，math
```

 解释：

​	captchaType 值可以为digit，string，math。

###### 返回值

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



##### 登录接口

###### URL地址

``` bash
/pd/login
```



###### method

```bash
GET
```



###### 传参说明

```bash
# 采用parameter query 方式
	username 	string  	用户名					 必传 
	password	string		用户密码				必传
	wcode		string		小程序code				 必传
	captchaId	string		图形验证码后台验证的唯一id 必传
	captchaValue string		图形验证码输入框值		 必传
```

​	解释：无

###### 返回值

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



##### 登出接口

###### URL地址

```bash
/pd/app/logout
```

###### method

```bash
POST
```

###### 传参说明

```
无需传参
```

###### 返回值

```bash
// 返回格式json
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



##### 修改密码接口

###### URL

```
/pd/app/modifyPwd
```

###### method

```
POST
```

###### 传参说明

```bash
# 包体传值，application/json
	username	string		用户名			 必传
	oldpwd		string		原始密码		必传
	newpwd		string		新密码			 必传
	wphone		int			手机号码		可选 		#如果用户不同意授权，也让用户修改密码
```

###### 返回值

```json
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

##### 管理后台登录接口

###### URL地址

```bash
/pd/admin/login
```

###### method

```bash
GET
```

###### 传参说明

```bash
# query param
username	string		用户名							必传
password	string		密码							 必传
captchaId	string		图形验证码后台验证的唯一id 		  必传
captchaValue string		图形验证码输入框值		 		必传
```

###### 返回值

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

##### 管理后台用户登出接口

###### URL

```
/pd/admin/logout
```



###### method

```
POST
```



###### 传参说明

无

###### 返回值说明

```bash
{
    // 返回格式json
    "code":200/503,
    "errMsg":"xxxx",   // 当返回503时返回错误信息
}
```



##### 用户账号信息查询接口

###### URL

```bash
/pd/admin/queryUser
```

###### method

```bash
POST
```

###### 传参说明

```bash
# 包体传值，application/json
query_type		string		查询用户信息的类型		必传		# 可以传role，user，all，			                                                                            # detail_user_by_role,detail_role_by_user
username	string			用户名					可选		
role		string			角色名					可选
```

​	解释：

```bash
# queryTrype
传参role时，返回所有角色名
传参user时,返回所有用户名
传参all时，返回所有的角色和用户，以及他们的对应关系
传参detail_user_by_role 根据参数role查询用户信息
传参detail_role_by_user 根据username查询用户所属的角色
# username
# role
```

###### 返回参数

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

```

##### 用户创建接口

###### URL

```bash
/pd/admin/createUser
```

###### method

```bash
POST
```

###### 传参说明

```bash
# body体传值 json
UserName 		string		用户名			必传
Password		string		密码			必传
Phone			string		手机号码		可选
Company			string		公司名称		必传
Department		string		部门名称		可选
Role			string		角色名称		必传
```

###### 返回参数说明

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

##### 用户更新接口

###### URL

```bash
/pd/admin/updateUser
```

###### method

```bash
POST
```

###### 传参说明

```bash
# body传参，json格式
UserName 		string		用户名			必传
Password		string		密码			 可选
Phone			string		手机号码		可选
Company			string		公司名称		可选
Department		string		部门名称		可选
Role			string		角色名称		可选

```

###### 返回参数说明

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

##### 删除用户接口

###### URL

```bash
/pd/admin/delUser
```

###### method

```bash
GET
```

###### 传参说明

```bash
# query param
UserName	string	用户名		必传
```

###### 返回参数说明

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

##### 创建活动接口

###### URL

```bash
/pd/admin/createActivity
```

###### method

```bash
POST
```

###### 传参说明

```bash
# body传值 json
ActivityName		string		活动名称	必传		// 主要于后台区分唯一活动的标识
ActivityContent		string		活动描述	必选	// 用于展示给用户看，可能存在活动相同的情况，所以这个字段类似别名
StartTime			string		开始时间	必传		// "2021-10-27T23:25:30.511+08:00"  RFC3339 
EndTime				string		结束时间	必传		//"2021-10-28T23:25:30.511+08:00"   RFC3339 
```

###### 返回参数说明

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
    "ActivityName":"一个测试活动2",
    "StartTime":"2021-10-27T23:25:30.511+08:00",
    "EndTime":"2021-10-28T23:25:30.511+08:00"
}
#出参
{
    "code": 200
}
```

##### 删除活动接口

###### URL

```bash
/pd/admin/delActivity
```



###### method

```bash
GET
```



###### 传参说明

```bash
# param query
ActivityName		string		活动名称		#必传
```



###### 返回参数说明

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

##### 活动信息修改接口

###### URL

```
/pd/admin/updateActivity
```



###### method

```
POST
```



###### 传参说明

```bash
# body传值 json
ActivityName		string		活动名称	必传		// 主要于后台区分唯一活动的标识
ActivityContent		string		活动描述	可选	// 用于展示给用户看，可能存在活动相同的情况，所以这个字段类似别名
StartTime			string		开始时间	可选		// "2021-10-27T23:25:30.511+08:00"  RFC3339 
EndTime				string		结束时间	可选		//"2021-10-28T23:25:30.511+08:00"   RFC3339 
```



###### 返回参数说明

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
    "StartTime":"2021-10-27T23:25:30.511+08:00", //结束时间早于开始时间
    "EndTime":"2021-10-26T23:25:30.511+08:00"
}
#出参
{
    "code": 606
}
# 案例3 同时修改别名和活动时间
# 入参
{
    "ActivityName":"一个测试活动2",
    "ActivityContent":"这是我的update测试2",
    "StartTime":"2021-10-27T23:25:30.511+08:00", // +08:00 代表北京时间前端一定要传过来，不然时间会偏差8小时
    "EndTime":"2021-10-29T23:25:30.511+08:00"
}
#出参
{
    "code": 200
}

```

##### 活动战队用户信息查询接口



##### 创建战队接口

###### URL

```bash
/pd/admin/createGroup
```

###### method

```bash
GET
```

###### 传参说明

```bash
# body  json
GroupName 		string		战队名称	必传
ActivityName		string		活动名称	必传
```

###### 返回参数说明

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

##### 添加用户到战队接口

###### URL

```bash
/pd/admin/addUsersToGroup
```

###### method

```bash
POST
```



###### 传参说明

```bash
# body json
GroupName 		string		战队名称	必传
ActivityName		string		活动名称	必传
Users			array			用户名列表	必传
```



###### 返回参数说明

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

##### 从战队中删除用户接口

###### URL

```
/pd/admin/delUserFromGroup
```



###### method

```
POST
```



###### 传参说明

```bash
# body json
GroupName 		string		战队名称	必传
ActivityName		string		活动名称	必传
Users			array			用户名列表	必传
```



###### 返回参数说明

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



##### 删除战队接口

###### URL

```bash
/pd/admin/delGroup
```



###### method

```bash
POST
```



###### 传参说明

```bash
# body json
GroupName 		string		战队名称	必传
ActivityName		string		活动名称	必传
```



###### 返回参数说明

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





​			







