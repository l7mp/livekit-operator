package renderer

import (
	"github.com/l7mp/livekit-operator/internal/store"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Renderer) renderLiveKitEgressConfigMap(renderContext *RenderContext) {
	log := r.logger.WithName("renderLiveKitEgressConfigMap")

	log.V(2).Info("trying to render LiveKit-Egress ConfigMap")

	lkMesh := renderContext.liveKitMesh
	configMap, err := createLiveKitEgressConfigMap(lkMesh)
	if err != nil {
		log.Error(err, "cannot create livekit-egress config map")
		return
	}
	if err := controllerutil.SetOwnerReference(lkMesh, configMap, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(configMap))
		return
	}

	renderContext.update.UpsertQueue.ConfigMaps.Upsert(configMap)

	log.V(2).Info("Upserted LiveKit-Egress ConfigMap into UpsertQueue", "cm", store.GetObjectKey(configMap))
}

func (r *Renderer) renderLiveKitEgressDeployment(renderContext *RenderContext) {
	log := r.logger.WithName("renderLiveKitEgressDeployment")

	log.V(2).Info("trying to render LiveKit-Egress Deployment")
	lkMesh := renderContext.liveKitMesh

	deployment := createLiveKitEgressDeployment(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, deployment, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(deployment))
		return
	}

	renderContext.update.UpsertQueue.Deployments.Upsert(deployment)
	log.V(2).Info("Upserted LiveKit-Egress Deployment into UpsertQueue", "deployment", store.GetObjectKey(deployment))
}
