package middle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"miaosha-system/utils/jwt"
	"net/http"
)

// AuthMiddleware 验证 Token 的中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Token
		tokenString := c.GetHeader("Authorization")
		fmt.Println(tokenString)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "未提供 Token",
			})
			c.Abort()
			return
		}

		// 去除 Token 前缀 "Bearer "
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// 解析 Token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "无效的 Token",
			})
			c.Abort()
			return
		}

		// 将用户 ID 存储到上下文中，方便后续处理
		fmt.Println("userid", claims.UserID)
		c.Set("user_id", claims.UserID)

		// 继续处理请求
		c.Next()
	}
}
