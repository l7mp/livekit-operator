package renderer

import (
	"fmt"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/store"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	stnrauthsvc "github.com/l7mp/stunner-auth-service/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"
	"strconv"
	"strings"
)

func createLiveKitConfigMap(lkMesh *lkstnv1a1.LiveKitMesh, iceConfig stnrauthsvc.IceConfig) (*corev1.ConfigMap, error) {
	dp := lkMesh.Spec.Components.LiveKit.Deployment
	name := getLiveKitServerConfigMapName(*dp.Name)
	config := dp.Config

	//TODO fix the below code, first turn address is taken, others are left there
	iceServers := *iceConfig.IceServers
	iceAuthenticationToken := iceServers[0]
	urls := *iceAuthenticationToken.Urls
	username := *iceAuthenticationToken.Username
	credential := *iceAuthenticationToken.Credential
	protocol := "udp"
	for i, url := range urls {
		host, port, err := getAddressAndPortFromTurnUrl(url)
		if err != nil {
			return nil, fmt.Errorf("this should only skip the current turn server TODO: %v", err)
		}
		config.Rtc.TurnServers[i].Username = &username
		config.Rtc.TurnServers[i].Credential = &credential
		config.Rtc.TurnServers[i].Host = host
		config.Rtc.TurnServers[i].Port = port
		config.Rtc.TurnServers[i].Protocol = &protocol
	}
	if lkMesh.Spec.Components.LiveKit.Deployment.Config.Redis == nil {
		redisAddress := fmt.Sprintf("%s.%s.svc:6379", getRedisName(lkMesh.GetName()), lkMesh.Namespace)
		config.Redis = &lkstnv1a1.Redis{
			Address: &redisAddress,
		}
	}
	if lkMesh.Spec.Components.Ingress != nil && lkMesh.Spec.Components.ApplicationExpose.HostName != nil {
		rtmp := fmt.Sprintf("rtmp://ingress.%s/rtmp", *lkMesh.Spec.Components.ApplicationExpose.HostName)
		whip := fmt.Sprintf("https://ingress.%s/whip", *lkMesh.Spec.Components.ApplicationExpose.HostName)
		config.IngressAddresses = &lkstnv1a1.IngressAddresses{}
		config.IngressAddresses.RtmpBaseUrl = ptr.To(rtmp)
		config.IngressAddresses.WhipBaseUrl = ptr.To(whip)
	}

	//config.Rtc.StunServers = fmt.Sprintf("%s:%s", )

	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return nil, err
	}
	yamlMap := make(map[string]string)
	yamlMap[opdefault.DefaultLiveKitConfigFileName] = string(yamlData)

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
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

func getAddressAndPortFromTurnUrl(url string) (*string, *int, error) {
	splitString := strings.Split(url, ":")
	address := splitString[1]
	portString := strings.Split(splitString[2], "?")[0]
	portInt, err := strconv.Atoi(portString)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert port string to int: %v", err)
	}
	return &address, &portInt, nil
}

func createLiveKitService(lkMesh *lkstnv1a1.LiveKitMesh) *corev1.Service {

	name := getLiveKitServiceName(*lkMesh.Spec.Components.LiveKit.Deployment.Name)

	labels := map[string]string{
		opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
		opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
		opdefault.RelatedComponent:      opdefault.ComponentLiveKit,
		"app.kubernetes.io/name":        *lkMesh.Spec.Components.LiveKit.Deployment.Name,
		"app.kubernetes.io/instance":    "livekit",
		"app.kubernetes.io/version":     fetchVersion(lkMesh.Spec.Components.LiveKit.Deployment.Container.Image),
	}

	if current := store.Services.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      name,
	}); current != nil {
		labels = mergeMaps(labels, current.Labels)
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: lkMesh.Namespace,
			Name:      name,
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
				"app.kubernetes.io/name":     *lkMesh.Spec.Components.LiveKit.Deployment.Name,
				"app.kubernetes.io/instance": "livekit",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					TargetPort: intstr.FromInt32(7880),
					Port:       443,
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       "rtc-tcp",
					TargetPort: intstr.FromInt32(7801),
					Port:       7801,
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	return svc
}

