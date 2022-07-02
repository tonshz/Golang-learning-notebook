package routers

import (
	"demo/ch02/global"
	"demo/ch02/internal/middleware"
	"demo/ch02/internal/routers/api"
	v1 "demo/ch02/internal/routers/api/v1"
	"demo/ch02/pkg/limiter"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"

	// 注意此处导包，不需要设置别名
	"github.com/swaggo/gin-swagger/swaggerFiles"
	// 初始化 docs 包
	_ "demo/ch02/docs"
)

// 对认证接口进行限流
var methodLimiters = limiter.NewMethodLimiter().AddBucket(limiter.LimiterBucketRule{
	Key:          "/auth",          // 对接受的 auth 请求进行限流
	FillInterval: 10 * time.Second, // 对 10s 内接受的请求进行限流
	Capacity:     10,               // 10s 内最多接受 10 个
	Quantum:      10,               // 每过 10s 将令牌桶的数量减少 10 #{Quantum}
})

func NewRouter() *gin.Engine {
	// 或者 r := gin.Default() 也可
	r := gin.New()
	// 根据不同的部署环境进行了应用中间件的设置
	// 使用了自定义的 Logger 和 Recovery 后就不需要使用 gin 原生提供的组件了
	if global.ServerSetting.RunMode == "debug" {
		// 在本地开发环境中，可能没有全应用生态圈，故作特殊处理
		r.Use(gin.Logger())
		// 在注册顺序上需要注意，类似 Recovery 这类的中间件应当尽可能早的注册
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}

	// 新增链路追踪中间件注册
	r.Use(middleware.Tracing())
	// 新增限流控制中间件的注册
	r.Use(middleware.RateLimiter(methodLimiters))
	// 新增统一超时控制中间件注册
	r.Use(middleware.ContextTimeout(global.AppSetting.DefaultContextTimeout))
	// 新增中间件 Translations 的注册
	r.Use(middleware.Translations())
	// 新增应用信息中间件注册
	r.Use(middleware.AppInfo())

	// 手动指定当前应用所启动的 swagger/doc.json 路径
	//url := ginSwagger.URL("http://127.0.0.1:8000/swagger/doc.json")
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	// 注册一个针对 swagger 的路由，默认指向当前应用所启动的域名下的 swagger/doc.json 路径
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	article := v1.NewArticle()
	tag := v1.NewTag()
	// 添加上传文件的对应路由
	upload := api.NewUpload()
	r.POST("/upload/file", upload.UploadFile)
	// 提供静态资源的访问
	// Static 只能展示文件，StaticFS 可以连目录也展示
	// http.Dir() 实现了 FileSystem接口，利用本地目录实现一个文件系统
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	// 新增 auth 相关路由
	r.POST("/auth", api.GetAuth)

	// 使用路由组设置访问路由的统一前缀 e.g. /api/v1
	// 此处定义了一个路由组 /api/v1
	apiv1 := r.Group("/api/v1")
	// apiv1 路由分组引入 JWT 中间件
	apiv1.Use(middleware.JWT())
	// 上面花括号是代表中间的语句属于一个空间内，不受外界干扰，可去掉
	{
		// 将实现的 Handler 方法注册到对应的路由规则上
		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.GET("/tags", tag.List)

		apiv1.POST("/articles", article.Create)
		apiv1.DELETE("/articles/:id", article.Delete)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("/articles/:id/state", article.Update)
		apiv1.GET("/articles/:id", article.Get)
		apiv1.GET("/articles", article.List)
	}
	return r
}
