package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type AppConfig struct {
	Log            Log                      `koanf:"log"`
	CloudProviders map[string]CloudProvider `koanf:"cloud-providers"`
}

type Log struct {
	EnableFile bool   `koanf:"enable-file"`
	FileDir    string `koanf:"file-dir"`
	Level      string `koanf:"level"`
}

type CloudProvider struct {
	Provider  string `koanf:"provider"`
	SecretId  string `koanf:"secret-id"`
	SecretKey string `koanf:"secret-key"`
}

func ReadConfig(configFile string) (*AppConfig, error) {
	var k = koanf.New(".")

	// default config
	_ = k.Set("log.enable-file", false)

	if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	r := &AppConfig{}
	if err := k.Unmarshal("", r); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return r, nil
}
