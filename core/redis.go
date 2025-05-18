package core

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"miaosha-system/global"
	"strconv"
	"time"
)

var rdb *redis.Client

// Background返回一个非空context，它永远不会取消，没有值
// 也没有期限，通常在main函数，初始化和测试时使用，并用作
// 传入请求的顶级上下文
var ctx = context.Background()

func InitRedis(addr, pwd string, db int) (client *redis.Client) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,  //默认db为0，指的是连接到哪一个数据库
		PoolSize: 100, //连接池大小
	})
	_, cancel := context.WithTimeout(ctx, 500*time.Microsecond)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
		return
	}
	global.Log.Println("初始化连接redis成功")
	return rdb
}

func UpdatePublish() {
	// 订阅更新缓存的消息
	pubsub := global.Redis.Subscribe(context.Background())
	defer pubsub.Close()

	// 监听所有商品的更新缓存消息
	channelPattern := "update_cache:*"
	err := pubsub.PSubscribe(context.Background(), channelPattern)
	if err != nil {
		global.Log.Fatalf("订阅更新缓存消息时出错: %v", err)
	}

	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			global.Log.Printf("接收更新缓存消息时出错: %v", err)
			continue
		}

		// 解析商品 ID
		var productID int
		fmt.Sscanf(msg.Channel, "update_cache:%d", &productID)

		// 解析新的库存数量
		newStock, err := strconv.Atoi(msg.Payload)
		if err != nil {
			global.Log.Printf("解析库存数量时出错: %v", err)
			continue
		}

		// 更新缓存
		stockKey := fmt.Sprintf("stock:%d", productID)
		err = global.Redis.Set(context.Background(), stockKey, newStock, 0).Err()
		if err != nil {
			global.Log.Printf("更新商品 %d 缓存库存时出错: %v", productID, err)
		} else {
			global.Log.Printf("商品 %d 缓存库存更新成功，新库存: %d", productID, newStock)
		}
	}
}
