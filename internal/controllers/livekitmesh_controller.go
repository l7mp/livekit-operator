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

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	ievent "github.com/l7mp/livekit-operator/internal/event"
	"github.com/l7mp/livekit-operator/internal/store"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
)

const (
	serviceLiveKitIndex    = "serviceLiveKitIndex"
	configMapLiveKitIndex  = "configMapLiveKitIndex"
	deploymentLiveKitIndex = "deploymentLiveKitIndex"
)

// LiveKitMeshReconciler reconciles a LiveKitMesh object
type LiveKitMeshReconciler struct {
	client.Client
	eventCh chan ievent.Event
	Scheme  *runtime.Scheme
	Log     logr.Logger
}

func RegisterLiveKitMeshController(mgr manager.Manager, ch chan ievent.Event, logger logr.Logger) error {
	log := logger.WithName("RegisterLiveKitMeshController")

	if err := (&LiveKitMeshReconciler{
		Client:  mgr.GetClient(),
		eventCh: ch,
		Scheme:  mgr.GetScheme(),
		Log:     logger,
	}).SetupWithManager(mgr); err != nil {
		log.Error(err, "unable to create controller", "controller", "LiveKitMesh")
		os.Exit(1)
	}

	return nil
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LiveKitMesh object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *LiveKitMeshReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("LiveKitMesh Reconciler", req.String())
	log.Info("reconciling")

	var liveKitMeshList []client.Object
	var configMapList []client.Object
	var serviceList []client.Object
	var deploymentList []client.Object

	//find liveKitMesh resources in the cluster
	liveKitMeshes := &lkstnv1a1.LiveKitMeshList{}
	if err := r.List(ctx, liveKitMeshes); err != nil {
		log.Error(err, "error obtaining LiveKitMesh objects")
		return ctrl.Result{}, err
	} else {
		for _, lkm := range liveKitMeshes.Items {
			lkm := lkm
			//TODO if this controller handles it (in case if multiple operators and controllers are running
			liveKitMeshList = append(liveKitMeshList, &lkm)
		}
	}

	//find configMap resources in the cluster
	configMaps := &corev1.ConfigMapList{}
	listOps := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			map[string]string{
				opdefault.DefaultLabelKeyForConfigMap: opdefault.DefaultLabelValueForConfigMap,
			}),
	}
	if err := r.List(ctx, configMaps, listOps); err != nil {
		log.Error(err, "error obtaining ConfigMap objects")
		return ctrl.Result{}, err
	} else {
		for _, cfgmp := range configMaps.Items {
			cfgmp := cfgmp
			//TODO if this controller handles it (in case if multiple operators and controllers are running
			if shouldEnqueueConfigMap(&cfgmp) {
				configMapList = append(configMapList, &cfgmp)
			}
		}
	}

	//find services resources in the cluster
	services := &corev1.ServiceList{}
	listOpsSvc1 := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			map[string]string{
				opdefault.OwnedByLabelKey: opdefault.OwnedByLabelValue,
			}),
	}
	listOpsSvc2 := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			map[string]string{
				"stunner.l7mp.io/owned-by": "stunner",
			}),
	}
	if err := r.List(ctx, services, listOpsSvc1, listOpsSvc2); err != nil {
		log.Error(err, "error obtaining Service objects")
		return ctrl.Result{}, err
	} else {
		for _, service := range services.Items {
			service := service
			//TODO if this controller handles it (in case if multiple operators and controllers are running
			serviceList = append(serviceList, &service)
		}
	}

	//find services resources in the cluster
	deployments := &v1.DeploymentList{}
	listOpsDeployment := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			map[string]string{
				opdefault.OwnedByLabelKey: opdefault.OwnedByLabelValue,
			}),
	}
	if err := r.List(ctx, deployments, listOpsDeployment); err != nil {
		log.Error(err, "error obtaining Deployment objects")
		return ctrl.Result{}, err
	} else {
		for _, deployment := range deployments.Items {
			deployment := deployment
			//TODO if this controller handles it (in case if multiple operators and controllers are running
			deploymentList = append(deploymentList, &deployment)
		}
	}

	store.LiveKitMeshes.Reset(liveKitMeshList)
	log.Info("reset LiveKitMesh store", "lkmeshes", store.LiveKitMeshes.String())

	store.ConfigMaps.Reset(configMapList)
	log.Info("reset ConfigMap store", "configmaps", store.ConfigMaps.String())

	store.Services.Reset(serviceList)
	log.Info("reset Service store", "services", store.Services.String())

	store.Deployments.Reset(deploymentList)
	log.Info("reset Deployment store", "deployment", store.Deployments.String())

	r.eventCh <- ievent.NewEventRender()

	objectWithUpdatedStatus := r.updateLiveKitMeshStatus(ctx, req)
	if objectWithUpdatedStatus != nil {
		err := r.Status().Update(ctx, objectWithUpdatedStatus)
		if err != nil {
			log.Error(err, "Error happened while updating status")
			return ctrl.Result{}, err
		}
		log.Info("Status updated on LiveKitMesh", "lkMesh", objectWithUpdatedStatus)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LiveKitMeshReconciler) SetupWithManager(mgr ctrl.Manager) error {

	ctx := context.Background()
	controller := ctrl.NewControllerManagedBy(mgr).
		For(&lkstnv1a1.LiveKitMesh{})

	if err := mgr.GetFieldIndexer().IndexField(ctx, &lkstnv1a1.LiveKitMesh{},
		serviceLiveKitIndex, r.serviceMeshIndexFunc); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(ctx, &lkstnv1a1.LiveKitMesh{},
		configMapLiveKitIndex, configMapMeshIndexFunc); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(ctx, &lkstnv1a1.LiveKitMesh{},
		deploymentLiveKitIndex, r.deploymentMeshIndexFunc); err != nil {
		return err
	}

	// a label-selector predicate to select the loadbalancer services we are interested in
	stunnerLoadBalancerPredicate, err := predicate.LabelSelectorPredicate(
		metav1.LabelSelector{
			MatchLabels: map[string]string{
				"stunner.l7mp.io/owned-by": "stunner",
			},
		})
	if err != nil {
		return err
	}

	// a label-selector predicate to select the loadbalancer services we are interested in
	ownedByPredicate, err := predicate.LabelSelectorPredicate(
		metav1.LabelSelector{
			MatchLabels: map[string]string{
				opdefault.OwnedByLabelKey: opdefault.OwnedByLabelValue,
			},
		})
	if err != nil {
		return err
	}

	controller.
		Watches(&corev1.Service{},
			&handler.EnqueueRequestForObject{},
			builder.WithPredicates(predicate.Or(stunnerLoadBalancerPredicate,
				ownedByPredicate))).
		Watches(&corev1.ConfigMap{},
			&handler.EnqueueRequestForObject{},
			builder.WithPredicates(predicate.NewPredicateFuncs(r.validateConfigMapForReconcile))).
		Watches(&v1.Deployment{},
			&handler.EnqueueRequestForObject{},
			builder.WithPredicates(ownedByPredicate))

	return controller.Complete(r)
}

