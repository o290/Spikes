package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"miaosha-system/core"
	"miaosha-system/global"
	"strconv"
	"time"
)

const messageQueueKey = "message_que"

// 锁的过期时间，避免死锁
const lockExpiration = 10 * time.Second

// Producer 生产者函数，将用户的秒杀请求加入消息队列
func Producer(userID, goodID int) {
	message := fmt.Sprintf("%d:%d", userID, goodID)
	err := global.Redis.RPush(context.Background(), messageQueueKey, message).Err()
	if err != nil {
		log.Printf("用户 %d 发起商品 %d 的秒杀请求时，加入队列失败: %v", userID, goodID, err)
		return
	}
	log.Printf("用户 %d 发起商品 %d 的秒杀请求，已加入队列", userID, goodID)
}

// Consumer 消费者函数，从消息队列中取出消息并处理
func Consumer() {
	for {
		result, err := global.Redis.BLPop(context.Background(), 10*time.Second, messageQueueKey).Result()
		if err != nil {
			if err == redis.Nil {
				// 队列中没有消息，继续等待
				continue
			}
			log.Printf("从队列中获取消息时出错: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		// 解析消息
		message := result[1]
		var userID, productID int
		fmt.Sscanf(message, "%d:%d", &userID, &productID)
		// 处理秒杀订单
		processSeckillOrder(userID, productID)
	}
}

// processSeckillOrder 处理秒杀订单，添加了分布式锁
func processSeckillOrder(userID, productID int) {
	log.Printf("开始处理用户 %d 对商品 %d 的秒杀订单", userID, productID)
	// 锁的键名，以商品 ID 为标识
	lockKey := fmt.Sprintf("product_lock:%d", productID)
	// 尝试获取分布式锁
	locked, err := global.Redis.SetNX(context.Background(), lockKey, "locked", lockExpiration).Result()
	if err != nil {
		log.Printf("获取商品 %d 的分布式锁时出错: %v", productID, err)
		return
	}
	if !locked {
		log.Printf("商品 %d 的分布式锁已被占用，用户 %d 的秒杀请求稍后重试", productID, userID)
		return
	}
	// 使用 defer 确保锁最终会被释放
	defer global.Redis.Del(context.Background(), lockKey)

	// 检查库存
	stockKey := fmt.Sprintf("product_stock:%d", productID)
	stock, err := global.Redis.Get(context.Background(), stockKey).Int64()
	if err != nil {
		if err == redis.Nil {
			log.Printf("商品 %d 库存信息不存在", productID)
		} else {
			log.Printf("获取商品 %d 库存信息时出错: %v", productID, err)
		}
		return
	}
	if stock <= 0 {
		log.Printf("商品 %d 已售罄，用户 %d 的秒杀请求失败", productID, userID)
		return
	}
	// 扣减库存
	newStock, err := global.Redis.Decr(context.Background(), stockKey).Result()
	if err != nil {
		log.Printf("扣减商品 %d 库存时出错: %v", productID, err)
		return
	}
	log.Printf("商品 %d 库存剩余: %d", productID, newStock)
	// 生成订单
	orderID := generateOrderID(userID, productID)
	log.Printf("用户 %d 对商品 %d 的秒杀订单生成成功，订单 ID: %s", userID, productID, orderID)
}

// generateOrderID 生成订单 ID
func generateOrderID(userID, productID int) string {
	// 简单示例，实际应用中可使用更复杂的方式生成唯一订单 ID
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	return fmt.Sprintf("%s_%d_%d", timestamp, userID, productID)
}

func main() {
	// 初始化配置
	core.InitConfig()
	// 初始化日志
	global.Log = core.InitLogger()
	// 连接数据库
	global.DB = core.Initgorm()
	global.Redis = core.InitRedis(global.Config.Redis.Addr, global.Config.Redis.Pwd, global.Config.Redis.DB)
	productID := 2
	stockKey := fmt.Sprintf("product_stock:%d", productID)
	err := global.Redis.Set(context.Background(), stockKey, 5, 0).Err()
	if err != nil {
		log.Fatalf("初始化商品 %d 库存时出错: %v", productID, err)
	}

	// 启动消费者 goroutine
	go Consumer()

	// 模拟多个用户发起秒杀请求
	for i := 1; i <= 10; i++ {
		Producer(i, productID)
	}

	// 保持主 goroutine 运行
	time.Sleep(100 * time.Second)
}
