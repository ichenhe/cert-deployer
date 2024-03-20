package config

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func ReadConfig(configFile string) (*domain.AppConfig, error) {
	var k = koanf.New(".")

	// default config
	_ = k.Set("log.enable-file", false)

	if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	r := &domain.AppConfig{}
	if err := k.Unmarshal("", r); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return r, nil
}
