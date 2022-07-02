package middleware

import (
	"demo/ch02/pkg/app"
	"demo/ch02/pkg/errcode"
	"demo/ch02/pkg/limiter"
	"github.com/gin-gonic/gin"
)

// 入参为 LimiterIface 接口类型，这样只要符合该接口类型的具体限流器实现都可以传入使用
func RateLimiter(l limiter.LimiterIface) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := l.Key(c)
		if bucket, ok := l.GetBucket(key); ok {
			/*
				TakeAvailable()
				占用存储桶中立即可用的令牌的数量。
				返回值为删除的令牌数，
				如果没有可用的令牌，则返回零。它不会阻塞。
			*/
			count := bucket.TakeAvailable(1)
			if count == 0 {
				response := app.NewResponse(c)
				response.ToErrorResponse(errcode.TooManyRequests)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
