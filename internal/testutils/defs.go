package testutils

import (
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
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
	TestConfigMapNamespacedName       = lkstnv1a1.NamespacedName{
		Namespace: &TestNsName,
		Name:      &TestConfigMapName,
	}
	TestGatewayNamespacedName = "testgateway"
)

// TestNs is a Namespace for testing purposes
var TestNs = corev1.Namespace{
	ObjectMeta: metav1.ObjectMeta{
		Name:   TestNsName,
		Labels: map[string]string{TestLabelName: TestLabelValue},
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
					ConfigMap: &TestConfigMapNamespacedName,
				},
			},
			Ingress: nil,
			Egress:  nil,
			Gateway: &lkstnv1a1.Gateway{
				RelatedStunnerGatewayAnnotations: &lkstnv1a1.NamespacedName{
					Namespace: &TestNsName,
					Name:      &TestGatewayNamespacedName,
				},
			},
			CertManager: nil,
			Monitoring:  nil,
		},
	},
	Status: lkstnv1a1.LiveKitMeshStatus{},
}
