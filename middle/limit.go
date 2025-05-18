package middle

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// 定义全局令牌桶管理器（支持按键区分不同限流维度）
var (
	buckets     = make(map[string]*ratelimit.Bucket)
	bucketsLock sync.Mutex
)

// RateLimitMiddleware 令牌桶限流中间件
// fillInterval: 令牌填充间隔（如 time.Second）
// capacity: 令牌桶容量
func RateLimitMiddleware(fillInterval time.Duration, capacity int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成限流键,按 IP 限流
		key := c.ClientIP()
		// 获取或创建令牌桶（并发安全）
		bucketsLock.Lock()
		bucket, exists := buckets[key]
		if !exists {
			// 新创建的桶：初始令牌数为容量值，填充间隔控制令牌生成速率
			bucket = ratelimit.NewBucketWithRate(
				float64(capacity)/fillInterval.Seconds(), // 每秒生成的令牌数（速率）
				capacity, // 桶容量
			)
			buckets[key] = bucket
		}
		bucketsLock.Unlock()

		// 尝试获取 1 个令牌
		if bucket.TakeAvailable(1) < 1 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": http.StatusTooManyRequests,
				"msg":  "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
