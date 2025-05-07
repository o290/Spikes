package core

import (
	"github.com/spf13/viper"
	"log"
	"miaosha-system/global"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		//log.Fatalf() 时，它不仅会输出错误信息，而且会 调用 os.Exit(1)，
		//使程序立即终止并退出，后续的代码（如 fmt.Println()）将不会被执行。
		log.Fatalf("config read error:%s", err)
	}
	log.Println("读取配置文件成功")
	if err := viper.UnmarshalKey("config", &global.Config); err != nil {
		log.Fatalf("Error Unmarshaling key, %s", err)
	}
	log.Println("反序列化配置文件成功")
}
