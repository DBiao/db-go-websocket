package config

type Etcd struct {
	Endpoints []string `mapstructure:"endpoints" yaml:"endpoints"` // etcd集群地址
}
