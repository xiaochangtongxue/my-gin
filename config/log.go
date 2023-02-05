package config

type Log struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"`
	Root       string `mapstructure:"root" json:"root" yaml:"root"`
	FileName   string `mapstructure:"filename" json:"filename" yaml:"filename"`
	Format     string `mapstructure:"format" json:"format" yaml:"format"`
	ShowLine   bool   `mapstructure:"showline" json:"showline" yaml:"showline"`
	MaxBackups int    `mapstructure:"maxbackups" json:"maxbackups" yaml:"maxbackups"`
	MaxSize    int    `mapstructure:"maxsize" json:"maxsize" yaml:"maxsize"`
	MaxAge     int    `mapstructure:"maxage" json:"maxage" yaml:"maxage"`
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}
