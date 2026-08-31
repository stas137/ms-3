package config

import "net"

type paymentConfig struct {
	Host string `yaml:"host" env:"GRPC_HOST_PAYMENT" env-default:"localhost"`
	Port string `yaml:"port" env:"GRPC_PORT_PAYMENT" env-default:"50052"`
}

func (c *paymentConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
