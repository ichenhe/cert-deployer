package domain

type AppConfig struct {
	Log            LogConfig                 `koanf:"log"`
	CloudProviders map[string]CloudProvider  `koanf:"cloud-providers"`
	Deployments    map[string]Deployment     `koanf:"deployments"`
	Triggers       map[string]TriggerDefiner `koanf:"-"` // this field is set manually
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

type Deployment struct {
	Name       string            // name of this deployment, infer from the map
	ProviderId string            `koanf:"provider-id"`
	Cert       string            `koanf:"cert"` // path to full chain pem
	Key        string            `koanf:"key"`  // path to private pem
	Assets     []DeploymentAsset `koanf:"assets"`
}

type DeploymentAsset struct {
	Type string `koanf:"type"`
	Id   string `koanf:"id"`
}

type TriggerDefiner interface {
	GetName() string
	GetType() string
	GetDeploymentIds() []string
}

type triggerBaseInfo struct {
	Name        string
	Type        string   `koanf:"type"`
	Deployments []string `koanf:"deployments"`
}

func (t triggerBaseInfo) GetName() string {
	return t.Name
}

func (t triggerBaseInfo) GetType() string {
	return t.Type
}

func (t triggerBaseInfo) GetDeploymentIds() []string {
	return t.Deployments
}

type FileMonitoringTriggerOptions struct {
	File string `koanf:"file"`
	Wait int    `koanf:"wait"`
}

type FileMonitoringTriggerDef struct {
	triggerBaseInfo `koanf:",squash"`
	Options         FileMonitoringTriggerOptions `koanf:"options"`
}
