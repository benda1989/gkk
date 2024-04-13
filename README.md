# gkk

基于 Gin Rpcx Gorm Redis NSQ 等开发 WebServer，封装各种常用组件，目的在于提高 Go 应用开发部署效率。
<br><br>

## 目前已支持

✅ 简化 gin
<br>
✅ 接口文档自动生成
<br>
✅ Cron 定时器
<br>
✅ 服务优雅关闭
<br>
✅ 日志抽象和标准字段统一
<br>
✅ 统一参数验证和返回
<br>
✅ 统一错误处理和提示
<br>
✅ 统一数据查询和更新
<br>
✅ 进程结束资源自动回收
<br>
✅ redis/mem 超时缓存
<br>
✅ redis/mem 定时限流 的接口验证码解锁
<br>
✅ 数字 Captcha 生成
<br>

<br>
<br><br>

# 目录结构

```
├── README.md
├── api
│   ├── auto
│   │   ├── auto.go
│   │   ├── common.go
│   │   └── params.go
│   ├── captcha
│   │   └── captcha.go
│   ├── doc
│   │   ├── doc.go
│   │   ├── parser.go
│   │   └── store.go
│   ├── gin.go
│   ├── limit
│   │   ├── limiter.go
│   │   ├── store.go
│   │   └── verify.go
│   ├── log.go
│   ├── params.go
│   └── response.go
├── cache
│   ├── cache.go
│   ├── gin.go
│   └── store.go
├── captcha
│   ├── captcha.go
│   └── store.go
├── code
│   └── code.go
├── compare
│   ├── diff.go
│   └── map.go
├── config
│   ├── config.go
│   └── store.go
├── cron
│   └── crontab.go
├── db
│   ├── gorm.go
│   └── model.go
├── expect
│   └── exception.go
├── go.mod
├── go.sum
├── http.go
├── js
│   ├── encode.go
│   ├── fold.go
│   ├── indent.go
│   ├── scanner.go
│   ├── tables.go
│   └── tags.go
├── logger
│   ├── errMsg.go
│   └── logger.go
├── middle
│   ├── cors.go
│   └── recovery.go
├── new.go
├── req
│   ├── req.go
│   ├── res.go
│   └── valid.go
├── str
│   └── string.go
├── tool
│   ├── encrypt.go
│   └── utils.go
└── upload.go


```

# 配置相关

使用 yaml 配置<br>
配置格式如下

```
app1:
  gin:
    mode: debug
    port: 8401
  db:
    name: x
  custom:
    domain: 
    num: 123
    title: 测试

default: //系统默认项目
  cache:
    host: "127.0.0.1:6379"  //tcp链接
    db: 1  //不配置默认使用0
  db:
    db: postgres
    user: postgres
    password: qwert
    host: localhost
    port: 5432
    MaxIdle: 20
    MaxOpen: 20
    PreferSimpleProtocol: true
  log:
    mode: debug // 前端打印 ，下面是mongo的配置
    host: "localhost:27017"
    db: log
    collection: gkk
    user: gkk
    password: "123456"
```

自定义 Custom 的使用

```
var CF struct {
	SendNum uint   `yaml:"num"`
	Title   string `yaml:"title"`
}
config.GetCustom(&CF)
```

# 统一响应

每个服务请求返回数据参数 Response：
<br>
Code： 0 时为正常，其他参考 code 定义
<br>
Msg： 响应说明
<br>
Data： 数据（请求正常时）
<br>
Path： 路径（请求错误时）
<br>

## 内置响应

```
func RD(c *gin.Context, data any)
func RDE(c *gin.Context, data any, err error)
func RDS(c *gin.Context, data any, count int64)
func RDB(c *gin.Context, data []byte) 
func R(c *gin.Context)
func RE(c *gin.Context, err error) 
func RM(c *gin.Context, msg string)
func RME(c *gin.Context, msg string, err error)
```

# 统一参数

参数校验使用 validator，可以自助添加校验, 项目添加了 Gin 错误优化和翻译器
gkk.Validate => github.com/go-playground/validator/v10
注册一个获取 json tag 的自定义方法:

```
gkk.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
    name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
    if re,_ := IsInSlice([]string{"-","_"},name);re[0]{
        return ""
    }
    return name
})
```

Gin 服务中验证器使用:

```
var s struct {
	Name string `json:"name" bing:"required,gte=3"`
	Mail string `json:"name" bing:"required,email"`
}
api.BindJson(c,&s)
token := api.TokenBindJson(c,&s)
```

# 统一错误

在 gin/rpc 的服务流程中任意位置，使用 panic 抛出 gkk.Expection，服务会相应错误返回给请求者，不会引起错误中断

```
抛出错误
expect.PBMC(pcb.TodaySign,code.SIGNED,"已签到")
接口返回：
{
    "code": 20200,
    "msg": "已签到",
    "path": "POST /api/v1/point/sign"
}

```

## 内置错误
```
type E struct {
	Code int   `json:"code"`
	Msg  any   `json:"msg"`
	Err  error `json:"-"`
}
func (r E) Error() string 
func (r E) Return(c *gin.Context)
```

