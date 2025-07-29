// Package config defines the necessary types to configure the application.
// An example config file config.yaml is provided in the repository.
package config

import (
	"time"

	"github.com/openkcm/common-sdk/pkg/commoncfg"
)

type ListenerType string

const (
	UNIXListener ListenerType = "unix"
	TCPListener  ListenerType = "tcp"
)

type Config struct {
	commoncfg.BaseConfig `mapstructure:",squash"`

	Listener Listener `yaml:"listener"`
}

type Listener struct {
	Type             ListenerType                   `yaml:"type" default:"tcp"`
	TCP              commoncfg.GRPCServer           `yaml:"tcp"`
	UNIX             UNIX                           `yaml:"unix"`
	ShutdownTimeout  time.Duration                  `yaml:"shutdownTimeout" default:"5s"`
	ClientAttributes commoncfg.GRPCClientAttributes `yaml:"clientAttributes"`
}

type UNIX struct {
	// SocketPath is the Unix Path to listen on for gRPC requests
	SocketPath string `yaml:"socketPath" default:"/etc/envoy/gateway/extension.sock"`
}
