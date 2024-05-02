package renderer

import (
	"github.com/l7mp/livekit-operator/internal/store"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Renderer) renderLiveKitIngressConfigMap(renderContext *RenderContext) {
	log := r.logger.WithName("renderLiveKitIngressConfigMap")

	log.V(2).Info("trying to render LiveKit-Ingress ConfigMap")

	lkMesh := renderContext.liveKitMesh
	configMap, err := createLiveKitIngressConfigMap(lkMesh)
	if err != nil {
		log.Error(err, "cannot create livekit-ingress config map")
		return
	}
	if err := controllerutil.SetOwnerReference(lkMesh, configMap, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(configMap))
		return
	}

	renderContext.update.UpsertQueue.ConfigMaps.Upsert(configMap)

	log.V(2).Info("Upserted LiveKit-Ingress ConfigMap into UpsertQueue", "cm", store.GetObjectKey(configMap))
}

func (r *Renderer) renderLiveKitIngressDeployment(renderContext *RenderContext) {
	log := r.logger.WithName("renderLiveKitIngressDeployment")

	log.V(2).Info("trying to render LiveKit-Ingress Deployment")
	lkMesh := renderContext.liveKitMesh

	deployment := createLiveKitIngressDeployment(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, deployment, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(deployment))
		return
	}

	renderContext.update.UpsertQueue.Deployments.Upsert(deployment)
	log.V(2).Info("Upserted LiveKit-Ingress Deployment into UpsertQueue", "deployment", store.GetObjectKey(deployment))
}

func (r *Renderer) renderLiveKitIngressService(renderContext *RenderContext) {
	log := r.logger.WithName("renderLiveKitIngressService")

	log.V(2).Info("trying to render LiveKit-Ingress Service")

	lkMesh := renderContext.liveKitMesh
	service := createLiveKitIngressService(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, service, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(service))
		return
	}

	renderContext.update.UpsertQueue.Services.Upsert(service)

	log.V(2).Info("Upserted LiveKit-Ingress Service into UpsertQueue", "cm", store.GetObjectKey(service))
}
