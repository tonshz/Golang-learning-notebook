# Go 语言编程之旅(二)：HTTP 应用(六) 

## 八、对接口进行访问控制

在完成了相关的业务接口的开发后，还有一个问题，这些 API 接口，没有鉴权功能，也就是所有知道地址的人都可以请求该项目的 API 接口和 Swagger 文档，甚至有可能会被网络上的端口扫描器扫描到后滥用，这非常的不安全，怎么办呢。实际上，应该要考虑做纵深防御，对 API 接口进行访问控制。

目前市场上比较常见的两种 API 访问控制方案，分别是 OAuth 2.0 和 JWT(JSON Web Token)，但实际上这两者并不能直接的进行对比，因为它们是两个完全不同的东西，对应的应用场景也不一样，可以先大致了解，如下：

- OAuth 2.0：**OAuth 2.0 是一种授权框架**，本质上是一个授权的行业标准协议，提供了一整套的授权机制的指导标准，常用于使用第三方登陆的情况，像是在网站登录时，会有提供其它第三方站点（例如用微信、QQ、Github 账号）关联登陆的，往往就是用 OAuth 2.0 的标准去实现的。并且 OAuth 2.0 会相对重一些，常常还会授予第三方应用去获取到对应账号的个人基本信息等等。在实现 OAuth 2.0 时可以将 JWT 作为一种认证机制使用。
- JWT：**JWT 是一种认证协议**，与 OAuth 2.0 完全不同，它常用于前后端分离的情况，能够非常便捷的给 API 接口提供安全鉴权，因此在本章节采用的就是 JWT 的方式，来实现 API 访问控制功能。

### 1. JWT 是什么

JSON Web Token（JWT）是一个开放标准（RFC7519），它定义了一种紧凑且自包含的方式，用于在各方之间作为 JSON 对象安全地传输信息。 由于此信息是经过数字签名的，因此可以被验证和信任。 可以使用使用 RSA 或 ECDSA 的公用/专用密钥对对 JWT 进行签名，其格式如下：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205152113517.jpeg)

JSON Web 令牌（JWT）是由紧凑的形式三部分组成，这些部分由点 “.“ 分隔，组成为 `”xxxxx.yyyyy.zzzzz“ `的格式，三个部分分别代表的意义如下：

- Header：头部。
- Payload：有效载荷。
- Signature：签名。

#### a. Header

Header（头部）通常由两部分组成，**分别是令牌的类型和所使用的签名算法**（HMAC SHA256、RSA 等），其会组成一个 JSON 对象用于描述其元数据，例如：

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

在上述 JSON 中 `alg` 字段表示所使用的签名算法，默认是 HMAC SHA256（HS256），而 type 字段表示所使用的令牌类型，使用的 JWT 令牌类型，在最后会对上面的 JSON 对象进行` base64UrlEncode `算法进行转换成为 JWT 的第一部分。

#### b. Payload

Payload（有效负载）也是一个 JSON 对象，**主要存储在 JWT 中实际传输的数据**，如下：

```json
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true
}
```

- `aud（Audience）`：受众，也就是接受 JWT 的一方。
- `exp（ExpiresAt）`：所签发的 JWT 过期时间，过期时间必须大于签发时间。
- `jti（JWT Id）`：JWT 的唯一标识。
- `iat（IssuedAt）`：签发时间。
- `iss（Issuer）`：JWT 的签发者。
- `nbf（Not Before）`：JWT 的生效时间，如果未到这个时间则为不可用。
- `sub（Subject）`：主题。

同样也会对该 JSON 对象进行 base64UrlEncode 算法将其转换为 JWT Token 的第二部分。

这时候需要注意一个问题点，也就是 JWT 在转换时用的 base64UrlEncode 算法，也就是它是可逆的，因此一些敏感信息不要放到 JWT 中，若有特殊情况一定要放，**也应当进行一定的加密处理。**

#### c. Signature

