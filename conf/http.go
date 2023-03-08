package config

type Http struct {
	Port int    `mapstructure:"port" yaml:"port"` // 端口值
	Mode string `mapstructure:"mode" yaml:"mode"` // gin调式模式
}
