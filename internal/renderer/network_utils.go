package renderer

import (
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/l7mp/livekit-operator/api/v1alpha1"
	"github.com/l7mp/livekit-operator/internal/store"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	stnrauthsvc "github.com/l7mp/stunner-auth-service/pkg/types"
	"io"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"strings"
	"time"
)

const (
	stunnerAuthNamespace = opdefault.StunnerGatewayChartNamespace
	service              = "turn"
	maxRetries           = 5
	initialDelay         = 1 * time.Second
)

var (
	//baseUrl              = fmt.Sprintf("http://stunner-auth.%s.svc.cluster.local:8088/ice?service=%s", stunnerAuthNamespace, service)
	baseUrlWithLbStunner = fmt.Sprintf("http://34.118.95.206:8088/ice?service=%s", service)
)

func getIceConfigurationFromStunnerAuth(lkMesh v1alpha1.LiveKitMesh, log logr.Logger) (*stnrauthsvc.IceConfig, error) {
	log.WithName("getIceConfigurationFromStunnerAuth")

	parameterList, logMsg := createParameterList(lkMesh)
	if logMsg != nil && parameterList == nil {
		log.V(2).Info("Could not create parameter list", "reason", logMsg)
		return nil, nil
	}
	url := fmt.Sprintf("%s%s", baseUrlWithLbStunner, *parameterList)

	var iceConfig stnrauthsvc.IceConfig
	var resp *http.Response
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = sendRequest(url)
		if err != nil {
			log.Error(err, "Error sending request")
			time.Sleep(initialDelay * time.Duration(attempt+1))
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err, "error reading response body")
			return nil, err
		}

		log.V(2).Info("response body", "body", string(body), "status", resp.StatusCode)
		if resp.StatusCode == http.StatusOK {
			if err := json.Unmarshal(body, &iceConfig); err != nil {
				log.Error(err, "unmarshal response failed")
				return nil, err
			} else if iceConfig == (stnrauthsvc.IceConfig{}) {
				return nil, fmt.Errorf("response was not an ice config")
			}
			log.Info("response body", "body", string(body), "status", resp.StatusCode, "iceConfig", iceConfig)
			iceServers := *iceConfig.IceServers
			urls := *iceServers[0].Urls
			turnUrl := urls[0]
			address := strings.Split(turnUrl, ":")[1]
			if address != "0.0.0.0" {
				return &iceConfig, nil
			} else {
				log.V(2).Info("Received 0.0.0.0 address. Retrying...")
				time.Sleep(initialDelay * time.Duration(attempt+1))
				continue
			}
		} else if resp.StatusCode == http.StatusNotFound {
			log.V(2).Info("Received 404 status code. Retrying...")
			time.Sleep(initialDelay * time.Duration(attempt+1))
			continue
		}
	}

	return nil, fmt.Errorf("maximum retries reached or unexpected response")
}

func sendRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func createParameterList(lkMesh v1alpha1.LiveKitMesh) (*string, *string) {
	parameterList := ""
	gwConfig := store.GatewayConfigs.GetObject(types.NamespacedName{
		Namespace: lkMesh.Namespace,
		Name:      getStunnerGatewayConfigName(lkMesh.Name),
	})
	if gwConfig == nil {
		log := fmt.Sprintf("GatewayConfig not found in the global store, probably not ready yet")
		return nil, &log
	}
	gwConfigSpec := gwConfig.Spec
	if *gwConfigSpec.AuthType == "longterm" && gwConfigSpec.AuthLifetime != nil {
		parameterList = fmt.Sprintf("%s&ttl=%d", parameterList, *gwConfigSpec.AuthLifetime)
	}
	if gwConfigSpec.Username != nil {
		parameterList = fmt.Sprintf("%s&username=%s", parameterList, *gwConfigSpec.Username)
	}
	gatewayNamespace := lkMesh.Namespace
	gatewayName := getStunnerGatewayName(lkMesh.Name)

	parameterList = fmt.Sprintf("%s&namespace=%s&gateway=%s", parameterList, gatewayNamespace, gatewayName)
	return &parameterList, nil
}
