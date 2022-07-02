package middleware

import (
	"bytes"
	"demo/ch02/global"
	"demo/ch02/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

// 针对访问日志的 Writer 结构体
type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 实现特定的 Write 方法
func (w AccessLogWriter) Write(p []byte) (int, error) {
	// 获取响应主体
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 初始化 AccessLogWriter
		bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		// 将其赋值给当前的 Writer 写入流
		c.Writer = bodyWriter
		beginTime := time.Now().Unix()
		c.Next()
		endTime := time.Now().Unix()
		fields := logger.Fields{
			"request":  c.Request.PostForm.Encode(),
			"response": bodyWriter.body.String(),
		}
		// 将 Context 传入日志方法中
		global.Logger.WithFields(fields).Infof(c, "access log: method: %s, status_code: %d, begin_time: %d, end_time: %d",
			c.Request.Method,
			bodyWriter.Status(),
			beginTime,
			endTime,
		)
	}
}
