package external

import (
	"context"
	"github.com/go-logr/logr"
	helmClient "github.com/mittwald/go-helm-client"
)

type Chart interface {
	InstallChart(ctx context.Context, logger logr.Logger) error
	UninstallChart() error
	SetChartSpec(spec helmClient.ChartSpec)
	GetChartSpec() *helmClient.ChartSpec
	SetClient(client helmClient.Client)
	GetClient() *helmClient.Client
	IsInstalled() bool
	SetInstalled(installed bool)
}

type chartImpl struct {
	chartSpec helmClient.ChartSpec
	client    helmClient.Client
	installed bool
}

func NewChartImpl() Chart {
	return &chartImpl{}
}

func (c *chartImpl) InstallChart(ctx context.Context, logger logr.Logger) error {
	return nil
}

func (c *chartImpl) UninstallChart() error {
	//fmt.Println("uninstalling", "chart", c.chartSpec.ChartName)
	//return c.client.UninstallRelease(&c.chartSpec)
	return nil
}

func (c *chartImpl) SetChartSpec(spec helmClient.ChartSpec) {
	c.chartSpec = spec
}

func (c *chartImpl) GetChartSpec() *helmClient.ChartSpec {
	return &c.chartSpec
}

func (c *chartImpl) SetClient(client helmClient.Client) {
	c.client = client
}

func (c *chartImpl) GetClient() *helmClient.Client {
	return &c.client
}

func (c *chartImpl) IsInstalled() bool {
	return c.installed
}

func (c *chartImpl) SetInstalled(installed bool) {
	c.installed = installed
}
