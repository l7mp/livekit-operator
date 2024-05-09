package config

import (
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"k8s.io/utils/ptr"
)

type LiveKitConfig struct {
	Keys             *map[string]string `yaml:"keys" json:"keys,omitempty"`
	LogLevel         *string            `yaml:"log_level" json:"log_level,omitempty"`
	Port             *int               `yaml:"port" json:"port,omitempty"`
	Redis            *lkstnv1a1.Redis   `yaml:"redis" json:"redis,omitempty"`
	Rtc              *Rtc               `yaml:"rtc" json:"rtc,omitempty"`
	IngressAddresses *IngressAddresses  `yaml:"ingress" json:"ingress,omitempty"`
}

type Rtc struct {
	PortRangeEnd   *int         `yaml:"port_range_end" json:"port_range_end,omitempty"`
	PortRangeStart *int         `yaml:"port_range_start" json:"port_range_start,omitempty"`
	TcpPort        *int         `yaml:"tcp_port" json:"tcp_port,omitempty"`
	StunServers    []string     `yaml:"stun_servers" json:"stun_servers,omitempty"`
	TurnServers    []TurnServer `yaml:"turn_servers" json:"turn_servers,omitempty"`
	// +kubebuilder:default=false
	// +optional
	UseExternalIp *bool `yaml:"use_external_ip" json:"use_external_ip,omitempty"`
}

type TurnServer struct {
	Credential *string `yaml:"credential" json:"credential,omitempty"`
	Host       *string `yaml:"host" json:"host,omitempty"`
	Port       *int    `yaml:"port" json:"port,omitempty"`
	Protocol   *string `yaml:"protocol" json:"protocol,omitempty"`
	Username   *string `yaml:"username" json:"username,omitempty"`
	AuthURI    *string `yaml:"uri,omitempty" json:"uri,omitempty"`
}

type IngressAddresses struct {
	RtmpBaseUrl *string `json:"rtmp_base_url,omitempty"`
	WhipBaseUrl *string `json:"whip_base_url,omitempty"`
}

func ConvertServerConfig(config lkstnv1a1.LiveKitConfig) *LiveKitConfig {
	liveKitConfig := &LiveKitConfig{}

	if config.Keys != nil {
		liveKitConfig.Keys = config.Keys
	}
	if config.LogLevel != nil {
		liveKitConfig.LogLevel = config.LogLevel
	}
	if config.Port != nil {
		liveKitConfig.Port = config.Port
	}
	if config.Redis != nil {
		liveKitConfig.Redis = &lkstnv1a1.Redis{}
		if config.Redis.Address != nil {
			liveKitConfig.Redis.Address = config.Redis.Address
		}
		if config.Redis.Password != nil {
			liveKitConfig.Redis.Password = config.Redis.Password
		}
		if config.Redis.Username != nil {
			liveKitConfig.Redis.Username = config.Redis.Username
		}
		if config.Redis.Db != nil {
			liveKitConfig.Redis.Db = config.Redis.Db
		}
	}

	liveKitConfig.Rtc = &Rtc{}
	if config.Rtc != nil {
		if config.Rtc.PortRangeEnd != nil {
			liveKitConfig.Rtc.PortRangeEnd = config.Rtc.PortRangeEnd
		}
		if config.Rtc.PortRangeStart != nil {
			liveKitConfig.Rtc.PortRangeStart = config.Rtc.PortRangeStart
		}
		if config.Rtc.TcpPort != nil {
			liveKitConfig.Rtc.TcpPort = config.Rtc.TcpPort
		}
		liveKitConfig.Rtc.UseExternalIp = ptr.To(false)
	}
	//if config.IngressAddresses != nil {
	//	liveKitConfig.IngressAddresses = &IngressAddresses{}
	//	if config.IngressAddresses.RtmpBaseUrl != nil {
	//		liveKitConfig.IngressAddresses.RtmpBaseUrl = config.IngressAddresses.RtmpBaseUrl
	//	}
	//}
	return liveKitConfig
}
