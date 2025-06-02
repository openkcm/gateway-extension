// Package config defines the necessary types to configure the application.
// An example config file config.yaml is provided in the repository.
package config

import (
	"github.com/openkcm/common-sdk/pkg/commoncfg"
)

type ListenerType string

const (
	UNIXListener ListenerType = "unix"
	TCPListener  ListenerType = "tcp"
)

type Config struct {
	commoncfg.BaseConfig `mapstructure:",squash"`

	Listener             Listener                       `yaml:"listener"`
	GRPCClientAttributes commoncfg.GRPCClientAttributes `yaml:"grpcClientAttributes"`
}

type Listener struct {
	Type ListenerType `yaml:"type"`
	TCP  *TCP         `yaml:"tcp"`
	UNIX *UNIX        `yaml:"unix"`
}

type TCP struct {
	// TCP.Address is the address to listen on for gRPC requests
	Address string `yaml:"address"`
}

type UNIX struct {
	// SocketPath is the Unix Path to listen on for gRPC requests
	SocketPath string `yaml:"socketPath"`
}
