package renderer

import (
	"fmt"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/store"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	stnrgwv1 "github.com/l7mp/stunner-gateway-operator/api/v1"
	stnrgwpkg "github.com/l7mp/stunner-gateway-operator/pkg/config"
	stnrpgkv1 "github.com/l7mp/stunner/pkg/apis/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func (r *Renderer) renderStunnerGatewayClass(renderContext *RenderContext) {
	log := r.logger.WithName("renderStunnerGatewayClass")

	log.V(2).Info("trying to render STUNner GatewayClass")
	lkMesh := renderContext.liveKitMesh

	gatewayClass := createStunnerGatewayClass(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, gatewayClass, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(gatewayClass))
		return
	}

	renderContext.update.UpsertQueue.GatewayClasses.Upsert(gatewayClass)
	log.V(2).Info("Upserted STUNner GatewayClass into UpsertQueue", "gwclass", store.GetObjectKey(gatewayClass))
}

func (r *Renderer) renderStunnerGatewayConfig(renderContext *RenderContext) {
	log := r.logger.WithName("renderStunnerGatewayConfig")

	log.V(2).Info("trying to render STUNner GatewayConfig")
	lkMesh := renderContext.liveKitMesh

	gatewayConfig := createStunnerGatewayConfig(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, gatewayConfig, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(gatewayConfig))
		return
	}

	renderContext.update.UpsertQueue.GatewayConfigs.Upsert(gatewayConfig)
	log.V(2).Info("Upserted STUNner GatewayConfig into UpsertQueue", "gwconfig", store.GetObjectKey(gatewayConfig))
}

func (r *Renderer) renderStunnerGateway(renderContext *RenderContext) {
	log := r.logger.WithName("renderStunnerGateway")

	log.V(2).Info("trying to render STUNner Gateway")
	lkMesh := renderContext.liveKitMesh

	gateway := createStunnerGateway(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, gateway, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(gateway))
		return
	}

	renderContext.update.UpsertQueue.Gateways.Upsert(gateway)
	log.V(2).Info("Upserted STUNner Gateway into UpsertQueue", "gw", store.GetObjectKey(gateway))
}

func (r *Renderer) renderStunnerUdpRoute(renderContext *RenderContext) {
	log := r.logger.WithName("renderStunnerUdpRoute")

	log.V(2).Info("trying to render STUNner UDPRoute")
	lkMesh := renderContext.liveKitMesh

	udpRoute := createStunnerUDPRoute(lkMesh)
	if err := controllerutil.SetOwnerReference(lkMesh, udpRoute, r.scheme); err != nil {
		log.Error(err, "cannot set owner reference", "owner",
			store.GetObjectKey(lkMesh), "reference",
			store.GetObjectKey(udpRoute))
		return
	}

	renderContext.update.UpsertQueue.UDPRoutes.Upsert(udpRoute)
	log.V(2).Info("Upserted STUNner UDPRoute into UpsertQueue", "udproute", store.GetObjectKey(udpRoute))
}

func createStunnerGateway(lkMesh *lkstnv1a1.LiveKitMesh) *gwapiv1.Gateway {

	name := GetStunnerGatewayName(lkMesh.Name)
	labels := map[string]string{
		opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
		opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
		opdefault.RelatedComponent:      opdefault.ComponentStunner,
	}

	if current := store.UDPRoutes.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      name,
	}); current != nil {
		labels = mergeMaps(labels, current.Labels)
	}

	return &gwapiv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: lkMesh.GetNamespace(),
			Labels:    labels,
			Annotations: map[string]string{
				opdefault.RelatedLiveKitMeshKey: types.NamespacedName{
					Namespace: lkMesh.GetNamespace(),
					Name:      lkMesh.GetName(),
				}.String(),
			},
		},
		Spec: gwapiv1.GatewaySpec{
			GatewayClassName: gwapiv1.ObjectName(GetStunnerGatewayClassName(lkMesh.Name)),
			Listeners:        lkMesh.Spec.Components.Stunner.GatewayListeners,
		},
		Status: gwapiv1.GatewayStatus{},
	}
}

