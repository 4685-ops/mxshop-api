package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
func InitConfig() {
	//debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// 这个对象如何在其他文件中使用 -全局变量
	// serverConfig := config.ServerConfig{}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}

	zap.S().Info("配置信息：&v", global.ServerConfig)

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件产生变化:%s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Info("配置信息：&v", global.ServerConfig)

	})
}
