package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// log_level: info
// server:
//   tcp: true
//   udp: true
//   udp_port: 3478
//   tcp_port: 3478
//   publicip: 127.0.0.1
//   realm: rtc.qcloud.com
//   password: pass

var DefaultConfig = &Config{}

type serverstruct struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	TCP      bool   `yaml:"tcp"`
	UDP      bool   `yaml:"udp"`
	TCPPort  int    `yaml:"tcp_port"`
	UDPPort  int    `yaml:"udp_port"`
	PublicIP string `yaml:"publicip"`
	Realm    string `yaml:"realm"`
	Password string `yaml:"password"`
}

type Config struct {
	LogLevel string        `yaml:"log_level"`
	Server   *serverstruct `yaml:"server"`
}

func LoadConfig(filePath string) (*Config, error) {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, DefaultConfig)
	if err != nil {
		return nil, err
	}

	return DefaultConfig, nil
}
