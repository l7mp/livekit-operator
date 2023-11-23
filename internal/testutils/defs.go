package testutils

import (
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	TestNsName          = "testnamespace"
	TestConfigMapName   = "testconfigmap"
	TestLabelName       = "testlabel"
	TestLabelValue      = "testvalue"
	TestDeploymentName  = "testdeployment"
	TestReplicaNumber   = int32(1)
	TestImage           = "testrepo/testimage"
	TestImagePullPolicy = corev1.PullAlways
	TestCPURequest      = resource.MustParse("250m")
	TestMemoryLimit     = resource.MustParse("10M")
	TestResourceRequest = corev1.ResourceList(map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU: TestCPURequest,
	})
	TestResourceLimit = corev1.ResourceList(map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceMemory: TestMemoryLimit,
	})
	TestResourceRequirements = corev1.ResourceRequirements{
		Limits:   TestResourceLimit,
		Requests: TestResourceRequest,
	}
	TestTerminationGracePeriodSeconds = int64(3600)
	TestGatewayNamespacedName         = "testgateway"
	TestConfigCredential              = "testcredential"
	TestAuthUri                       = "http://localhost:8080?service=turn"
	TestAccessToken                   = "testtoken"
	TestLogLevel                      = "info"
	TestPort                          = 1234
	TestRedisAddress                  = "dummy_address"
	TestPortRangeStart                = 11111
	TestPortRangeEnd                  = 22222
	TestTCPPort                       = 1235
	TestUseExternalIP                 = false
)

// TestNs is a Namespace for testing purposes
var TestNs = corev1.Namespace{
	ObjectMeta: metav1.ObjectMeta{
		Name:   TestNsName,
		Labels: map[string]string{TestLabelName: TestLabelValue},
	},
}

// TestConfigMap is a ConfigMap for testing purposes which holds a dummy configuration for the LiveKit Deployment
var TestConfigMap = corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      TestConfigMapName,
		Namespace: TestNsName,
		Labels: map[string]string{
			opdefault.DefaultLabelKeyForConfigMap: opdefault.DefaultLabelValueForConfigMap,
		},
	},
	Data: map[string]string{
		"dummydatakey": "dummydatavalue",
	},
}

// TestLkMesh is a LiveKitMesh for testing purposes
var TestLkMesh = lkstnv1a1.LiveKitMesh{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "testlivekitmesh",
		Namespace: TestNsName,
	},
	Spec: lkstnv1a1.LiveKitMeshSpec{
		Components: &lkstnv1a1.Component{
			LiveKit: &lkstnv1a1.LiveKit{
				Type: "",
				Deployment: &lkstnv1a1.Deployment{
					Name:     &TestDeploymentName,
					Replicas: &TestReplicaNumber,
					Container: &lkstnv1a1.Container{
						Image:                         TestImage,
						ImagePullPolicy:               &TestImagePullPolicy,
						Command:                       []string{"testcommand-1"},
						Args:                          []string{"testarg-1", "testarg-2"},
						Env:                           nil,
						Resources:                     &TestResourceRequirements,
						TerminationGracePeriodSeconds: &TestTerminationGracePeriodSeconds,
						HostNetwork:                   false,
						Affinity:                      nil,
						SecurityContext:               nil,
						Tolerations:                   nil,
						HealthCheckPort:               nil,
					},
					Config: &lkstnv1a1.LiveKitConfig{
						Keys: &lkstnv1a1.Keys{
							AccessToken: &TestAccessToken,
						},
						LogLevel: &TestLogLevel,
						Port:     &TestPort,
						Redis: &lkstnv1a1.Redis{
							Address: &TestRedisAddress,
						},
						Rtc: &lkstnv1a1.Rtc{
							PortRangeEnd:   &TestPortRangeEnd,
							PortRangeStart: &TestPortRangeStart,
							TcpPort:        &TestTCPPort,
							StunServers:    nil,
							TurnServers: []lkstnv1a1.TurnServer{{
								AuthURI: &TestAuthUri,
								//Credential: TestConfigCredential,
								//Host:
							}},
							UseExternalIp: &TestUseExternalIP,
						},
					},
				},
			},
			Ingress: nil,
			Egress:  nil,
			//Gateway: &lkstnv1a1.Gateway{
			//	RelatedStunnerGatewayAnnotations: &lkstnv1a1.NamespacedName{
			//		Namespace: &TestNsName,
			//		Name:      &TestGatewayNamespacedName,
			//	},
			//},
			Gateway:     nil,
			CertManager: nil,
			Monitoring:  nil,
		},
	},
	Status: lkstnv1a1.LiveKitMeshStatus{},
}

/*// TestService is representing the service created by stunner
var TestService = corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Namespace: "testnamespace",
		Name:      "testservice-ok",
		Annotations: map[string]string{
			opdefault.RelatedGatewayKey: "testnamespace/gateway-1",
		},
	},
	Spec: corev1.ServiceSpec{
		Type:     corev1.ServiceTypeLoadBalancer,
		Selector: map[string]string{"app": "dummy"},
		Ports: []corev1.ServicePort{
			{
				Name:     "udp-ok",
				Protocol: corev1.ProtocolUDP,
				Port:     1,
			},
		},
	},
	Status: corev1.ServiceStatus{
		LoadBalancer: corev1.LoadBalancerStatus{
			Ingress: []corev1.LoadBalancerIngress{{
				IP: "1.2.3.4",
				Ports: []corev1.PortStatus{{
					Port:     1,
					Protocol: corev1.ProtocolUDP,
				}},
			}, {
				IP: "5.6.7.8",
				Ports: []corev1.PortStatus{{
					Port:     2,
					Protocol: corev1.ProtocolTCP,
				}},
			}},
		}},
}
*/