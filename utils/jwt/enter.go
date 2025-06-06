package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"miaosha-system/global"
	"net/http"
	"time"
)

// 定义 JWT 的签名密钥和过期时间
var (
	SigningKey []byte
	ExpireTime time.Duration
)

// InitJWT 初始化 JWT 配置
func InitJWT() {
	if global.Config == nil {
		global.Log.Fatalf("配置信息未初始化")
	}
	if global.Config.JWT.SignKey == "" {
		global.Log.Fatalf("JWT 签名密钥未配置")
	}
	SigningKey = []byte(global.Config.JWT.SignKey)
	ExpireTime = time.Duration(global.Config.JWT.ExpireTime) * time.Minute
	global.Log.Info("JWT 配置初始化成功")
}

// Claims 定义 JWT 的 Claims 结构体
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint) (string, error) {
	if SigningKey == nil || ExpireTime == 0 {
		return "", errors.New("JWT 配置未初始化")
	}
	// 创建一个新的 Claims
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ExpireTime).Unix(),
			Issuer:    global.Config.JWT.Issuer,
		},
	}

	// 创建一个新的 JWT Token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用签名密钥对 Token 进行签名
	tokenString, err := token.SignedString(SigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	if SigningKey == nil {
		return nil, errors.New("JWT 配置未初始化")
	}
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GetUserID(c *gin.Context) (uint, error) {
	//从token中解析出的userid
	token := c.GetHeader("Authorization")
	if token == "" {
		global.Log.Info("无效token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未提供 Token",
		})
		return 0, errors.New("获取id失败")
	}
	// 去除 Token 前缀 "Bearer "
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// 解析 Token
	claims, err := ParseToken(token)
	if err != nil {
		global.Log.Info("token解析失败")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "无效的 Token",
		})
		return 0, errors.New("获取id失败")
	}
	return claims.UserID, nil
}
