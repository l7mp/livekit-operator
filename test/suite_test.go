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

package test

import (
	"context"
	cert "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	//"github.com/l7mp/livekit-operator/internal/controllers"
	"github.com/l7mp/livekit-operator/internal/operator"
	"github.com/l7mp/livekit-operator/internal/renderer"
	"github.com/l7mp/livekit-operator/internal/updater"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/testutils"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

const (
	logLevel = -4
	timeout  = time.Second * 10
	interval = time.Millisecond * 250
)

var (
	// Resources
	testNs     *corev1.Namespace
	testLkMesh *lkstnv1a1.LiveKitMesh
	//testConfigMap *corev1.ConfigMap
	// Globals
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
	scheme    *runtime.Scheme
)

func InitResources() {
	ctrl.Log.Info("testns", "ns", testNs)
	testNs = testutils.TestNs.DeepCopy()
	testLkMesh = testutils.TestLkMesh.DeepCopy()
	//testConfigMap = testutils.TestConfigMap.DeepCopy()
	scheme = runtime.NewScheme()
	ctx, cancel = context.WithCancel(context.Background())
}

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	ctrl.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true), func(o *zap.Options) {
		o.TimeEncoder = zapcore.RFC3339NanoTimeEncoder
	}, zap.Level(zapcore.Level(logLevel))))
	setupLog := ctrl.Log.WithName("setup")

	By("bootstrapping test environment")
	InitResources()
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = clientgoscheme.AddToScheme(scheme)

	// LiveKitMesh CRD scheme
	err = lkstnv1a1.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	// Cert-manager scheme
	err = cert.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	setupLog.Info("creating a testing namespace")
	Expect(k8sClient.Create(ctx, testNs)).Should(Succeed())
	setupLog.Info("created a testing namespace")

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	setupLog.Info("setting up LikeKitMesh config renderer")
	r := renderer.NewRenderer(renderer.Config{
		Scheme: scheme,
		Logger: ctrl.Log,
	})
	Expect(r).NotTo(BeNil())

	setupLog.Info("setting up updater client")
	u := updater.NewUpdater(updater.Config{
		Manager: k8sManager,
		Logger:  ctrl.Log,
	})
	Expect(u).NotTo(BeNil())

	setupLog.Info("setting up operator")
	op := operator.NewOperator(operator.Config{
		ControllerName:      opdefault.DefaultControllerName,
		Manager:             k8sManager,
		RenderCh:            r.GetRenderChannel(),
		UpdaterCh:           u.GetUpdaterChannel(),
		ShouldInstallCharts: false,
		Logger:              ctrl.Log,
	})
	Expect(op).NotTo(BeNil())

	r.SetOperatorChannel(op.GetOperatorChannel())

	setupLog.Info("Start renderer thread")
	err = r.Start(ctx)
	Expect(err).NotTo(HaveOccurred())

	setupLog.Info("Start updater thread")
	err = u.Start(ctx)
	Expect(err).NotTo(HaveOccurred())

	setupLog.Info("Start operator thread")
	err = op.Start(ctx)
	Expect(err).NotTo(HaveOccurred())

	setupLog.Info("starting manager")
	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
})

