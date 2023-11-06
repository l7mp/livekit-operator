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
	// +kubebuilder:default=false
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

// NamespacedName is the namespaced name of the configmap that stores
// the base configuration for a given LiveKit container
type NamespacedName struct {
	// Namespace is the namespace of the configMap resource
	Namespace *string `json:"namespace"`
	// Name is the name of the configMap resource
	Name *string `json:"name"`
}

type Deployment struct {

	// Name is the name of the Deployment that will be created.
	// Optional, if not filled default name 'livekit' will be used
	// Note that the same namespace will be used as the CR was deployed into.
	//
	// +kubebuilder:default=livekit-server
	// +optional
	Name *string `json:"name"`

	// Replicas Number of desired pods. This is a pointer to distinguish between explicit zero and not
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
	// TODO in the future we should make a copy from the configmap into the namespace the lkmesh was deployed to
	//
	//ConfigMap *NamespacedName `json:"configMap"`
	ConfigMap *string `json:"configMap"`
}

type LiveKit struct {

	// +kubebuilder:validation:Enum=livekit
	// +kubebuilder:default=livekit
	// +optional
	Type string `json:"type,omitempty"`

	// Deployment holds the configuration for the future Deployment manifest that will be created
	// by the operator.
	//
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

type Gateway struct {
	//Defaults bool `json:"defaults"`

	// RelatedStunnerGatewayAnnotations is the name of the related gateway name for STUNner
	// When deploying the LiveKit server pod we need to know the external IP of the LB SVC
	// that was created based on the very given GW
	// The value of this filed will be present in the SVC's annotation list
	RelatedStunnerGatewayAnnotations *NamespacedName `json:"relatedStunnerGatewayAnnotations"`
}

//type StunnerGatewayOperator struct {
//
//	// Enable sets whether the Chart should be installed or not
//
//	Enable *bool `json:"enable"`
//
//}
//
//type Helm struct {
//
//	// StunnerGatewayOperator holds the values for the STUNner Gateway Operator (https://github.com/l7mp/stunner)
//	// chart.
//	//
//	//
//	// +kubebuilder:validation:Required
//	StunnerGatewayOperator *StunnerGatewayOperator `json:"stunnerGatewayOperator"`
//
//}

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

	// Gateway field should hold the configuration for ANY Gateway deployments in the cluster (STUNner and Envoy)
	Gateway *Gateway `json:"gateway"`

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

	//// Helm holds a configuration for the desired Helm charts in the cluster.
	//// In case the user installs the operator in a cluster that has already one or more of the
	//// Helm charts installed, they can disable the installation to prevent collision.
	////
	//// +kubebuilder:validation:Required
	//Helm *Helm `json:"helm"`
}

// LiveKitMeshSpec defines the desired state of LiveKitMesh
type LiveKitMeshSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	Components *Component `json:"components"`
}

// InstallStatus is the status of the component.
// Enum with the following possible values:
// -NONE meaning the component is not present
// -UPDATING meaning the component is being updated to a different version
// -RECONCILING meaning the controller has started but not yet completed reconciliation loop for the component
// -HEALTHY meaning the component is healthy
// -ERROR meaning a critical error happened to the component
// -ACTION_REQUIRED meaning there is a user action needed in order to proceed
// +kubebuilder:validation:Enum=NONE;UPDATING;RECONCILING;HEALTHY;ERROR;ACTION_REQUIRED
type InstallStatus string

// LiveKitMeshStatus defines the observed state of LiveKitMesh
type LiveKitMeshStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ComponentStatus is a key-value store to signal the components' status after installation
	// The map will give a brief overview for the user which component was successful or failed etc.
	// THE FIELD IS POPULATED BY THE OPERATOR NOT BY THE USER. IT WILL BE OVERWRITTEN
	ComponentStatus map[string]InstallStatus `json:"componentStatus"`

	// OverallStatus of all components controlled by the operator.
	//
	// * If all components have status `NONE`, overall status is `NONE`.
	// * If all components are `HEALTHY`, overall status is `HEALTHY`.
	// * If one or more components are `RECONCILING` and others are `HEALTHY`, overall status is `RECONCILING`.
	// * If one or more components are `UPDATING` and others are `HEALTHY`, overall status is `UPDATING`.
	// * If components are a mix of `RECONCILING`, `UPDATING` and `HEALTHY`, overall status is `UPDATING`.
	// * If any component is in `ERROR` state, overall status is `ERROR`.
	// * If further action is needed for reconciliation to proceed, overall status is `ACTION_REQUIRED`.
	//
	OverallStatus *InstallStatus `json:"overallStatus"`

	// ConfigStatus holds the current configuration for the LiveKit component
	// if it is available in the cluster. nil meaning the ConfigMap provided is not present.
	//
	ConfigStatus *string `json:"configStatus"`
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
