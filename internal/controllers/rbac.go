package controllers

// RBAC for directly watched resources.
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io,resources=livekitmeshes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io,resources=livekitmeshes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io,resources=livekitmeshes/finalizers,verbs=update

// RBAC for references in watched resources.
// +kubebuilder:rbac:groups=core,resources=serviceaccounts;services;secrets;configmaps;,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles;clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways;gatewayclasses;tcproutes;httproutes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stunner.l7mp.io,resources=udproutes;gatewayconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=issuers,verbs=get;list;watch;create;update;patch;delete

// Helm chart related RBAC
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles;clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stunner.l7mp.io,resources=dataplanes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.k8s.io,resources=jobs,verbs=get;list;watch;create;update;patch;delete
