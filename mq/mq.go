package mq

import "miaosha-system/inter"

func Init() {
	CreateMQ = &CreateOderMQ{
		Order: inter.GetOrder(),
	}
	CloseMQ = &CloseOrderMQ{
		Order: inter.GetOrder(),
	}
	//StockMQ = &UpdateStockMQ{}
	Refresh = &RefreshTask{
		Good: inter.GetGood(),
	}
}

func Run() {
	go CreateMQ.Receive()
	go CloseMQ.Receive()
	//go StockMQ.PeriodicUpdateStock()
	go Refresh.Start()
}
