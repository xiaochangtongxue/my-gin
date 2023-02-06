package core

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/xiaochangtongxue/my-gin/core/internal"
	"github.com/xiaochangtongxue/my-gin/global"
)

func Viper() *viper.Viper {
	var config string

	if configEev := os.Getenv(internal.ConfigEnv); configEev == "" {
		switch gin.Mode() {
		case gin.DebugMode:
			config = internal.ConfigDefaultFile
			fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigDefaultFile)
		case gin.TestMode:
			config = internal.ConfigDevFile
			fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigDevFile)
		case gin.ReleaseMode:
			config = internal.ConfigReleaseFile
			fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigReleaseFile)
		}
	} else {
		config = configEev
		fmt.Printf("您正在使用%s变量,config的路径为%s\n", internal.ConfigEnv, configEev)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config is change", e.Name)
		if err := v.Unmarshal(&global.MGIN_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&global.MGIN_CONFIG); err != nil {
		fmt.Println(err)
	}
	return v
}
