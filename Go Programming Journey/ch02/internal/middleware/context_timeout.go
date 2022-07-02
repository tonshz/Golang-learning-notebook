package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func ContextTimeout(t time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		/*
			before c.Next()
		*/
		// 使用 context.WithTimeout() 设置当前 context 的超时时间
		// 返回值 Context, CancelFunc，即 cancel 值是一个取消函数
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer cancel()
		// 将设置了超时时间的 ctx 赋值给当前的 context
		c.Request = c.Request.WithContext(ctx)
		// c.Next() 使用后会使得在其之后的代码执行在 main() 之后
		c.Next()
		/*
			after c.Next()
			此处的代码在 main() 中代码处理完成后再执行
			即执行顺序为： before => main => after
			若不使用 c.Next() 执行顺序为： before => after => main
		*/
	}
}
