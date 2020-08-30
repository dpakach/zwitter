package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ServiceConfig struct {
  Server struct {
    Name        string `yaml:"Name"`
    GrpcAddr    string `yaml:"GrpcAddr"`
    RestAddr    string `yaml:"RestAddr"`
    CertFile    string `yaml:"CertFile"`
    KeyFile     string `yaml:"KeyFile"`
    ServerName  string `yaml:"ServerName"`
  } `yaml:"Server"`
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

