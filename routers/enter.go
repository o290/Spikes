package routers

import (
	"github.com/gin-gonic/gin"
	"miaosha-system/controller/good"
	"miaosha-system/controller/order"
	"miaosha-system/mq"
)

func initService() {
	defer mq.Init()
	order.OrderInit()
	good.GoodInit()
	good.GoodControllerr{}.Init()
}
func Init() (r *gin.Engine) {
	initService()
	r = gin.Default()
	UserRouterInit(r)
	GoodRouterInit(r)
	OrderRouterInit(r)
	return
}