func (r *LiveKitMeshReconciler) validateServiceForReconcile(object client.Object) bool {
	key := ""

	if svc, ok := object.(*corev1.Service); ok {
		key = store.GetObjectKey(svc)
	} else {
		return false
	}

	lkMeshList := &lkstnv1a1.LiveKitMeshList{}

	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(serviceLiveKitIndex, key),
	}

	if err := r.List(context.Background(), lkMeshList, listOps); err != nil {
		r.Log.Error(err, "unable to find associated livekit meshes")
	}
	if len(lkMeshList.Items) > 0 {
		r.Log.Info("Service validation", "lkmeshes list", len(lkMeshList.Items))
	}
	return len(lkMeshList.Items) > 0
}

func (r *LiveKitMeshReconciler) validateConfigMapForReconcile(object client.Object) bool {
	key := ""

	if cm, ok := object.(*corev1.ConfigMap); ok {
		key = store.GetObjectKey(cm)
	} else {
		return false
	}

	lkMeshList := &lkstnv1a1.LiveKitMeshList{}

	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configMapLiveKitIndex, key),
	}

	if err := r.List(context.Background(), lkMeshList, listOps); err != nil {
		r.Log.Error(err, "unable to find associated livekit meshes")
	}
	if len(lkMeshList.Items) > 0 {
		r.Log.Info("configmap validation", "lkmeshes list", len(lkMeshList.Items))
	}
	return len(lkMeshList.Items) > 0
}

