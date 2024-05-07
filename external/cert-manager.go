package external

import (
	"context"
	"github.com/go-logr/logr"
	opdefault "github.com/l7mp/livekit-operator/pkg/config"
	helmClient "github.com/mittwald/go-helm-client"
	"github.com/mittwald/go-helm-client/values"
	"helm.sh/helm/v3/pkg/repo"
	"time"
)

var CertManagerChart = NewCertManager()

type CertManager struct {
	Chart
}

func NewCertManager() *CertManager {
	return &CertManager{
		Chart: NewChartImpl(),
	}
}

func (e *CertManager) InstallChart(ctx context.Context, logger logr.Logger) error {
	log := logger.WithName("Cert-Manager")

	e.SetInstalled(false)

	chartRepo := repo.Entry{
		Name: "jetstack",
		URL:  "https://charts.jetstack.io",
	}

	opt := &helmClient.Options{
		//TODO
		Namespace:        opdefault.CertManagerChartNamespace, // Change this to the namespace you wish the client to operate in.
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
		ReleaseName:     "cert-manager",
		ChartName:       "jetstack/cert-manager",
		CreateNamespace: true,
		Namespace:       opdefault.CertManagerChartNamespace,
		UpgradeCRDs:     true,
		Wait:            true,
		WaitForJobs:     true,
		Timeout:         5 * time.Minute,
		ValuesOptions: values.Options{Values: []string{
			"installCRDs=true",
			"featureGates=ExperimentalGatewayAPISupport=true",
		}},
	})
	// Install or upgrade a chart.
	go func() {

		egwRelease, err := (*e.GetClient()).InstallOrUpgradeChart(ctx, e.GetChartSpec(), nil)
		if err != nil {
			// Rollback to the previous version of the release.
			if err := (*e.GetClient()).RollbackRelease(e.GetChartSpec()); err != nil {
				// In case rollback also failed throw hands in the air and then die
				panic(err)
			}
		}
		if egwRelease.Info.Status == "deployed" {
			log.Info("chart installed", "release name", egwRelease.Name, "status", egwRelease.Info.Status)
			e.SetInstalled(true)
		} else {
			log.Error(nil, "installation of the chart was NOT successful", "release name", egwRelease.Name, "status", egwRelease.Info.Status)
			err := e.UninstallChart()
			if err != nil {
				log.Info("could not uninstall chart when its' installation was NOT successful")
				return
			}
		}
	}()

	return err
}
