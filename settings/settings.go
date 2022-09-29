package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Init() 加载配置
func Init() (err error) {
	// 指定文件路径
	viper.SetConfigFile("./conf/config.yaml")
	// 监控
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改...")
	})
	// 读配置信息
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("ReadInConfig failed, err: %v", err))
	}
	return err
}
