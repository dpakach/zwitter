package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ServiceConfig struct {
	Server struct {
		Name       string `yaml:"Name"`
		GrpcAddr   string `yaml:"GrpcAddr"`
		RestAddr   string `yaml:"RestAddr"`
		CertFile   string `yaml:"CertFile"`
		KeyFile    string `yaml:"KeyFile"`
		ServerName string `yaml:"ServerName"`
		Nodes      []struct {
			Name string `yaml:"Name"`
			Host string `yaml:"Host"`
			Port string `yaml:"Port"`
		} `yaml:"Nodes"`
	} `yaml:"Server"`
}

func (sc *ServiceConfig) GetNodeAddr(name string) (error, string) {
	for _, node := range sc.Server.Nodes {
		if node.Name == name {
			return nil, fmt.Sprintf("%v:%v", node.Host, node.Port)
		}
	}
	return fmt.Errorf("Node with name %v not found", name), ""
}

func NewServerconfig(configFile string) (*ServiceConfig, error) {
	f, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	sc := &ServiceConfig{}

	err = yaml.Unmarshal(f, &sc)
	if err != nil {
		return nil, err
	}

	return sc, nil
}
