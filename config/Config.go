package config

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type AppConfig struct {
	Log            Log             `yaml:"log"`
	CloudProviders []CloudProvider `yaml:"cloud-providers"`
}

type Log struct {
	EnableFile bool   `yaml:"enable-file"`
	FileDir    string `yaml:"file-dir"`
	Level      string `yaml:"level"`
}

type CloudProvider struct {
	Provider  string `yaml:"provider"`
	SecretId  string `yaml:"secret-id"`
	SecretKey string `yaml:"secret-key"`
}

func newDefaultAppConfig() *AppConfig {
	return &AppConfig{
		Log: Log{
			EnableFile: false,
		},
		CloudProviders: make([]CloudProvider, 0),
	}
}

func ReadConfig(configFile string) (*AppConfig, error) {
	if !utils.IsFile(configFile) {
		return nil, fmt.Errorf("config file does not exist: %s", configFile)
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	r := newDefaultAppConfig()
	if err = yaml.Unmarshal(data, r); err != nil {
		return nil, err
	}
	return r, nil
}
