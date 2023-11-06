package renderer

import (
	"fmt"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func livekitServiceSkeleton(lkMesh *lkstnv1a1.LiveKitMesh) (*corev1.Service, error) {

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: lkMesh.Namespace,
			Name:      fmt.Sprintf("%s-service", *lkMesh.Spec.Components.LiveKit.Deployment.Name),
			Labels: map[string]string{
				opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
				opdefault.RelatedComponent:      opdefault.ComponentLiveKit,
				"app.kubernetes.io/name":        *lkMesh.Spec.Components.LiveKit.Deployment.Name,
				"app.kubernetes.io/instance":    "livekit",
				"app.kubernetes.io/version":     "v1.4.2",
			},
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

	err := controllerutil.SetOwnerReference(lkMesh, svc, nil)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func livekitDeploymentSkeleton(lkMesh *lkstnv1a1.LiveKitMesh, cm *corev1.ConfigMap) (*v1.Deployment, error) {

	containerSpec := lkMesh.Spec.Components.LiveKit.Deployment.Container
	var envList []corev1.EnvVar

	for _, env := range containerSpec.Env {
		env := env
		envList = append(envList, env)
	}
	/*
	   - name: LIVEKIT_CONFIG
	     valueFrom:
	       configMapKeyRef:
	         name: livekit-server
	         key: config.yaml
	*/
	envList = append(envList, corev1.EnvVar{
		Name: "LIVEKIT_CONFIG",
		ValueFrom: &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: *lkMesh.Spec.Components.LiveKit.Deployment.ConfigMap,
				},
				Key: "config.yaml",
			},
		},
	})

	dp := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *lkMesh.Spec.Components.LiveKit.Deployment.Name,
			Namespace: lkMesh.Namespace,
			Labels: map[string]string{
				opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
				opdefault.RelatedComponent:      opdefault.ComponentLiveKit,
				"app.kubernetes.io/name":        *lkMesh.Spec.Components.LiveKit.Deployment.Name,
				"app.kubernetes.io/instance":    "livekit",
				"app.kubernetes.io/version":     "v1.4.2",
			},
			Annotations: map[string]string{
				opdefault.RelatedLiveKitMeshKey: types.NamespacedName{
					Namespace: lkMesh.GetNamespace(),
					Name:      lkMesh.GetName(),
				}.String(),
			},
		},
		Spec: v1.DeploymentSpec{
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
					//TODO	Annotations:                nil,
					//TODO	OwnerReferences:            nil,
					//TODO	Finalizers:                 nil,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:            "default",
					TerminationGracePeriodSeconds: containerSpec.TerminationGracePeriodSeconds,
					Containers: []corev1.Container{{
						Name:    *lkMesh.Spec.Components.LiveKit.Deployment.Name,
						Image:   containerSpec.Image,
						Command: containerSpec.Command,
						Args:    containerSpec.Args,
						Ports: []corev1.ContainerPort{{
							Name:          "http",
							ContainerPort: 7880,
							Protocol:      corev1.ProtocolTCP,
						}},
						Env:       containerSpec.Env,
						Resources: *containerSpec.Resources,
					},
					},
					HostNetwork: containerSpec.HostNetwork,
				},
			},
		},
	}

	err := controllerutil.SetOwnerReference(lkMesh, dp, nil)
	if err != nil {
		return nil, err
	}

	return dp, nil
}
