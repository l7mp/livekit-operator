package external

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/l7mp/livekit-operator/internal/defaults"
	helmClient "github.com/mittwald/go-helm-client"
	"github.com/mittwald/go-helm-client/values"
	"helm.sh/helm/v3/pkg/repo"
	"time"
)

var StunnerGatewayOperatorChart = NewStunnerGatewayOperator()

type StunnerGatewayOperator struct {
	Chart
}

func NewStunnerGatewayOperator() *StunnerGatewayOperator {
	return &StunnerGatewayOperator{
		Chart: NewChartImpl(),
	}
}

func (e *StunnerGatewayOperator) InstallChart(ctx context.Context, logger logr.Logger) error {

	log := logger.WithName("Stunner-GW-operator")

	chartRepo := repo.Entry{
		Name: "stunner",
		URL:  "https://l7mp.io/stunner",
	}

	opt := &helmClient.Options{
		//TODO
		Namespace:        defaults.StunnerGatewayChartNamespace, // Change this to the namespace you wish the client to operate in.
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

	// Add a chart-repository to the client.
	if err := c.AddOrUpdateChartRepo(chartRepo); err != nil {
		panic(err)
	}

	e.SetChartSpec(helmClient.ChartSpec{
		//TODO
		ReleaseName:     "stunner-gateway-operator",
		ChartName:       "stunner/stunner-gateway-operator",
		CreateNamespace: true,
		Namespace:       defaults.StunnerGatewayChartNamespace,
		UpgradeCRDs:     true,
		Wait:            true,
		WaitForJobs:     true,
		Timeout:         5 * time.Minute,
		ValuesOptions: values.Options{Values: []string{
			"stunnerGatewayOperator.dataplane.mode=managed",
		}},
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
