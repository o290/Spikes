package mq

import "miaosha-system/inter"

func Init() {
	CreateMQ = &CreateOderMQ{
		Order: inter.GetOrder(),
	}
	CloseMQ = &CloseOrderMQ{
		Order: inter.GetOrder(),
	}
	StockMQ = &UpdateStockMQ{}
}
func Run() {
	go CreateMQ.Receive()
	go CloseMQ.Receive()
	go StockMQ.PeriodicUpdateStock()
}
