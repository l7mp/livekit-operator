package renderer

import (
	"fmt"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"
)

func createLiveKitEgressConfigMap(lkMesh *lkstnv1a1.LiveKitMesh) (*corev1.ConfigMap, error) {

	egress := lkMesh.Spec.Components.Egress
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

	egressConfig := config.ConvertEgressConfig(*egress.Config)
	egressConfig.APIKey = &apiKey
	egressConfig.APISecret = &apiSecret
	egressConfig.Redis = &redis
	egressConfig.WSURL = &wsUrl

	yamlData, err := yaml.Marshal(egressConfig)
	if err != nil {
		return nil, err
	}

	yamlMap := make(map[string]string)
	yamlMap[config.DefaultLiveKitConfigFileName] = string(yamlData)

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getEgressName(lkMesh.Name),
			Namespace: lkMesh.GetNamespace(),
			Labels: map[string]string{
				config.OwnedByLabelKey:             config.OwnedByLabelValue,
				config.RelatedLiveKitMeshKey:       lkMesh.GetName(),
				config.DefaultLabelKeyForConfigMap: config.DefaultLabelValueForConfigMap,
				config.RelatedComponent:            config.ComponentEgress,
			},
		},
		Data: yamlMap,
	}

	return cm, nil
}

func createLiveKitEgressDeployment(lkMesh *lkstnv1a1.LiveKitMesh) *appsv1.Deployment {

	var envList []corev1.EnvVar

	envList = append(envList, corev1.EnvVar{
		Name: "EGRESS_CONFIG_BODY",
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: getEgressName(lkMesh.Name),
				},
				Key: "config.yaml",
			},
		},
	})

	egressName := getEgressName(lkMesh.Name)

	dp := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      egressName,
			Namespace: lkMesh.Namespace,
			Labels: map[string]string{
				config.OwnedByLabelKey:       config.OwnedByLabelValue,
				config.RelatedLiveKitMeshKey: lkMesh.GetName(),
				config.RelatedComponent:      config.ComponentEgress,
			},
			Annotations: map[string]string{
				config.RelatedLiveKitMeshKey: types.NamespacedName{
					Namespace: lkMesh.GetNamespace(),
					Name:      lkMesh.GetName(),
				}.String(),
				config.RelatedConfigMapKey: egressName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To(int32(1)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":     egressName,
					"app.kubernetes.io/instance": egressName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":     egressName,
						"app.kubernetes.io/instance": egressName,
					},
					Annotations: map[string]string{
						config.DefaultConfigMapResourceVersionKey: "",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:            "default",
					TerminationGracePeriodSeconds: ptr.To(int64(3600)),
					Containers: []corev1.Container{{
						Name:            egressName,
						Image:           "livekit/egress:v1.8",
						ImagePullPolicy: "IfNotPresent",
						Ports: []corev1.ContainerPort{
							{
								Name:          "metrics",
								ContainerPort: 7889,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								Name:          "health",
								ContainerPort: 8080,
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
