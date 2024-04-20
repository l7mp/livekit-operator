package renderer

import (
	envygwapiv1 "github.com/envoyproxy/gateway/api/v1alpha1"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/store"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func createEnvoyGatewayClass(lkMesh *lkstnv1a1.LiveKitMesh) *gwapiv1.GatewayClass {
	name := getEnvoyGatewayClassName(lkMesh.Name)

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
			ControllerName: envygwapiv1.GatewayControllerName,
		},
	}
}

func createEnvoyGateway(lkMesh *lkstnv1a1.LiveKitMesh) *gwapiv1.Gateway {

	name := getEnvoyGatewayName(lkMesh.Name)
	ns := gwapiv1.Namespace(lkMesh.Namespace)
	labels := map[string]string{
		opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
		opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
		opdefault.RelatedComponent:      opdefault.ComponentStunner,
	}

	if current := store.Gateways.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      name,
	}); current != nil {
		labels = mergeMaps(labels, current.Labels)
	}
	hostName := gwapiv1.Hostname(getHostNameWithSubDomain("*", *lkMesh.Spec.Components.ApplicationExpose.HostName))
	mode := gwapiv1.TLSModeTerminate
	kind := gwapiv1.Kind("Secret")

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
				"cert-manager.io/issuer": "cloudflare-issuer",
			},
		},
		Spec: gwapiv1.GatewaySpec{
			GatewayClassName: gwapiv1.ObjectName(getEnvoyGatewayClassName(lkMesh.Name)),
			Listeners: []gwapiv1.Listener{{
				Name:     gwapiv1.SectionName(getEnvoyGatewayListenerName(lkMesh.Name)),
				Protocol: gwapiv1.HTTPSProtocolType,
				Hostname: &hostName,
				Port:     gwapiv1.PortNumber(443),
				TLS: &gwapiv1.GatewayTLSConfig{
					Mode: &mode,
					CertificateRefs: []gwapiv1.SecretObjectReference{{
						Kind:      &kind,
						Name:      gwapiv1.ObjectName(getEnvoyGatewayListenerSecretName(lkMesh.Name)),
						Namespace: &ns,
					},
					},
				},
			},
			},
		},
		Status: gwapiv1.GatewayStatus{},
	}
}

func createEnvoyHTTPRoute(lkMesh *lkstnv1a1.LiveKitMesh) *gwapiv1.HTTPRoute {

	name := getEnvoyHTTPRouteName(lkMesh.Name)
	ns := gwapiv1.Namespace(lkMesh.Namespace)

	labels := map[string]string{
		opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
		opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
		opdefault.RelatedComponent:      opdefault.ComponentStunner,
	}

	if current := store.HTTPRoutes.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      name,
	}); current != nil {
		labels = mergeMaps(labels, current.Labels)
	}

	parentRefObjectName := gwapiv1.ObjectName(getEnvoyGatewayName(lkMesh.Name))
	specifiedHostName := getHostNameWithSubDomain("server", *lkMesh.Spec.Components.ApplicationExpose.HostName)
	hostnames := []gwapiv1.Hostname{gwapiv1.Hostname(specifiedHostName)}

	pathMatchType := gwapiv1.PathMatchPathPrefix
	pathMatchValue := "/"
	weight := int32(1)
	backendRefSvcName := ServiceNameFormat(*lkMesh.Spec.Components.LiveKit.Deployment.Name)
	kind := gwapiv1.Kind("Service")
	port := gwapiv1.PortNumber(443)

	return &gwapiv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: lkMesh.Namespace,
			Labels:    labels,
			Annotations: map[string]string{
				opdefault.RelatedLiveKitMeshKey: types.NamespacedName{
					Namespace: lkMesh.GetNamespace(),
					Name:      lkMesh.GetName(),
				}.String(),
				opdefault.HostnameAnnotationKey: specifiedHostName,
			},
		},
		Spec: gwapiv1.HTTPRouteSpec{
			CommonRouteSpec: gwapiv1.CommonRouteSpec{
				ParentRefs: []gwapiv1.ParentReference{{
					Name:      parentRefObjectName,
					Namespace: &ns,
				},
				},
			},
			Hostnames: hostnames,
			Rules: []gwapiv1.HTTPRouteRule{{
				Matches: []gwapiv1.HTTPRouteMatch{{
					Path: &gwapiv1.HTTPPathMatch{
						Type:  &pathMatchType,
						Value: &pathMatchValue,
					},
				}},
				BackendRefs: []gwapiv1.HTTPBackendRef{{
					BackendRef: gwapiv1.BackendRef{
						BackendObjectReference: gwapiv1.BackendObjectReference{
							Name:      gwapiv1.ObjectName(backendRefSvcName),
							Namespace: &ns,
							Kind:      &kind,
							Port:      &port,
						},
						Weight: &weight,
					},
				}},
			}},
		},
	}
}
