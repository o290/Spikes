package good

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"miaosha-system/common"
	"miaosha-system/common/req"
	"miaosha-system/common/res"
	"miaosha-system/global"
	"miaosha-system/inter"
	"miaosha-system/model"
	"net/http"
	"sync"
	"time"
)

var once sync.Once

func GoodInit() {
	once.Do(func() {
		inter.GoodController = &GoodControllerr{}
	})
}

type GoodControllerr struct {
}

//12.修改商品信息 在更新数据库时，不直接更新缓存，而是让缓存失效
//11.删除商品

// InitStock 商品库存预热功能 初始化商品库存信息到 Redis 缓存中。
func (m GoodControllerr) Init() (err error) {
	var goodList []model.GoodModel
	err = global.DB.Find(&goodList).Error
	if err != nil {
		global.Log.Error("初始化库存失败", err)
		return
	}
	for _, goodModel := range goodList {
		if err = m.SetStock(goodModel.ID, goodModel.Stock); err != nil {
			global.Log.Error("初始化库存失败", err)
			return
		}
	}
	return
}

// GoodAdd 添加商品
func (m GoodControllerr) GoodAdd(c *gin.Context) {
	var req req.GoodAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Fatalln("解析失败")
		c.JSON(http.StatusOK, gin.H{
			"message": "数据解析失败",
		})
		return
	}
	start, _ := time.Parse("2006-01-02 15:05", req.StartTime)
	end, _ := time.Parse("2006-01-02 15:05", req.EndTime)
	if req.Price <= 0 || req.OriginPrice <= 0 || req.Stock < 0 || start.After(end) || start.Before(time.Now()) {
		global.Log.Error("商品参数设置错误")
		c.JSON(http.StatusOK, gin.H{
			"message": "商品参数设置错误",
		})
		return
	}
	good := model.GoodModel{
		Name:        req.GoodName,
		Img:         req.Img,
		OriginPrice: req.OriginPrice,
		Price:       req.Price,
		Stock:       req.Stock,
		StartTime:   start,
		EndTime:     end,
		Status:      1,
	}
	if err := global.DB.Create(&good).Error; err != nil {
		global.Log.Fatalln("添加商品失败")
		c.JSON(http.StatusOK, gin.H{
			"message": "添加商品失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "添加商品成功",
	})
	return
}

// GoodList 商品列表查询 从数据库中获取商品信息列表
func (m GoodControllerr) GoodList(c *gin.Context) {
	var req req.GoodListRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Log.Fatalln("解析失败")
		c.JSON(http.StatusOK, gin.H{
			"message": "get good list  error",
		})
		return
	}
	var info common.PageInfo
	if req.Page <= 0 {
		info.Page = 1
	}
	if req.Limit != -1 { //-1查全部
		if req.Limit <= 0 {
			info.Limit = 10
		} else {
			info.Limit = req.Limit
		}
	}
	offset := (info.Page - 1) * info.Limit
	//根据开始时间降序排序
	var goodList []model.GoodModel
	err := global.DB.Order("start_time desc").Limit(info.Limit).Offset(offset).Find(&goodList).Error
	if err != nil {
		global.Log.Error("获取商品列表失败", err)
		c.JSON(http.StatusOK, gin.H{
			"message": "get good list  error",
		})
		return
	}
	//求总数
	var count int64
	err = global.DB.Model(model.GoodModel{}).Count(&count).Error
	if err != nil {
		global.Log.Error("获取商品数目失败", err)
		c.JSON(http.StatusOK, gin.H{
			"message": "get good list  error",
		})
		return
	}

	var list = make([]res.GoodInfoResponse, 0)
	for _, goodModel := range goodList {
		good := res.GoodInfoResponse{
			GoodID:      goodModel.ID,
			GoodName:    goodModel.Name,
			Img:         goodModel.Img,
			OriginPrice: goodModel.OriginPrice,
			Price:       goodModel.Price,
			Stock:       goodModel.Stock,
			StartTime:   goodModel.StartTime,
			EndTime:     goodModel.EndTime,
			Status:      goodModel.Status,
		}
		list = append(list, good)
	}

	response := res.GoodInfoListResponse{
		List:  list,
		Count: int(count),
	}
	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

