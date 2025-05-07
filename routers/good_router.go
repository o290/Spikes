package routers

import (
	"github.com/gin-gonic/gin"
	"miaosha-system/controller/good"
)

func GoodRouterInit(r *gin.Engine) {
	goodRouters := r.Group("/good")
	{
		goodRouters.GET("/list", good.GoodControllerr{}.GoodList)
		//goodRouters.GET("/good", good.GoodControllerr{}.GetGood)
		goodRouters.POST("/good", good.GoodControllerr{}.GoodAdd)
		goodRouters.GET("/detail/:id", good.GoodControllerr{}.GetGoodDetail)
	}
}
