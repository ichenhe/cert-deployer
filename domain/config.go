package domain

type AppConfig struct {
	Log            LogConfig                `koanf:"log"`
	CloudProviders map[string]CloudProvider `koanf:"cloud-providers"`
}

type LogConfig struct {
	EnableFile bool   `koanf:"enable-file"`
	FileDir    string `koanf:"file-dir"`
	Level      string `koanf:"level"`
}

type CloudProvider struct {
	Provider  string `koanf:"provider"`
	SecretId  string `koanf:"secret-id"`
	SecretKey string `koanf:"secret-key"`
}
