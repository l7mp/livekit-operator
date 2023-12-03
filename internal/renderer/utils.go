package renderer

import (
	"fmt"
	"github.com/go-logr/logr"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/store"
)

func (r *Renderer) getLoadBalancerIP(logger logr.Logger, gw *lkstnv1a1.Gateway) *string {
	log := logger.WithName("getLoadBalancerIP")

	serviceList := store.Services.GetAll()
	for _, svc := range serviceList {
		if val, ok := svc.Annotations["stunner.l7mp.io/related-gateway-name"]; ok {
			if val == fmt.Sprintf("%s/%s", *gw.RelatedStunnerGatewayAnnotations.Namespace, *gw.RelatedStunnerGatewayAnnotations.Name) {
				if len(svc.Status.LoadBalancer.Ingress) > 0 {
					log.Info("LoadBalancerIP", "ip", svc.Status.LoadBalancer.Ingress[0].IP)
					return &svc.Status.LoadBalancer.Ingress[0].IP
				}
			}
		}
	}
	return nil
}

func mergeMaps(maps ...map[string]string) map[string]string {
	mergedMap := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			mergedMap[k] = v
		}
	}

	return mergedMap
}

func ServiceNameFormat(lkMeshName string) string {
	return fmt.Sprintf("%s-service", lkMeshName)
}

func ConfigMapNameFormat(lkDeploymentName string) string {
	return fmt.Sprintf("%s-config", lkDeploymentName)
}

func RedisNameFormat(lkMeshName string) string {
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
