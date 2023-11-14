package renderer

import (
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
	name     string
	lkMeshes []lkstnv1a1.LiveKitMesh
	svcs     []corev1.Service
	dps      []appv1.Deployment
	cms      []corev1.ConfigMap
	prep     func(c *renderTestConfig)
	tester   func(t *testing.T, r *Renderer)
}
