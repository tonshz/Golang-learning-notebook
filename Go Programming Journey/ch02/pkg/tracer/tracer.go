package tracer

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

func NewJaegerTrace(serviceName, agentHostPort string) (opentracing.Tracer, io.Closer, error) {
	// config.Configuration 为 jaeger client 的配置项，主要设置应用的基本信息
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const", // 固定采样，对所有数据都进行采样
			Param: 1,       // 1 表示 true,0 表示 false
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,            // 是否启用 LoggingReporter
			BufferFlushInterval: 1 * time.Second, // 刷新缓冲区的频率
			LocalAgentHostPort:  agentHostPort,   // 上报的 Agent 地址
		},
	}
	// cfg.NewTracer() 根据配置项初始化 Tracer 对象，返回 opentracing.Tracer，而不是特定供应商的追踪系统对象
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}
	// opentracing.SetGlobalTracer() 设置全局的 Tracer 对象
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}
