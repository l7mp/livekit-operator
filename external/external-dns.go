package external

import (
	"context"
	"github.com/go-logr/logr"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	helmClient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
	"time"
)

var ExternalDNSChart = NewExternalDNS()

type ExternalDNS struct {
	Chart
}

func NewExternalDNS() *ExternalDNS {
	return &ExternalDNS{
		Chart: NewChartImpl(),
	}
}

func (e *ExternalDNS) InstallChart(ctx context.Context, logger logr.Logger) error {

	log := logger.WithName("ExternalDNS")

	chartRepo := repo.Entry{
		Name: "external-dns",
		URL:  "https://kubernetes-sigs.github.io/external-dns/",
	}

	opt := &helmClient.Options{
		Namespace:        opdefault.ExternalDNSNamespace,
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
		ReleaseName:     "external-dns",
		ChartName:       "external-dns/external-dns",
		Version:         "1.14.4",
		CreateNamespace: true,
		Namespace:       opdefault.ExternalDNSNamespace,
		UpgradeCRDs:     true,
		Wait:            true,
		WaitForJobs:     true,
		Timeout:         5 * time.Minute,
		ValuesYaml: `
sources:
- gateway-httproute
- gateway-grpcroute
- gateway-tcproute
- gateway-tlsroute
- gateway-udproute
provider=cloudflare
env:
- name: CF_API_KEY
  value: "V5g7yndcMqkzJh4-cOUmIQfeowaBPh2UwvUAUBIb"
- name: CF_API_EMAIL
  value: "info@l7mp.io"`,
	})

	// Install or upgrade a chart.
	go func() {

		externalDNSRelease, err := (*e.GetClient()).InstallOrUpgradeChart(ctx, e.GetChartSpec(), nil)
		if err != nil {
			// Rollback to the previous version of the release.
			if err := (*e.GetClient()).RollbackRelease(e.GetChartSpec()); err != nil {
				// In case rollback also failed throw hands in the air and then die
				panic(err)
			}
		}
		if externalDNSRelease.Info.Status == "deployed" {
			log.Info("chart installed", "release name", externalDNSRelease.Name, "status", externalDNSRelease.Info.Status)
		} else {
			log.Error(nil, "installation of the ExternalDNS chart was NOT successful", "status", externalDNSRelease.Info.Status)
			err := e.UninstallChart()
			if err != nil {
				log.Info("could not uninstall chart when it its' installation was NOT successful")
				return
			}
		}
	}()

	return err
}
