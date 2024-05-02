package renderer

import (
	"fmt"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func createExternalDNSCloudFlareDeployment(lkMesh *lkstnv1a1.LiveKitMesh) *appsv1.Deployment {

	var envList []corev1.EnvVar
	var argList []string

	envList = append(envList,
		corev1.EnvVar{
			Name:  "CF_API_TOKEN",
			Value: *lkMesh.Spec.Components.ApplicationExpose.ExternalDNS.CloudFlare.Token,
		},
		corev1.EnvVar{
			Name:  "CF_API_EMAIL",
			Value: *lkMesh.Spec.Components.ApplicationExpose.ExternalDNS.CloudFlare.Email,
		})

	argList = append(argList,
		"--source=gateway-httproute",
		"--source=gateway-tcproute",
		"--provider=cloudflare",
		"--cloudflare-dns-records-per-page=1000",
		fmt.Sprintf("--namespace=%s", lkMesh.Namespace),
	)

	dp := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getExternalDNSDeploymentName(lkMesh.Name),
			Namespace: lkMesh.Namespace,
			Labels: map[string]string{
				opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
				opdefault.RelatedComponent:      opdefault.ComponentApplicationExpose,
			},
			Annotations: map[string]string{
				opdefault.RelatedLiveKitMeshKey: types.NamespacedName{
					Namespace: lkMesh.GetNamespace(),
					Name:      lkMesh.GetName(),
				}.String(),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					opdefault.ExternalDNSLabelKey: getExternalDNSDeploymentName(lkMesh.Name),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						opdefault.ExternalDNSLabelKey: getExternalDNSDeploymentName(lkMesh.Name),
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: getExternalDNSServiceAccountName(lkMesh.Name),
					Containers: []corev1.Container{{
						Name:  getExternalDNSDeploymentName(lkMesh.Name),
						Image: "registry.k8s.io/external-dns/external-dns:v0.14.1",
						Args:  argList,
						Env:   envList,
					},
					},
				},
			},
		},
	}

	return dp
}

func createExternalDNSServiceAccount(lkMesh *lkstnv1a1.LiveKitMesh) *corev1.ServiceAccount {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getExternalDNSServiceAccountName(lkMesh.Name),
			Namespace: lkMesh.Namespace,
			Labels: map[string]string{
				opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
				opdefault.RelatedComponent:      opdefault.ComponentApplicationExpose,
			},
		},
	}

	return serviceAccount
}

func createExternalDNSClusterRole(lkMesh *lkstnv1a1.LiveKitMesh) *rbacv1.ClusterRole {
	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: getExternalDNSClusterRoleName(lkMesh.Name, lkMesh.Namespace),
			Labels: map[string]string{
				opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
				opdefault.RelatedComponent:      opdefault.ComponentApplicationExpose,
			},
		},
		Rules: []rbacv1.PolicyRule{{
			APIGroups: []string{"gateway.networking.k8s.io"},
			Resources: []string{"gateways", "httproutes", "tcproutes"},
			Verbs:     []string{"get", "list", "watch"},
		},
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"get", "list", "watch"},
			}},
	}

	return role
}

func createExternalDNSClusterRoleBinding(lkMesh *lkstnv1a1.LiveKitMesh) *rbacv1.ClusterRoleBinding {
	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: getExternalDNSClusterRoleBindingName(lkMesh.Name, lkMesh.Namespace),
			Labels: map[string]string{
				opdefault.OwnedByLabelKey:       opdefault.OwnedByLabelValue,
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
				opdefault.RelatedComponent:      opdefault.ComponentApplicationExpose,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     getExternalDNSClusterRoleName(lkMesh.Name, lkMesh.Namespace),
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      getExternalDNSServiceAccountName(lkMesh.Name),
			Namespace: lkMesh.Namespace,
		}},
	}

	return roleBinding
}
