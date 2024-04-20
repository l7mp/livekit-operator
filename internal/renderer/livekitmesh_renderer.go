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

			//gw := renderContext.liveKitMesh.Spec.Components.Gateway
			// if the PR goes in then this whole if-else block should be removed
			//if gw != nil {
			//	log.V(2).Info("Gateway is configured, looking for loadbalancerip for the LiveKit config")
			//	if addr := r.getLoadBalancerIP(r.logger, gw); addr == nil {
			//		log.Info("LoadBalancerIP is not present yet for", "Gateway",
			//			types.NamespacedName{
			//				Namespace: *gw.RelatedStunnerGatewayAnnotations.Namespace,
			//				Name:      *gw.RelatedStunnerGatewayAnnotations.Name,
			//			}.String())
			//		//continue
			//	} else {
			//		log.Info("LoadBalancerIP is present for", "Gateway",
			//			types.NamespacedName{
			//				Namespace: *gw.RelatedStunnerGatewayAnnotations.Namespace,
			//				Name:      *gw.RelatedStunnerGatewayAnnotations.Name,
			//			}.String(),
			//			"addr", addr)
			//		renderContext.turnServerPublicAddress = addr
			//	}
			//}

			r.renderStunnerComponentResources(renderContext)

			//var iceConfig *types.IceConfig
			//go func(log logr.Logger, mesh *v1alpha1.LiveKitMesh, iceConfig *types.IceConfig) {
			//	var err error
			//	iceConfig, err = getIceConfigurationFromStunnerAuth(*mesh)
			//	if err != nil {
			//		log.Error(err, "Failed to get ICE configuration from STUNner auth")
			//	}
			//}(log, lkMesh, iceConfig)

			// this will be supported way to render the configmap however
			r.renderLiveKitComponentResources(renderContext)

			r.renderCertManagerComponentResources(renderContext)

			r.renderEnvoyGatewayResources(renderContext)

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

func (r *Renderer) renderEnvoyGatewayResources(renderContext *RenderContext) {
	r.renderEnvoyGatewayClass(renderContext)
	r.renderEnvoyGateway(renderContext)
	r.renderEnvoyHTTPRouteForLiveKitServer(renderContext)
}

//func (r *Renderer) renderExternalDNSResources(renderContext *RenderContext) {
//	r.renderExternalDNSDeployment
//}

func (r *Renderer) renderCertManagerComponentResources(renderContext *RenderContext) {
	r.renderCertManagerIssuerAndSecret(renderContext)
}
