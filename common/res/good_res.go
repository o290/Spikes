package res

import "time"

type GoodInfoResponse struct {
	GoodID      uint      `json:"goodID"`
	GoodName    string    `json:"goodName"`
	Img         string    `json:"img"`
	OriginPrice float32   `json:"originPrice"`
	Price       float32   `json:"price"`
	Stock       int       `json:"stock"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Status      int8      `json:"status"`
}
type GoodInfoListResponse struct {
	List  []GoodInfoResponse `json:"list"`
	Count int                `json:"count"`
}
type GetGoodResponse struct {
	GoodID      uint      `json:"goodID"`
	GoodName    string    `json:"goodName"`
	Img         string    `json:"img"`
	OriginPrice float32   `json:"originPrice"`
	Price       float32   `json:"price"`
	Stock       int       `json:"stock"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Status      int8      `json:"status"`
}
