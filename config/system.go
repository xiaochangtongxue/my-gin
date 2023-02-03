package config

type System struct {
	Env     string `mapstructure:"env" json:"env" yaml:"env"`
	Port    string `mapstructure:"port" json:"port" yaml:"port"`
	AppName string `mapstructure:"appname" json:"appname" yaml:"appname"`
	AppUrl  string `mapstructure:"appurl" json:"appurl" yaml:"appurl"`
}
