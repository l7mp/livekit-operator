/*
Copyright 2023 Kornel David.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	stnrgwv1 "github.com/l7mp/stunner-gateway-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type TurnServer struct {
	Credential *string `yaml:"credential" json:"credential,omitempty"`
	Host       *string `yaml:"host" json:"host,omitempty"`
	Port       *int    `yaml:"port" json:"port,omitempty"`
	Protocol   *string `yaml:"protocol" json:"protocol,omitempty"`
	Username   *string `yaml:"username" json:"username,omitempty"`
	AuthURI    *string `yaml:"uri,omitempty" json:"uri,omitempty"`
}

type Container struct {
	// Container image name.
	//
	// +optional
	Image string `json:"image,omitempty"`

	// Image pull policy. One of Always, Never, IfNotPresent.
	//
	// +optional
	ImagePullPolicy *corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// Entrypoint array.
	//
	// +optional
	Command []string `json:"command,omitempty"`

	// Arguments to the entrypoint.
	//
	// +optional
	Args []string `json:"args,omitempty"`

	// List of environment variables to set in the stunnerd container.
	//
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`

	// Resources required by stunnerd.
	//
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// Optional duration in seconds the stunnerd needs to terminate gracefully. Defaults to 3600 seconds.
	//
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`

	// Host networking requested for the stunnerd pod to use the host's network namespace.
	// Can be used to implement public TURN servers with Kubernetes.  Defaults to false.
	//
	// +kubebuilder:default=false
	// +optional
	HostNetwork bool `json:"hostNetwork,omitempty"`

	// Scheduling constraints.
	//
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// SecurityContext holds pod-level security attributes and common container settings.
	//
	// +optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// If specified, the pod's tolerations.
	//
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// If specified, the health-check port.
	//
	// +optional
	HealthCheckPort *int `json:"healthCheckPort,omitempty"`

	// // If specified, the metrics collection port.
	// //
	// // +optional
	// MetricsEndpointPort *int `json:"metricsEndpointPort,omitempty"`
}

// NamespacedName is the namespaced name of the configmap that stores
// the base configuration for a given LiveKit container
type NamespacedName struct {
	// Namespace is the namespace of the configMap resource
	Namespace *string `json:"namespace"`
	// Name is the name of the configMap resource
	Name *string `json:"name"`
}

type Keys struct {
	AccessToken *map[string]string `yaml:"access_token" json:"access_token,omitempty"`

	//ApiKey *string `yaml:"api_key" json:"api_key,omitempty"`
	//
	//ApiSecret *string `yaml:"api_secret" json:"api_secret,omitempty"`
}

// Redis holds the configuration for the Redis deployment
// If the field is configured, no redis will be created by the operator, in this case
// user deployed Redis deployment is required
// If it is omitted, a default Redis will be created.
type Redis struct {
	Address  *string `yaml:"address" json:"address,omitempty"`
	Username *string `yaml:"username" json:"username,omitempty"`
	Password *string `yaml:"password" json:"password,omitempty"`
	Db       *string `yaml:"db" json:"db,omitempty"`
}

type Rtc struct {
	PortRangeEnd   *int `yaml:"port_range_end" json:"port_range_end,omitempty"`
	PortRangeStart *int `yaml:"port_range_start" json:"port_range_start,omitempty"`
	TcpPort        *int `yaml:"tcp_port" json:"tcp_port,omitempty"`
}

//type IngressAddresses struct {
//	RtmpBaseUrl *string `json:"rtmp_base_url,omitempty"`
//	WhipBaseUrl *string `json:"whip_base_url,omitempty"`
//}

type LiveKitConfig struct {
	Keys     *map[string]string `yaml:"keys" json:"keys,omitempty"`
	LogLevel *string            `yaml:"log_level" json:"log_level,omitempty"`
	Port     *int               `yaml:"port" json:"port,omitempty"`

	// Redis in case redis is configured no Redis resources will be created by the operator
	// The value of this field inherited to the Ingress and Egress config as well
	Redis *Redis `yaml:"redis" json:"redis,omitempty"`
	Rtc   *Rtc   `yaml:"rtc" json:"rtc,omitempty"`
	//IngressAddresses *IngressAddresses `yaml:"ingress" json:"ingress,omitempty"`
}

type Deployment struct {

	// Name is the name of the Deployment that will be created.
	// Optional, if not filled default name 'livekit' will be used
	// Note that the same namespace will be used as the CR was deployed into.
	//
	// +kubebuilder:default=livekit-server
	// +optional
	Name *string `json:"name"`

	// Replicas Number of desired pods. This is a pointer to distinguish between explicit zero and not
	// specified. Defaults to 1.
	//
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Container template for the containers created in each Pod in the replicaset.
	// If omitted the default template will be used. Which spawns a single container
	//TODO
	//
	// +optional
	Container *Container `json:"container"`

	// Config holds the configuration for the livekit server
	Config *LiveKitConfig `json:"config"`
}

type LiveKit struct {

	// +kubebuilder:validation:Enum=livekit
	// +kubebuilder:default=livekit
	// +optional
	Type string `json:"type,omitempty"`

	// Deployment holds the configuration for the future Deployment manifest that will be created
	// by the operator.
	//
	// +kubebuilder:validation:Required
	Deployment *Deployment `json:"deployment"`
}

type Issuer struct {

	/*	// Name of the issuer resource
		// +kubebuilder:default=livekitmesh-issuer
		// +optional
		Name *string `json:"name"`*/

	// Email is the email address to be associated with the ACME account.
	// This field is optional, but it is strongly recommended to be set.
	// It will be used to contact you in case of issues with your account or
	// certificates, including expiry notification emails.
	// This field may be updated after the account is initially registered.
	// +optional
	Email *string `json:"email,omitempty"`

	// ChallengeSolver Used to configure a DNS01 challenge provider to be used when solving DNS01 challenges.
	//
	// +kubebuilder:validation:Enum=cloudflare;route53;clouddns;digitalocean;azuredns
	// +kubebuilder:validation:Required
	ChallengeSolver *string `json:"challengeSolver"`

	// ApiToken is the API token for the CloudFlare account that owns the challenged DnsZone
	//
	// +kubebuilder:validation:Required
	ApiToken *string `json:"apiToken"`
}

