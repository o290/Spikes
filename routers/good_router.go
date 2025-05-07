package routers

import (
	"github.com/gin-gonic/gin"
	"miaosha-system/controller/good"
	"miaosha-system/middle"
)

func GoodRouterInit(r *gin.Engine) {
	goodRouters := r.Group("/good")
	{
		goodRouters.Use(middle.AuthMiddleware())
		goodRouters.GET("/test", good.GoodControllerr{}.TestToken)
		goodRouters.GET("/list", good.GoodControllerr{}.GoodList)
		//goodRouters.GET("/good", good.GoodControllerr{}.GetGood)
		goodRouters.POST("/good", good.GoodControllerr{}.GoodAdd)
		goodRouters.GET("/detail/:id", good.GoodControllerr{}.GetGoodDetail)
	}
}
