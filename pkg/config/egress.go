package config

import lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"

type EgressConfig struct {
	APIKey         *string          `yaml:"api_key" json:"api_key"`
	APISecret      *string          `yaml:"api_secret" json:"api_secret"`
	WSURL          *string          `yaml:"ws_url" json:"ws_url"`
	Redis          *lkstnv1a1.Redis `yaml:"redis" json:"redis"`
	HealthPort     *int             `yaml:"health_port" json:"health_port"`
	TemplatePort   *int             `yaml:"template_port" json:"template_port"`
	PrometheusPort *int             `yaml:"prometheus_port" json:"prometheus_port"`
	LogLevel       *string          `yaml:"log_level" json:"log_level"`
	//TemplateBase
	//EnableChromeSandbox
	Insecure *bool  `yaml:"insecure" json:"insecure"`
	S3       *S3    `yaml:"s3" json:"s3"`
	Azure    *Azure `yaml:"azure" json:"azure"`
	Gcp      *Gcp   `yaml:"gcp" json:"gcp"`
}

type S3 struct {
	AccessKey *string `yaml:"access_key" json:"access_key"`
	Secret    *string `yaml:"secret" json:"secret"`
	Region    *string `yaml:"region" json:"region"`
	Endpoint  *string `yaml:"endpoint" json:"endpoint"`
	Bucket    *string `yaml:"bucket" json:"bucket"`
}

type Azure struct {
	AccountName   *string `yaml:"account_name" json:"account_name"`
	AccountKey    *string `yaml:"account_key" json:"account_key"`
	ContainerName *string `yaml:"container_name" json:"container_name"`
}

type Gcp struct {
	CredentialsJson *string `yaml:"credentials_json" json:"credentials_json"`
	Bucket          *string `yaml:"bucket" json:"bucket"`
}

func ConvertEgressConfig(config lkstnv1a1.EgressConfig) *EgressConfig {
	egressConfig := &EgressConfig{}
	if config.HealthPort != nil {
		egressConfig.HealthPort = config.HealthPort
	}
	if config.TemplatePort != nil {
		egressConfig.TemplatePort = config.TemplatePort
	}
	if config.PrometheusPort != nil {
		egressConfig.PrometheusPort = config.PrometheusPort
	}
	if config.LogLevel != nil {
		egressConfig.LogLevel = config.LogLevel
	}
	if config.Insecure != nil {
		egressConfig.Insecure = config.Insecure
	}
	if config.S3 != nil {
		egressConfig.S3 = &S3{
			AccessKey: config.S3.AccessKey,
			Secret:    config.S3.Secret,
			Region:    config.S3.Region,
			Endpoint:  config.S3.Endpoint,
			Bucket:    config.S3.Bucket,
		}
	} else if config.Azure != nil {
		egressConfig.Azure = &Azure{
			AccountName:   config.Azure.AccountName,
			AccountKey:    config.Azure.AccountKey,
			ContainerName: config.Azure.ContainerName,
		}
	} else if config.Gcp != nil {
		egressConfig.Gcp = &Gcp{
			CredentialsJson: config.Gcp.CredentialsJson,
			Bucket:          config.Gcp.Bucket,
		}
	}
	return egressConfig
}