```
func NE(msg any) error
func NEC(msg any, code int) error
func NEE(err error, msg string) error
func NEEC(err error, msg string, code int) error
//主动抛出错误
func PM(msg any)
func PMC(code int, msg string)
func PBM(flage bool, msg any)
func PBMC(flage bool, code int, msg string)
func PEM(err error, msg any )
func PEMC(err error, code int, msg string)
func PDM(db *gorm.DB, msgs ...string)
```

# 统一日志

基于 logrus 二次封装，抽象统一接口、RPC 调用日志结构，便于后期日志收集和搜索

```
{"code":200,"ip":"192.168.10.98","level":"info","method":"GET","msg":"Gin","period":"31.449944ms","time":"2022-02-17T11:19:44+08:00","url":"/api/v1/setting"}
{"code":200,"ip":"192.168.10.98","level":"info","method":"GET","msg":"Gin","period":"27.957655ms","time":"2022-02-17T11:19:44+08:00","url":"/api/v1/order/price?type=会员"}
{"code":200,"ip":"192.168.10.98","level":"info","method":"GET","msg":"Gin","period":"91.255677ms","time":"2022-02-17T11:19:44+08:00","url":"/api/v1/user"}
```

# 统一数据查询(范型)
抽象类型：get[T any]  
具备的方法：
- FormPSO(c *gin.Context) (id string, re M, p *PageSizeOrder)
- Form(c *gin.Context) M
- DB() *gorm.DB
- Model() *T
- Slice() []*T
- First(w M, preloads ...string) *T
- Find(w M, preloads ...string) []*T
- List(w M, p *PageSizeOrder, preloads ...string) ([]*T, int64)
- Select(w any, args ...any) []IdName
- Get(c *gin.Context)
- Comment() MS  

实例方法：   
NewFind[T any] (order string) *get[T]

备注：  
查询标签：cru "-" 为忽略, 前端传送字段使用form标签
查询条件： =, >, >=, <, <=, like(英文逗号分隔条件表示or), in (英文逗号分隔条件)
特例查询：cru: "create_at > ?"

# 统一数据更新(范型)
抽象类型：post[T any]  继承get[T any]
具备的方法：

- Req(c *gin.Context) *T
- Create(data *T) *T
- Update(new *T, w M) (*T, *T)
- Active(id any, w M)
- Delete(ids []uint, where M) []*T

实例方法：   
NewModel[T any](order string, unUpdate ...string) *post[T]

备注：  
unUpdate：用于更新时忽略字段控制，json标签为——的也会被忽略 


# GIN 服务(范型)

根据项目使用习惯和规则，设计规范，简化整体环节  
依赖统一数据查询/更新

## 1 模型添加标签：cru
## 2 编写路由的方法

map[string]GinHandler[T]  
key解析：(权限:)方式:路径

- 权限如果有，使用：auth/admin
- 方法有：GET;POST;PUT;DELETE;ACTIVE;SELECT
- SELECT: 用户下拉 返回 id 和 name 组成的切片
- ACTIVE: 用来更改 status 状态，1-2 之间的变换
- GET 方法不传 id 则是 list

备注：使用内置方法时，字典的值传递nil ，方法参考 getDefaultMethod


## 3 编写对应路由方法 GinHandler

内置方法不满足业务需求，则可以自定义处理方法  
func () GET(c gkk.GC[T])
gkk.GC 有如下属性： 
- Context  *gin.Context
- Auth     注册的auth
- IfReturn 

gkk.GC 有如下方法：  

- BindJson(ptr)
- BindParam(ptr)
- Id(): query 中 id
- Ids(): json 中的 ids
- Form(): query 查询参数
- FormFSO() id,map,pageSizeOrder: 主键，参数，分页排序
- Req(): json 中全部参数T
- 
- DB() *gorm.DB
- DBM() (*gorm.DB, *T)
- DBMS() (*gorm.DB, []*T)
- Comment() MS
- 
- RD(any,...error) 单条返回
- RDS(any,int64,...error) list 返回 （当 int64 为负数时：等同于 RD 返回）
- 
- GetAuth() (w M)
- Get(where M, preloads ...string)
- First(w M, preloads ...string) *T
- Find(w M, preloads ...string) []*T
- List(w M, p *PageSizeOrder, preloads ...string) ([]*T, int64)

## 4 注册
GinRegister[T any](  
gin.IRoutes,     
*post[T],   
map[string]GinHandler[T],   
...*authHandler  
)


## 5 操作记录Record
内置方法在新增/修改/删除时会调用注册的 Record方法
record func(Auth, *Record[T])  
```
type Record[T any] struct {
	Data   *T   //操作前的数据
	Diff   MSS  //更新世数据不一致的数据
	Method string
}
```
## 6 路由打印

