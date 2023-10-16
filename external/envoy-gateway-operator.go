package external

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/l7mp/livekit-operator/internal/defaults"
	helmClient "github.com/mittwald/go-helm-client"
	"time"
)

var EnvoyGatewayOperatorChart = NewEnvoyGatewayOperator()

type EnvoyGatewayOperator struct {
	Chart
}

func NewEnvoyGatewayOperator() *EnvoyGatewayOperator {
	return &EnvoyGatewayOperator{
		Chart: NewChartImpl(),
	}
}

func (e *EnvoyGatewayOperator) InstallChart(ctx context.Context, logger logr.Logger) error {

	log := logger.WithName("Envoy-GW")

	opt := &helmClient.Options{
		Namespace:        defaults.EnvoyGatewayChartNamespace, // Change this to the namespace you wish the client to operate in.
		RepositoryCache:  "/tmp/.helmcache",
		RepositoryConfig: "/tmp/.helmrepo",
		Debug:            true,
		Linting:          true,
		DebugLog:         func(format string, v ...interface{}) {},
	}

	c, err := helmClient.New(opt)
	if err != nil {
		panic(err)
	} else {
		e.Chart.SetClient(c)
	}

	e.SetChartSpec(helmClient.ChartSpec{
		ReleaseName:     "eg",
		Version:         "v0.0.0-latest",
		ChartName:       "oci://docker.io/envoyproxy/gateway-helm",
		CreateNamespace: true,
		Namespace:       defaults.EnvoyGatewayChartNamespace,
		UpgradeCRDs:     true,
		Wait:            true,
		WaitForJobs:     true,
		Timeout:         3 * time.Minute,
	})
	// Install or upgrade a chart.
	go func() {

		egwRelease, err := (*e.GetClient()).InstallOrUpgradeChart(ctx, e.GetChartSpec(), nil)
		if err != nil {
			panic(err)
		}
		if egwRelease.Info.Status == "deployed" {
			log.Info("chart installed", "release name", egwRelease.Name, "status", egwRelease.Info.Status)
		} else {
			log.Error(nil, "installation of the Envoy Gateway Operator was NOT successful", "status", egwRelease.Info.Status)
			err := e.UninstallChart()
			if err != nil {
				log.Info("could not uninstall chart when it its' installation was NOT successful")
				return
			}
		}
	}()

	return err
}
