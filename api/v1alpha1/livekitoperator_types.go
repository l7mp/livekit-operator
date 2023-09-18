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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type TURNServer struct {
	Host       *string `json:"host"`
	Port       *int    `json:"port"`
	Protocol   *string `json:"protocol"`
	Username   *string `json:"username,omitempty"`
	Credential *string `json:"credential,omitempty"`
}

type RTCConfig struct {
	PortRangeStart *int32 `json:"port_range_start"`
	PortRangeEnd   *int32 `json:"port_range_end"`
	TcpPort        *int32 `json:"tcp_port"`

	// TURNServers holds the configuration for the user defined
	// TURN servers. In case users want to define themselves this they can.
	// However, it is advised to omit this to let the operator configure it.
	// +optional
	TURNServers []*TURNServer `json:"turn_servers,omitempty"`
}

type Config struct {
	Keys     *map[string]string `json:"keys"`
	LogLevel *string            `json:"log_level"`
	Port     *int32             `json:"port"`
	Redis    *map[string]string `json:"redis"`
	RTC      *RTCConfig         `json:"rtc"`
}

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

	// Config holds the configuration for the livekit server that is executed.
	//
	//
	Config *Config `json:"config"`
}

type MediaServer struct {

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

	// MediaServer is the main resource that the operator manages. By default, it supports
	// only the LiveKit server as a media server but in the future it might support other
	// media servers as well.
	// +kubebuilder:validation:Required
	MediaServer *MediaServer `json:"mediaServer"`

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

// LiveKitOperatorSpec defines the desired state of LiveKitOperator
type LiveKitOperatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of LiveKitOperator. Edit livekitoperator_types.go to remove/update
	Foo string `json:"foo,omitempty"`

	// +kubebuilder:validation:Required
	Components *Component `json:"components"`
}

// LiveKitOperatorStatus defines the observed state of LiveKitOperator
type LiveKitOperatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LiveKitOperator is the Schema for the livekitoperators API
type LiveKitOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LiveKitOperatorSpec   `json:"spec,omitempty"`
	Status LiveKitOperatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LiveKitOperatorList contains a list of LiveKitOperator
type LiveKitOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LiveKitOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LiveKitOperator{}, &LiveKitOperatorList{})
}
