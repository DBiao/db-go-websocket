package config

type Config struct {
	System System `mapstructure:"system" yaml:"system"`
	Http   Http   `mapstructure:"http" yaml:"http"`
	Zap    Zap    `mapstructure:"zap" yaml:"zap"`
	Etcd   Etcd   `mapstructure:"etcd" yaml:"etcd"`
	Grpc   Grpc   `mapstructure:"grpc" yaml:"grpc"`
}
