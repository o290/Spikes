package model

import (
	"miaosha-system/common"
)

type OrderModel struct {
	common.Model
	UserID        uint      `json:"userID"`
	UserModel     UserModel `gorm:"foreignKey:UserID" json:"-"`
	GoodID        uint      `json:"goodID"`
	GoodModel     GoodModel `gorm:"foreignKey:GoodID" json:"-"`
	GoodNumber    int8      `json:"goodNumber"` //商品数量
	OrderNumber   string    `gorm:"size:256" json:"orderNumber"`
	ActualPayment float32   `json:"actualPayment"` //实际付款，显示到小数点后两位
	Status        int8      `json:"status"`        // 2 完成  1 未完成 0 关闭
}

//订单号可以使用uuid或者分布式id
