package middleware

import "github.com/gin-gonic/gin"

func AppInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("app_name", "blog-service")
		c.Set("app_version", "1.0.0")
		/*
			Next 应该只在中间件内部使用。
			它执行调用处理程序内链中的待处理处理程序。
			请参阅 GitHub 中的示例
		*/
		c.Next()
	}
}
