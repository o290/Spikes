package model

import (
	"miaosha-system/common"
	"time"
)

type GoodModel struct {
	common.Model
	Name        string       `gorm:"size:32" json:"name"`
	Img         string       `gorm:"size:256" json:"img"`
	OriginPrice float32      `json:"originPrice"`
	Price       float32      `json:"price"`
	Stock       int          `json:"stock"`
	StartTime   time.Time    `json:"startTime"`
	EndTime     time.Time    `json:"endTime"`
	Orders      []OrderModel `gorm:"foreignKey:GoodID" json:"orders"`
	Status      int8         `json:"status"` //0 已下架 1 未开始 2 进行中
}

// CheckPrice 对price合法性检查
func (g *GoodModel) CheckPrice() bool {
	if g.Price > 0 {
		return true
	}
	return false
}

// CheckStock 对库存数量检查
func (g *GoodModel) CheckStock() bool {
	if g.Stock >= 0 {
		return true
	}
	return false
}

// CheckTime 添加对 StartTime 和 EndTime 的验证，确保 StartTime 早于 EndTime。
func (g *GoodModel) CheckTime() bool {
	if g.StartTime.Before(g.EndTime) {
		return true
	}
	return false
}

// Check 判断当前时间是否处于活动的有效时间范围内，以及商品是否还有库存
func (g *GoodModel) Check() bool {
	now := time.Now().Unix()
	start := g.StartTime.Unix()
	end := g.EndTime.Unix()
	if now < start || now > end {
		return false
	}
	return true
}
