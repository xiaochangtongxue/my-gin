package config

type Mysql struct {
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
	PassWord string `mapstructure:"passWord" json:"passWord" yaml:"passWord"`
	DbName   string `mapstructure:"dbName" json:"dbName" yaml:"dbName"`
}
