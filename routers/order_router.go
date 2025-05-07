package routers

import (
	"github.com/gin-gonic/gin"
	"miaosha-system/controller/order"
)

func OrderRouterInit(r *gin.Engine) {
	orderRouters := r.Group("/order")
	{
		//秒杀前要限流
		orderRouters.POST("/spikes", (&order.OrderController{}).Spikes)
	}
}
