package config

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func createDefaultConfig() (k *koanf.Koanf) {
	k = koanf.New(".")
	_ = k.Set("log.enable-file", false)
	return
}

func unmarshal(k *koanf.Koanf) (*domain.AppConfig, error) {
	r := &domain.AppConfig{}
	err := k.Unmarshal("", r)
	return r, err
}

func DefaultConfig() *domain.AppConfig {
	k := createDefaultConfig()
	r, _ := unmarshal(k)
	return r
}

func ReadConfig(configFile string) (*domain.AppConfig, error) {
	var k = createDefaultConfig()

	if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if r, err := unmarshal(k); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	} else {
		return r, nil
	}
}
