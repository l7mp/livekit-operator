package renderer

import (
	"fmt"
	"github.com/l7mp/livekit-operator/internal/store"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Renderer) renderLiveKitConfigMap(renderContext *RenderContext) {
	log := r.logger.WithName("renderLiveKitConfigMap")

	log.V(2).Info("trying to render LiveKit-Server ConfigMap")

	lkMesh := renderContext.liveKitMesh

	iceConfig, err := getIceConfigurationFromStunnerAuth(*lkMesh, log)
	if err != nil {
		log.Error(err, "Failed to get ICE configuration from STUNner auth")
		return
	} else if iceConfig != nil {
		address := getAddressFromIceConfig(iceConfig)
		if validateIPAddress(address) {
			log.V(1).Info("Valid IP address found for STUNner", "address", address)
		} else {
			err := fmt.Errorf("invalid turn address found: %s", address)
			log.Error(err, "Failed to render LiveKit-Server ConfigMap")
			return
		}
	} else {
		//iceConfig is nil
		return
	}

	cm, err := createLiveKitConfigMap(lkMesh, *iceConfig)
	if err != nil {
		log.Error(err, "cannot create livekit config map")
		return
	}
	if err := controllerutil.SetOwnerReference(lkMesh, cm, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(cm))
		return
	}

	renderContext.update.UpsertQueue.ConfigMaps.Upsert(cm)

	log.V(2).Info("Upserted LiveKit-Server ConfigMap into UpsertQueue", "cm", store.GetObjectKey(cm))
}

func (r *Renderer) renderLiveKitService(renderContext *RenderContext) {
	log := r.logger.WithName("renderLiveKitService")

	log.V(2).Info("trying to render LiveKit-Server Service")

	lkMesh := renderContext.liveKitMesh
	service := createLiveKitService(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, service, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(service))
		return
	}

	renderContext.update.UpsertQueue.Services.Upsert(service)

	log.V(2).Info("Upserted LiveKit-Server Service into UpsertQueue", "cm", store.GetObjectKey(service))
}

func (r *Renderer) renderLiveKitDeployment(renderContext *RenderContext) {
	log := r.logger.WithName("renderLiveKitDeployment")

	log.V(2).Info("trying to render LiveKit-Server Deployment")

	lkMesh := renderContext.liveKitMesh
	deployment := createLiveKitDeployment(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, deployment, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(deployment))
		return
	}

	renderContext.update.UpsertQueue.Deployments.Upsert(deployment)

	log.V(2).Info("Upserted LiveKit-Server Deployment into UpsertQueue", "cm", store.GetObjectKey(deployment))
}

func (r *Renderer) renderLiveKitRedis(renderContext *RenderContext) {
	log := r.log.WithName("renderLiveKitRedis")

	log.V(2).Info("trying to render Redis StatefulSets and Service")

	lkMesh := renderContext.liveKitMesh
	redis := lkMesh.Spec.Components.LiveKit.Deployment.Config.Redis

	if redis == nil {
		log.V(2).Info("creation of a Redis deployment is required due to empty configuration")
		redisStatefulSet, redisService, redisConfigMap := createLiveKitRedis(lkMesh)
		if err := controllerutil.SetOwnerReference(lkMesh, redisStatefulSet, r.scheme); err != nil {
			log.Error(err, "cannot set owner reference", "owner",
				store.GetObjectKey(lkMesh), "reference",
				store.GetObjectKey(redisStatefulSet))
			return
		}
		if err := controllerutil.SetOwnerReference(lkMesh, redisService, r.scheme); err != nil {
			log.Error(err, "cannot set owner reference", "owner",
				store.GetObjectKey(lkMesh), "reference",
				store.GetObjectKey(redisService))
			return
		}

		if err := controllerutil.SetOwnerReference(lkMesh, redisConfigMap, r.scheme); err != nil {
			log.Error(err, "cannot set owner reference", "owner",
				store.GetObjectKey(lkMesh), "reference",
				store.GetObjectKey(redisConfigMap))
			return
		}

		renderContext.update.UpsertQueue.StatefulSets.Upsert(redisStatefulSet)

		log.V(2).Info("Upserted Redis StatefulSets into UpsertQueue", "ss", store.GetObjectKey(redisStatefulSet))

		renderContext.update.UpsertQueue.Services.Upsert(redisService)

		log.V(2).Info("Upserted Redis Service into UpsertQueue", "svc", store.GetObjectKey(redisService))

		renderContext.update.UpsertQueue.ConfigMaps.Upsert(redisConfigMap)

		log.V(2).Info("Upserted Redis ConfigMap into UpsertQueue", "cm", store.GetObjectKey(redisConfigMap))

		//TODO why was the below added on the first hand?
		/*	renderContext.update.UpsertQueue.ConfigMaps.Get(types.NamespacedName{
			Namespace: lkMesh.GetNamespace(),
			Name:      getLiveKitServerConfigMapName(*lkMesh.Spec.Components.LiveKit.Deployment.Name),
		})*/

	} else {
		log.V(2).Info("creation of a Redis deployment is NOT required due to configuration", "redis", redis)
	}
}
