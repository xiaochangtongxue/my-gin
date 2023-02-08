package config

type Log struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"`
	Root       string `mapstructure:"root" json:"root" yaml:"root"`
	FileName   string `mapstructure:"fileName" json:"fileName" yaml:"fileName"`
	Format     string `mapstructure:"format" json:"format" yaml:"format"`
	ShowLine   bool   `mapstructure:"showLine" json:"showLine" yaml:"showLine"`
	MaxBackups int    `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"`
	MaxSize    int    `mapstructure:"maxSize" json:"maxSize" yaml:"maxSize"`
	MaxAge     int    `mapstructure:"maxAge" json:"maxAge" yaml:"maxAge"`
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}
