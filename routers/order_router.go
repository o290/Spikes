package routers

import (
	"github.com/gin-gonic/gin"
	"miaosha-system/controller/order"
)

func OrderRouterInit(r *gin.Engine) {
	orderRouters := r.Group("/order")
	{
		//orderRouters.Use(middle.Auth())
		orderRouters.GET("/list", (&order.OrderController{}).GetOrderList)
		orderRouters.POST("/spikes", (&order.OrderController{}).Spikes)
		orderRouters.GET("/detail", (&order.OrderController{}).GetOrderInfo)
		orderRouters.POST("/close", (&order.OrderController{}).CloseOrder)
	}
}
