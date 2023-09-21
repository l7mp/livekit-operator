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
	"fmt"
	"github.com/go-logr/logr"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
)

// LiveKitMeshReconciler reconciles a LiveKitMesh object
type LiveKitMeshReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io.l7mp.io,resources=livekitmeshes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io.l7mp.io,resources=livekitmeshes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io.l7mp.io,resources=livekitmeshes/finalizers,verbs=update

func RegisterLiveKitMeshController(mgr manager.Manager, logger logr.Logger) error {
	//ctx := context.Background()
	log := logger.WithName("RegisterLiveKitMeshController")

	if err := (&LiveKitMeshReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    logger,
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
	log := r.Log.WithValues("livekit", req.String())
	log.Info("reconciling")

	var liveKitList []client.Object
	defaultConfigMap := &corev1.ConfigMap{}
	//find all configMaps purposed for the operator
	if err := r.Get(ctx, req.NamespacedName, defaultConfigMap); err != nil {
		if errors.IsNotFound(err) {
			log.Info("no default config map found")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	} else {
		log.Info(defaultConfigMap.String())
	}

	// find all LiveKitMesh
	lkList := &lkstnv1a1.LiveKitMeshList{}
	if err := r.List(ctx, lkList); err != nil {
		log.Info("no LiveKitMesh resource found")
		return ctrl.Result{}, err
	}

	for _, i := range lkList.Items {
		lk := i
		log.V(1).Info("processing LiveKitMesh")

		liveKitList = append(liveKitList, &lk)
	}

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
				fmt.Println("update func", configMap)
				if ok {
					// Return true if you want to enqueue the event, false otherwise
					return shouldEnqueue(configMap)
				}
				return true
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				// Add your custom logic to filter ConfigMap deletion events
				configMap, ok := e.Object.(*corev1.ConfigMap)
				fmt.Println("delete func", shouldEnqueue(configMap))
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
