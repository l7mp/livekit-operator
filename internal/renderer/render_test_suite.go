package renderer

import (
	"context"
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/l7mp/livekit-operator/internal/store"
	stnrgwv1 "github.com/l7mp/stunner-gateway-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	"testing"

	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"go.uber.org/zap/zapcore"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var testerLogLevel = zapcore.ErrorLevel

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(lkstnv1a1.AddToScheme(scheme))
}

type renderTestConfig struct {
	name          string
	lkMeshes      []lkstnv1a1.LiveKitMesh
	svcs          []corev1.Service
	dps           []appv1.Deployment
	cms           []corev1.ConfigMap
	gwcs          []gwapiv1.GatewayClass
	stnrGwconfigs []stnrgwv1.GatewayConfig
	gws           []gwapiv1.Gateway
	udpRoutes     []stnrgwv1.UDPRoute
	prep          func(c *renderTestConfig)
	tester        func(t *testing.T, r *Renderer)
}

func renderTester(t *testing.T, testConfig []renderTestConfig) {
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	z, err := zc.Build()
	assert.NoError(t, err, "logger created")
	zlogger := zapr.NewLogger(z)
	log := zlogger.WithName("tester")

	for i := range testConfig {
		c := testConfig[i]
		t.Run(c.name, func(t *testing.T) {
			log.V(1).Info(fmt.Sprintf("-------------- Running test: %s -------------", c.name))

			c.prep(&c)

			log.V(1).Info("setting up config renderer")
			r := NewRenderer(Config{
				Scheme: scheme,
				Logger: log.WithName("renderer"),
			})

			log.V(1).Info("preparing local storage")

			store.LiveKitMeshes.Flush()
			for i := range c.lkMeshes {
				store.LiveKitMeshes.Upsert(&c.lkMeshes[i])

				log.Info("lkmesh config added to store")
			}

			store.GatewayClasses.Flush()
			for i := range c.gwcs {
				store.GatewayClasses.Upsert(&c.gwcs[i])
			}

			store.GatewayConfigs.Flush()
			for i := range c.stnrGwconfigs {
				store.GatewayConfigs.Upsert(&c.stnrGwconfigs[i])
				log.Info("gateway config added to store")
			}

			store.Gateways.Flush()
			for i := range c.gws {
				store.Gateways.Upsert(&c.gws[i])
			}

			store.UDPRoutes.Flush()
			for i := range c.udpRoutes {
				store.UDPRoutes.Upsert(&c.udpRoutes[i])
			}

			store.Services.Flush()
			for i := range c.svcs {
				store.Services.Upsert(&c.svcs[i])
			}

			log.V(1).Info("starting renderer thread")
			ctx, cancel := context.WithCancel(context.Background())
			err := r.Start(ctx)
			assert.NoError(t, err, "renderer thread started")
			defer cancel()

			c.tester(t, r)

		})
	}
}