var _ = AfterSuite(func() {
	By("removing test namespace")
	Expect(k8sClient.Delete(ctx, testNs)).Should(Succeed())

	cancel()

	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

type LiveKitMeshMutator func(lkMesh *lkstnv1a1.LiveKitMesh)

func createOrUpdateLiveKitMesh(template *lkstnv1a1.LiveKitMesh, f LiveKitMeshMutator) {
	current := &lkstnv1a1.LiveKitMesh{
		ObjectMeta: metav1.ObjectMeta{
			Name:      template.GetName(),
			Namespace: template.GetNamespace(),
		}}

	_, err := createOrUpdate(ctx, k8sClient, current, func() error {
		template.Spec.DeepCopyInto(&current.Spec)
		if f != nil {
			f(current)
		}
		return nil
	})
	Expect(err).Should(Succeed())
}

// createOrUpdate will retry when ctrlutil.CreateOrUpdate fails 5 times
func createOrUpdate(ctx context.Context, c client.Client, obj client.Object, f ctrlutil.MutateFn) (ctrlutil.OperationResult, error) {
	var res ctrlutil.OperationResult
	var err error

	for i := 1; i < 5; i++ {
		if res, err = ctrlutil.CreateOrUpdate(ctx, c, obj, f); err == nil {
			return res, err
		}
	}
	return res, err
}

var _ = Describe("Integration test:", func() {

	Context("When creating a LiveKitMesh", func() {
		It("should survive creating a minimal config", func() {
			current := &lkstnv1a1.LiveKitMesh{ObjectMeta: metav1.ObjectMeta{
				Name:      "testlivekitmesh",
				Namespace: testutils.TestNsName,
			}}
			_, err := ctrlutil.CreateOrUpdate(ctx, k8sClient, current, func() error {
				testLkMesh.Spec.DeepCopyInto(&current.Spec)
				return nil
			})
			Expect(err).Should(Succeed())

		})

		/*		It("should render configmap", func() {
					lookUpKey := types.NamespacedName{
						Namespace: testutils.TestNsName,
						Name:      renderer.ConfigMapNameFormat(*testutils.TestLkMesh.Spec.Components.LiveKit.Deployment.Name),
					}

					cm := &corev1.ConfigMap{}
					Eventually(func() error {
						err := k8sClient.Get(ctx, lookUpKey, cm)
						if err != nil {
							return err
						}
						return nil
					}, timeout, interval).Should(BeNil())

					Expect(cm.Labels[opdefault.RelatedComponent]).To(Equal(opdefault.ComponentLiveKit))
				})

				It("should render deployment", func() {
					lookUpKey := types.NamespacedName{
						Namespace: testutils.TestNsName,
						Name:      *testutils.TestLkMesh.Spec.Components.LiveKit.Deployment.Name,
					}
					dp := &v1.Deployment{}
					Eventually(func() error {
						err := k8sClient.Get(ctx, lookUpKey, dp)
						if err != nil {
							return err
						}
						return nil
					}, timeout, interval).Should(BeNil())

					Expect(dp.Labels[opdefault.RelatedComponent]).To(Equal(opdefault.ComponentLiveKit))

				})

				It("should render service", func() {
					lookUpKey := types.NamespacedName{
						Namespace: testutils.TestNsName,
						Name:      renderer.ServiceNameFormat(*testutils.TestLkMesh.Spec.Components.LiveKit.Deployment.Name),
					}

					svc := &corev1.Service{}
					Eventually(func() error {
						err := k8sClient.Get(ctx, lookUpKey, svc)
						if err != nil {
							return err
						}
						return nil
					}, timeout, interval).Should(BeTrue())

					Expect(svc.Labels[opdefault.RelatedComponent]).To(Equal(opdefault.ComponentLiveKit))

				})

				It("should create issuer and secret", func() {
					lookUpKey := types.NamespacedName{
						Namespace: testutils.TestNsName,
						Name:      "cloudflare-issuer",
					}

					issuer := &cert.Issuer{}
					Eventually(func() error {
						err := k8sClient.Get(ctx, lookUpKey, issuer)
						if err != nil {
							return err
						}
						return nil
					}, timeout, interval).Should(BeTrue())

					Expect(issuer.Labels[opdefault.RelatedComponent]).To(Equal(opdefault.ComponentCertManager))

					lookUpKey = types.NamespacedName{
						Namespace: testutils.TestNsName,
						Name:      "cloudflare-api-token-secret",
					}

					secret := &corev1.Secret{}
					Eventually(func() error {
						err := k8sClient.Get(ctx, lookUpKey, secret)
						if err != nil {
							return err
						}
						return nil
					}, timeout, interval).Should(BeTrue())

					Expect(secret.Labels[opdefault.RelatedComponent]).To(Equal(opdefault.ComponentCertManager))

				})*/

	})
})
