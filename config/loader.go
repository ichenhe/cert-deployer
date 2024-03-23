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

	var config *domain.AppConfig
	config, err := unmarshal(k)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err = parseTriggers(k, config); err != nil {
		return nil, fmt.Errorf("failed to load triggers: %w", err)
	}

	if err = verifyDeploymentsReferencedInTriggerExist(config); err != nil {
		return nil, err
	}

	return config, nil
}

func verifyDeploymentsReferencedInTriggerExist(config *domain.AppConfig) error {
	for name, trigger := range config.Triggers {
		for _, deploymentName := range trigger.GetDeploymentIds() {
			if _, ex := config.Deployments[deploymentName]; !ex {
				return fmt.Errorf("deployment '%s' referenced in trigger '%s' does not exist", deploymentName, name)
			}
		}
	}
	return nil
}