Signature（签名）部分是对前面两个部分组合（Header+Payload）进行约定算法和规则的签名，**而签名将会用于校验消息在整个过程中有没有被篡改**，并且对有使用私钥进行签名的令牌，它还可以验证 JWT 的发送者是否它的真实身份。

在签名的生成上，在应用程序指定了密钥（secret）后，会使用传入的指定签名算法（默认是 HMAC SHA256），然后通过下述的签名方式来完成 Signature（签名）部分的生成，如下：

```lua
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret)
```

可以看出 JWT 的第三部分是由 Header、Payload 以及 Secret 的算法组成而成的，因此它最终可达到用于校验消息是否被篡改的作用之一，因为如果一旦被篡改，Signature 就会无法对上。

#### d. Base64UrlEncode

实际上 Base64UrlEncode 是 Base64 算法的变种，为什么要变呢，原因是在实际开发过程中经常可以看到 JWT 令牌会被放入 Header 或 Query Param 中（也就是 URL）。

而在 URL 中，一些个别字符是有特殊意义的，例如：“+”、“/”、“=” 等等，因此在 Base64UrlEncode 算法中，会对其进行替换，例如：“+” 替换为 “-”、“/” 替换成 “_”、“=” 会被进行忽略处理，以此来保证 JWT 令牌的在 URL 中的可用性和准确性。

### 2. JWT 的使用场景

通常会先在内部约定好 JWT 令牌的交流方式，像是存储在 Header、Query Param、Cookie、Session 都有，但最常见的是存储在 Header 中。然后服务端提供一个获取 JWT 令牌的接口方法，返回而客户端去使用，在客户端请求其余的接口时需要带上所签发的 JWT 令牌，然后服务端接口也会到约定位置上获取 JWT 令牌来进行鉴权处理，以此流程来鉴定是否合法。

### 3. 安装 JWT

拉取 `jwt-go`，该库提供了 JWT 的 Go 实现，能够便捷的提供 JWT 支持，不需要自己去实现。

```bash
$ go get -u github.com/dgrijalva/jwt-go@v3.2.0
```

### 4. 配置 JWT

#### a. 创建认证表

在介绍 JWT 和其使用场景时，了解了实际上需要一个服务端的接口来提供 JWT 令牌的签发，并且可以将自定义的私有信息存入其中，那么必然需要一个地方来存储签发的凭证，否则谁来都签发，似乎不大符合实际的业务需求，因此要创建一个新的数据表，用于存储签发的认证信息，表 SQL 语句如下：

```sql
CREATE TABLE `blog_auth` (
                             `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                             `app_key` varchar(20) DEFAULT '' COMMENT 'Key',
                             `app_secret` varchar(50) DEFAULT '' COMMENT 'Secret',
                             `created_on` int(10) unsigned DEFAULT '0' COMMENT '创建时间',
                             `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
                             `modified_on` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
                             `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
                             `deleted_on` int(10) unsigned DEFAULT '0' COMMENT '删除时间',
                             `is_del` tinyint(3) unsigned DEFAULT '0' COMMENT '是否删除 0 为未删除、1 为已删除',
                              PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='认证管理';
```

上述表 SQL 语句的主要作用是创建了一张名为` blog_auth` 的表，其核心是 `app_key `和 `app_secret` 字段，用于签发的认证信息，接下来默认插入一条认证的 SQL 语句（也可以做一个接口），便于认证接口的后续使用，插入的 SQL 语句如下：

```sql
INSERT INTO `ch02`.`blog_auth`(`id`, `app_key`, `app_secret`, `created_on`, `created_by`, `modified_on`, `modified_by`, `deleted_on`, `is_del`) VALUES (1, 'admin', 'go-learning', 0, 'test', 0, '', 0, 0);
```

该条语句的主要作用是新增了一条`app_key` 为 admin以及 `app_secret `为 `go-learning`的数据。

#### b. 新建 model 对象

接下来打开项目的 `internal/model` 目录下的` auth.go` 文件，写入对应刚刚新增的` blog_auth `表的数据模型，如下：

```go
package model

