package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"miaosha-system/global"
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
		log.Fatalf("配置信息未初始化")
	}
	if global.Config.JWT.SignKey == "" {
		log.Fatalf("JWT 签名密钥未配置")
	}
	SigningKey = []byte(global.Config.JWT.SignKey)
	ExpireTime = time.Duration(global.Config.JWT.ExpireTime) * time.Second
	fmt.Println("JWT 配置初始化成功")
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
