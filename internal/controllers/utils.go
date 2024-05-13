package controllers

import (
	"context"
	"github.com/go-logr/logr"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/store"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// TODO
func (r *LiveKitMeshReconciler) updateLiveKitMeshStatus(ctx context.Context, req ctrl.Request) *lkstnv1a1.LiveKitMesh {
	log := r.Log.WithName("updateLiveKitMeshStatus")

	log.Info("trying to update LiveKitMesh status")

	lkMesh, _, _ := r.fetchLiveKitMesh(log, ctx, req)

	if lkMesh == nil {
		return nil
	}
	log.Info("trying to update LiveKitMesh status", "status", lkMesh.Status)
	if lkMesh.Status.ComponentStatus == nil {
		log.Info("Unprocessed LiveKitMesh, initializing its status")
		lkMesh.Status.ComponentStatus = make(map[string]lkstnv1a1.InstallStatus)
	}

	relatedObjects := store.FetchAllObjectsBasedOnLabelFromAllStores(lkMesh.Name)

	livekitComponentRelatedResources := fetchByComponent(relatedObjects, opdefault.ComponentLiveKit)
	if len(livekitComponentRelatedResources) == 0 {
		lkMesh.Status.ComponentStatus[opdefault.ComponentLiveKit] = opdefault.StatusNone
	} else {
		for _, lkc := range livekitComponentRelatedResources {
			if _, ok := lkc.(*corev1.Service); ok {
				lkMesh.Status.ComponentStatus[opdefault.ComponentLiveKit] = opdefault.StatusHealthy
			} else if dp, ok := lkc.(*v1.Deployment); ok {
				for _, condition := range dp.Status.Conditions {
					if condition.Type == v1.DeploymentAvailable && condition.Status == corev1.ConditionTrue {
						lkMesh.Status.ComponentStatus[opdefault.ComponentLiveKit] = opdefault.StatusHealthy
						break
					} else {
						lkMesh.Status.ComponentStatus[opdefault.ComponentLiveKit] = opdefault.StatusError
					}
				}
			}
			//else {
			//	lkMesh.Status.ComponentStatus[opdefault.ComponentLiveKit] = opdefault.StatusError
			//}
		}
	}
	ingressComponents := fetchByComponent(relatedObjects, opdefault.ComponentIngress)
	if len(ingressComponents) == 0 {
		lkMesh.Status.ComponentStatus[opdefault.ComponentIngress] = opdefault.StatusNone
	}
	egressComponents := fetchByComponent(relatedObjects, opdefault.ComponentEgress)
	if len(egressComponents) == 0 {
		lkMesh.Status.ComponentStatus[opdefault.ComponentEgress] = opdefault.StatusNone
	}

	//TODO set all other component statuses

	for _, v := range lkMesh.Status.ComponentStatus {
		if v == opdefault.StatusHealthy || v == opdefault.StatusNone {
			overallStatus := lkstnv1a1.InstallStatus(opdefault.StatusHealthy)
			lkMesh.Status.OverallStatus = &overallStatus
		} else {
			overallStatus := lkstnv1a1.InstallStatus(opdefault.StatusError)
			lkMesh.Status.OverallStatus = &overallStatus
		}
	}

	// TODO remove configstatus
	dummystatus := "dev"
	lkMesh.Status.ConfigStatus = &dummystatus

	return lkMesh
}

func (r *LiveKitMeshReconciler) fetchLiveKitMesh(log logr.Logger, ctx context.Context, req ctrl.Request) (*lkstnv1a1.LiveKitMesh, *corev1.Service, *gwapiv1.Gateway) {
	log = log.WithName("fetchLiveKitMesh")

	lkMesh := &lkstnv1a1.LiveKitMesh{}
	svc := &corev1.Service{}
	gw := &gwapiv1.Gateway{}

	if err := r.Get(ctx, req.NamespacedName, lkMesh); err == nil {
		// reconciliation triggering req is a livekitMesh
		log.Info("LiveKit mesh triggered")

	} else if err := r.Get(ctx, req.NamespacedName, svc); err == nil {
		// reconciliation triggering req is a service
		if value, ok := svc.GetLabels()[opdefault.RelatedLiveKitMeshKey]; ok {
			lkMeshName := types.NamespacedName{
				Namespace: svc.GetNamespace(),
				Name:      value,
			}
			if err := r.Get(ctx, lkMeshName, lkMesh); err != nil {
				log.Info("could not fetch lkmesh from svc",
					"lkmesh", lkMeshName,
					"svc", svc.GetName(), "err", err)
			}
		} else {
			return nil, nil, nil
		}
	} else if err := r.Get(ctx, req.NamespacedName, gw); err == nil {
		// reconciliation triggering req is a gateway
		if value, ok := gw.GetLabels()[opdefault.RelatedLiveKitMeshKey]; ok {
			lkMeshName := types.NamespacedName{
				Namespace: gw.GetNamespace(),
				Name:      value,
			}
			if err := r.Get(ctx, lkMeshName, lkMesh); err != nil {
				log.Info("could not fetch lkmesh from gateway",
					"lkmesh", lkMeshName,
					"gw", gw.GetName(), "err", err)
			}
		} else {
			return nil, nil, nil
		}
	} else {
		log.Info("Reconcile was triggered by configmap, not updating status")
		return nil, nil, nil
	}

	return lkMesh, svc, gw
}

func fetchByComponent(ro []client.Object, component string) []client.Object {
	var objects []client.Object
	for _, ro := range ro {
		if ro.GetLabels()[opdefault.RelatedComponent] == component {
			objects = append(objects, ro)
		}
	}
	return objects
}
