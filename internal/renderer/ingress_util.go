package renderer

import (
	"fmt"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/store"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"
)

type IngressConfig struct {
	APIKey         string  `yaml:"api_key" json:"api_key"`
	APISecret      string  `yaml:"api_secret" json:"api_secret"`
	CPUCost        CPUCost `yaml:"cpu_cost" json:"cpu_cost"`
	HealthPort     int     `yaml:"health_port" json:"health_port"`
	HTTPRelayPort  int     `yaml:"http_relay_port" json:"http_relay_port"`
	Logging        Logging `yaml:"logging" json:"logging"`
	PrometheusPort int     `yaml:"prometheus_port" json:"prometheus_port"`
	Redis          Redis   `yaml:"redis" json:"redis"`
	RTMPPort       int     `yaml:"rtmp_port" json:"rtmp_port"`
	WSURL          string  `yaml:"ws_url" json:"ws_url"`
}

type CPUCost struct {
	RTMPCPUCost int `yaml:"rtmp_cpu_cost" json:"rtmp_cpu_cost"`
}

type Logging struct {
	Level string `yaml:"level" json:"level"`
}

type Redis struct {
	Address string `yaml:"address" json:"address"`
}

func createLiveKitIngressConfigMap(lkMesh *lkstnv1a1.LiveKitMesh) (*corev1.ConfigMap, error) {
	redisAddress := ""

	if lkMesh.Spec.Components.LiveKit.Deployment.Config.Redis != nil {
		redisAddress = *lkMesh.Spec.Components.LiveKit.Deployment.Config.Redis.Address
	} else {
		redisAddress = fmt.Sprintf("%s.%s.%d", getRedisName(lkMesh.Name), lkMesh.Namespace, 6379)
	}

	//TODO fix service name
	wsUrl := fmt.Sprintf("ws://%s.%s:%d", getLiveKitServiceName(*lkMesh.Spec.Components.LiveKit.Deployment.Name), lkMesh.Namespace, 443)

	config := &IngressConfig{
		APIKey:    "access_token",
		APISecret: *lkMesh.Spec.Components.LiveKit.Deployment.Config.Keys.AccessToken,
		CPUCost: CPUCost{
			RTMPCPUCost: 2,
		},
		HealthPort:    7888,
		HTTPRelayPort: 9090,
		Logging: Logging{
			Level: "debug",
		},
		PrometheusPort: 7889,
		Redis: Redis{
			Address: redisAddress,
		},
		RTMPPort: 1935,
		WSURL:    wsUrl,
	}

	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	yamlMap := make(map[string]string)
	yamlMap[opdefault.DefaultLiveKitConfigFileName] = string(yamlData)

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getIngressName(lkMesh.Name),
			Namespace: lkMesh.GetNamespace(),
			Labels: map[string]string{
				opdefault.OwnedByLabelKey:             opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey:       lkMesh.GetName(),
				opdefault.DefaultLabelKeyForConfigMap: opdefault.DefaultLabelValueForConfigMap,
				opdefault.RelatedComponent:            opdefault.ComponentLiveKit,
			},
		},
		Data: yamlMap,
	}

	return cm, nil
}

func createLiveKitIngressService(lkMesh *lkstnv1a1.LiveKitMesh) *corev1.Service {

	labels := map[string]string{
		opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
		opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
		opdefault.RelatedComponent:      opdefault.ComponentIngress,
	}

	if current := store.Services.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      getIngressName(lkMesh.Name),
	}); current != nil {
		labels = mergeMaps(labels, current.Labels)
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: lkMesh.Namespace,
			Name:      getIngressName(lkMesh.Name),
			Labels:    labels,
			Annotations: map[string]string{
				opdefault.RelatedLiveKitMeshKey: types.NamespacedName{
					Namespace: lkMesh.GetNamespace(),
					Name:      lkMesh.GetName(),
				}.String(),
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app.kubernetes.io/name":     getIngressName(lkMesh.Name),
				"app.kubernetes.io/instance": getIngressName(lkMesh.Name),
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "ws",
					TargetPort: intstr.FromInt32(7888),
					Port:       7888,
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       "rtmp",
					TargetPort: intstr.FromInt32(1935),
					Port:       1935,
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       "whip",
					TargetPort: intstr.FromInt32(8080),
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	return svc
}

func createLiveKitIngressDeployment(lkMesh *lkstnv1a1.LiveKitMesh) *appsv1.Deployment {

	var envList []corev1.EnvVar

	envList = append(envList, corev1.EnvVar{
		Name: "INGRESS_CONFIG_BODY",
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: getIngressName(lkMesh.Name),
				},
				Key: "config.yaml",
			},
		},
	})

	dp := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getIngressName(lkMesh.Name),
			Namespace: lkMesh.Namespace,
			Labels: map[string]string{
				opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
				opdefault.RelatedComponent:      opdefault.ComponentIngress,
			},
			Annotations: map[string]string{
				opdefault.RelatedLiveKitMeshKey: types.NamespacedName{
					Namespace: lkMesh.GetNamespace(),
					Name:      lkMesh.GetName(),
				}.String(),
				opdefault.RelatedConfigMapKey: getIngressName(lkMesh.Name),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To(int32(1)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":     getIngressName(lkMesh.Name),
					"app.kubernetes.io/instance": getIngressName(lkMesh.Name),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":     getIngressName(lkMesh.Name),
						"app.kubernetes.io/instance": getIngressName(lkMesh.Name),
					},
					Annotations: map[string]string{
						opdefault.DefaultConfigMapResourceVersionKey: "",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:            "default",
					TerminationGracePeriodSeconds: ptr.To(int64(3600)),
					Containers: []corev1.Container{{
						Name:            getIngressName(lkMesh.Name),
						Image:           "livekit/ingress:v1.4",
						ImagePullPolicy: "IfNotPresent",
						Ports: []corev1.ContainerPort{
							{
								Name:          "http-relay",
								ContainerPort: 9090,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								Name:          "http",
								ContainerPort: 7888,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								Name:          "rtmp-port",
								ContainerPort: 1935,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								Name:          "metrics",
								ContainerPort: 7889,
								Protocol:      corev1.ProtocolTCP,
							},
						},
						Env: envList,
						//TODO Resources: ,
					},
					},
				},
			},
		},
	}

	return dp
}
