package res

type OrderInfoResponse struct {
	ID          uint    `json:"ID"`
	OrderNumber string  `json:"orderNumber"`
	BuyerID     uint    `json:"buyerID"`
	BuyerName   string  `json:"buyerName"`
	GoodID      uint    `json:"goodID"`
	GoodName    string  `json:"goodName"`
	Img         string  `json:"img"`
	GoodPrice   float32 `json:"goodPrice"`
	ActualPay   float32 `json:"actualPay"`
	PayWay      uint    `json:"payWay"` //1 微信 2 支付宝 3 其他
	Number      int8    `json:"number"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	Status      int8    `json:"status"`
}

type OrderInfoListResponse struct {
	List  []OrderInfoResponse `json:"list"`
	Count int                 `json:"count"`
}
