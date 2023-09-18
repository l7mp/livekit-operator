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
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// LiveKitOperatorReconciler reconciles a LiveKitOperator object
type LiveKitOperatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io.l7mp.io,resources=livekitoperators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io.l7mp.io,resources=livekitoperators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=livekit.stunner.l7mp.io.l7mp.io,resources=livekitoperators/finalizers,verbs=update

func RegisterLiveKitOperatorController(mgr manager.Manager, logger logr.Logger) error {
	//ctx := context.Background()
	log := logger.WithName("RegisterLiveKitOperatorController")

	if err := (&LiveKitOperatorReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    logger,
	}).SetupWithManager(mgr); err != nil {
		log.Error(err, "unable to create controller", "controller", "LiveKitOperator")
		os.Exit(1)
	}

	//c, err := controller.New("livekitoperator", mgr, controller.Options{Reconciler: r})
	//if err != nil {
	//	return err
	//}
	//r.log.Info("created livekitoperator controller")
	//
	//if err := c.Watch(
	//	source.Kind(mgr.GetCache(), &corev1.ConfigMap{}),
	//	&handler.EnqueueRequestForObject{},
	//	predicate.GenerationChangedPredicate{},
	//); err != nil {
	//	return err
	//}
	return nil
}

func (r *LiveKitOperatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("livekit", req.String())
	log.Info("reconciling")

	var liveKitList []client.Object

	// find all LiveKitOperator
	lkList := &lkstnv1a1.LiveKitOperatorList{}
	if err := r.List(ctx, lkList); err != nil {
		r.Log.Info("no LiveKitOperator resource found")
		return reconcile.Result{}, err
	}

	for _, i := range lkList.Items {
		lk := i
		r.Log.V(1).Info("processing LiveKitOperator")

		liveKitList = append(liveKitList, &lk)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LiveKitOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lkstnv1a1.LiveKitOperator{}).
		Watches(&corev1.ConfigMap{}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
