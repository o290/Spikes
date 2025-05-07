package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"miaosha-system/common/req"
	"miaosha-system/global"
	"miaosha-system/model"
	"miaosha-system/utils"
	"net/http"
)

type UserController struct {
}

func (m UserController) Login(c *gin.Context) {
	var req req.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		global.Log.Fatalln("解析json失败")
		c.JSON(http.StatusOK, gin.H{
			"message": "login error",
		})
		return
	}
	var u model.UserModel
	if err := global.DB.Where("id =?", req.ID).First(&u).Error; err != nil {
		global.Log.Println("不存在该用户,无法直接登录")
		c.JSON(http.StatusOK, gin.H{
			"message": "不存在该用户，请前往注册",
		})
		return
	}

	if !utils.CheckPwd(u.Pwd, req.Password) {
		global.Log.Println("用户密码错误")
		c.JSON(http.StatusOK, gin.H{
			"message": "用户名或密码错误",
		})
		return
	}

	//重定向
	c.JSON(200, gin.H{
		"message":  "登录成功",
		"nickname": u.Nickname,
	})
	return
}
