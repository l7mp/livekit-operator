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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
)

// LiveKitMeshReconciler reconciles a LiveKitMesh object
type LiveKitMeshReconciler struct {
	client.Client
	eventCh chan ievent.Event
	Scheme  *runtime.Scheme
	Log     logr.Logger
}

//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io.l7mp.io,resources=livekitmeshes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io.l7mp.io,resources=livekitmeshes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io.l7mp.io,resources=livekitmeshes/finalizers,verbs=update

func RegisterLiveKitMeshController(mgr manager.Manager, ch chan ievent.Event, logger logr.Logger) error {
	//ctx := context.Background()
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
	if err := r.List(ctx, configMaps); err != nil {
		log.Error(err, "error obtaining ConfigMap objects")
		return ctrl.Result{}, err
	} else {
		for _, cfgmp := range configMaps.Items {
			//TODO if this controller handles it (in case if multiple operators and controllers are running
			configMapList = append(configMapList, &cfgmp)
		}
	}

	store.LiveKitMeshes.Reset(liveKitMeshList)
	log.Info("reset LiveKitMesh store", "lkmeshes", store.LiveKitMeshes.String())

	store.ConfigMaps.Reset(configMapList)
	log.Info("reset ConfigMap store", "configmaps", store.ConfigMaps.String())

	r.eventCh <- ievent.NewEventRender()
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LiveKitMeshReconciler) SetupWithManager(mgr ctrl.Manager) error {

	controller := ctrl.NewControllerManagedBy(mgr).
		For(&lkstnv1a1.LiveKitMesh{})
	controller = controller.
		Watches(&corev1.ConfigMap{}, &handler.EnqueueRequestForObject{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				// Add your custom logic to filter ConfigMap creation events
				configMap, ok := e.Object.(*corev1.ConfigMap)
				if ok {
					// Return true if you want to enqueue the event, false otherwise
					return shouldEnqueue(configMap)
				}
				return true
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				// Add your custom logic to filter ConfigMap update events
				configMap, ok := e.ObjectNew.(*corev1.ConfigMap)
				//fmt.Println("update func", configMap)
				if ok {
					// Return true if you want to enqueue the event, false otherwise
					return shouldEnqueue(configMap)
				}
				return true
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				// Add your custom logic to filter ConfigMap deletion events
				configMap, ok := e.Object.(*corev1.ConfigMap)
				//fmt.Println("delete func", shouldEnqueue(configMap))
				if ok {
					// Return true if you want to enqueue the event, false otherwise
					return shouldEnqueue(configMap)
				}
				return true
			},
		})

	return controller.Complete(r)
}

// Add your custom filtering logic here
func shouldEnqueue(configMap *corev1.ConfigMap) bool {
	// For example, you can check some condition on the ConfigMap
	// and decide whether to enqueue it or not.
	// Return true to enqueue, false to skip.
	// Modify this logic according to your requirements.
	return configMap.Labels[opdefault.DefaultLabelKeyForConfigMap] == opdefault.DefaultLabelValueForConfigMap
}
