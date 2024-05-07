package external

import (
	"context"
	"github.com/go-logr/logr"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
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

	e.SetInstalled(false)

	opt := &helmClient.Options{
		Namespace:        opdefault.EnvoyGatewayChartNamespace, // Change this to the namespace you wish the client to operate in.
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
		ReleaseName: "eg",
		//Version:         "v0.0.0-latest",
		//FIXME v0.0.0 uses Gateway API v1 instead of v1b1 so we need to stuck to the v0.5.0 (latest stable)
		Version:         "v1.0.1",
		ChartName:       "oci://docker.io/envoyproxy/gateway-helm",
		CreateNamespace: true,
		Namespace:       opdefault.EnvoyGatewayChartNamespace,
		UpgradeCRDs:     true,
		SkipCRDs:        false,
		Wait:            true,
		WaitForJobs:     true,
		Timeout:         3 * time.Minute,
	})
	// Install or upgrade a chart.
	go func() {

		egwRelease, err := (*e.GetClient()).InstallOrUpgradeChart(ctx, e.GetChartSpec(), nil)
		if err != nil {
			// Rollback to the previous version of the release.
			log.Error(err, "failed to install chart")
			if err := (*e.GetClient()).RollbackRelease(e.GetChartSpec()); err != nil {
				// In case rollback also failed throw hands in the air and then die
				panic(err)
			}
		}
		if egwRelease.Info != nil {
			if egwRelease.Info.Status == "deployed" {
				log.Info("chart installed", "release name", egwRelease.Name, "status", egwRelease.Info.Status)
				e.SetInstalled(true)
			} else {
				log.Error(nil, "installation of the Envoy Gateway Operator was NOT successful", "status", egwRelease.Info.Status)
				err := e.UninstallChart()
				if err != nil {
					log.Info("could not uninstall chart when it its' installation was NOT successful")
					return
				}
			}
		}
	}()

	return err
}
