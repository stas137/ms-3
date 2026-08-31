package config

import "net"

type inventoryConfig struct {
	Host string `yaml:"host" env:"GRPC_HOST_INVENTORY" env-default:"localhost"`
	Port string `yaml:"port" env:"GRPC_PORT_INVENTORY" env-default:"50051"`
}

func (c *inventoryConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
