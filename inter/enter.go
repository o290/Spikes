package inter

import (
	"github.com/gin-gonic/gin"
	"miaosha-system/model"
)

var (
	OrderController Order
	GoodController  Good
)

func GetOrder() Order {
	return OrderController
}

type Order interface {
	CreateOrder(userID, goodID uint) (err error)
	GenerateOrderID(userID, productID uint) string
	GetOrder(orderID string) (order model.OrderModel, err error)
	CloseUpdateStock(order model.OrderModel) (err error)
}
type Good interface {
	//init() (err error)
	GoodAdd(c *gin.Context)
	GoodList(c *gin.Context)
	GetGoodDetail(c *gin.Context)
	SetGoodToRedis(good model.GoodModel) error
}
