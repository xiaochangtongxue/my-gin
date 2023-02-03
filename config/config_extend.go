package config

type ExtendConfig struct {
	Mysql MysqlConfig `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis RedisConfig `mapstructure:"redis" json:"redis" yaml:"redis"`
}

type MysqlConfig struct {
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
	PassWord string `mapstructure:"password" json:"password" yaml:"password"`
	DbName   string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port string `mapstructure:"port" json:"port" yaml:"port"`
}