type Rtmp struct {
	Port *int `yaml:"port" json:"port,omitempty"`
}

type Whip struct {
	Port *int `yaml:"port" json:"port,omitempty"`
}

type IngressConfig struct {
	// +optional
	CPUCost *CPUCost `yaml:"cpu_cost" json:"cpu_cost"`
	// +optional
	HealthPort *int `yaml:"health_port" json:"health_port"`
	// +optional
	PrometheusPort *int `yaml:"prometheus_port" json:"prometheus_port"`
	// Rtmp holds the configuration for the RTMP service
	// If specified the all the necessary rescources (e.g.: gatewayclass, gateway listener, tcp route) will be created.
	// If omitted completely no resources will be created for it.
	// +kubebuilder:default=1935
	// +optional
	RTMPPort *int `yaml:"rtmp_port" json:"rtmp_port"`
	// Whip holds the configuration for the WHIP service
	// If specified the all the necessary rescources (e.g.: gatewayclass, gateway listener, http route) will be created.
	// If omitted completely no resources will be created for it.
	// +kubebuilder:default=8080
	// +optional
	WHIPPort *int `yaml:"whip_port" json:"whip_port"`
	// +optional
	HTTPRelayPort *int `yaml:"http_relay_port" json:"http_relay_port"`
	// +optional
	Logging *Logging `yaml:"logging" json:"logging"`
}

type CPUCost struct {
	// +optional
	RTMPCPUCost *int `yaml:"rtmp_cpu_cost" json:"rtmp_cpu_cost"`
	// +optional
	WHIPCPUCost *int `yaml:"whip_cpu_cost" json:"whip_cpu_cost"`
	// +optional
	WHIPBypassTranscodingCPUCost *int `yaml:"whip_bypass_transcoding_cpu_cost" json:"whip_bypass_transcoding_cpu_cost"`
}