type Auth struct {
   *Model
   AppKey    string `json:"app_key"`
   AppSecret string `json:"app_secret"`
}

func (a Auth) TableName() string{
   return "blog_auth"
}
```

#### c. 初始化配置

接下来需要针对 JWT 的一些相关配置进行设置，修改项目的 `configs/config.yaml` 配置文件，写入新的配置项，如下：

```yaml
# JWT 初始化配置
JWT:
  Secret: admin
  Issuer: blog-service
  Expire: 7200
```

然后对 JWT 的配置进行初始化操作，修改项目的启动文件` main.go`，修改其 `setupSetting` 方法，如下：

```go
func setupSetting() error {
	...
	err = settings.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	global.JWTSetting.Expire *= time.Second
	...
}
```

在上述配置中，设置了 JWT 令牌的 Secret（密钥）为 `admin`，签发者（Issuer）是 `blog-service`，有效时间（Expire）为 7200 秒，这里需要注意的是 Secret 千万不要暴露给外部，只能有服务端知道，否则是可以解密出来的，非常危险。

### 5. 处理 JWT 令牌

虽然 `jwt-go `库能够帮助开发者快捷的处理 JWT 令牌相关的行为，但是还是需要根据项目特性对其进行设计的，简单来讲，就是组合其提供的 API，设计鉴权场景。

在 `pkg/app` 并创建` jwt.go` 文件，写入第一部分的代码：

```go
package app

import (
   "demo/ch02/global"
   "github.com/dgrijalva/jwt-go"
)

type Claims struct {
   AppKey    string `json:"app_key"`
   AppSecret string `json:"app_secret"`
   jwt.StandardClaims
}

