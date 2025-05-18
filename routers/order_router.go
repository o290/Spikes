package routers

import (
	"github.com/gin-gonic/gin"
	"miaosha-system/controller/order"
	"miaosha-system/middle"
	"time"
)

func OrderRouterInit(r *gin.Engine) {
	orderRouters := r.Group("/order")
	{
		// 普通限流,每秒 100 请求
		orderRouters.Use(middle.RateLimitMiddleware(time.Second, 100))
		orderRouters.GET("/list", (&order.OrderController{}).GetOrderList)
		orderRouters.GET("/detail", (&order.OrderController{}).GetOrderInfo)
		orderRouters.POST("/close", (&order.OrderController{}).CloseOrder)

		// 秒杀单独限流,每秒 10 请求
		spikesRouters := orderRouters.Group("/spikes")
		spikesRouters.Use(middle.RateLimitMiddleware(100*time.Millisecond, 20)) // 每秒 10 令牌（100ms/个）
		spikesRouters.POST("", (&order.OrderController{}).Spikes)
	}
}
