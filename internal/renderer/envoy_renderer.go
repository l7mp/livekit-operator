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
	log.V(2).Info("Upserted STUNner HTTPRoute into UpsertQueue", "httproute", store.GetObjectKey(httpRoute))
}