func GetJWTSecret() []byte {
   return []byte(global.JWTSetting.Secret)
}
```

这块主要涉及 JWT 的一些基本属性，第一个是` GetJWTSecret` 方法，用于获取该项目的 JWT Secret，目前是直接使用配置所配置的 Secret，第二个是 Claims 结构体，分为两大块，第一块是项目嵌入的 `AppKey` 和 `AppSecret`，用于自定义的认证信息，第二块是 `jwt.StandardClaims` 结构体，它是` jwt-go` 库中预定义的，也是 JWT 的规范，其涉及字段如下：

```go
// Structured version of Claims Section, as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1
// See examples for how to use this with your own claim types
type StandardClaims struct {
   Audience  string `json:"aud,omitempty"`
   ExpiresAt int64  `json:"exp,omitempty"`
   Id        string `json:"jti,omitempty"`
   IssuedAt  int64  `json:"iat,omitempty"`
   Issuer    string `json:"iss,omitempty"`
   NotBefore int64  `json:"nbf,omitempty"`
   Subject   string `json:"sub,omitempty"`
}
```

它对应的其实是本章节中 Payload 的相关字段，这些字段都是非强制性但官方建议使用的预定义权利要求，能够提供一组有用的，可互操作的约定。

接下来在 `jwt.go`中写入第二部分代码。

```go
// 生成 JWT
func GenerateToken(appKey, appSecret string) (string, error) {
   nowTime := time.Now()
   expireTime := nowTime.Add(global.JWTSetting.Expire)
   claims := Claims{
      AppKey:    util.EncodeMD5(appKey),
      AppSecret: util.EncodeMD5(appSecret),
      StandardClaims: jwt.StandardClaims{
         ExpiresAt: expireTime.Unix(),
         Issuer:    global.JWTSetting.Issuer,
      },
   }
   // 根据 Claims 结构体创建 Token 实例，jwt.NewWithClaims() 包含两个形参
   // SigningMethod，其包含 SigningMethodHS256、SigningMethodHS384、SigningMethodHS512 三种 crypto.Hash 加密算法的方案
   // 第二个参数为 Claims 主要用于传递用户所预定义的一些权限要求，方便后续的加密、校验
   tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
   // SignedString() 生成签名后的 token 字符串
   token, err := tokenClaims.SignedString(GetJWTSecret())
   return token, err
}
```

在 `GenerateToken` 方法中，它承担了整个流程中比较重要的职责，也就是生成 JWT Token 的行为，主体的函数流程逻辑是根据客户端传入的` AppKey `和 `AppSecret `以及在项目配置中所设置的签发者（Issuer）和过期时间（`ExpiresAt`），根据指定的算法生成签名后的 Token。这其中涉及两个的内部方法，如下：

- `jwt.NewWithClaims`：根据 Claims 结构体创建 Token 实例，它一共包含两个形参，第一个参数是 `SigningMethod`，其包含 SigningMethodHS256、SigningMethodHS384、SigningMethodHS512 三种 `crypto.Hash `加密算法的方案。第二个参数是 Claims，主要是用于传递用户所预定义的一些权利要求，便于后续的加密、校验等行为。
- `tokenClaims.SignedString`：生成签名字符串，根据所传入 Secret 不同，进行签名并返回标准的 Token。

接下来继续在` jwt.go` 文件中写入第三部分代码，如下：

```go
// 解析和校验 Token
func ParseToken(token string) (*Claims, error) {
   // jwt.ParseWithClaims() 用于解析鉴权的声明，最终返回 *Token
   tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
      return GetJWTSecret(), nil
   })
   if err != nil {
      return nil, err
   }
   if tokenClaims != nil {
      // Token.Valid 当转换与核实 token 时填充该值
      if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
         return claims, nilgo
      }
   }
   return nil, err
}
```

```go
// A JWT Token.  Different fields will be used depending on whether you're
// creating or parsing/verifying a token.
type Token struct {
   Raw       string                 // The raw token.  Populated when you Parse a token
   Method    SigningMethod          // The signing method used or to be used
   Header    map[string]interface{} // The first segment of the token
   Claims    Claims                 // The second segment of the token
   Signature string                 // The third segment of the token.  Populated when you Parse a tokengo
   Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
}
```

在 `ParseToken `方法中，它主要的功能是解析和校验 Token，承担着与 `GenerateToken` 相对的功能，其函数流程主要是解析传入的 Token，然后根据 Claims 的相关属性要求进行校验。这其中涉及两个的内部方法，如下：

- `ParseWithClaims`：用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回 `*Token`。
- `Valid`：验证基于时间的声明，例如：过期时间（`ExpiresAt`）、签发者（Issuer）、生效时间（Not Before），需要注意的是，如果在令牌中没有任何声明，仍然会被认为是有效的。

至此完成了 JWT 令牌的生成、解析、校验的方法编写，在后续的应用中间件中对其进行调用，使其能够在应用程序中将一整套的动作给串联起来。

### 6. 获取 JWT 令牌

#### a. 新建 model 方法

修改`internal/model` 下的 `auth.go `文件。

```go
// 通过传入的 app_key 和 app_secret 获取认证信息
func (a Auth) Get(db *gorm.DB) (Auth, error) {
   var auth Auth
   db = db.Where("app_key = ? AND app_secret = ? AND is_del = ?", a.AppKey, a.AppSecret, 0)
   err := db.First(&auth).Error
   if err != nil && err != gorm.ErrRecordNotFound {
      return auth, err
   }
   return auth, nil
}
```

上述方法主要是用于服务端在获取客户端所传入的 app_key 和 app_secret 后，根据所传入的认证信息进行获取，以此判别是否真的存在这一条数据。

#### b. 新建 dao 方法

在 `internal/dao` 下新建`auth.go `文件，并编写针对获取认证信息的方法。

```go
package dao

import "demo/ch02/internal/model"

func (d *Dao) GetAuth(appKey, appSecret string) (model.Auth, error) {
   auth := model.Auth{AppKey: appKey, AppSecret: appSecret}
   return auth.Get(d.engine)
}
```

#### c. 新建 service 方法

在 `internal/service` 下新建`auth.go `文件，针对一些相应的基本逻辑进行处理。

```go
package service

import "errors"

