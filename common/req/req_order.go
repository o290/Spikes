package req

// SpikesRequest 秒杀请求
type SpikesRequest struct {
	GoodID uint `json:"goodID"`
}

// CloseRequest 手动关闭订单请求
type CloseRequest struct {
	UserID  uint `json:"userID"`
	OrderID uint `json:"orderID"`
}

type GetOrderListRequest struct {
	Page  int `form:"page,optional"`
	Limit int `form:"limit,optional"`
}
