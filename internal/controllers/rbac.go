package controllers

// RBAC for directly watched resources.
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io,resources=livekitmeshes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io,resources=livekitmeshes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io,resources=livekitmeshes/finalizers,verbs=update

// RBAC for references in watched resources.
// +kubebuilder:rbac:groups=core,resources=serviceaccounts;services;secrets;configmaps;,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways;gatewayclasses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stunner.l7mp.io,resources=udproutes;gatewayconfigs,verbs=get;list;watch;create;update;patch;delete
