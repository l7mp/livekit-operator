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
					//continue
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

			r.renderStunnerComponentResources(renderContext)
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

func (r *Renderer) renderStunnerComponentResources(renderContext *RenderContext) {
	r.renderStunnerGatewayClass(renderContext)
	r.renderStunnerGatewayConfig(renderContext)
	r.renderStunnerGateway(renderContext)
	r.renderStunnerUdpRoute(renderContext)
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
