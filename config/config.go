package config

type BaseConfig struct {
	App AppConfig `mapstructure:"app" json:"app" yaml:"app"`
	Log LogConfig `mapstructure:"log" json:"log" yaml:"log"`
}

type AppConfig struct {
	Env     string `mapstructure:"env" json:"env" yaml:"env"`
	Port    string `mapstructure:"port" json:"port" yaml:"port"`
	AppName string `mapstructure:"appname" json:"appname" yaml:"appname"`
	AppUrl  string `mapstructure:"appurl" json:"appurl" yaml:"appurl"`
}

type LogConfig struct {
	Address string `mapstructure:"address" json:"address" yaml:"address"`
}
