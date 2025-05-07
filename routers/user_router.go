package routers

import (
	"github.com/gin-gonic/gin"
	"miaosha-system/controller/user"
)

func UserRouterInit(r *gin.Engine) {
	userRouters := r.Group("/user")
	{
		userRouters.POST("/login", user.UserController{}.Login)
		//userRouters.POST("/register",user.UserController{}.)
	}
}
