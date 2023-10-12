/*
Copyright 2023 Kornel David.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Container struct {
	// Container image name.
	//
	// +optional
	Image string `json:"image,omitempty"`

	// Image pull policy. One of Always, Never, IfNotPresent.
	//
	// +optional
	ImagePullPolicy *corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// Entrypoint array. Defaults: "stunnerd".
	//
	// +optional
	Command []string `json:"command,omitempty"`

	// Arguments to the entrypoint.
	//
	// +optional
	Args []string `json:"args,omitempty"`

	// List of environment variables to set in the stunnerd container.
	//
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`

	// Resources required by stunnerd.
	//
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// Optional duration in seconds the stunnerd needs to terminate gracefully. Defaults to 3600 seconds.
	//
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`

	// Host networking requested for the stunnerd pod to use the host's network namespace.
	// Can be used to implement public TURN servers with Kubernetes.  Defaults to false.
	//
	// +optional
	HostNetwork bool `json:"hostNetwork,omitempty"`

	// Scheduling constraints.
	//
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// SecurityContext holds pod-level security attributes and common container settings.
	//
	// +optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// If specified, the pod's tolerations.
	//
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// If specified, the health-check port.
	//
	// +optional
	HealthCheckPort *int `json:"healthCheckPort,omitempty"`

	// // If specified, the metrics collection port.
	// //
	// // +optional
	// MetricsEndpointPort *int `json:"metricsEndpointPort,omitempty"`
}

// ConfigMapNamespacedName is the namespaced name of the configmap that stores
// the base configuration for a given LiveKit container
type ConfigMapNamespacedName struct {
	// Namespace is the namespace of the configMap resource
	Namespace *string `json:"namespace"`
	// Name is the name of the configMap resource
	Name *string `json:"name"`
}

type Deployment struct {

	// Number of desired pods. This is a pointer to distinguish between explicit zero and not
	// specified. Defaults to 1.
	//
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Container template for the containers created in each Pod in the replicaset.
	// If omitted the default template will be used. Which spawns a single container
	//TODO
	//
	// +optional
	Container *Container `json:"container"`

	// ConfigMap holds the configuration for the livekit server that is executed.
	//
	//
	ConfigMap *ConfigMapNamespacedName `json:"configMap"`
}

type LiveKit struct {

	// +kubebuilder:validation:Enum=livekit
	// +kubebuilder:default=livekit
	// +optional
	Type string `json:"type,omitempty"`

	// +kubebuilder:validation:Required
	Deployment *Deployment `json:"deployment"`
}

// Ingress is the LiveKit tool not the gateway resource to ingest traffic into the cluster
type Ingress struct {
	//TODO
}

type Egress struct {
	//TODO
}

type CertManager struct {
	//TODO
}

type Monitoring struct {
	//TODO
}

type Component struct {

	// LiveKit is the main resource that the operator manages. By default, it supports
	// only the LiveKit server as a media server but in the future it might support other
	// media servers as well.
	// +kubebuilder:validation:Required
	LiveKit *LiveKit `json:"liveKit"`

	// LiveKit's Ingress resource descriptor.
	// This resource makes it possible to stream videos(prerecorded or live) into
	// the Kubernetes cluster and further into a chosen room. Note that this resource
	// enables a one-way communication between the client and the media server.
	//
	// +optional
	Ingress *Ingress `json:"ingress,omitempty"`

	// LiveKit's Egress resource descriptor.
	// Egress makes it possible to stream a single user's or any number of users'
	// streams out of a room onto an RTMP port.
	//
	// +optional
	Egress *Egress `json:"egress,omitempty"`

	// CertManager manages the cert
	//TODO
	//
	//
	// +optional
	CertManager *CertManager `json:"certManager,omitempty"`

	// Monitoring enables the Prometheus metric exposition, installs
	// a Prometheus operator and Grafana operator with the corresponding resources
	//
	// +optional
	Monitoring *Monitoring `json:"monitoring,omitempty"`
}

// LiveKitMeshSpec defines the desired state of LiveKitMesh
type LiveKitMeshSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	Components *Component `json:"components"`
}

// LiveKitMeshStatus defines the observed state of LiveKitMesh
type LiveKitMeshStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LiveKitMesh is the Schema for the livekitmeshes API
type LiveKitMesh struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LiveKitMeshSpec   `json:"spec,omitempty"`
	Status LiveKitMeshStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LiveKitMeshList contains a list of LiveKitMesh
type LiveKitMeshList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LiveKitMesh `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LiveKitMesh{}, &LiveKitMeshList{})
}