// GetGoodDetail 根据商品 ID 获取商品信息 首先从缓存中获取，没有从数据库中获取
func (m GoodControllerr) GetGoodDetail(c *gin.Context) {
	goodIDStr := c.Query("goodID")
	if goodIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "goodID parameter is missing",
		})
		return
	}
	var goodID uint
	_, err := fmt.Sscanf(goodIDStr, "%d", &goodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid goodID parameter format",
		})
		return
	}

	// 从缓存中获取
	good, err := m.GetGoodFromRedis(goodID)
	if err != nil && err != redis.Nil {
		global.Log.Error("获取商品信息失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取商品信息失败",
		})
		return
	}
	if good.ID <= 0 && err == redis.Nil {
		// 从数据库中获取
		err = global.DB.Where("id=?", goodID).First(&good).Error
		// 存入缓存
		err = m.SetGoodToRedis(good)
		if err != nil {
			global.Log.Error("商品信息获取失败:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get good error",
			})
			return
		}
	}
	c.JSON(http.StatusOK, good)
}
func (m GoodControllerr) GetGood(goodID uint) (good model.GoodModel, err error) {
	//从缓存中获取
	good, err = m.GetGoodFromRedis(goodID)
	if err != nil && err != redis.Nil {
		global.Log.Error("获取商品信息失败:", err)
		return
	}
	if good.ID <= 0 && err == redis.Nil {
		//从数据库中获取
		err = global.DB.Where("id=?", goodID).First(&good).Error
		//存入缓存
		err = m.SetGoodToRedis(good)
		if err != nil {
			global.Log.Error("商品信息获取失败:", err)
			return
		}
	}
	return good, nil
}

// SetGoodToRedis 将商品信息存入 Redis 缓存 序列化为 JSON 字符串，并存储到 Redis 中
func (m GoodControllerr) SetGoodToRedis(good model.GoodModel) error {
	data, err := json.Marshal(good)
	if err != nil {
		global.Log.Error("商品信息序列化失败")
		return errors.New("商品信息序列化失败")
	}
	key := fmt.Sprintf("good_%d", good.ID)
	err = global.Redis.Set(context.Background(), key, data, 24*time.Hour).Err()
	if err != nil {
		global.Log.Error("商品信息缓存失败")
		return errors.New("商品信息缓存失败")
	}
	global.Log.Infof("%s商品存入缓存成功", key)
	return nil
}

// GetGoodFromRedis 从 Redis 缓存中获取商品信息，并反序列化
func (m GoodControllerr) GetGoodFromRedis(id uint) (good model.GoodModel, err error) {
	key := fmt.Sprintf("good_%d", id)
	data, err1 := global.Redis.Get(context.Background(), key).Result()
	if err1 != nil {
		err = err1
		return
	}
	if data == "" {
		return
	}

	err2 := json.Unmarshal([]byte(data), &good)
	if err2 != nil {
		err = err2
		return
	}
	global.Log.Infof("%s商品获取缓存成功", key)
	return good, err
}

// SetStock 将商品的库存信息存储到 Redis 中，不设置过期时间
func (m GoodControllerr) SetStock(id uint, stock int) (err error) {
	key := fmt.Sprintf("stock:%d", id)
	if err = global.Redis.Set(context.Background(), key, stock, time.Minute*5).Err(); err != nil {
		global.Log.Error("设置库存失败")
		global.Log.Errorf("设置库存失败，键: %s，值: %d，错误信息: %v", key, stock, err)
		return err
	}
	return
}

// DecrStock 减少指定商品 ID 的库存数量，通过 Redis 的 DECR 命令实现
func (m GoodControllerr) DecrStock(id uint) (stock int64, err error) {
	key := fmt.Sprintf("stock:%d", id)
	stock, err = global.Redis.Decr(context.Background(), key).Result()
	if err != nil {
		global.Log.Error("减库存失败")
		return
	}
	return
}

// IncrStock 增加指定商品 ID 的库存数量，通过 Redis 的 INCR 命令实现
func (m GoodControllerr) IncrStock(id uint) (stock int64, err error) {
	key := fmt.Sprintf("stock:%d", id)
	stock, err = global.Redis.Incr(context.Background(), key).Result()
	if err != nil {
		global.Log.Error("加库存失败")
		return
	}
	return
}

func (m GoodControllerr) TestToken(c *gin.Context) {
	id, exists := c.Get("user_id")
	if !exists {
		global.Log.Error("未获取到用户 ID")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "未获取到用户 ID",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
	fmt.Println("用户id：", id)
}