type Logging struct {
	Level string `yaml:"level" json:"level"`
}

// Ingress is the LiveKit tool not the gateway resource to ingest traffic into the cluster
type Ingress struct {

	// Config holds configuration for the LiveKit Ingress
	Config *IngressConfig `yaml:"config" json:"config,omitempty"`
}

type EgressConfig struct {
	// +optional
	HealthPort *int `yaml:"health_port" json:"health_port"`
	// +optional
	TemplatePort *int `yaml:"template_port" json:"template_port"`
	// +optional
	PrometheusPort *int `yaml:"prometheus_port" json:"prometheus_port"`
	// +optional
	LogLevel *string `yaml:"log_level" json:"log_level"`
	//TemplateBase
	//EnableChromeSandbox
	// +optional
	Insecure *bool `yaml:"insecure" json:"insecure"`
	// +optional
	S3 *S3 `yaml:"s3" json:"s3"`
	// +optional
	Azure *Azure `yaml:"azure" json:"azure"`
	// +optional
	Gcp *Gcp `yaml:"gcp" json:"gcp"`
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

type Egress struct {
	// Config holds configuration for the LiveKit Ingress
	Config *EgressConfig `yaml:"config" json:"config,omitempty"`

	// Container template for the containers created in each Pod in the replicaset.
	// If omitted the default template will be used. Which spawns a single container
	//TODO
	//
	// +optional
	Container *Container `json:"container"`
}

type CertManager struct {
	// Issuer holds the necessary configuration about the used Issuer
	//
	// +optional
	Issuer *Issuer `json:"issuer,omitempty"`
}

type Monitoring struct {
	//TODO
}

type Stunner struct {

	// GatewayConfig is the configuration for the STUNner deployment's GatewayConfig.spec object
	GatewayConfig *stnrgwv1.GatewayConfigSpec `yaml:"gatewayConfig" json:"gatewayConfig"`

	// GatewayListeners is the configuration of the STUNner deployment's Gateway object
	GatewayListeners []gwapiv1.Listener `yaml:"gatewayListeners" json:"gatewayListeners"`
}

type CloudFlare struct {
	// Token is the API Token to authenticate the external DNS requests towards the CloudFlare servers
	Token *string `yaml:"token" json:"token"`

	// Email is the email address of the account associated with the Token provided
	Email *string `yaml:"email" json:"email"`
}

type ExternalDNS struct {

	// CloudFlare holds the necessary configuration for the CloudFlare ExternalDNS integration.
	// Currently, only CloudFlare is supported by the operator.
	CloudFlare *CloudFlare `yaml:"cloudFlare" json:"cloudFlare"`
}

type ApplicationExpose struct {

	// HostName is the DNS host name that will be used by both cert-manager and Envoy GW.
	// DnsZone Certificate requests will be issues against this HostName.
	// This ChallengeSolver will use this to solve the challenge.
	//
	// +kubebuilder:validation:Required
	HostName *string `yaml:"hostName" json:"hostName"`

	// CertManager will obtain certificates from a variety of Issuers, both popular public Issuers and private Issuers,
	// and ensure the certificates are valid and up-to-date, and will attempt to renew certificates at a configured time before expiry.
	//
	//
	// +optional
	CertManager *CertManager `yaml:"certManager" json:"certManager,omitempty"`

	// ExternalDNS creates the A record in the DNS provider configured.
	// Currently, only CloudFlare is support. To use other DNS services please configure your own External DNS deployment
	// and omit this field.
	// Make sure to enable RBAC resources to watch HTTPRoutes and/or services.
	// See more: https://github.com/kubernetes-sigs/external-dns/tree/master/docs/tutorials
	//
	// +optional
	ExternalDNS *ExternalDNS `yaml:"externalDNS" json:"externalDNS"`
}

type Component struct {

	// LiveKit is the main resource that the operator manages. By default, it supports
	// only the LiveKit server as a media server but in the future it might support other
	// media servers as well.
	// +kubebuilder:validation:Required
	LiveKit *LiveKit `json:"liveKit"`

	// LiveKit's Ingress resource descriptor.
	// This resource makes it possible to stream videos(prerecorded or live) into
	// the Kubernetes cluster and further into a chosen room. Note that this resource
	// enables a one-way communication between the client and the media server.
	//
	// +optional
	Ingress *Ingress `json:"ingress,omitempty"`

	// LiveKit's Egress resource descriptor.
	// Egress makes it possible to stream a single user's or any number of users'
	// streams out of a room onto an RTMP port.
	//
	// +optional
	Egress *Egress `json:"egress,omitempty"`

	// ApplicationExpose is the component that contains all required subcomponents that are accountable for exposing the
	// application to the internet on a secure, encrypted way. This includes the Cert-manager, Envoy GW, and ExternalDns.
	// +kubebuilder:validation:Required
	ApplicationExpose *ApplicationExpose `json:"applicationExpose"`

	// Monitoring enables the Prometheus metric exposition, installs
	// a Prometheus operator and Grafana operator with the corresponding resources
	//
	// +optional
	Monitoring *Monitoring `json:"monitoring,omitempty"`

	// Stunner holds the configuration for the to-be created STUNner resources.
	// These resources provide the necessary resources for a possible/successful TURN connection
	//
	// +optional
	Stunner *Stunner `json:"stunner"`
}

// LiveKitMeshSpec defines the desired state of LiveKitMesh
type LiveKitMeshSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	Components *Component `json:"components"`
}

