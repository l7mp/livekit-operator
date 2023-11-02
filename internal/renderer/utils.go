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
