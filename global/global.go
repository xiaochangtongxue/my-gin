package global

import (
	"github.com/spf13/viper"
	"github.com/xiaochangtongxue/my-gin/config"
	"go.uber.org/zap"
)

var (
	MGIN_CONFIG config.Server
	MGIN_VIP    *viper.Viper
	MGIN_ZAP    *zap.Logger
)
