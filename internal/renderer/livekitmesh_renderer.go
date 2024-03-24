package renderer

import (
	"github.com/l7mp/livekit-operator/internal/event"
	"github.com/l7mp/livekit-operator/internal/store"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Renderer) RenderLiveKitMesh(e *event.Render) {
	log := r.logger.WithName("RenderLiveKitMesh")
	log.Info("Trying to render LiveKitMeshes")
	//TODO render each component in the livekitmesh
	r.gen += 1

	if liveKitMeshes := store.LiveKitMeshes.GetAll(); len(liveKitMeshes) == 0 {
		log.Info("no LiveKitMesh objects found", "event", e.String())
	} else {

		for _, lkMesh := range liveKitMeshes {
			log.Info("Found in store", "lk", lkMesh.Name)
			lkMesh := lkMesh
			renderContext := NewRenderContext(e, r, lkMesh)
			gw := renderContext.liveKitMesh.Spec.Components.Gateway

			// if the PR goes in then this whole if-else block should be removed
			if gw != nil {
				log.V(2).Info("Gateway is configured, looking for loadbalancerip for the LiveKit config")
				if addr := r.getLoadBalancerIP(r.logger, gw); addr == nil {
					log.Info("LoadBalancerIP is not present yet for", "Gateway",
						types.NamespacedName{
							Namespace: *gw.RelatedStunnerGatewayAnnotations.Namespace,
							Name:      *gw.RelatedStunnerGatewayAnnotations.Name,
						}.String())
					continue
				} else {
					log.Info("LoadBalancerIP is present for", "Gateway",
						types.NamespacedName{
							Namespace: *gw.RelatedStunnerGatewayAnnotations.Namespace,
							Name:      *gw.RelatedStunnerGatewayAnnotations.Name,
						}.String(),
						"addr", addr)
					renderContext.turnServerPublicAddress = addr
				}
			}

			// gw is not configured in the LiveKitMesh
			// this will be supported way to render the configmap however
			r.renderLiveKitComponentResources(renderContext)

			r.renderCertManagerComponentResources(renderContext)

			//renderContext.liveKit
			log.Info("event to channel")
			r.operatorCh <- renderContext.update
		}
	}
}

func (r *Renderer) renderLiveKitComponentResources(renderContext *RenderContext) {
	r.renderLiveKitConfigMap(renderContext)
	r.renderLiveKitDeployment(renderContext)
	r.renderLiveKitService(renderContext)
	r.renderLiveKitRedis(renderContext)
}

func (r *Renderer) renderCertManagerComponentResources(renderContext *RenderContext) {
	r.renderCertManagerIssuerAndSecret(renderContext)
}

func (r *Renderer) renderCertManagerIssuerAndSecret(context *RenderContext) {
	log := r.logger.WithName("renderCertManagerIssuer")

	log.V(2).Info("trying to render Issuer", "component", "Cert-Manager")
	lkMesh := context.liveKitMesh

	issuer, secret := createIssuer(*lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, issuer, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(issuer))
		return
	}
	if err := controllerutil.SetOwnerReference(lkMesh, secret, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(secret))
		return
	}

	context.update.UpsertQueue.Issuer.Upsert(issuer)

	log.V(2).Info("Upserted Cert-Manager Issuer into UpsertQueue", "cm", store.GetObjectKey(issuer))

	context.update.UpsertQueue.Secrets.Upsert(secret)

	log.V(2).Info("Upserted Cert-Manager Secrets into UpsertQueue", "cm", store.GetObjectKey(secret))
}

func (r *Renderer) renderLiveKitConfigMap(context *RenderContext) {
	log := r.logger.WithName("renderLiveKitConfigMap")

	log.V(2).Info("trying to render LiveKit-Server ConfigMap")

	lkMesh := context.liveKitMesh

	cm, err := createLiveKitConfigMap(lkMesh, context.turnServerPublicAddress)
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

	context.update.UpsertQueue.ConfigMaps.Upsert(cm)

	log.V(2).Info("Upserted LiveKit-Server ConfigMap into UpsertQueue", "cm", store.GetObjectKey(cm))
}

func (r *Renderer) renderLiveKitService(context *RenderContext) {
	log := r.logger.WithName("renderLiveKitService")

	log.V(2).Info("trying to render LiveKit-Server Service")

	lkMesh := context.liveKitMesh
	service := createLiveKitService(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, service, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(service))
		return
	}

	context.update.UpsertQueue.Services.Upsert(service)

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
			Name:      ConfigMapNameFormat(*lkMesh.Spec.Components.LiveKit.Deployment.Name),
		})*/

	} else {
		log.V(2).Info("creation of a Redis deployment is NOT required due to configuration", "redis", redis)
	}
}
