package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

// 声明了 LimiterIface 接口，用于定义当前限流器所必须要的方法
// 由于限流器的策略不同，需要定义通用的接口以保证接口的设计
type LimiterIface interface {
	// 获取对应限流器的键值对名称
	Key(c *gin.Context) string
	// 获取令牌桶
	GetBucket(key string) (*ratelimit.Bucket, bool)
	// 新增多个令牌桶
	AddBucket(rules ...LimiterBucketRule) LimiterIface
}

type Limiter struct {
	// 存储令牌桶和键值对名称的映射关系
	limiterBuckets map[string]*ratelimit.Bucket
}

// 定义 LimiterBucketRule 结构体用于存储令牌桶的规则属性
type LimiterBucketRule struct {
	// 自定义键值对名称
	Key string
	// 间隔多长时间释放 N 个令牌
	FillInterval time.Duration
	// 令牌桶的容量
	Capacity int64
	// 每次到达间隔时间后所释放的具体令牌数量
	Quantum int64
}
