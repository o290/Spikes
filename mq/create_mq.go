package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"miaosha-system/common/msg"
	"miaosha-system/global"
	"miaosha-system/inter"
	"miaosha-system/utils/lock"
	"time"
)

// CreateOderMQ 创建订单消息队列
type CreateOderMQ struct {
	Order inter.Order
}

var CreateMQ *CreateOderMQ

func (m *CreateOderMQ) Send(msg msg.CreateMsg) error {
	data, _ := json.Marshal(msg)
	err := global.Redis.LPush(context.Background(), "create_mq", data).Err()
	return err
}
func (m *CreateOderMQ) Receive() {
	for {
		result, err := global.Redis.BRPop(context.Background(), 2*time.Second, "create_mq").Result()
		if err != nil {
			if err == redis.Nil {
				//fmt.Println("创建订单消息队列中暂无订单，等待中...")
			} else {
				global.Log.Printf("处理订单时出错: %v", err)
			}
			time.Sleep(1 * time.Second)
			continue
		}
		orderInfo := result[1]
		var msg msg.CreateMsg
		if err = json.Unmarshal([]byte(orderInfo), &msg); err != nil {
			global.Log.Printf("json.Unmarshal() failed, err: %v", err)
			continue
		}
		global.Log.Printf("开始处理用户%d，商品%d的订单", msg.UserID, msg.GoodID)
		//创建订单
		m.Order.CreateOrder(msg.UserID, msg.GoodID)
		fmt.Printf("用户%d创建商品%d订单\n", msg.UserID, msg.GoodID)
		//获取分布式锁值
		value, _ := global.Redis.Get(context.Background(), fmt.Sprintf("lock:%d", msg.GoodID)).Result()
		//释放分布式锁
		_, err = lock.ReleaseLock(context.Background(), fmt.Sprintf("lock:%d", msg.GoodID), value)
		if err != nil {
			fmt.Printf("用户 %d 释放锁失败: %v\n", msg.UserID, err)
		}
		lock.StopRenew()
		fmt.Printf("用户 %d 释放锁成功\n", msg.UserID)

		global.Log.Printf("用户%d，商品%d的创建订单处理完成", msg.UserID, msg.GoodID)
	}

}