func createLiveKitDeployment(lkMesh *lkstnv1a1.LiveKitMesh) *appsv1.Deployment {

	containerSpec := lkMesh.Spec.Components.LiveKit.Deployment.Container
	var envList []corev1.EnvVar

	for _, env := range containerSpec.Env {
		env := env
		envList = append(envList, env)
	}

	envList = append(envList, corev1.EnvVar{
		Name: "LIVEKIT_CONFIG",
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: getLiveKitServerConfigMapName(*lkMesh.Spec.Components.LiveKit.Deployment.Name),
				},
				Key: "config.yaml",
			},
		},
	})

	dp := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *lkMesh.Spec.Components.LiveKit.Deployment.Name,
			Namespace: lkMesh.Namespace,
			Labels: map[string]string{
				opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
				opdefault.RelatedComponent:      opdefault.ComponentLiveKit,
				"app.kubernetes.io/name":        *lkMesh.Spec.Components.LiveKit.Deployment.Name,
				"app.kubernetes.io/instance":    "livekit",
				"app.kubernetes.io/version":     fetchVersion(containerSpec.Image),
			},
			Annotations: map[string]string{
				opdefault.RelatedLiveKitMeshKey: types.NamespacedName{
					Namespace: lkMesh.GetNamespace(),
					Name:      lkMesh.GetName(),
				}.String(),
				opdefault.RelatedConfigMapKey: getLiveKitServerConfigMapName(*lkMesh.Spec.Components.LiveKit.Deployment.Name),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: lkMesh.Spec.Components.LiveKit.Deployment.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":     *lkMesh.Spec.Components.LiveKit.Deployment.Name,
					"app.kubernetes.io/instance": "livekit",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":     *lkMesh.Spec.Components.LiveKit.Deployment.Name,
						"app.kubernetes.io/instance": "livekit",
					},
					Annotations: map[string]string{
						opdefault.DefaultConfigMapResourceVersionKey: "",
					},
					//TODO	Finalizers:                 nil,
				},
				Spec: corev1.PodSpec{
					//ServiceAccountName:            "default",
					TerminationGracePeriodSeconds: containerSpec.TerminationGracePeriodSeconds,
					Containers: []corev1.Container{{
						Name:            *lkMesh.Spec.Components.LiveKit.Deployment.Name,
						Image:           containerSpec.Image,
						ImagePullPolicy: *containerSpec.ImagePullPolicy,
						Command:         containerSpec.Command,
						Args:            containerSpec.Args,
						Ports: []corev1.ContainerPort{{
							Name:          "http",
							ContainerPort: 7880,
							Protocol:      corev1.ProtocolTCP,
						}},
						Env:       envList,
						Resources: *containerSpec.Resources,
					},
					},
					HostNetwork:     containerSpec.HostNetwork,
					Affinity:        containerSpec.Affinity,
					SecurityContext: containerSpec.SecurityContext,
					Tolerations:     containerSpec.Tolerations,
				},
			},
		},
	}

	return dp
}

func createLiveKitRedis(lkMesh *lkstnv1a1.LiveKitMesh) (*appsv1.StatefulSet, *corev1.Service, *corev1.ConfigMap) {

	replicasValue := int32(1)

	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getRedisName(lkMesh.GetName()),
			Namespace: lkMesh.GetNamespace(),
			Labels: map[string]string{
				opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
				opdefault.RelatedComponent:      opdefault.ComponentLiveKit,
				"app":                           "redis",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: getRedisName(lkMesh.GetName()),
			Replicas:    &replicasValue,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
					opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
					opdefault.RelatedComponent:      opdefault.ComponentLiveKit,
					"app":                           "redis",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
						opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
						opdefault.RelatedComponent:      opdefault.ComponentLiveKit,
						"app":                           "redis",
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "redis-config",
									},
									Items: []corev1.KeyToPath{{
										Key:  "redis-config",
										Path: "redis.conf",
									}},
								},
							},
						},
					},
					//ServiceAccountName:            "default",
					//TerminationGracePeriodSeconds: containerSpec.TerminationGracePeriodSeconds,
					Containers: []corev1.Container{{
						Name:            "redis",
						Image:           "redis",
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command: []string{
							"redis-server",
							"/redis-master/redis.conf",
						},
						Env: []corev1.EnvVar{
							{
								Name:  "MASTER",
								Value: "true",
							},
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 6379,
						}},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/redis-master-data",
								Name:      "data",
							},
							{
								MountPath: "/redis-master",
								Name:      "config",
							},
						},
					}},
				},
			},
		},
	}

	// MUST be headless svc
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getRedisName(lkMesh.GetName()),
			Namespace: lkMesh.GetNamespace(),
			Labels: map[string]string{
				opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
				opdefault.RelatedComponent:      opdefault.ComponentLiveKit,
				"app":                           "redis",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     6379,
					Protocol: corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						IntVal: 6379,
						StrVal: "6379",
					},
				},
			},
			Selector: map[string]string{
				opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
				opdefault.RelatedComponent:      opdefault.ComponentLiveKit,
				"app":                           "redis",
			},
			ClusterIP: "None",
		},
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis-config",
			Namespace: lkMesh.GetNamespace(),
			Labels: map[string]string{
				opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
				opdefault.RelatedComponent:      opdefault.ComponentLiveKit,
			},
		},
		Data: map[string]string{
			"redis-config": "bind 0.0.0.0",
		},
	}

	return ss, svc, cm
}

// fetchVersion fetches the version from the specified image
func fetchVersion(image string) string {
	// Find the last ":" in the input string
	lastIndex := strings.LastIndex(image, ":")
	if lastIndex != -1 {
		// Version tag found
		version := image[lastIndex+1:]
		return version
	}
	// default is latest
	// in case no tag has been provided
	return "latest"
}
