package renderer

import (
	"encoding/json"
	"fmt"
	lkstnv1a1 "github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/store"
	"github.com/l7mp/livekit-operator/internal/testutils"
	stnrgwv1 "github.com/l7mp/stunner-gateway-operator/api/v1"
	"github.com/stretchr/testify/assert"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"net/http/httptest"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	"testing"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/types"
	// "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestNetworkUtils(t *testing.T) {
	renderTester(t, []renderTestConfig{
		{
			name:          "creating parameter list ok",
			stnrGwconfigs: []stnrgwv1.GatewayConfig{testutils.TestGatewayConfig},
			lkMeshes:      []lkstnv1a1.LiveKitMesh{testutils.TestLkMesh},
			svcs:          []corev1.Service{},
			dps:           []appv1.Deployment{},
			cms:           []corev1.ConfigMap{},
			gwcs:          []gwapiv1.GatewayClass{},
			gws:           []gwapiv1.Gateway{},
			udpRoutes:     []stnrgwv1.UDPRoute{},
			prep:          func(c *renderTestConfig) {},
			tester: func(t *testing.T, r *Renderer) {
				lkMesh := store.LiveKitMeshes.GetAll()[0]
				parameterList, log := createParameterList(*lkMesh)
				r.log.WithName("test").Info("log", "log", log)
				assert.Nil(t, log, "log should be nil")
				assert.Equal(t, "&username=testuser&namespace=testnamespace&gateway=testlivekitmesh-stunner-udp-gateway", *parameterList, "parameter list")
			},
		},
	})
	renderTester(t, []renderTestConfig{
		{
			name:          "get ice config from stunner auth ok",
			stnrGwconfigs: []stnrgwv1.GatewayConfig{testutils.TestGatewayConfig},
			lkMeshes:      []lkstnv1a1.LiveKitMesh{testutils.TestLkMesh},
			svcs:          []corev1.Service{},
			dps:           []appv1.Deployment{},
			cms:           []corev1.ConfigMap{},
			gwcs:          []gwapiv1.GatewayClass{},
			gws:           []gwapiv1.Gateway{},
			udpRoutes:     []stnrgwv1.UDPRoute{},
			prep:          func(c *renderTestConfig) {},
			tester: func(t *testing.T, r *Renderer) {
				mockLiveKitMesh := store.LiveKitMeshes.GetAll()[0]
				// Start a mock HTTP server
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Mock response JSON
					iceConfig := testutils.TestTurnIceConfig
					responseJSON, _ := json.Marshal(iceConfig)

					// Set response headers
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)

					// Write response body
					_, _ = w.Write(responseJSON)
				}))
				defer mockServer.Close()

				// Replace baseUrlWithLbStunner with mock server URL
				baseUrlWithLbStunner = fmt.Sprintf("%s?ice", mockServer.URL)

				// Run the function under test
				result, err := getIceConfigurationFromStunnerAuth(*mockLiveKitMesh, r.log)
				assert.Nil(t, err, "no error from getIceConfigurationFromStunnerAuth")

				assert.Equal(t, testutils.TestTurnIceConfig, *result, "ice config is ok")

				// Additional assertions if needed
			},
		},
	})
	renderTester(t, []renderTestConfig{
		{
			name:          "get ice config from stunner auth results in 404 max tried reached",
			stnrGwconfigs: []stnrgwv1.GatewayConfig{testutils.TestGatewayConfig},
			lkMeshes:      []lkstnv1a1.LiveKitMesh{testutils.TestLkMesh},
			svcs:          []corev1.Service{},
			dps:           []appv1.Deployment{},
			cms:           []corev1.ConfigMap{},
			gwcs:          []gwapiv1.GatewayClass{},
			gws:           []gwapiv1.Gateway{},
			udpRoutes:     []stnrgwv1.UDPRoute{},
			prep: func(c *renderTestConfig) {

			},
			tester: func(t *testing.T, r *Renderer) {
				mockLiveKitMesh := store.LiveKitMeshes.GetAll()[0]
				// Start a mock HTTP server
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Mock response JSON
					//iceConfig := testutils.TestTurnIceConfig
					//responseJSON, _ := json.Marshal(iceConfig)

					// Set response headers
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusNotFound)

					// Write response body
					//_, _ = w.Write()
				}))
				defer mockServer.Close()
				baseUrlWithLbStunner = fmt.Sprintf("%s?ice", mockServer.URL)

				// Run the function under test
				result, err := getIceConfigurationFromStunnerAuth(*mockLiveKitMesh, r.log)
				assert.NotNil(t, err, "error from getIceConfigurationFromStunnerAuth")
				assert.EqualErrorf(t, err, "maximum retries reached or unexpected response", err.Error(), "error message")
				assert.Nil(t, nil, result, "ice config is nil")

			},
		},
	})
	renderTester(t, []renderTestConfig{
		{
			name:          "get ice config from stunner auth results in response was not an ice config",
			stnrGwconfigs: []stnrgwv1.GatewayConfig{testutils.TestGatewayConfig},
			lkMeshes:      []lkstnv1a1.LiveKitMesh{testutils.TestLkMesh},
			svcs:          []corev1.Service{},
			dps:           []appv1.Deployment{},
			cms:           []corev1.ConfigMap{},
			gwcs:          []gwapiv1.GatewayClass{},
			gws:           []gwapiv1.Gateway{},
			udpRoutes:     []stnrgwv1.UDPRoute{},
			prep: func(c *renderTestConfig) {

			},
			tester: func(t *testing.T, r *Renderer) {
				mockLiveKitMesh := store.LiveKitMeshes.GetAll()[0]
				// Start a mock HTTP server
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Mock response JSON
					testResponse := make(map[string]string)
					testResponse["should"] = "fail"
					responseJSON, _ := json.Marshal(testResponse)

					// Set response headers
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)

					//Write response body
					_, _ = w.Write(responseJSON)
				}))
				defer mockServer.Close()
				baseUrlWithLbStunner = fmt.Sprintf("%s?ice", mockServer.URL)

				// Run the function under test
				result, err := getIceConfigurationFromStunnerAuth(*mockLiveKitMesh, r.log)
				assert.NotNil(t, err, "error from getIceConfigurationFromStunnerAuth")
				assert.EqualErrorf(t, err, "response was not an ice config", err.Error(), "error message: response was not an ice config")
				assert.Nil(t, nil, result, "ice config is nil")

			},
		},
	})
}
