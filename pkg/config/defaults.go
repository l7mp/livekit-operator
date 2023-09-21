package config

const (
	// DefaultControllerName is a unique identifier which indicates this operator's name.
	DefaultControllerName = "l7mp.io/livekit-operator"

	// DefaultLabelKeyForConfigMap is a label key
	DefaultLabelKeyForConfigMap = "l7mp.io/livekit-config"

	// DefaultLabelValueForConfigMap is a label value
	DefaultLabelValueForConfigMap = "livekit-config"

	// DefaultNamespaceForConfigMap is the default namespaced name of base livekit config map
	DefaultNamespaceForConfigMap = "default"

	// DefaultNameForConfigMap is the default namespaced name of base livekit config map
	DefaultNameForConfigMap = "default-livekit-config"
)
