package config

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/domain"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

func createDefaultConfig() (k *koanf.Koanf) {
	k = koanf.New(".")
	logDrivers := []domain.LogDriver{
		{Driver: "stdout", Level: "info"},
	}
	_ = k.Set("log-drivers", logDrivers)
	return
}

func unmarshal(k *koanf.Koanf) (*domain.AppConfig, error) {
	r := &domain.AppConfig{}
	err := k.Unmarshal("", r)
	return r, err
}

func resolveProfilePath(c *cli.Context) string {
	if c != nil {
		if profilePath := c.Path("profile"); profilePath != "" {
			return profilePath
		}
	}

	searchDirs := make([]string, 0)
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		searchDirs = append(searchDirs, dir, filepath.Join(dir, "config"))
	}
	if dir, err := os.Getwd(); err == nil {
		searchDirs = append(searchDirs, dir, filepath.Join(dir, "config"))
	}

	fileName := "cert-deployer"
	exts := []string{".yaml", ".yml"}
	for _, dir := range searchDirs {
		for _, ext := range exts {
			path := filepath.Join(dir, fileName+ext)
			if isFile(path) {
				return path
			}
		}
	}

	return ""
}

// CreateWithModifier creates a default configuration, tries to resolve profile path and load it if valid.
// Optional modifier will be applied before parsing and unmarshalling the config.
func CreateWithModifier(c *cli.Context, modifier func(k *koanf.Koanf)) (*domain.AppConfig, error) {
	k := createDefaultConfig()

	// try to load from profile
	if profile := resolveProfilePath(c); profile != "" {
		// found valid profile, load it
		if err := k.Load(file.Provider(profile), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("failed to load config file '%s': %w", profile, err)
		}
	}

	if modifier != nil {
		modifier(k)
	}

	return parseAndVerifyConfig(k)
}

func parseAndVerifyConfig(k *koanf.Koanf) (*domain.AppConfig, error) {
	var config *domain.AppConfig
	config, err := unmarshal(k)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// infer the deployment's name
	for name, deployment := range config.Deployments {
		deployment.Name = name
		config.Deployments[name] = deployment
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
