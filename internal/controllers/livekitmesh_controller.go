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
	"github.com/l7mp/livekit-operator/internal/renderer"
	"github.com/l7mp/livekit-operator/internal/store"
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

	liveKitMesh := &lkstnv1a1.LiveKitMesh{}
	defaultConfigMap := &corev1.ConfigMap{}

	if err := r.Get(ctx, req.NamespacedName, liveKitMesh); err != nil {
		if errors.IsNotFound(err) {
			log.Info("no LiveKitMesh resource found")
		} else {
			return ctrl.Result{}, err
		}
	} else {
		log.Info("LiveKitMesh resource found", "name", req.Name)
		log.Info("store length", "len", store.LiveKitMeshes.Len())
		isEqual := store.GetNamespacedName(liveKitMesh) == req.NamespacedName
		log.Info("Namespaced name", "store", store.GetNamespacedName(liveKitMesh),
			"request", req.NamespacedName,
			"isEqual", isEqual)
		if !isEqual {
			panic("store namespacedname does not equal to request's namespacedname")
			//TODO delete this later and solve the issue
		}

		log.Info("livekitmeshes.get", "key", store.LiveKitMeshes.Get(store.GetNamespacedName(liveKitMesh)))
		log.Info("store length", "len", store.LiveKitMeshes.Len())

		//if it does not exist store it
		if ok := store.LiveKitMeshes.Get(store.GetNamespacedName(liveKitMesh)); ok == nil {
			log.Info("New LiveKitMesh found, storing it", "name", liveKitMesh.Name)
			store.LiveKitMeshes.Upsert(liveKitMesh)
			if store.LiveKitMeshes.IsConfigMapReadyForMesh(liveKitMesh) {
				renderer.RenderLiveKitMesh(liveKitMesh)
				// renderCH <- livekitMesh
				// return ctrl.Result{}, nil
				// TODO
			} else {
				return ctrl.Result{}, nil
			}
		} else {
			//TODO lkmesh was already stored, handle this as well
			//TODO check if has been changed etc
			return ctrl.Result{}, nil
		}
	}

	log.Info("If we reach this LoC livekitmesh var should be nil", "liveKitMesh", liveKitMesh)
	log.Info("Trying to get get the corresponding configMap", "configmap", req.NamespacedName)
	if err := r.Get(ctx, req.NamespacedName, defaultConfigMap); err != nil {
		if errors.IsNotFound(err) {
			log.Info("no default config map found")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	} else {
		isEqual := store.GetNamespacedName(defaultConfigMap) == req.NamespacedName
		log.Info("Namespaced name", "store", store.GetNamespacedName(liveKitMesh),
			"request", req.NamespacedName,
			"isEqual", isEqual)
		if !isEqual {
			panic("store namespacedname does not equal to request's namespacedname")
			//TODO delete this later and solve the issue
		}
		//if it does not exist store it
		if ok := store.ConfigMaps.Get(req.NamespacedName); ok == nil {
			log.Info("New ConfigMap found, storing it", "name", defaultConfigMap.Name)
			store.ConfigMaps.Upsert(defaultConfigMap)
			if liveKitMeshes := store.ConfigMaps.GetLiveKitMeshesBasedOnConfigMap(defaultConfigMap); liveKitMeshes != nil {
				for _, mesh := range liveKitMeshes {
					mesh := mesh
					renderer.RenderLiveKitMesh(mesh)
				}
			}
		} else {
			//TODO configmap was already stored, handle this as well
			//TODO check if has been changed etc
			return ctrl.Result{}, nil
		}

		//TODO set the freshly found cm as the config
	}

	// if a configmap change has triggered the reconciliation loop then read it
	// and try to find the corresponding lkmesh in the store
	if err := r.Get(ctx, req.NamespacedName, defaultConfigMap); err != nil {
		if errors.IsNotFound(err) {
			log.Info("no default config map found")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	} else {
		log.Info(defaultConfigMap.String())
		//TODO solve what happens when a new configmap has been found
		//TODO fetch corresponding livekitmeshes from store
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
