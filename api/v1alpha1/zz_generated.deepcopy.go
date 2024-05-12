//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	apiv1 "github.com/l7mp/stunner-gateway-operator/api/v1"
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	apisv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApplicationExpose) DeepCopyInto(out *ApplicationExpose) {
	*out = *in
	if in.HostName != nil {
		in, out := &in.HostName, &out.HostName
		*out = new(string)
		**out = **in
	}
	if in.CertManager != nil {
		in, out := &in.CertManager, &out.CertManager
		*out = new(CertManager)
		(*in).DeepCopyInto(*out)
	}
	if in.ExternalDNS != nil {
		in, out := &in.ExternalDNS, &out.ExternalDNS
		*out = new(ExternalDNS)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApplicationExpose.
func (in *ApplicationExpose) DeepCopy() *ApplicationExpose {
	if in == nil {
		return nil
	}
	out := new(ApplicationExpose)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Azure) DeepCopyInto(out *Azure) {
	*out = *in
	if in.AccountName != nil {
		in, out := &in.AccountName, &out.AccountName
		*out = new(string)
		**out = **in
	}
	if in.AccountKey != nil {
		in, out := &in.AccountKey, &out.AccountKey
		*out = new(string)
		**out = **in
	}
	if in.ContainerName != nil {
		in, out := &in.ContainerName, &out.ContainerName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Azure.
func (in *Azure) DeepCopy() *Azure {
	if in == nil {
		return nil
	}
	out := new(Azure)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUCost) DeepCopyInto(out *CPUCost) {
	*out = *in
	if in.RTMPCPUCost != nil {
		in, out := &in.RTMPCPUCost, &out.RTMPCPUCost
		*out = new(int)
		**out = **in
	}
	if in.WHIPCPUCost != nil {
		in, out := &in.WHIPCPUCost, &out.WHIPCPUCost
		*out = new(int)
		**out = **in
	}
	if in.WHIPBypassTranscodingCPUCost != nil {
		in, out := &in.WHIPBypassTranscodingCPUCost, &out.WHIPBypassTranscodingCPUCost
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUCost.
func (in *CPUCost) DeepCopy() *CPUCost {
	if in == nil {
		return nil
	}
	out := new(CPUCost)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertManager) DeepCopyInto(out *CertManager) {
	*out = *in
	if in.Issuer != nil {
		in, out := &in.Issuer, &out.Issuer
		*out = new(Issuer)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertManager.
func (in *CertManager) DeepCopy() *CertManager {
	if in == nil {
		return nil
	}
	out := new(CertManager)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudFlare) DeepCopyInto(out *CloudFlare) {
	*out = *in
	if in.Token != nil {
		in, out := &in.Token, &out.Token
		*out = new(string)
		**out = **in
	}
	if in.Email != nil {
		in, out := &in.Email, &out.Email
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudFlare.
func (in *CloudFlare) DeepCopy() *CloudFlare {
	if in == nil {
		return nil
	}
	out := new(CloudFlare)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Component) DeepCopyInto(out *Component) {
	*out = *in
	if in.LiveKit != nil {
		in, out := &in.LiveKit, &out.LiveKit
		*out = new(LiveKit)
		(*in).DeepCopyInto(*out)
	}
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(Ingress)
		(*in).DeepCopyInto(*out)
	}
	if in.Egress != nil {
		in, out := &in.Egress, &out.Egress
		*out = new(Egress)
		(*in).DeepCopyInto(*out)
	}
	if in.ApplicationExpose != nil {
		in, out := &in.ApplicationExpose, &out.ApplicationExpose
		*out = new(ApplicationExpose)
		(*in).DeepCopyInto(*out)
	}
	if in.Monitoring != nil {
		in, out := &in.Monitoring, &out.Monitoring
		*out = new(Monitoring)
		**out = **in
	}
	if in.Stunner != nil {
		in, out := &in.Stunner, &out.Stunner
		*out = new(Stunner)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Component.
func (in *Component) DeepCopy() *Component {
	if in == nil {
		return nil
	}
	out := new(Component)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Container) DeepCopyInto(out *Container) {
	*out = *in
	if in.ImagePullPolicy != nil {
		in, out := &in.ImagePullPolicy, &out.ImagePullPolicy
		*out = new(v1.PullPolicy)
		**out = **in
	}
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.TerminationGracePeriodSeconds != nil {
		in, out := &in.TerminationGracePeriodSeconds, &out.TerminationGracePeriodSeconds
		*out = new(int64)
		**out = **in
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.SecurityContext != nil {
		in, out := &in.SecurityContext, &out.SecurityContext
		*out = new(v1.PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.HealthCheckPort != nil {
		in, out := &in.HealthCheckPort, &out.HealthCheckPort
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Container.
func (in *Container) DeepCopy() *Container {
	if in == nil {
		return nil
	}
	out := new(Container)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Deployment) DeepCopyInto(out *Deployment) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.Container != nil {
		in, out := &in.Container, &out.Container
		*out = new(Container)
		(*in).DeepCopyInto(*out)
	}
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(LiveKitConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Deployment.
func (in *Deployment) DeepCopy() *Deployment {
	if in == nil {
		return nil
	}
	out := new(Deployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Egress) DeepCopyInto(out *Egress) {
	*out = *in
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(EgressConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Container != nil {
		in, out := &in.Container, &out.Container
		*out = new(Container)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Egress.
func (in *Egress) DeepCopy() *Egress {
	if in == nil {
		return nil
	}
	out := new(Egress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressConfig) DeepCopyInto(out *EgressConfig) {
	*out = *in
	if in.HealthPort != nil {
		in, out := &in.HealthPort, &out.HealthPort
		*out = new(int)
		**out = **in
	}
	if in.TemplatePort != nil {
		in, out := &in.TemplatePort, &out.TemplatePort
		*out = new(int)
		**out = **in
	}
	if in.PrometheusPort != nil {
		in, out := &in.PrometheusPort, &out.PrometheusPort
		*out = new(int)
		**out = **in
	}
	if in.LogLevel != nil {
		in, out := &in.LogLevel, &out.LogLevel
		*out = new(string)
		**out = **in
	}
	if in.Insecure != nil {
		in, out := &in.Insecure, &out.Insecure
		*out = new(bool)
		**out = **in
	}
	if in.S3 != nil {
		in, out := &in.S3, &out.S3
		*out = new(S3)
		(*in).DeepCopyInto(*out)
	}
	if in.Azure != nil {
		in, out := &in.Azure, &out.Azure
		*out = new(Azure)
		(*in).DeepCopyInto(*out)
	}
	if in.Gcp != nil {
		in, out := &in.Gcp, &out.Gcp
		*out = new(Gcp)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressConfig.
func (in *EgressConfig) DeepCopy() *EgressConfig {
	if in == nil {
		return nil
	}
	out := new(EgressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNS) DeepCopyInto(out *ExternalDNS) {
	*out = *in
	if in.CloudFlare != nil {
		in, out := &in.CloudFlare, &out.CloudFlare
		*out = new(CloudFlare)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNS.
func (in *ExternalDNS) DeepCopy() *ExternalDNS {
	if in == nil {
		return nil
	}
	out := new(ExternalDNS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Gcp) DeepCopyInto(out *Gcp) {
	*out = *in
	if in.CredentialsJson != nil {
		in, out := &in.CredentialsJson, &out.CredentialsJson
		*out = new(string)
		**out = **in
	}
	if in.Bucket != nil {
		in, out := &in.Bucket, &out.Bucket
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Gcp.
func (in *Gcp) DeepCopy() *Gcp {
	if in == nil {
		return nil
	}
	out := new(Gcp)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ingress) DeepCopyInto(out *Ingress) {
	*out = *in
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(IngressConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ingress.
func (in *Ingress) DeepCopy() *Ingress {
	if in == nil {
		return nil
	}
	out := new(Ingress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressConfig) DeepCopyInto(out *IngressConfig) {
	*out = *in
	if in.CPUCost != nil {
		in, out := &in.CPUCost, &out.CPUCost
		*out = new(CPUCost)
		(*in).DeepCopyInto(*out)
	}
	if in.HealthPort != nil {
		in, out := &in.HealthPort, &out.HealthPort
		*out = new(int)
		**out = **in
	}
	if in.PrometheusPort != nil {
		in, out := &in.PrometheusPort, &out.PrometheusPort
		*out = new(int)
		**out = **in
	}
	if in.RTMPPort != nil {
		in, out := &in.RTMPPort, &out.RTMPPort
		*out = new(int)
		**out = **in
	}
	if in.WHIPPort != nil {
		in, out := &in.WHIPPort, &out.WHIPPort
		*out = new(int)
		**out = **in
	}
	if in.HTTPRelayPort != nil {
		in, out := &in.HTTPRelayPort, &out.HTTPRelayPort
		*out = new(int)
		**out = **in
	}
	if in.Logging != nil {
		in, out := &in.Logging, &out.Logging
		*out = new(Logging)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressConfig.
func (in *IngressConfig) DeepCopy() *IngressConfig {
	if in == nil {
		return nil
	}
	out := new(IngressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Issuer) DeepCopyInto(out *Issuer) {
	*out = *in
	if in.Email != nil {
		in, out := &in.Email, &out.Email
		*out = new(string)
		**out = **in
	}
	if in.ChallengeSolver != nil {
		in, out := &in.ChallengeSolver, &out.ChallengeSolver
		*out = new(string)
		**out = **in
	}
	if in.ApiToken != nil {
		in, out := &in.ApiToken, &out.ApiToken
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Issuer.
func (in *Issuer) DeepCopy() *Issuer {
	if in == nil {
		return nil
	}
	out := new(Issuer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Keys) DeepCopyInto(out *Keys) {
	*out = *in
	if in.AccessToken != nil {
		in, out := &in.AccessToken, &out.AccessToken
		*out = new(map[string]string)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]string, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Keys.
func (in *Keys) DeepCopy() *Keys {
	if in == nil {
		return nil
	}
	out := new(Keys)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LiveKit) DeepCopyInto(out *LiveKit) {
	*out = *in
	if in.Deployment != nil {
		in, out := &in.Deployment, &out.Deployment
		*out = new(Deployment)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LiveKit.
func (in *LiveKit) DeepCopy() *LiveKit {
	if in == nil {
		return nil
	}
	out := new(LiveKit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LiveKitConfig) DeepCopyInto(out *LiveKitConfig) {
	*out = *in
	if in.Keys != nil {
		in, out := &in.Keys, &out.Keys
		*out = new(map[string]string)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]string, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
	if in.LogLevel != nil {
		in, out := &in.LogLevel, &out.LogLevel
		*out = new(string)
		**out = **in
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int)
		**out = **in
	}
	if in.Redis != nil {
		in, out := &in.Redis, &out.Redis
		*out = new(Redis)
		(*in).DeepCopyInto(*out)
	}
	if in.Rtc != nil {
		in, out := &in.Rtc, &out.Rtc
		*out = new(Rtc)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LiveKitConfig.
func (in *LiveKitConfig) DeepCopy() *LiveKitConfig {
	if in == nil {
		return nil
	}
	out := new(LiveKitConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LiveKitMesh) DeepCopyInto(out *LiveKitMesh) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LiveKitMesh.
func (in *LiveKitMesh) DeepCopy() *LiveKitMesh {
	if in == nil {
		return nil
	}
	out := new(LiveKitMesh)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LiveKitMesh) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LiveKitMeshList) DeepCopyInto(out *LiveKitMeshList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LiveKitMesh, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LiveKitMeshList.
func (in *LiveKitMeshList) DeepCopy() *LiveKitMeshList {
	if in == nil {
		return nil
	}
	out := new(LiveKitMeshList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LiveKitMeshList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LiveKitMeshSpec) DeepCopyInto(out *LiveKitMeshSpec) {
	*out = *in
	if in.Components != nil {
		in, out := &in.Components, &out.Components
		*out = new(Component)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LiveKitMeshSpec.
func (in *LiveKitMeshSpec) DeepCopy() *LiveKitMeshSpec {
	if in == nil {
		return nil
	}
	out := new(LiveKitMeshSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LiveKitMeshStatus) DeepCopyInto(out *LiveKitMeshStatus) {
	*out = *in
	if in.ComponentStatus != nil {
		in, out := &in.ComponentStatus, &out.ComponentStatus
		*out = make(map[string]InstallStatus, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.OverallStatus != nil {
		in, out := &in.OverallStatus, &out.OverallStatus
		*out = new(InstallStatus)
		**out = **in
	}
	if in.ConfigStatus != nil {
		in, out := &in.ConfigStatus, &out.ConfigStatus
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LiveKitMeshStatus.
func (in *LiveKitMeshStatus) DeepCopy() *LiveKitMeshStatus {
	if in == nil {
		return nil
	}
	out := new(LiveKitMeshStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Logging) DeepCopyInto(out *Logging) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Logging.
func (in *Logging) DeepCopy() *Logging {
	if in == nil {
		return nil
	}
	out := new(Logging)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Monitoring) DeepCopyInto(out *Monitoring) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Monitoring.
func (in *Monitoring) DeepCopy() *Monitoring {
	if in == nil {
		return nil
	}
	out := new(Monitoring)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespacedName) DeepCopyInto(out *NamespacedName) {
	*out = *in
	if in.Namespace != nil {
		in, out := &in.Namespace, &out.Namespace
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespacedName.
func (in *NamespacedName) DeepCopy() *NamespacedName {
	if in == nil {
		return nil
	}
	out := new(NamespacedName)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Redis) DeepCopyInto(out *Redis) {
	*out = *in
	if in.Address != nil {
		in, out := &in.Address, &out.Address
		*out = new(string)
		**out = **in
	}
	if in.Username != nil {
		in, out := &in.Username, &out.Username
		*out = new(string)
		**out = **in
	}
	if in.Password != nil {
		in, out := &in.Password, &out.Password
		*out = new(string)
		**out = **in
	}
	if in.Db != nil {
		in, out := &in.Db, &out.Db
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Redis.
func (in *Redis) DeepCopy() *Redis {
	if in == nil {
		return nil
	}
	out := new(Redis)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rtc) DeepCopyInto(out *Rtc) {
	*out = *in
	if in.PortRangeEnd != nil {
		in, out := &in.PortRangeEnd, &out.PortRangeEnd
		*out = new(int)
		**out = **in
	}
	if in.PortRangeStart != nil {
		in, out := &in.PortRangeStart, &out.PortRangeStart
		*out = new(int)
		**out = **in
	}
	if in.TcpPort != nil {
		in, out := &in.TcpPort, &out.TcpPort
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rtc.
func (in *Rtc) DeepCopy() *Rtc {
	if in == nil {
		return nil
	}
	out := new(Rtc)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rtmp) DeepCopyInto(out *Rtmp) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rtmp.
func (in *Rtmp) DeepCopy() *Rtmp {
	if in == nil {
		return nil
	}
	out := new(Rtmp)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *S3) DeepCopyInto(out *S3) {
	*out = *in
	if in.AccessKey != nil {
		in, out := &in.AccessKey, &out.AccessKey
		*out = new(string)
		**out = **in
	}
	if in.Secret != nil {
		in, out := &in.Secret, &out.Secret
		*out = new(string)
		**out = **in
	}
	if in.Region != nil {
		in, out := &in.Region, &out.Region
		*out = new(string)
		**out = **in
	}
	if in.Endpoint != nil {
		in, out := &in.Endpoint, &out.Endpoint
		*out = new(string)
		**out = **in
	}
	if in.Bucket != nil {
		in, out := &in.Bucket, &out.Bucket
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new S3.
func (in *S3) DeepCopy() *S3 {
	if in == nil {
		return nil
	}
	out := new(S3)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Stunner) DeepCopyInto(out *Stunner) {
	*out = *in
	if in.GatewayConfig != nil {
		in, out := &in.GatewayConfig, &out.GatewayConfig
		*out = new(apiv1.GatewayConfigSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.GatewayListeners != nil {
		in, out := &in.GatewayListeners, &out.GatewayListeners
		*out = make([]apisv1.Listener, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Stunner.
func (in *Stunner) DeepCopy() *Stunner {
	if in == nil {
		return nil
	}
	out := new(Stunner)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TurnServer) DeepCopyInto(out *TurnServer) {
	*out = *in
	if in.Credential != nil {
		in, out := &in.Credential, &out.Credential
		*out = new(string)
		**out = **in
	}
	if in.Host != nil {
		in, out := &in.Host, &out.Host
		*out = new(string)
		**out = **in
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int)
		**out = **in
	}
	if in.Protocol != nil {
		in, out := &in.Protocol, &out.Protocol
		*out = new(string)
		**out = **in
	}
	if in.Username != nil {
		in, out := &in.Username, &out.Username
		*out = new(string)
		**out = **in
	}
	if in.AuthURI != nil {
		in, out := &in.AuthURI, &out.AuthURI
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TurnServer.
func (in *TurnServer) DeepCopy() *TurnServer {
	if in == nil {
		return nil
	}
	out := new(TurnServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Whip) DeepCopyInto(out *Whip) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Whip.
func (in *Whip) DeepCopy() *Whip {
	if in == nil {
		return nil
	}
	out := new(Whip)
	in.DeepCopyInto(out)
	return out
}
