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
		authGroup := goodRouters.Group("")
		authGroup.Use(middle.AuthMiddleware())
		authGroup.GET("/detail", good.GoodControllerr{}.GetGoodDetail)
		authGroup.POST("/good", good.GoodControllerr{}.GoodAdd)

	}
}
