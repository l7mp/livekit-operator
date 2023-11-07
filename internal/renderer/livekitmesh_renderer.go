package renderer

import (
	"github.com/go-logr/logr"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/event"
	"github.com/l7mp/livekit-operator/internal/store"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
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
					r.renderLiveKitDeployment(renderContext)
					r.renderLiveKitService(renderContext)
					//updateLiveKitMeshStatus(r.logger, renderContext)
				}
			}
			//renderContext.liveKit
			log.Info("event to channel")
			r.operatorCh <- renderContext.update
		}
	}
}

// TODO
func updateLiveKitMeshStatus(logger logr.Logger, context *RenderContext) {
	log := logger.WithName("updateLiveKitMeshStatus")

	log.Info("trying to update LiveKitMesh status")

	lkMesh := context.liveKitMesh

	if lkMesh.Status.OverallStatus == nil {
		//TODO init components' map cuz it is nil in the beginning
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

func (r *Renderer) renderLiveKitService(context *RenderContext) {
	log := r.log.WithName("renderLiveKitService")

	log.Info("trying to render LiveKit-Server Service")

	lkMesh := context.liveKitMesh
	service := createLiveKitService(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, service, r.scheme); err != nil {
		r.log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(service))
	}

	context.update.UpsertQueue.Services.Upsert(service)

	log.Info("Upserted LiveKit-Server Service into UpsertQueue")
}

func (r *Renderer) renderLiveKitDeployment(context *RenderContext) {
	log := r.log.WithName("renderLiveKitDeployment")

	log.Info("trying to render LiveKit-Server Deployment")

	lkMesh := context.liveKitMesh
	deployment := createLiveKitDeployment(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, deployment, r.scheme); err != nil {
		r.log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(deployment))
	}

	context.update.UpsertQueue.Deployments.Upsert(deployment)

	log.Info("Upserted LiveKit-Server Deployment into UpsertQueue")
}
