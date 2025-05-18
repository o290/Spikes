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

// RefreshTask 定时刷新任务结构体
type RefreshTask struct {
	Good inter.Good
}

var Refresh *RefreshTask

// Start 启动定时刷新任务
func (t *RefreshTask) Start() {
	// 每 30 秒刷新一次
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 刷新库存
			t.refreshStock()
			// 刷新商品信息
			t.refreshGoodInfo()
		case <-context.Background().Done():
			return
		}
	}
}

// refreshStock 刷新库存
func (t *RefreshTask) refreshStock() {
	// 获取所有商品信息
	var goodList []model.GoodModel
	err := global.DB.Find(&goodList).Error
	if err != nil {
		global.Log.Error("获取商品列表失败", err)
		return
	}

	for _, good := range goodList {
		// 获取缓存中的库存
		stock, err := global.Redis.Get(context.Background(), fmt.Sprintf("stock:%d", good.ID)).Int64()
		if err != nil {
			global.Log.Error("获取库存缓存失败", err)
			continue
		}

		// 更新数据库中的库存
		err = global.DB.Transaction(func(tx *gorm.DB) error {
			return tx.Model(&model.GoodModel{}).Where("id=?", good.ID).Update("stock", stock).Error
		})
		if err != nil {
			global.Log.Error("更新数据库库存失败", err)
		}
	}
}

// refreshGoodInfo 刷新商品信息
func (t *RefreshTask) refreshGoodInfo() {
	// 获取所有商品信息
	var goodList []model.GoodModel
	err := global.DB.Find(&goodList).Error
	if err != nil {
		global.Log.Error("获取商品列表失败", err)
		return
	}

	for _, good := range goodList {
		// 将商品信息存入 Redis 缓存
		err = t.Good.SetGoodToRedis(good)
		if err != nil {
			global.Log.Error("更新商品信息缓存失败", err)
		}
	}
}
