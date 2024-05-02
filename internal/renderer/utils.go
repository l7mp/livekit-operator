package renderer

import (
	"fmt"
	stnrauthsvc "github.com/l7mp/stunner-auth-service/pkg/types"
	"net"
	"strings"
)

//func (r *Renderer) getLoadBalancerIP(logger logr.Logger, gw *lkstnv1a1.Gateway) *string {
//	log := logger.WithName("getLoadBalancerIP")
//
//	serviceList := store.Services.GetAll()
//	for _, svc := range serviceList {
//		if val, ok := svc.Annotations["stunner.l7mp.io/related-gateway-name"]; ok {
//			if val == fmt.Sprintf("%s/%s", *gw.RelatedStunnerGatewayAnnotations.Namespace, *gw.RelatedStunnerGatewayAnnotations.Name) {
//				if len(svc.Status.LoadBalancer.Ingress) > 0 {
//					log.Info("LoadBalancerIP", "ip", svc.Status.LoadBalancer.Ingress[0].IP)
//					return &svc.Status.LoadBalancer.Ingress[0].IP
//				}
//			}
//		}
//	}
//	return nil
//}

func mergeMaps(maps ...map[string]string) map[string]string {
	mergedMap := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			mergedMap[k] = v
		}
	}

	return mergedMap
}

func getLiveKitServiceName(lkMeshName string) string {
	return fmt.Sprintf("%s-service", lkMeshName)
}

func getLiveKitServerConfigMapName(lkDeploymentName string) string {
	return fmt.Sprintf("%s-config", lkDeploymentName)
}

func getRedisName(lkMeshName string) string {
	return fmt.Sprintf("%s-redis", lkMeshName)
}

//func ParseLiveKitConfigMap(cm v1.ConfigMap) (lkstnv1a1.LiveKitConfig, error) {
//	var lkConfig lkstnv1a1.LiveKitConfig
//
//	yamlConf, found := cm.Data[opdefault.DefaultLiveKitConfigFileName]
//	if !found {
//		return lkConfig, fmt.Errorf("error unpacking configmap data: %s not found",
//			opdefault.DefaultLiveKitConfigFileName)
//	}
//
//	if err := yaml.Unmarshal([]byte(yamlConf), &lkConfig); err != nil {
//		return lkConfig, err
//	}
//
//	return lkConfig, nil
//}

// STUNner related utils
func getStunnerGatewayName(lkMeshName string) string {
	return fmt.Sprintf("%s-%s", lkMeshName, "stunner-gateway")
}

func getStunnerGatewayConfigName(lkMeshName string) string {
	return fmt.Sprintf("%s-%s", lkMeshName, "stunner-gatewayconfig")
}

func getStunnerGatewayClassName(lkMeshName string) string {
	return fmt.Sprintf("%s-%s", lkMeshName, "stunner-gatewayclass")
}

func getStunnerUDPRouteName(lkMeshName string) string {
	return fmt.Sprintf("%s-%s", lkMeshName, "stunner-udproute")
}

// Envoy related utils
func getEnvoyGatewayClassName(lkMeshName string) string {
	return fmt.Sprintf("%s-%s", lkMeshName, "envoy-gatewayclass")
}

func getEnvoyLiveKitServerGatewayName(lkMeshName string) string {
	return fmt.Sprintf("%s-%s", lkMeshName, "server-envoy-gateway")
}

func getEnvoyLiveKitServerGatewayListenerName(lkMeshName string) string {
	return fmt.Sprintf("%s-%s", lkMeshName, "server-envoy-gateway-https")
}

func getEnvoyLiveKitServerGatewayListenerSecretName(lkMeshName string) string {
	return fmt.Sprintf("%s-%s", lkMeshName, "server-envoy-gateway-https-secret")
}

func getEnvoyLiveKitServerHTTPRouteName(lkMeshName string) string {
	return fmt.Sprintf("%s-%s", lkMeshName, "server-envoy-httproute")
}

func getHostNameWithSubDomain(subDomain string, hostName string) string {
	return fmt.Sprintf("%s.%s", subDomain, hostName)
}

func getEnvoyLiveKitIngressGatewayName(lkMeshName string) string {
	return fmt.Sprintf("%s-%s", lkMeshName, "ingress-envoy-gateway")
}

func getEnvoyLiveKitServerTCPRouteName(lkMeshName string) string {
	return fmt.Sprintf("%s-%s", lkMeshName, "ingress-envoy-tcproute")
}

func getEnvoyLiveKitIngressGatewayListenerName(lkMeshName string, mode string) string {
	return fmt.Sprintf("%s-%s-%s", lkMeshName, "ingress-envoy-gateway-tcp", mode)
}

// Network related utils
func validateIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

func getAddressFromIceConfig(iceConfig *stnrauthsvc.IceConfig) string {
	iceServers := *iceConfig.IceServers
	urls := *iceServers[0].Urls
	turnUrl := urls[0]
	address := strings.Split(turnUrl, ":")[1]
	return address
}

// ExternalDNS related utils
func getExternalDNSDeploymentName(name string) string {
	return fmt.Sprintf("%s-external-dns", name)
}

func getExternalDNSServiceAccountName(name string) string {
	return fmt.Sprintf("%s-external-dns-service-account", name)
}

func getExternalDNSClusterRoleName(name string, namespace string) string {
	return fmt.Sprintf("%s-%s-external-dns-role", namespace, name)
}

func getExternalDNSClusterRoleBindingName(name string, namespace string) string {
	return fmt.Sprintf("%s-%s-external-dns-role-binding", namespace, name)
}

// Ingress related utils

func getIngressName(name string) string {
	return fmt.Sprintf("%s-ingress", name)
}