// InstallStatus is the status of the component.
// Enum with the following possible values:
// -NONE meaning the component is not present
// -UPDATING meaning the component is being updated to a different version
// -RECONCILING meaning the controller has started but not yet completed reconciliation loop for the component
// -HEALTHY meaning the component is healthy
// -ERROR meaning a critical error happened to the component
// -ACTION_REQUIRED meaning there is a user action needed in order to proceed
// +kubebuilder:validation:Enum=NONE;UPDATING;RECONCILING;HEALTHY;ERROR;ACTION_REQUIRED
type InstallStatus string

// LiveKitMeshStatus defines the observed state of LiveKitMesh
type LiveKitMeshStatus struct {

	// ComponentStatus is a key-value store to signal the components' status after installation
	// The map will give a brief overview for the user which component was successful or failed etc.
	// THE FIELD IS POPULATED BY THE OPERATOR NOT BY THE USER. IT WILL BE OVERWRITTEN
	ComponentStatus map[string]InstallStatus `json:"componentStatus"`

	// OverallStatus of all components controlled by the operator.
	//
	// * If all components have status `NONE`, overall status is `NONE`.
	// * If all components are `HEALTHY`, overall status is `HEALTHY`.
	// * If one or more components are `RECONCILING` and others are `HEALTHY`, overall status is `RECONCILING`.
	// * If one or more components are `UPDATING` and others are `HEALTHY`, overall status is `UPDATING`.
	// * If components are a mix of `RECONCILING`, `UPDATING` and `HEALTHY`, overall status is `UPDATING`.
	// * If any component is in `ERROR` state, overall status is `ERROR`.
	// * If further action is needed for reconciliation to proceed, overall status is `ACTION_REQUIRED`.
	//
	OverallStatus *InstallStatus `json:"overallStatus"`

	// ConfigStatus holds the current configuration for the LiveKit component
	// if it is available in the cluster. nil meaning the ConfigMap provided is not present.
	//
	ConfigStatus *string `json:"configStatus"`
}

// LiveKitMesh is the Schema for the livekitmeshes API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type LiveKitMesh struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LiveKitMeshSpec   `json:"spec,omitempty"`
	Status LiveKitMeshStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LiveKitMeshList contains a list of LiveKitMesh
type LiveKitMeshList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LiveKitMesh `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LiveKitMesh{}, &LiveKitMeshList{})
}
