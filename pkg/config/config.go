// Utils to load configurations for server from yaml file
package config

import (
	"os"

	p "github.com/liangLouise/http_server/pkg/httpProto"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Server struct {
		Port    int                     `yaml:"SERVER_PORT"`
		Host    string                  `yaml:"SERVER_HOST"`
		Version p.HTTP_PROTOCOL_VERSION `yaml:"HTTP_VERSION"`
	} `yaml:"Server"`

	RunTime struct {
		MaxConnections  int  `yaml:"MAX_CONCURRENT_CONNECTIONS"`
		HasPersistant   bool `yaml:"ENABLE_PESISTANT"`
		HasPipelining   bool `yaml:"ENABLE_PIPELINING"`
		MaxPipelining   int  `yaml:"MAX_PIPELINING_NUMBER"`
		TimeoutDuration int  `yaml:"TIMEOUT_DURATION"`
	} `yaml:"RunTime"`
}

// function to parse the config data from config yaml file
func LoadConfig(path string) (config *ServerConfig, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg ServerConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