```
default注册
[GIN-debug] GET   /api/apps/vip_link/agent/select  --> gkk.newGinHandlerFunc.func1 (4 handlers)
[GIN-debug] GET    /api/apps/vip_link/agent/admin --> gkk.newGinHandlerFunc.func1 (4 handlers)
[GIN-debug] POST   /api/apps/vip_link/agent/admin --> gkk.newGinHandlerFunc.func1 (4 handlers)
[GIN-debug] PUT    /api/apps/vip_link/agent/admin --> gkk.newGinHandlerFunc.func1 (4 handlers)
[GIN-debug] DELETE /api/apps/vip_link/agent/admin --> gkk.newGinHandlerFunc.func1 (4 handlers)
[GIN-debug] GET    /api/apps/vip_link/asset  --> gkk.newGinHandlerFunc.func1 (4 handlers)
[GIN-debug] GET    /api/apps/vip_link/asset/admin --> gkk.newGinHandlerFunc.func1 (4 handlers)
path注册
[GIN-debug] GET    /api/apps/vip_link/asset/total --> gkk.newGinHandlerFunc.func1 (4 handlers)
[GIN-debug] GET    /api/apps/vip_link/asset/total/admin --> gkk.newGinHandlerFunc.func1 (4 handlers)

```

# 接口文档生成

POST：/api/doc 开始记录   
GET： /api/doc 获取结果

- 使用BindJson/BindParam绑定参数
- 使用gkk.R*等 处理返回数据
- 接口说明：注释写在调用Bind*的gin.HandlerFunc前
- 接口需要访问一次才能生成请求/返回参数
- 根据路由自动划分目录， 编写路由要做到有逻辑
- 本地测试环境使用，保存生成的json


# 数据缓存

默认开启内存缓存，配置 default: cache: host: 可以使用 redis 缓存

```
gin.RouterGroup.GET("/", cache.Gin(time.Hour*24), func(c *gin.context){c.JSON("Hello")})
定制缓存的key
gin.RouterGroup.GET("/", cache.Gin(time.Hour*24,func(c *gin.context)string{return c.Request.RequestURI}), func(c *gin.context){c.JSON("Hello")})
```

自定义缓存

```
func add(a string,b *string){
  *b = a  //b是输出指针类型，a是值类型变量(局域)，a的值和指针在函数结束会被销毁。
}
a := "1234"
var b string

cache.Cache(10*time.second,  a,  &b,  add)
//            保存时间,    输入，输出，函数

```

### 缓存分组预置：

limit： 限速  
code： 验证码  
cache： 数据缓存  
func： 函数缓存

# 限流(存储在缓存)

使用方法

```
限制单个ip每分钟100次访问
*gin.Engine.Use( limit.Gin( gkk.NewLimiter(time.Minute,100) ) )
```

定制使用

```
定制ip白名单
func ss(s string)(re bool){
	if s == "127.0.0.1" {
	    re = true
	 }
	return
}
定制缓存key
func cc(c *gin.Context) string{
    return c.Request.URL.Path
}
定制限速处理
func vv(key string, c *gin.Context){
    c.Header("X-Limit-Reset", l.Reset)
	expect.Panic(c, &gkk.Exception{429,gkk.IP_LINMIT, "超过响应次数限制"})
}
*gin.Engine.Use( limit.Gin( gkk.NewLimiter(time.Minute, 100, gkk.LimiterIP(ss)),
                                                                gkk.LimiterKey(cc)),
                                                                gkk.LimiterLimited(vv)))
```

预置定制

```
限制接口 每分钟100次访问
*gin.Engine.Use( limit.Path( time.Minute, 100) )

限制接口 每个ip, 每分钟10次访问
*gin.Engine.Use( limit.PathIp( time.Minute, 10) )
限制接口 每个ip, 每分钟1次访问自动验证码
*gin.Engine.Use( limit.PathIPCode(time.Minute, 1) )
前端的数据：
        {
            "code": 10004,
            "data": {
                "key": "code:192.168.30.24:/api/v1/translate",
                "base64": "data:image/png;base64,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"
            },
            "msg": "超过响应次数限制"
        }
```


# Captcha

### value, item := gkk.Captcha.DrawCaptcha()

value：验证码的字符串，5 位数字

### 验证和重置

*gin.Engine.POST("captcha", captcha.Verify) //验证<br>
*gin.Engine.GET("captcha", captcha.Reset) //重置

#### 转 base64

item.EncodeB64()

#### 转二进制

item.EncodeBinary()

#### 存储本地

item.WriteTo(os.OpenFile("1.png"))

## 定制

captcha := captcha.NewDigitCaptcha(80, 240, 5, 0.7, 80) <br>
宽/高/字数/倾斜度/星星数量

# 七牛云上传

两种方式使用  
1 HandFunc 接口使用:  
router.POST("upload", gkk.UploadQNGin("auth/avatar/"))  
2 从网络链接上传获取地址：  
Avatar := gkk.UploadQNUrl(user.HeadImgURL, conf.QnBase)

# 项目运行

```
func main(){
    gkk.Run(
        code.SSHTunnel,                     //函数1
                    "8306",                //参数1
                    "11313",
                    "7000",
                    "6379",
        app1.Run, false,   //函数2，参数
        app2.Run,                     //函数3
    )
}
```

# 项目引用

go.mod

```
require gkk v0.0.0

replace (
	gkk => ../gkk
)
```

go mod tidy -compat=1.18

<br> <br> <br>

# 框架使用的 golang 版本：1.18
#### 20231101,GH,fixed