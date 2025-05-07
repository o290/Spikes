package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"miaosha-system/common/msg"
	"miaosha-system/common/req"
	"miaosha-system/controller/good"
	"miaosha-system/global"
	"miaosha-system/inter"
	"miaosha-system/model"
	"miaosha-system/mq"
	"miaosha-system/utils/lock"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var once sync.Once

func OrderInit() {
	once.Do(func() {
		inter.OrderController = &OrderController{}
	})
}

type OrderController struct {
}

//2.获取订单列表
//3.获取订单信息

//4.超时关闭订单

//5.获取秒杀结果：preload预加载商品信息
//6.将订单信息添加到缓存
//7.删除订单信息的缓存
//8加锁，释放锁

// Spikes 秒杀
func (o *OrderController) Spikes(c *gin.Context) {
	var req req.SpikesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("解析失败")
		c.JSON(http.StatusOK, gin.H{
			"message": "数据解析失败",
		})
		return
	}
	//判断用户是否存在
	var user model.UserModel
	err := global.DB.Model(&model.UserModel{}).Where("id = ?", req.UserID).Take(&user).Error
	if err != nil {
		global.Log.Error("不存在该用户")
		c.JSON(http.StatusOK, gin.H{
			"message": "不存在该用户",
		})
		return
	}
	//判断用户是否可以进行秒杀，限制秒杀次数。判断商品是否在秒杀期间，
	//获取商品信息
	var goodController good.GoodControllerr
	good, err := goodController.GetGood(req.GoodID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "商品信息获取失败",
		})
		return
	}
	//判断是否符合秒杀条件
	//if !(good.Check() && good.Stock > 0) {
	//fmt.Println(good.StartTime)
	if !(good.Check()) {
		global.Log.Error("商品不满足秒杀条件")
		c.JSON(http.StatusOK, gin.H{
			"message": "商品不满足秒杀条件",
		})
		return
	}
	//判断用户是否重复秒杀
	err = global.Redis.Get(context.Background(), fmt.Sprintf("order:%d:%d", req.UserID, req.GoodID)).Err()
	if err == nil {
		global.Log.Infof("用户秒杀次数已达上限")
		c.JSON(http.StatusOK, gin.H{
			"message": "用户秒杀次数已达上限",
		})
		return
	}
	//获取锁
	_, acquired, err := lock.AcquireLock(context.Background(), fmt.Sprintf("lock:%d", req.GoodID), 5*time.Second, req.UserID, req.GoodID)
	if err != nil || !acquired {
		fmt.Printf("用户 %d 获取商品%d的分布式锁失败: %v\n", req.UserID, req.GoodID, err)
		return
	}
	fmt.Printf("用户 %d 获取锁成功\n", req.UserID)
	//预减库存
	stock, _ := global.Redis.Decr(context.Background(), fmt.Sprintf("stock:%d", req.GoodID)).Result()
	if stock < 0 {
		global.Log.Infof("商品%d已经售罄", req.GoodID)
		c.JSON(http.StatusOK, gin.H{
			"message": "商品已经售罄",
		})
		return
	}
	//构造创建订单消息
	msg := msg.CreateMsg{
		UserID: req.UserID,
		GoodID: req.GoodID,
	}
	mq.CreateMQ.Send(msg)
	//fmt.Println("当前库存信息：", stock)
}

func (o *OrderController) CreateOrder(userID, goodID uint) (err error) {
	//创建订单
	order := model.OrderModel{
		UserID:     userID,
		GoodID:     goodID,
		GoodNumber: 1, //购买的商品数量，先写1
		Status:     1, //未完成
	}
	order.OrderNumber = o.GenerateOrderID(userID, goodID)
	//创建订单信息 数据库
	//开启事务
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		var good model.GoodModel
		result := tx.Where("id=? and stock>0", goodID).First(&good).Update("stock", gorm.Expr("stock-1"))
		if result.Error != nil {
			return fmt.Errorf("减库存失败: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("库存不足")
		}
		// 创建订单
		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("创建订单失败: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 创建订单信息缓存
	err = global.Redis.Set(context.Background(), fmt.Sprintf("order:%d:%d", userID, goodID), order.OrderNumber, 2*time.Minute).Err()
	if err != nil {
		global.Log.Printf("创建订单信息缓存失败: %v", err)
		return err
	}
	// 加入订单超时延迟队列
	mq.CloseMQ.Send(order.OrderNumber)
	return
}

// GenerateOrderID 生成订单id
func (o *OrderController) GenerateOrderID(userID, productID uint) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	return fmt.Sprintf("%s_%d_%d", timestamp, userID, productID)
}

// GetOrderInfo 获取订单信息
func (o *OrderController) GetOrderInfo(orderID string) (order model.OrderModel, err error) {
	err = global.DB.Where("order_number=?", orderID).Take(&order).Error
	return
}

// CloseOrder 关闭订单
func (o *OrderController) CloseOrder(order model.OrderModel) {
	//判断订单的状态
	if order.Status != 1 {
		fmt.Println("订单已完成")
		//移除订单
		mq.CloseMQ.Remove(order.OrderNumber)
	}
	//关闭订单
	//增加库存
	//删除订单缓存
	//移除订单
}
func Close(order model.OrderModel) (err error) {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&struct{}{}).Where("id=?", order.GoodID).UpdateColumn("stock", gorm.Expr("stock+?", 1))
		if result.Error != nil {
			global.Log.Errorf("回滚库存时出错, err: %v, 商品ID: %v", result.Error, order.GoodID)
			return result.Error
		}
		if result.RowsAffected == 0 {
			global.Log.Errorf("回滚库存时未影响任何行, 商品ID: %v", order.GoodID)
			return errors.New("加库存失败")
		}

		// 修改订单信息状态
		result = tx.Model(&model.OrderModel{}).Where("order_id = ? AND status = ?", order.ID, order.Status).Update("status", 0)
		if result.Error != nil {
			log.Printf("修改订单状态时出错, err: %v", result.Error)
			return result.Error
		}
		if result.RowsAffected == 0 {
			log.Printf("修改订单状态时未影响任何行, 订单ID: %v", order.ID)
			return errors.New("修改订单信息状态失败")
		}
		global.Log.Infof("订单%d关闭订单回滚库存和修改订单成功", order.ID)
		return nil
	})
}

// PayOrder 支付订单，修改订单的状态
//func (o *OrderController) PayOrder() (err error) {
//	//修改订单状态
//	return global.DB.Transaction(func(tx *gorm.DB) error {
//		//
//	})
//	//移除订单
//	mq.CloseMQ.Remove(order.OrderNumber)
//}
