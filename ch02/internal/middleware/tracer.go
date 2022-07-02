package middleware

import (
	"context"
	"demo/ch02/global"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

// 返回入参为一个 gin 上下文的函数
func Tracing() func(c *gin.Context) {
	return func(c *gin.Context) {
		var newCtx context.Context
		var span opentracing.Span
		// Extract() 返回给定 `format` 和 `carrier` 的 SpanContext 实例。
		spanCtx, err := opentracing.GlobalTracer().Extract(
			// HTTPHeaders 将 SpanContexts 表示为 HTTP 标头字符串对。
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header))

		if err != nil {
			/*
				StartSpanFromContextWithTracer
				使用在上下文中找到的跨度作为 ChildOfRef 启动并返回一个带有 `operationName` 的跨度。
				如果不存在，它会创建一个根跨度。它还返回一个围绕返回的范围构建的 context.Context 对象。
			*/
			span, newCtx = opentracing.StartSpanFromContextWithTracer(
				c.Request.Context(),
				global.Tracer,
				c.Request.URL.Path)
		} else {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(
				c.Request.Context(),
				global.Tracer,
				// 相对路径
				c.Request.URL.Path,
				// ChildOf 返回一个指向依赖父跨度的 StartSpanOption
				opentracing.ChildOf(spanCtx),
				// Tag 可以作为 StartSpanOption 传递以将标签添加到新 Span，或者其 Set 方法可用于将标签应用于现有 Span
				/*
					StartSpanOptions 允许 Tracer.StartSpan() 调用者和实现者使用一种机制来覆盖开始时间戳、
					指定 Span 引用并在 Span 开始时使单个标签或多个标签可用。
				*/
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"})
		}

		defer span.Finish()
		// 添加链路 TraceID 与 SpanID 的记录
		var traceID string
		var spanID string
		var spanContext = span.Context()
		// 此方法为类型断言
		// interface{}.(type) 使用类型断言断定某个接口是否是指定的类型
		// 该写法必须与switch case联合使用，case中列出实现该接口的类型。
		switch spanContext.(type) {
		case jaeger.SpanContext:
			// 将 spanContext 转换为 jaeger.SpanContext 类型并赋值给 jaegerContext
			// jaegerContext, flag := spanContext.(jaeger.SpanContext)
			// 其中 flag 为是否转换成功，该参数可省略
			jaegerContext := spanContext.(jaeger.SpanContext)
			traceID = jaegerContext.TraceID().String()
			spanID = jaegerContext.SpanID().String()
		}

		c.Set("X-Trace-ID", traceID)
		c.Set("X-Span-ID", spanID)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
