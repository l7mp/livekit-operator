package controllers

// RBAC for directly watched resources.
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io.l7mp.io,resources=livekitmeshes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io.l7mp.io,resources=livekitmeshes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io.l7mp.io,resources=livekitmeshes/finalizers,verbs=update

// RBAC for references in watched resources.
// +kubebuilder:rbac:groups=core,resources=services;secrets,verbs=get;list;watch;create;update;patch;delete
