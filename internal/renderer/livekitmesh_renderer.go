package renderer

import (
	"github.com/l7mp/livekit-operator/internal/event"
	"github.com/l7mp/livekit-operator/internal/store"
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

			r.renderStunnerComponentResources(renderContext)

			r.renderLiveKitComponentResources(renderContext)

			r.renderCertManagerComponentResources(renderContext)

			r.renderEnvoyGatewayResources(renderContext)

			r.renderExternalDNSResources(renderContext)

			r.renderLiveKitIngressResources(renderContext)

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

func (r *Renderer) renderEnvoyGatewayResources(renderContext *RenderContext) {
	r.renderEnvoyGatewayClass(renderContext)
	r.renderEnvoyGateway(renderContext)
	r.renderEnvoyHTTPRouteForLiveKitServer(renderContext)
	if renderContext.liveKitMesh.Spec.Components.Ingress != nil {
		r.renderEnvoyGatewayForLiveKitIngress(renderContext)
		r.renderEnvoyTCPRoutesForLiveKitIngress(renderContext)
	}
}

func (r *Renderer) renderExternalDNSResources(renderContext *RenderContext) {
	if renderContext.liveKitMesh.Spec.Components.ApplicationExpose.ExternalDNS != nil {
		if renderContext.liveKitMesh.Spec.Components.ApplicationExpose.ExternalDNS.CloudFlare != nil {
			r.renderExternalDNSCloudFlareDNSClusterRole(renderContext)
			r.renderExternalDNSCloudFlareDNSClusterRoleBinding(renderContext)
			r.renderExternalDNSCloudFlareServiceAccount(renderContext)
			r.renderExternalDNSCloudFlareDeployment(renderContext)
		}
	}
}

func (r *Renderer) renderCertManagerComponentResources(renderContext *RenderContext) {
	r.renderCertManagerIssuerAndSecret(renderContext)
}

func (r *Renderer) renderLiveKitIngressResources(renderContext *RenderContext) {
	if renderContext.liveKitMesh.Spec.Components.Ingress != nil {
		r.renderLiveKitIngressConfigMap(renderContext)
		r.renderLiveKitIngressDeployment(renderContext)
		r.renderLiveKitIngressService(renderContext)
	}
}
