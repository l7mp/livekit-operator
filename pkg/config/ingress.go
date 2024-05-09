package config

import (
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"k8s.io/utils/ptr"
)

// IngressConfig This ingressConfig is used to create the configmap for ingress
type IngressConfig struct {
	APIKey         *string          `yaml:"api_key" json:"api_key"`
	APISecret      *string          `yaml:"api_secret" json:"api_secret"`
	WSURL          *string          `yaml:"ws_url" json:"ws_url"`
	Redis          *lkstnv1a1.Redis `yaml:"redis" json:"redis"`
	CPUCost        *CPUCost         `yaml:"cpu_cost" json:"cpu_cost"`
	HealthPort     *int             `yaml:"health_port" json:"health_port"`
	PrometheusPort *int             `yaml:"prometheus_port" json:"prometheus_port"`
	RTMPPort       *int             `yaml:"rtmp_port" json:"rtmp_port"`
	WHIPPort       *int             `yaml:"whip_port" json:"whip_port"`
	HTTPRelayPort  *int             `yaml:"http_relay_port" json:"http_relay_port"`
	Logging        *Logging         `yaml:"logging" json:"logging"`
}

type CPUCost struct {
	RTMPCPUCost                  *int `yaml:"rtmp_cpu_cost" json:"rtmp_cpu_cost"`
	WHIPCPUCost                  *int `yaml:"whip_cpu_cost" json:"whip_cpu_cost"`
	WHIPBypassTranscodingCPUCost *int `yaml:"whip_bypass_transcoding_cpu_cost" json:"whip_bypass_transcoding_cpu_cost"`
}

type Logging struct {
	Level *string `yaml:"level" json:"level"`
}

func ConvertIngressConfig(config lkstnv1a1.IngressConfig) *IngressConfig {
	ingressConfig := &IngressConfig{}

	if config.CPUCost != nil {
		ingressConfig.CPUCost = &CPUCost{}
		if config.CPUCost.RTMPCPUCost != nil {
			ingressConfig.CPUCost.RTMPCPUCost = config.CPUCost.RTMPCPUCost
		}
		if config.CPUCost.WHIPCPUCost != nil {
			ingressConfig.CPUCost.WHIPCPUCost = config.CPUCost.WHIPCPUCost
		}
		if config.CPUCost.WHIPBypassTranscodingCPUCost != nil {
			ingressConfig.CPUCost.WHIPBypassTranscodingCPUCost = config.CPUCost.WHIPBypassTranscodingCPUCost
		}
	}

	if config.HealthPort != nil {
		ingressConfig.HealthPort = config.HealthPort
	}
	if config.PrometheusPort != nil {
		ingressConfig.PrometheusPort = config.PrometheusPort
	}
	if config.RTMPPort != nil {
		ingressConfig.RTMPPort = config.RTMPPort
	} else {
		ingressConfig.RTMPPort = ptr.To(1935)
	}
	if config.WHIPPort != nil {
		ingressConfig.WHIPPort = config.WHIPPort
	} else {
		ingressConfig.WHIPPort = ptr.To(8080)
	}
	if config.HTTPRelayPort != nil {
		ingressConfig.HTTPRelayPort = config.HTTPRelayPort
	} else {
		ingressConfig.HTTPRelayPort = ptr.To(9090)
	}
	if config.Logging != nil {
		ingressConfig.Logging = &Logging{
			Level: &config.Logging.Level,
		}
	}
	return ingressConfig
}
