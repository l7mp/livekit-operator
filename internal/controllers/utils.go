package controllers

import (
	"context"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/store"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO
func (r *LiveKitMeshReconciler) updateLiveKitMeshStatus(ctx context.Context, req ctrl.Request) *lkstnv1a1.LiveKitMesh {
	log := r.Log.WithName("updateLiveKitMeshStatus")

	log.Info("trying to update LiveKitMesh status")

	lkMesh := &lkstnv1a1.LiveKitMesh{}
	svc := &corev1.Service{}
	dp := &v1.Deployment{}

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
			return nil
		}
	} else if err := r.Get(ctx, req.NamespacedName, dp); err == nil {
		// reconciliation triggering req is a deployment
		if value, ok := dp.GetLabels()[opdefault.RelatedLiveKitMeshKey]; ok {
			lkMeshName := types.NamespacedName{
				Namespace: dp.GetNamespace(),
				Name:      value,
			}
			if err := r.Get(ctx, lkMeshName, lkMesh); err != nil {
				log.Info("could not fetch lkmesh from deployment",
					"lkmesh", lkMeshName,
					"dp", dp.GetName(), "err", err)
			}
		} else {
			return nil
		}
	} else {
		log.Info("Reconcile was triggered by configmap, not updating status")
		return nil
	}

	if lkMesh.Status.ComponentStatus == nil {
		log.Info("Unprocessed LiveKitMesh, initializing its status")
		lkMesh.Status.ComponentStatus = make(map[string]lkstnv1a1.InstallStatus)
	}

	relatedObjects := store.FetchAllObjectsBasedOnLabelFromAllStores(lkMesh.Name)

	livekitComponents := fetchByComponent(relatedObjects, opdefault.ComponentLiveKit)
	if len(livekitComponents) == 0 {
		lkMesh.Status.ComponentStatus[opdefault.ComponentLiveKit] = opdefault.StatusNone
	}
	for _, lkc := range livekitComponents {
		if _, ok := lkc.(*corev1.Service); ok {
			//TODO how to check a clusterIP service properly???
			lkMesh.Status.ComponentStatus[opdefault.ComponentLiveKit] = opdefault.StatusHealthy
		} else {
			lkMesh.Status.ComponentStatus[opdefault.ComponentLiveKit] = opdefault.StatusError
		}
		if dp, ok := lkc.(*v1.Deployment); ok {
			for _, condition := range dp.Status.Conditions {
				if condition.Type == v1.DeploymentAvailable && condition.Status == corev1.ConditionTrue {
					lkMesh.Status.ComponentStatus[opdefault.ComponentLiveKit] = opdefault.StatusHealthy
					break
				} else {
					lkMesh.Status.ComponentStatus[opdefault.ComponentLiveKit] = opdefault.StatusError
				}
			}
		} else {
			lkMesh.Status.ComponentStatus[opdefault.ComponentLiveKit] = opdefault.StatusError
		}
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

	//setting configstatus

	if ok, cm := store.LiveKitMeshes.IsConfigMapReadyForMesh(lkMesh); ok {
		cmData := cm.Data["config.yaml"]
		lkMesh.Status.ConfigStatus = &cmData
	}

	return lkMesh
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
