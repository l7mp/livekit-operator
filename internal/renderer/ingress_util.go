package renderer

import (
	"fmt"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/store"
	"github.com/l7mp/livekit-operator/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"
)

func createLiveKitIngressConfigMap(lkMesh *lkstnv1a1.LiveKitMesh) (*corev1.ConfigMap, error) {

	ingress := lkMesh.Spec.Components.Ingress
	redis := lkstnv1a1.Redis{}
	if lkMesh.Spec.Components.LiveKit.Deployment.Config.Redis != nil {
		redis = *lkMesh.Spec.Components.LiveKit.Deployment.Config.Redis
	} else {
		redis.Address = ptr.To(fmt.Sprintf("%s.%s:%d", getRedisName(lkMesh.Name), lkMesh.Namespace, 6379))
	}

	//TODO fix service name
	wsUrl := fmt.Sprintf("ws://%s.%s:%d", getLiveKitServiceName(*lkMesh.Spec.Components.LiveKit.Deployment.Name), lkMesh.Namespace, 443)

	//get the first key-value
	apiKey, apiSecret := "", ""
	for k, v := range *lkMesh.Spec.Components.LiveKit.Deployment.Config.Keys {
		apiKey = k
		apiSecret = v
		break
	}

	ingressConfig := config.ConvertIngressConfig(*ingress.Config)
	ingressConfig.APIKey = &apiKey
	ingressConfig.APISecret = &apiSecret
	ingressConfig.Redis = &redis
	ingressConfig.WSURL = &wsUrl

	yamlData, err := yaml.Marshal(ingressConfig)
	if err != nil {
		return nil, err
	}

	yamlMap := make(map[string]string)
	yamlMap[config.DefaultLiveKitConfigFileName] = string(yamlData)

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getIngressName(lkMesh.Name),
			Namespace: lkMesh.GetNamespace(),
			Labels: map[string]string{
				config.OwnedByLabelKey:             config.OwnedByLabelValue,
				config.RelatedLiveKitMeshKey:       lkMesh.GetName(),
				config.DefaultLabelKeyForConfigMap: config.DefaultLabelValueForConfigMap,
				config.RelatedComponent:            config.ComponentIngress,
			},
		},
		Data: yamlMap,
	}

	return cm, nil
}

func createLiveKitIngressService(lkMesh *lkstnv1a1.LiveKitMesh) *corev1.Service {

	labels := map[string]string{
		config.OwnedByLabelKey:       config.OwnedByLabelValue,
		config.RelatedLiveKitMeshKey: lkMesh.GetName(),
		config.RelatedComponent:      config.ComponentIngress,
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
				config.RelatedLiveKitMeshKey: types.NamespacedName{
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
				config.OwnedByLabelKey:       config.OwnedByLabelValue,
				config.RelatedLiveKitMeshKey: lkMesh.GetName(),
				config.RelatedComponent:      config.ComponentIngress,
			},
			Annotations: map[string]string{
				config.RelatedLiveKitMeshKey: types.NamespacedName{
					Namespace: lkMesh.GetNamespace(),
					Name:      lkMesh.GetName(),
				}.String(),
				config.RelatedConfigMapKey: getIngressName(lkMesh.Name),
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
						config.DefaultConfigMapResourceVersionKey: "",
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
