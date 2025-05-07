package mq

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"miaosha-system/global"
	"miaosha-system/inter"
	"strconv"
	"time"
)

// CloseOderMQ 创建订单消息队列
type CloseOrderMQ struct {
	Order inter.Order
}

var CloseMQ *CloseOrderMQ

const closeKey = "timeout_close"

func (m *CloseOrderMQ) Send(orderID string) error {
	err := global.Redis.ZAdd(context.Background(), closeKey, redis.Z{
		Score:  float64(time.Now().Unix() + 60*4),
		Member: orderID,
	}).Err()
	if err != nil {
		global.Log.Errorf("订单%v加入延迟队列失败", err)
	} else {
		global.Log.Infof("订单%s加入延迟队列成功", orderID)
	}
	return err
}

// Remove 移除订单
func (m *CloseOrderMQ) Remove(orderID string) {
	err := global.Redis.ZRem(context.Background(), closeKey, orderID).Err()
	if err != nil {
		global.Log.Errorf("订单%v移除延迟队列失败", err)
	} else {
		global.Log.Infof("订单%s移除延迟队列成功", orderID)
	}
	return
}

// Receive 处理订单关闭
func (m *CloseOrderMQ) Receive() {
	for {
		//获取分数（过期时间）在 0 到当前时间之间的订单 ID
		list, err := global.Redis.ZRangeByScore(context.Background(), closeKey, &redis.ZRangeBy{
			Min:    "0",
			Max:    strconv.FormatInt(time.Now().Unix(), 10),
			Offset: 0,
			//获取一个
			Count: 1,
		}).Result()
		if err != nil {
			log.Printf("receive处理错误: %v", err)
			continue
		}
		if len(list) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		//获取订单信息
		orderInfo, err1 := m.Order.GetOrderInfo(list[0])
		fmt.Println(list)
		if err1 != nil {
			global.Log.Errorf("获取订单信息失败: %v", err)
			return
		}
		//订单关闭处理
		//判断订单状态
		if orderInfo.Status != 2 {
			global.Log.Infof("订单%s超时，需移除超时队列", orderInfo.OrderNumber)
			//修改数据库库存
			err = m.Order.CloseUpdateStock(orderInfo)
			if err != nil {
				return
			}

			//删除订单信息缓存
			err = global.Redis.Del(context.Background(), fmt.Sprintf("order:%d:%d", orderInfo.UserID, orderInfo.GoodID)).Err()
			if err != nil {
				global.Log.Printf("删除订单信息缓存【%s】失败: %v", orderInfo.OrderNumber, err)
				return
			}
			m.Remove(orderInfo.OrderNumber)
		}

		return
	}

}
