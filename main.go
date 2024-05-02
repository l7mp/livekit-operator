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

package main

import (
	"flag"
	"fmt"
	renderer "github.com/l7mp/livekit-operator/internal/renderer"
	"github.com/l7mp/livekit-operator/internal/updater"
	stnrgwv1 "github.com/l7mp/stunner-gateway-operator/api/v1"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/l7mp/livekit-operator/internal/operator"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	cert "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	livekitstunnerl7mpiov1alpha1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(livekitstunnerl7mpiov1alpha1.AddToScheme(scheme))
	utilruntime.Must(gwapiv1.AddToScheme(scheme))
	utilruntime.Must(gwapiv1a2.AddToScheme(scheme))
	utilruntime.Must(stnrgwv1.AddToScheme(scheme))
	utilruntime.Must(cert.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr, probeAddr, controllerName string
	var enableLeaderElection bool

	flag.StringVar(&controllerName, "controller-name", opdefault.DefaultControllerName,
		"The controller name to be used in the GatewayClass resource to bind it to this operator.")
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	logger := zap.New(zap.UseFlagOptions(&opts), func(o *zap.Options) {
		o.TimeEncoder = zapcore.RFC3339NanoTimeEncoder
	})
	gracefulShutdown := time.Duration(15000000000)
	ctrl.SetLogger(logger.WithName("ctrl-runtime"))
	setupLog := ctrl.Log.WithName("setup")
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: metricsAddr,
		},
		HealthProbeBindAddress:  probeAddr,
		LeaderElection:          enableLeaderElection,
		LeaderElectionID:        "0386a07e.l7mp.io",
		GracefulShutdownTimeout: &gracefulShutdown,
	})

	// Add your custom runnable to the manager.
	if err := mgr.Add(manager.RunnableFunc(operator.HandleCleanup)); err != nil {
		panic(fmt.Sprintf("Failed to add runnable: %v", err))
	}

	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	setupLog.Info("setting up renderer")
	r := renderer.NewRenderer(renderer.Config{
		Scheme: scheme,
		//Manager: mgr,
		Logger: logger,
	})

	setupLog.Info("setting up updater")
	u := updater.NewUpdater(updater.Config{
		Manager: mgr,
		Logger:  logger,
	})

	setupLog.Info("setting up operator")
	op := operator.NewOperator(operator.Config{
		ControllerName:      controllerName,
		RenderCh:            r.GetRenderChannel(),
		UpdaterCh:           u.GetUpdaterChannel(),
		ShouldInstallCharts: true,
		Manager:             mgr,
		Logger:              logger,
	})

	r.SetOperatorChannel(op.GetOperatorChannel())

	ctx := ctrl.SetupSignalHandler()

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting renderer thread")
	if err := r.Start(ctx); err != nil {
		setupLog.Error(err, "problem running renderer")
		os.Exit(1)
	}

	setupLog.Info("starting updater thread")
	if err := u.Start(ctx); err != nil {
		setupLog.Error(err, "problem running updater")
		os.Exit(1)
	}

	setupLog.Info("starting operator thread")
	if err := op.Start(ctx); err != nil {
		setupLog.Error(err, "problem running operator")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
