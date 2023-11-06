package renderer

import (
	"github.com/go-logr/logr"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/event"
	"github.com/l7mp/livekit-operator/internal/store"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	"k8s.io/apimachinery/pkg/types"
)

func (r *Renderer) RenderLiveKitMesh(e *event.Render) {
	log := r.logger.WithName("RenderLiveKitMesh")
	log.Info("Trying to render LiveKitMeshes")
	//TODO render each component in the livekitmesh
	r.gen += 1

	if liveKitMeshes := store.LiveKitMeshes.GetAll(); len(liveKitMeshes) == 0 {
		log.Info("no LiveKitMesh objects found", "event", e.String())
		//TODO maybe return here?
	} else {

		for _, lkMesh := range liveKitMeshes {
			log.Info("Found in store", "lk", lkMesh.Name)
			lkMesh := lkMesh
			renderContext := NewRenderContext(e, r, lkMesh)
			if ok, cm := store.LiveKitMeshes.IsConfigMapReadyForMesh(lkMesh); ok {
				log.Info("ConfigMap is present")
				renderContext.liveKitConfig = cm
				gw := renderContext.liveKitMesh.Spec.Components.Gateway
				if addr := r.getLoadBalancerIP(r.logger, gw); addr == nil {
					log.Info("LoadBalancerIP is not present yet for", "Gateway",
						types.NamespacedName{
							Namespace: *gw.RelatedStunnerGatewayAnnotations.Namespace,
							Name:      *gw.RelatedStunnerGatewayAnnotations.Name,
						}.String())
				} else {
					log.Info("LoadBalancerIP is present for", "Gateway",
						types.NamespacedName{
							Namespace: *gw.RelatedStunnerGatewayAnnotations.Namespace,
							Name:      *gw.RelatedStunnerGatewayAnnotations.Name,
						}.String(),
						"addr", addr)
					renderContext.stunnerPublicAddress = addr
					renderLiveKitDeployment(r.logger, renderContext)
					renderLiveKitService(r.logger, renderContext)
					updateLiveKitMeshStatus(r.logger, renderContext)
				}
			}
			//renderContext.liveKit
			r.operatorCh <- renderContext.update
		}
	}
}

func updateLiveKitMeshStatus(logger logr.Logger, context *RenderContext) {
	log := logger.WithName("updateLiveKitMeshStatus")

	log.Info("trying to update LiveKitMesh status")

	lkMesh := context.liveKitMesh

	if lkMesh.Status.OverallStatus == nil {
		log.Info("Unprocessed LiveKitMesh, initializing its status")
		overallStatus := lkstnv1a1.InstallStatus(opdefault.StatusReconciling)
		lkMesh.Status.OverallStatus = &overallStatus

		if lkMesh.Spec.Components.LiveKit != nil {
			lkMesh.Status.ComponentStatus[opdefault.ComponentLiveKit] = opdefault.StatusReconciling
		} else {
			// SHOULD BE NONE in all other cases but here we need to raise panic
			panic("LiveKit Component has not been found")
		}
	}
}

func renderLiveKitService(logger logr.Logger, context *RenderContext) {
	log := logger.WithName("renderLiveKitService")

	log.Info("trying to render LiveKit-Server Service")

	service, err := livekitServiceSkeleton(context.liveKitMesh)
	if err != nil {
		log.Error(err, "Error while creating LiveKit Service")
	}

	context.update.UpsertQueue.Services.Upsert(service)

	log.Info("Upserted LiveKit-Server Service into UpsertQueue")
}

func renderLiveKitDeployment(logger logr.Logger, context *RenderContext) {
	log := logger.WithName("renderLiveKitDeployment")

	log.Info("trying to render LiveKit-Server Deployment")

	deployment, err := livekitDeploymentSkeleton(context.liveKitMesh, context.liveKitConfig)
	if err != nil {
		log.Error(err, "Error while creating LiveKit Deployment")
	}
	context.update.UpsertQueue.Deployments.Upsert(deployment)

	log.Info("Upserted LiveKit-Server Deployment into UpsertQueue")
}
