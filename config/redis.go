package config

type Redis struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port string `mapstructure:"port" json:"port" yaml:"port"`
}
