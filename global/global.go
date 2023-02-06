package global

import (
	"github.com/spf13/viper"
	"github.com/xiaochangtongxue/my-gin/config"
)

var (
	MGIN_CONFIG config.Server
	MGIN_VIP    *viper.Viper
)
