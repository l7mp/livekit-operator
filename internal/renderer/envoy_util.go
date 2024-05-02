package renderer

import (
	envygwapiv1 "github.com/envoyproxy/gateway/api/v1alpha1"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/store"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
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

	name := getEnvoyLiveKitServerGatewayName(lkMesh.Name)
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

	gw := &gwapiv1.Gateway{
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
				Name:     gwapiv1.SectionName(getEnvoyLiveKitServerGatewayListenerName(lkMesh.Name)),
				Protocol: gwapiv1.HTTPSProtocolType,
				Hostname: &hostName,
				Port:     gwapiv1.PortNumber(443),
				TLS: &gwapiv1.GatewayTLSConfig{
					Mode: &mode,
					CertificateRefs: []gwapiv1.SecretObjectReference{{
						Kind:      &kind,
						Name:      gwapiv1.ObjectName(getEnvoyLiveKitServerGatewayListenerSecretName(lkMesh.Name)),
						Namespace: &ns,
					},
					},
				},
			},
			},
		},
		Status: gwapiv1.GatewayStatus{},
	}

	return gw
}

func createEnvoyHTTPRoute(lkMesh *lkstnv1a1.LiveKitMesh) *gwapiv1.HTTPRoute {

	name := getEnvoyLiveKitServerHTTPRouteName(lkMesh.Name)
	ns := gwapiv1.Namespace(lkMesh.Namespace)

	labels := map[string]string{
		opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
		opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
		opdefault.RelatedComponent:      opdefault.ComponentApplicationExpose,
	}

	if current := store.HTTPRoutes.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      name,
	}); current != nil {
		labels = mergeMaps(labels, current.Labels)
	}

	parentRefObjectName := gwapiv1.ObjectName(getEnvoyLiveKitServerGatewayName(lkMesh.Name))
	specifiedHostName := getHostNameWithSubDomain("server", *lkMesh.Spec.Components.ApplicationExpose.HostName)
	hostnames := []gwapiv1.Hostname{gwapiv1.Hostname(specifiedHostName)}

	backendRefSvcName := getLiveKitServiceName(*lkMesh.Spec.Components.LiveKit.Deployment.Name)

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
						Type:  ptr.To(gwapiv1.PathMatchPathPrefix),
						Value: ptr.To("/"),
					},
				}},
				BackendRefs: []gwapiv1.HTTPBackendRef{{
					BackendRef: gwapiv1.BackendRef{
						BackendObjectReference: gwapiv1.BackendObjectReference{
							Name:      gwapiv1.ObjectName(backendRefSvcName),
							Namespace: &ns,
							Kind:      ptr.To(gwapiv1.Kind("Service")),
							Port:      ptr.To(gwapiv1.PortNumber(443)),
						},
						Weight: ptr.To(int32(1)),
					},
				}},
			}},
		},
	}
}

func createEnvoyLiveKitIngressGateway(lkMesh *lkstnv1a1.LiveKitMesh) *gwapiv1.Gateway {

	name := getEnvoyLiveKitIngressGatewayName(lkMesh.Name)
	labels := map[string]string{
		opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
		opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
		opdefault.RelatedComponent:      opdefault.ComponentApplicationExpose,
	}

	if current := store.Gateways.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      name,
	}); current != nil {
		labels = mergeMaps(labels, current.Labels)
	}
	hostName := gwapiv1.Hostname(getHostNameWithSubDomain("ingress", *lkMesh.Spec.Components.ApplicationExpose.HostName))
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
				"cert-manager.io/issuer":        "cloudflare-issuer",
				opdefault.HostnameAnnotationKey: string(hostName),
			},
		},
		Spec: gwapiv1.GatewaySpec{
			GatewayClassName: gwapiv1.ObjectName(getEnvoyGatewayClassName(lkMesh.Name)),
			Listeners: []gwapiv1.Listener{
				{
					Name:     gwapiv1.SectionName(getEnvoyLiveKitIngressGatewayListenerName(lkMesh.Name, "rtmp")),
					Protocol: gwapiv1.TLSProtocolType,
					Hostname: &hostName,
					Port:     gwapiv1.PortNumber(1935),
					TLS: &gwapiv1.GatewayTLSConfig{
						Mode: &mode,
						CertificateRefs: []gwapiv1.SecretObjectReference{{
							Kind:      &kind,
							Name:      gwapiv1.ObjectName(getEnvoyLiveKitServerGatewayListenerSecretName(lkMesh.Name)),
							Namespace: ptr.To(gwapiv1.Namespace(lkMesh.Namespace)),
						},
						},
					},
				},
				{
					Name:     gwapiv1.SectionName(getEnvoyLiveKitIngressGatewayListenerName(lkMesh.Name, "whip")),
					Protocol: gwapiv1.TLSProtocolType,
					Hostname: &hostName,
					Port:     gwapiv1.PortNumber(8080),
					TLS: &gwapiv1.GatewayTLSConfig{
						Mode: &mode,
						CertificateRefs: []gwapiv1.SecretObjectReference{{
							Kind:      &kind,
							Name:      gwapiv1.ObjectName(getEnvoyLiveKitServerGatewayListenerSecretName(lkMesh.Name)),
							Namespace: ptr.To(gwapiv1.Namespace(lkMesh.Namespace)),
						},
						},
					},
				},
			},
		},
	}
}

