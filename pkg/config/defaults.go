package config

const (
	// DefaultControllerName is a unique identifier which indicates this operator's name.
	DefaultControllerName = "l7mp.io/livekit-operator"

	// DefaultLabelKeyForConfigMap is a label key
	DefaultLabelKeyForConfigMap = "l7mp.io/livekit-config"

	// DefaultLabelValueForConfigMap is a label value
	DefaultLabelValueForConfigMap = "livekitmesh-config"

	// DefaultNamespaceForConfigMap is the default namespaced name of base livekit config map
	DefaultNamespaceForConfigMap = "default"

	// DefaultNameForConfigMap is the default namespaced name of base livekit config map
	DefaultNameForConfigMap = "default-livekit-config"

	// OwnedByLabelKey is the name of the label that is used to mark resources (Services,
	// ConfigMaps, and Deployments) dynamically created and maintained by the operator. Note
	// that the Deployments and Services created by the operator will have both the AppLabelKey
	// and the OwnedByLabelKey labels set.
	OwnedByLabelKey = "livekit.stunner.l7mp.io/owned-by"

	// OwnedByLabelValue is the value of OwnedByLabelKey to indicate that a resource is
	// maintained by the operator.
	OwnedByLabelValue = "livekitmesh-operator"

	// RelatedLiveKitMeshKey is the name of the label that is used to mark resources (Services,
	// ConfigMaps, and Deployments) dynamically created and maintained by the operator. Note
	// that the Deployments and Services created by the operator will have both the AppLabelKey
	// and the OwnedByLabelKey labels set.
	RelatedLiveKitMeshKey = "livekit.stunner.l7mp.io/livekit-mesh-name"

	// RelatedComponent is the name of the label that is used to determine which component this resource belongs to
	RelatedComponent = "livekit.stunner.l7mp.io/mesh-component"

	// DefaultLiveKitConfigFileName is the key of the livekit config in the config map data field
	DefaultLiveKitConfigFileName = "config.yaml"

	// DefaultClusterIssuerSecretApiTokenKey is the default api token key in the cluster issuer's secret
	DefaultClusterIssuerSecretApiTokenKey = "api-token"
)

// Statuses for the LiveKitMesh
const (
	// StatusNone Component is not present.
	StatusNone = "NONE"
	// StatusUpdating Component is being updated to a different version.
	StatusUpdating = "UPDATING"
	// StatusReconciling Controller has started but not yet completed reconciliation loop for the component.
	StatusReconciling = "RECONCILING"
	// StatusHealthy Component is healthy.
	StatusHealthy = "HEALTHY"
	// StatusError Component is in an error state.
	StatusError = "ERROR"
	// StatusActionRequired Action is needed from the user for reconciliation to proceed
	StatusActionRequired = "ACTION_REQUIRED"
)

// Component names
const (
	ComponentLiveKit     = "LIVEKIT"
	ComponentIngress     = "INGRESS"
	ComponentEgress      = "EGRESS"
	ComponentGateway     = "GATEWAY"
	ComponentCertManager = "CERTMANAGER"
	ComponentMonitoring  = "MONITORING"
	ComponentStunner     = "STUNNER"
)

// Issuer challenge provider types
const (
	IssuerCloudFlare   = "cloudflare"
	IssuerCloudDNS     = "clouddns"
	IssuerRoute53      = "route53"
	IssuerDigitalOcean = "digitalocean"
	IssuerAzureDNS     = "azuredns"
)

// Helm default values
const (
	//STUNNER

	// StunnerGatewayChartNamespace is the default namespace where the STUNner chart should be deployed
	StunnerGatewayChartNamespace = "stunner-gateway-system"

	//ENVOY

	// EnvoyGatewayChartNamespace is the default namespace where the Envoy chart should be deployed
	EnvoyGatewayChartNamespace = "envoy-gateway-system"

	//CERT-MANAGER

	// CertManagerChartNamespace is the default namespace where the Cert-Manager chart should be deployed
	CertManagerChartNamespace = "cert-manager"
)
