package mq

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"miaosha-system/global"
	"miaosha-system/inter"
	"miaosha-system/model"
	"time"
)

//后台异步更新库存到数据库中

// UpdateStockMQ 创建订单消息队列
type UpdateStockMQ struct {
	Order inter.Order
}

var StockMQ *UpdateStockMQ

const stockKey = "stock_update"
const interval = 5 * time.Second

// PeriodicUpdateStock 定期更新库存的函数
func (m *UpdateStockMQ) PeriodicUpdateStock() {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 获取所有需要更新的商品 ID
			members, err := global.Redis.ZRange(context.Background(), stockKey, 0, -1).Result()
			if err != nil {
				global.Log.Error("获取需要更新的商品id失败")
				continue
			}

			for _, member := range members {
				var goodID uint
				// 将成员转换为 uint 类型
				fmt.Sscanf(member, "%d", &goodID)

				// 获取缓存中的库存
				stock, err := global.Redis.Get(context.Background(), fmt.Sprintf("stock:%d", goodID)).Int64()
				if err != nil {
					global.Log.Error("获取库存缓存失败")
					continue
				}
				global.Log.Println("更新数据库库存成功")
				// 更新数据库中的库存
				err = global.DB.Transaction(func(tx *gorm.DB) error {

					return tx.Model(&model.GoodModel{}).Where("id=?", goodID).Update("stock", stock).Error
				})
				if err != nil {
					global.Log.Error("更新缓存失败")
				}

				// 更新成功后，从有序集合中移除该商品 ID
				global.Redis.ZRem(context.Background(), stockKey, goodID)
			}
		case <-context.Background().Done():
			return
		}
	}
}
