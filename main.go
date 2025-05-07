package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"miaosha-system/controller/order"
	"miaosha-system/core"
	"miaosha-system/global"
	"miaosha-system/inter"
	"miaosha-system/model"
	"miaosha-system/mq"
	"miaosha-system/routers"
)

type Options struct {
	DB bool
}

func main() {
	//读取配置文件
	core.InitConfig()
	//初始化日志
	global.Log = core.InitLogger()
	//连接数据库
	global.DB = core.Initgorm()
	//连接redis
	global.Redis = core.InitRedis(global.Config.Redis.Addr, global.Config.Redis.Pwd, global.Config.Redis.DB)

	fmt.Println(global.Config.System)
	//表结构迁移
	var opt Options
	flag.BoolVar(&opt.DB, "db", false, "db")
	flag.Parse()
	if opt.DB {
		err := global.DB.AutoMigrate(
			&model.UserModel{},
			&model.GoodModel{},
			&model.OrderModel{},
		)
		if err != nil {
			fmt.Println("表结构生成失败", err)
			return
		}
		fmt.Println("表结构生成成功")
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	inter.OrderController = &order.OrderController{}
	////表示加载templates目录下的文件夹下的文件
	//r.LoadHTMLGlob("./templates/**.html")
	////配置静态web目录，第一个参数表示路由，第二个参数表示映射目录
	//r.Static("/static", "./static")
	//routers.DefaultRouterInit(r)
	routers.Init()
	routers.UserRouterInit(r)
	routers.GoodRouterInit(r)
	routers.OrderRouterInit(r)
	mq.Run()
	//routers.FileRouter(r)

	addr := global.Config.System.Addr()
	global.Log.Printf("%s运行在: %s", global.Config.System.Name, addr)
	r.Run(addr)
}
