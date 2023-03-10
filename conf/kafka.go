package config

type Kafka struct {
	Brokers []string `mapstructure:"brokers" yaml:"brokers"` //
}
