package req

// 秒杀请求
type SpikesRequest struct {
	UserID uint `json:"userID"`
	GoodID uint `json:"goodID"`
}
