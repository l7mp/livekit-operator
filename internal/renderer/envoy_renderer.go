package renderer

import (
	"github.com/l7mp/livekit-operator/internal/store"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *Renderer) renderEnvoyGatewayClass(renderContext *RenderContext) {
	log := r.logger.WithName("renderEnvoyGatewayClass")

	log.V(2).Info("trying to render Envoy GatewayClass")
	lkMesh := renderContext.liveKitMesh

	gatewayClass := createEnvoyGatewayClass(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, gatewayClass, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(gatewayClass))
		return
	}

	renderContext.update.UpsertQueue.GatewayClasses.Upsert(gatewayClass)
	log.V(2).Info("Upserted Envoy GatewayClass into UpsertQueue", "gwclass", store.GetObjectKey(gatewayClass))

}

func (r *Renderer) renderEnvoyGateway(renderContext *RenderContext) {
	log := r.logger.WithName("renderEnvoyGateway")

	log.V(2).Info("trying to render Envoy Gateway")
	lkMesh := renderContext.liveKitMesh

	gateway := createEnvoyGateway(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, gateway, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(gateway))
		return
	}

	renderContext.update.UpsertQueue.Gateways.Upsert(gateway)
	log.V(2).Info("Upserted Envoy Gateway into UpsertQueue", "gw", store.GetObjectKey(gateway))
}

func (r *Renderer) renderEnvoyHTTPRouteForLiveKitServer(renderContext *RenderContext) {
	log := r.logger.WithName("renderEnvoyHTTPRouteForLiveKitServer")

	log.V(2).Info("trying to render Envoy HTTPRoute")
	lkMesh := renderContext.liveKitMesh

	httpRoute := createEnvoyHTTPRoute(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, httpRoute, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(httpRoute))
		return
	}

	renderContext.update.UpsertQueue.HTTPRoutes.Upsert(httpRoute)
	log.V(2).Info("Upserted Envoy HTTPRoute into UpsertQueue", "httproute", store.GetObjectKey(httpRoute))
}

func (r *Renderer) renderEnvoyTCPRoutesForLiveKitIngress(renderContext *RenderContext) {
	log := r.logger.WithName("renderEnvoyTCPRoutesForLiveKitIngress")

	log.V(2).Info("trying to render Envoy LiveKitIngress TCPRoute")
	lkMesh := renderContext.liveKitMesh

	if renderContext.liveKitMesh.Spec.Components.Ingress.Rtmp != nil {
		rtmpTcpRoute := createEnvoyLiveKitIngressTCPRouteRtmp(lkMesh)
		if err := controllerutil.SetOwnerReference(lkMesh, rtmpTcpRoute, r.scheme); err != nil {
			log.Error(err, "cannot set owner reference", "owner",
				store.GetObjectKey(lkMesh), "reference",
				store.GetObjectKey(rtmpTcpRoute))
			return
		}

		renderContext.update.UpsertQueue.TCPRoutes.Upsert(rtmpTcpRoute)
		log.V(2).Info("Upserted Envoy TCPRoute into UpsertQueue", "tcproute", store.GetObjectKey(rtmpTcpRoute))
	}

	if renderContext.liveKitMesh.Spec.Components.Ingress.Whip != nil {
		// FIXME change this tcp route to HTTP route
		whipTcpRoute := createEnvoyLiveKitIngressTCPRouteWhip(lkMesh)
		if err := controllerutil.SetOwnerReference(lkMesh, whipTcpRoute, r.scheme); err != nil {
			log.Error(err, "cannot set owner reference", "owner",
				store.GetObjectKey(lkMesh), "reference",
				store.GetObjectKey(whipTcpRoute))
			return
		}

		renderContext.update.UpsertQueue.TCPRoutes.Upsert(whipTcpRoute)
		log.V(2).Info("Upserted Envoy TCPRoute into UpsertQueue", "tcproute", store.GetObjectKey(whipTcpRoute))
	}
}

//func (r *Renderer) renderEnvoyGatewayForLiveKitIngress(renderContext *RenderContext) {
//	log := r.logger.WithName("renderEnvoyGatewayForLiveKitIngress")
//
//	log.V(2).Info("trying to render Envoy LiveKitIngress Gateway")
//	lkMesh := renderContext.liveKitMesh
//
//	gateway := createEnvoyLiveKitIngressGateway(lkMesh)
//	if err := controllerutil.SetOwnerReference(lkMesh, gateway, r.scheme); err != nil {
//		log.Error(err, "cannot set owner reference", "owner",
//			store.GetObjectKey(lkMesh), "reference",
//			store.GetObjectKey(gateway))
//		return
//	}
//
//	renderContext.update.UpsertQueue.Gateways.Upsert(gateway)
//	log.V(2).Info("Upserted Envoy LiveKit Ingress Gateway into UpsertQueue", "gw", store.GetObjectKey(gateway))
//}