type AuthRequest struct {
   AppKey    string `form:"app_key" binding:"required"`
   AppSecret string `form:"app_secret" binding:"required"`
}

func (svc *Service) CheckAuth(param *AuthRequest) error {
   auth, err := svc.dao.GetAuth(param.AppKey, param.AppSecret)
   if err != nil {
      return err
   }
   if auth.ID > 0 {
      return nil
   }
   return errors.New("auth info does not exist.")
}
```

在上述代码中，声明了 `AuthRequest` 结构体用于接口入参的校验，`AppKey `和 `AppSecret` 都设置为了必填项，在 `CheckAuth` 方法中，使用客户端所传入的认证信息作为筛选条件获取数据行，以此根据是否取到认证信息 ID 来进行是否存在的判定。

#### d. 新增路由方法

在 `internal/routers/api` 在新建`auth.go `文件。

```go
package api

import (
   "demo/ch02/global"
   "demo/ch02/internal/service"
   "demo/ch02/pkg/app"
   "demo/ch02/pkg/errcode"
   "github.com/gin-gonic/gin"
)

func GetAuth(c *gin.Context) {
   // 入参绑定与校验
   param := service.AuthRequest{}
   response := app.NewResponse(c)
   valid, errs := app.BindAndValid(c, &param)
   if !valid {
      global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
      response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
      return
   }

   // 判断认证信息
   svc := service.New(c.Request.Context())
   err := svc.CheckAuth(&param)
   if err != nil {
      global.Logger.Errorf(c, "svc.CheckAuth err: %v", err)
      response.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
      return
   }

   // 生成 token
   token, err := app.GenerateToken(param.AppKey, param.AppSecret)
   if err != nil {
      global.Logger.Errorf(c, "app.GenerateToken err: %v", err)
      response.ToErrorResponse(errcode.UnauthorizedTokenGenerate)
      return
   }

   // 返回生成的 token
   response.ToResponse(gin.H{
      "token": token,
   })
}
```

这块的逻辑主要是校验及获取入参后，绑定并获取到的 `app_key` 和 `app_secrect `进行数据库查询，检查认证信息是否存在，若存在则进行 Token 的生成并返回。

接下来修改 `internal/routers` 的 `router.go `文件，新增`auth`路由。至此，就完成了获取 Token 的整套流程。

```go
package routers

import (
   ...
)

func NewRouter() *gin.Engine {
   ...
   // 新增 auth 相关路由
   r.POST("/auth", api.GetAuth)
   ...
}
```

#### e. 接口验证

![image-20220515223328408](https://raw.githubusercontent.com/tonshz/test/master/img/202205152233457.png)

![image-20220515223350469](https://raw.githubusercontent.com/tonshz/test/master/img/202205152233508.png)

### 7. 处理应用中间件

#### a. 编写 JWT 中间件

在完成了获取 Token 的接口后，能获取 Token 了，但是对于其它的业务接口，它还没产生任何作用。涉及特定类别的接口统一处理，那必然是选择应用中间件的方式，接下来在`internal/middleware` 下新建 `jwt.go `文件，写入如下代码：

```go
package middleware

