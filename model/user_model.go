package model

import (
	"miaosha-system/common"
)

type UserModel struct {
	common.Model
	Pwd      string       `gorm:"size:64" json:"-"`
	Role     int8         `json:"role"` //1 管理员 2 卖家 3 买家
	Nickname string       `gorm:"size:32" json:"nickname"`
	Orders   []OrderModel `gorm:"foreignKey:UserID" json:"orders"`
}
