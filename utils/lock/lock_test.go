package lock

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"miaosha-system/global"
	"sync"
	"testing"
	"time"
)

// 初始化 Redis 客户端
func init() {
	global.Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // 根据实际情况修改 Redis 地址和端口
		Password: "123456",
		DB:       1,
	})
}

func TestSeckillConcurrency(t *testing.T) {
	ctx := context.Background()
	// 假设锁的过期时间
	expiration := 5 * time.Second
	// 假设参与秒杀的用户数量
	numUsers := 5

	var wg sync.WaitGroup
	wg.Add(numUsers)

	for i := 1; i <= numUsers; i++ {
		go func(userID uint) {
			defer wg.Done()
			key := fmt.Sprintf("lock:%d", 1)
			value, acquired, err := AcquireLock(ctx, key, expiration, userID, 1)
			if err != nil {
				fmt.Printf("用户 %d 获取锁失败: %v\n", userID, err)
				return
			}

			if acquired {
				fmt.Printf("用户 %d 获取锁成功\n", userID)
				defer func() {
					_, err := ReleaseLock(ctx, key, value)
					if err != nil {
						fmt.Printf("用户 %d 释放锁失败: %v\n", userID, err)
					}
					StopRenew()
					fmt.Printf("用户 %d 释放锁成功\n", userID)
				}()

				// 模拟检查库存和扣减库存的逻辑
				stock, err := global.Redis.Get(ctx, fmt.Sprintf("stock:%d", 1)).Int()
				if err != nil {
					if err != redis.Nil {
						fmt.Printf("用户 %d 获取库存失败: %v\n", userID, err)
					}
					return
				}
				if stock > 0 {
					_, err = global.Redis.Decr(ctx, fmt.Sprintf("stock:%d", 1)).Result()
					if err != nil {
						fmt.Printf("用户 %d 扣减库存失败: %v\n", userID, err)
						return
					}
					fmt.Printf("用户 %d 秒杀成功\n", userID)
					time.Sleep(2 * time.Second)
				} else {
					fmt.Printf("用户 %d 秒杀失败，商品已售罄\n", userID)
				}
			} else {
				fmt.Printf("用户 %d 获取锁失败，可能有其他用户正在秒杀\n", userID)
			}
		}(uint(i))
	}

	wg.Wait()
	time.Sleep(time.Second * 60)
	// 检查最终库存是否正确
	stock, err := global.Redis.Get(ctx, fmt.Sprintf("stock:%d", 1)).Int()
	if err != nil {
		if err != redis.Nil {
			t.Errorf("获取最终库存失败: %v", err)
		}
	} else {
		fmt.Printf("最终商品库存: %d\n", stock)
	}
}
