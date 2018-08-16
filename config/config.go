// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period     time.Duration `config:"period"`
	ClientCert string        `config:"client_cert"`
	ClientKey  string        `config:"client_key"`
	ServerCert string        `config:"server_cert"`
	Hosts      []string      `config:"hosts"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
}