func (r *LiveKitMeshReconciler) validateDeploymentForReconcile(object client.Object) bool {
	key := ""

	if dp, ok := object.(*v1.Deployment); ok {
		key = store.GetObjectKey(dp)
	} else {
		return false
	}

	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(deploymentLiveKitIndex, key),
	}

	lkMeshList := &lkstnv1a1.LiveKitMeshList{}

	if err := r.List(context.Background(), lkMeshList, listOps); err != nil {
		r.Log.Error(err, "unable to find associated livekit meshes")
	}
	if len(lkMeshList.Items) > 0 {
		r.Log.Info("deployment validation", "lkmeshes list", len(lkMeshList.Items))
	}
	return len(lkMeshList.Items) > 0
}

func configMapMeshIndexFunc(object client.Object) []string {
	if lkMesh, ok := object.(*lkstnv1a1.LiveKitMesh); ok {
		if lkMesh.Spec.Components.LiveKit.Deployment.ConfigMap == nil {
			return nil
		}
		cm := types.NamespacedName{
			Namespace: *lkMesh.Spec.Components.LiveKit.Deployment.ConfigMap.Namespace,
			Name:      *lkMesh.Spec.Components.LiveKit.Deployment.ConfigMap.Name,
		}.String()
		return []string{cm}
	}
	return nil
}

func (r *LiveKitMeshReconciler) serviceMeshIndexFunc(object client.Object) []string {
	var svcs []string
	services := &corev1.ServiceList{}
	lkMesh := &lkstnv1a1.LiveKitMesh{}

	lkMesh = object.(*lkstnv1a1.LiveKitMesh)

	listOps := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			map[string]string{
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
			}),
	}
	if err := r.List(context.Background(), services, listOps); err != nil {
		r.Log.Error(err, "error obtaining Service objects")
	}

	for _, svc := range services.Items {
		svc := svc
		svcs = append(svcs, store.GetObjectKey(&svc))
	}
	r.Log.Info("indexed services", "svcs", svcs)
	return svcs
}

func (r *LiveKitMeshReconciler) deploymentMeshIndexFunc(object client.Object) []string {
	var dps []string
	deployments := &v1.DeploymentList{}
	lkMesh := &lkstnv1a1.LiveKitMesh{}

	lkMesh = object.(*lkstnv1a1.LiveKitMesh)

	listOps := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			map[string]string{
				opdefault.RelatedLiveKitMeshKey: lkMesh.GetName(),
			}),
	}

	if err := r.List(context.Background(), deployments, listOps); err != nil {
		r.Log.Error(err, "error obtaining Deployment objects")
	}

	for _, dp := range deployments.Items {
		dp := dp
		dps = append(dps, store.GetObjectKey(&dp))
		r.Log.Info("deployments when indexing", "dps", dps)
	}
	return dps
}

// Add your custom filtering logic here
func shouldEnqueueConfigMap(configMap *corev1.ConfigMap) bool {
	// For example, you can check some condition on the ConfigMap
	// and decide whether to enqueue it or not.
	// Return true to enqueue, false to skip.
	// Modify this logic according to your requirements.
	//fmt.Println("shouldEnqueueConfigMap", "configmap", configMap.Name, "bool", configMap.Labels[opdefault.DefaultLabelKeyForConfigMap] == opdefault.DefaultLabelValueForConfigMap)
	return configMap.Labels[opdefault.DefaultLabelKeyForConfigMap] == opdefault.DefaultLabelValueForConfigMap
}