func createEnvoyLiveKitIngressTCPRouteRtmp(lkMesh *lkstnv1a1.LiveKitMesh) *gwapiv1a2.TCPRoute {

	name := getEnvoyLiveKitServerTCPRouteName(lkMesh.Name, "rtmp")
	ns := gwapiv1.Namespace(lkMesh.Namespace)

	labels := map[string]string{
		opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
		opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
		opdefault.RelatedComponent:      opdefault.ComponentApplicationExpose,
	}

	if current := store.HTTPRoutes.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      name,
	}); current != nil {
		labels = mergeMaps(labels, current.Labels)
	}

	parentRefObjectName := gwapiv1.ObjectName(getEnvoyLiveKitIngressGatewayName(lkMesh.Name))
	specifiedHostName := getHostNameWithSubDomain("ingress", *lkMesh.Spec.Components.ApplicationExpose.HostName)

	backendRefSvcName := getIngressName(lkMesh.Name)

	return &gwapiv1a2.TCPRoute{
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
		Spec: gwapiv1a2.TCPRouteSpec{
			CommonRouteSpec: gwapiv1.CommonRouteSpec{
				ParentRefs: []gwapiv1.ParentReference{{
					Name:        parentRefObjectName,
					Namespace:   &ns,
					SectionName: ptr.To(gwapiv1a2.SectionName(getEnvoyLiveKitIngressGatewayListenerName(lkMesh.Name, "rtmp"))),
				},
				},
			},
			Rules: []gwapiv1a2.TCPRouteRule{
				{
					BackendRefs: []gwapiv1a2.BackendRef{
						{
							BackendObjectReference: gwapiv1.BackendObjectReference{
								Name:      gwapiv1.ObjectName(backendRefSvcName),
								Namespace: &ns,
								Kind:      ptr.To(gwapiv1.Kind("Service")),
								Port:      ptr.To(gwapiv1.PortNumber(1935)),
							},
							Weight: ptr.To(int32(1)),
						},
					},
				},
			},
		},
	}
}

func createEnvoyLiveKitIngressTCPRouteWhip(lkMesh *lkstnv1a1.LiveKitMesh) *gwapiv1a2.TCPRoute {

	name := getEnvoyLiveKitServerTCPRouteName(lkMesh.Name, "whip")
	ns := gwapiv1.Namespace(lkMesh.Namespace)

	labels := map[string]string{
		opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
		opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
		opdefault.RelatedComponent:      opdefault.ComponentApplicationExpose,
	}

	if current := store.HTTPRoutes.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      name,
	}); current != nil {
		labels = mergeMaps(labels, current.Labels)
	}

	parentRefObjectName := gwapiv1.ObjectName(getEnvoyLiveKitIngressGatewayName(lkMesh.Name))
	specifiedHostName := getHostNameWithSubDomain("ingress", *lkMesh.Spec.Components.ApplicationExpose.HostName)

	backendRefSvcName := getIngressName(lkMesh.Name)

	return &gwapiv1a2.TCPRoute{
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
		Spec: gwapiv1a2.TCPRouteSpec{
			CommonRouteSpec: gwapiv1.CommonRouteSpec{
				ParentRefs: []gwapiv1.ParentReference{{
					Name:        parentRefObjectName,
					Namespace:   &ns,
					SectionName: ptr.To(gwapiv1a2.SectionName(getEnvoyLiveKitIngressGatewayListenerName(lkMesh.Name, "whip"))),
				},
				},
			},
			Rules: []gwapiv1a2.TCPRouteRule{
				{
					BackendRefs: []gwapiv1a2.BackendRef{
						{
							BackendObjectReference: gwapiv1.BackendObjectReference{
								Name:      gwapiv1.ObjectName(backendRefSvcName),
								Namespace: &ns,
								Kind:      ptr.To(gwapiv1.Kind("Service")),
								Port:      ptr.To(gwapiv1.PortNumber(8080)),
							},
							Weight: ptr.To(int32(1)),
						},
					},
				},
			},
		},
	}
}
