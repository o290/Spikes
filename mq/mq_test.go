package mq

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"miaosha-system/global"
	"testing"
)

// 初始化 Redis 客户端
func init() {
	global.Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // 根据实际情况修改 Redis 地址和端口
		Password: "123456",
		DB:       1,
	})
}
func TestMQ(t *testing.T) {
	m1 := CreateMsg{
		UserID: 1,
		GoodID: 1,
	}
	m2 := CreateMsg{2, 1}
	m3 := CreateMsg{1, 2}
	CreateMQ.Send(m1)
	CreateMQ.Send(m2)
	CreateMQ.Send(m3)
	fmt.Println(global.Redis.LRange(context.Background(), "create_mq", 0, -1).Result())
}