func createStunnerGatewayConfig(lkMesh *lkstnv1a1.LiveKitMesh) *stnrgwv1.GatewayConfig {

	name := GetStunnerGatewayConfigName(lkMesh.Name)

	labels := map[string]string{
		opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
		opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
		opdefault.RelatedComponent:      opdefault.ComponentStunner,
	}

	if current := store.GatewayClasses.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      name,
	}); current != nil {
		labels = mergeMaps(labels, current.Labels)
	}

	return &stnrgwv1.GatewayConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: lkMesh.Namespace,
			Labels:    labels,
			Annotations: map[string]string{
				opdefault.RelatedLiveKitMeshKey: types.NamespacedName{
					Namespace: lkMesh.GetNamespace(),
					Name:      lkMesh.GetName(),
				}.String(),
			},
		},
		Spec: *lkMesh.Spec.Components.Stunner.GatewayConfig,
	}
}

func createStunnerGatewayClass(lkMesh *lkstnv1a1.LiveKitMesh) *gwapiv1.GatewayClass {

	name := GetStunnerGatewayClassName(lkMesh.Name)
	ns := gwapiv1.Namespace(lkMesh.Namespace)
	desc := fmt.Sprintf("GatewayClass-LiveKitMesh-%s", lkMesh.GetName())

	labels := map[string]string{
		opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
		opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
		opdefault.RelatedComponent:      opdefault.ComponentStunner,
	}

	if current := store.GatewayClasses.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      name,
	}); current != nil {
		labels = mergeMaps(labels, current.Labels)
	}

	return &gwapiv1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: lkMesh.Namespace,
			Labels:    labels,
			Annotations: map[string]string{
				opdefault.RelatedLiveKitMeshKey: types.NamespacedName{
					Namespace: lkMesh.GetNamespace(),
					Name:      lkMesh.GetName(),
				}.String(),
			},
		},
		Spec: gwapiv1.GatewayClassSpec{
			ControllerName: stnrgwpkg.DefaultControllerName,
			ParametersRef: &gwapiv1.ParametersReference{
				Group:     stnrpgkv1.DefaultRealm,
				Kind:      "GatewayConfig",
				Name:      GetStunnerGatewayConfigName(lkMesh.Name),
				Namespace: &ns,
			},
			Description: &desc,
		},
	}
}

func createStunnerUDPRoute(lkMesh *lkstnv1a1.LiveKitMesh) *stnrgwv1.UDPRoute {

	backendRefSvcName := ServiceNameFormat(*lkMesh.Spec.Components.LiveKit.Deployment.Name)
	name := GetStunnerUDPRouteName(lkMesh.Name)
	ns := gwapiv1.Namespace(lkMesh.Namespace)

	labels := map[string]string{
		opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
		opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
		opdefault.RelatedComponent:      opdefault.ComponentStunner,
	}

	if current := store.UDPRoutes.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      name,
	}); current != nil {
		labels = mergeMaps(labels, current.Labels)
	}

	return &stnrgwv1.UDPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: lkMesh.Namespace,
			Labels:    labels,
			Annotations: map[string]string{
				opdefault.RelatedLiveKitMeshKey: types.NamespacedName{
					Namespace: lkMesh.GetNamespace(),
					Name:      lkMesh.GetName(),
				}.String(),
			},
		},
		Spec: stnrgwv1.UDPRouteSpec{
			CommonRouteSpec: gwapiv1.CommonRouteSpec{
				ParentRefs: []gwapiv1.ParentReference{{
					Name:      gwapiv1.ObjectName(GetStunnerGatewayName(lkMesh.Name)),
					Namespace: &ns,
				}},
			},
			Rules: []stnrgwv1.UDPRouteRule{{
				BackendRefs: []stnrgwv1.BackendRef{{
					BackendObjectReference: stnrgwv1.BackendObjectReference{
						Name:      gwapiv1.ObjectName(backendRefSvcName),
						Namespace: &ns,
					},
				}},
			}},
		},
	}
}
