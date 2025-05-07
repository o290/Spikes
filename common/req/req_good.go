package req

// GoodListRequest 获取商品列表  所有人都可以
type GoodListRequest struct {
	Page  int `form:"page,optional"`
	Limit int `form:"limit,optional"`
	//Key   string `form:"key,optional"`
}

// 添加商品 只有商家才可以
type GoodAddRequest struct {
	GoodName    string  `json:"goodName"`
	Img         string  `json:"img"`
	OriginPrice float32 `json:"originPrice"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
	StartTime   string  `json:"startTime"` //2006-01-02 15:00
	EndTime     string  `json:"endTime"`
}

// GetGoodRequest 查询商品 所有人都可以
type GetGoodRequest struct {
	GoodID uint `form:"goodID"`
}
type GoodUpdateRequest struct {
	GoodID      uint    `json:"goodID"`
	GoodName    string  `json:"goodName,optional"`
	Img         string  `json:"img,optional"`
	OriginPrice float32 `json:"originPrice,optional"`
	Price       float32 `json:"price,optional"`
	Stock       int     `json:"stock,optional"`
	StartTime   string  `json:"startTime,optional"` //2006-01-02 15:00
	EndTime     string  `json:"endTime,optional"`
}
