package core

import (
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"miaosha-system/global"
	"time"
)

// gorm 连数据库
func Initgorm() *gorm.DB {
	if global.Config.Mysql.Host == "" {
		global.Log.Warnln("未配置mysql，取消gorm连接")
		return nil
	}
	dsn := global.Config.Mysql.Dsn()

	var mysqlLogger logger.Interface
	if global.Config.System.Env == "debug" {
		//开发环境显示所有的sql,打印所有的日志
		mysqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		mysqlLogger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: mysqlLogger,
	})
	if err != nil {
		//连接失败，直接退出
		global.Log.Fatalf(fmt.Sprintf("[%s]mysql连接失败", dsn))
	}
	global.Log.Println("数据库连接成功")
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)               //最大空闲连接数
	sqlDB.SetMaxIdleConns(100)              //最多可容纳
	sqlDB.SetConnMaxIdleTime(time.Hour * 4) //连接最大复用时间
	return db
}