import (
   "demo/ch02/pkg/app"
   "demo/ch02/pkg/errcode"
   "github.com/dgrijalva/jwt-go"
   "github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
   return func(c *gin.Context) {
      var (
         token string
         ecode = errcode.Success
      )
      // 获取 token
      if s, exist := c.GetQuery("token"); exist {
         token = s
      } else {
         token = c.GetHeader("token")
      }
      if token == "" {
         ecode = errcode.InvalidParams
      } else {
         // ParseToken() 解析 token
         _, err := app.ParseToken(token)
         if err != nil {
            switch err.(*jwt.ValidationError).Errors {
            case jwt.ValidationErrorExpired:
               ecode = errcode.UnauthorizedTokenTimeout
            default:
               ecode = errcode.UnauthorizedTokenError
            }
         }
      }
      if ecode != errcode.Success {
         response := app.NewResponse(c)
         response.ToErrorResponse(ecode)
         /*
            Abort() 可防止调用挂起的处理程序
            请注意，这不会停止当前处理程序
            假设您有一个授权中间件来验证当前请求是否已获得授权
            如果授权失败（例如：密码不匹配）
            请调用 Abort 以确保不调用此请求的其余处理程序
         */
         c.Abort()
         return
      }
      c.Next()
   }
}
```

在上述代码中，通过` GetHeader `方法从 Header 中获取 token 参数，并调用` ParseToken` 对其进行解析，再根据返回的错误类型进行断言判定。

#### b. 接入 JWT 中间件

在完成了 JWT 的中间件编写后，需要将其接入到应用流程中，但是需要注意的是，并非所有的接口都需要用到 JWT 中间件，因此需要利用 gin 中的分组路由的概念，只针对 apiv1 的路由分组进行 JWT 中间件的引用，也就是只有 apiv1 路由分组里的路由方法会受此中间件的约束，修改`internal/routers`下的`router.go`。

```go
package routers

import (
   ...
)

func NewRouter() *gin.Engine {
   ...
   // 使用路由组设置访问路由的统一前缀 e.g. /api/v1
   // 此处定义了一个路由组 /api/v1
   apiv1 := r.Group("/api/v1")
   // apiv1 路由分组引入 JWT 中间件
   apiv1.Use(middleware.JWT())
   // 上面花括号是代表中间的语句属于一个空间内，不受外界干扰，可去掉
   {
      ...
   }
   return r
}
```

#### c. 验证接口

##### 没有获取 Token

![image-20220515225527386](https://raw.githubusercontent.com/tonshz/test/master/img/202205152255448.png)

##### Token 错误

![image-20220515230152513](https://raw.githubusercontent.com/tonshz/test/master/img/202205152301590.png)

##### Token 超时

![image-20220515230342431](https://raw.githubusercontent.com/tonshz/test/master/img/202205152303507.png)

### 8. 小结

通过本章节的学习，可以得知 JWT 令牌的内容是非严格加密的，大体上只是进行` base64UrlEncode `的处理，也就是对 JWT 令牌机制有一定了解的人可以进行反向解密，可以编写 base64 的解码代码，也可通过`jwt.io`网站直接进行解码。首先先调用 `/auth` 接口获取一个全新 token，例如：

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBfa2V5IjoiMjEyMzJmMjk3YTU3YTVhNzQzODk0YTBlNGE4MDFmYzMiLCJhcHBfc2VjcmV0IjoiNjgyYjU1NGRiYmQ5NGE3NDQ0NDU5NDJlOGMyZDk3Y2YiLCJleHAiOjE2NTI2MjcyMTIsImlzcyI6ImJsb2ctc2VydmljZSJ9.siQ-JLv3PZGUtn5OvLGzTOTV69PhkHCrmn1zfwb0dKE"
}
```

接下来针对新获取的 Token 值，只需要手动复制中间那一段（也就是 Payload），然后编写一个测试 Demo 来进行 base64 的解码，Demo 代码如下：

```go
func main() {
    payload, _ := base64.StdEncoding.DecodeString("eyJhcHBfa....DM5MTcsImlzcyI6ImJsb2ctc2VydmljZSJ9")
    fmt.Println(string(payload))
}
```

最终的输出结果如下：

```bash
{"app_key":"21232f297a57a5a743894a0e4a801fc3","app_secret":"682b554dbbd94a744445942e8c2d97cf","exp":1652627212,"iss":"blog-service"}
```

可以看到，假设有人拦截到 Token 后，是可以通过 Token 去解密并获取到 Payload 信息，也就是至少在 Payload 中不应该明文存储重要的信息，若非要存，就必须要进行不可逆加密，这样子才可以确保一定的安全性。

同时也可以发现，过期时间（`ExpiresAt`）是存储在 Payload 中的，也就是 JWT 令牌一旦签发，在没有做特殊逻辑的情况下，过期时间是不可以再度变更的，因此务必根据自己的实际项目情况进行设计和思考。

---------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)



