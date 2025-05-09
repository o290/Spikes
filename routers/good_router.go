package routers

import (
	"github.com/gin-gonic/gin"
	"miaosha-system/controller/good"
	"miaosha-system/middle"
)

func GoodRouterInit(r *gin.Engine) {
	goodRouters := r.Group("/good")
	{
		goodRouters.GET("/test", good.GoodControllerr{}.TestToken)
		goodRouters.GET("/list", good.GoodControllerr{}.GoodList)
		goodRouters.Use(middle.AuthMiddleware())
		goodRouters.GET("/detail", good.GoodControllerr{}.GetGoodDetail)
		goodRouters.POST("/good", good.GoodControllerr{}.GoodAdd)

	}
}
