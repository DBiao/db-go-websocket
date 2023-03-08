package config

type System struct {
	SystemId  uint16 `mapstructure:"system-id" yaml:"system-id"`   // 是否集群
	IsCluster bool   `mapstructure:"is-cluster" yaml:"is-cluster"` // 是否集群
}
