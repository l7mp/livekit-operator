package operator

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/l7mp/livekit-operator/external"
	"github.com/l7mp/livekit-operator/internal/controllers"
	"github.com/l7mp/livekit-operator/internal/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sync"
)

const channelBufferSize = 200

type ChartInstallConfig struct {
	ShouldInstallStunnerGatewayChart bool
	ShouldInstallEnvoyGatewayChart   bool
	ShouldInstallCertManagerChart    bool
}

type Config struct {
	ControllerName      string
	Manager             manager.Manager
	RenderCh            chan event.Event
	UpdaterCh           chan event.Event
	ShouldInstallCharts ChartInstallConfig
	Logger              logr.Logger
}

type Operator struct {
	ctx                             context.Context
	mgr                             manager.Manager
	renderCh, operatorCh, updaterCh chan event.Event
	shouldInstallCharts             ChartInstallConfig
	log, logger                     logr.Logger
}

func NewOperator(config Config) *Operator {
	return &Operator{
		mgr:                 config.Manager,
		renderCh:            config.RenderCh,
		operatorCh:          make(chan event.Event, channelBufferSize),
		updaterCh:           config.UpdaterCh,
		shouldInstallCharts: config.ShouldInstallCharts,
		logger:              config.Logger,
	}
}

func (o *Operator) Start(ctx context.Context) error {
	log := o.logger.WithName("operator")
	o.log = log
	o.ctx = ctx
	log.Info("Starting operator")

	if o.shouldInstallCharts.ShouldInstallStunnerGatewayChart {
		log.Info("start installing Stunner-Gateway-Operator")
		if err := external.StunnerGatewayOperatorChart.InstallChart(o.ctx, o.logger); err != nil {
			return fmt.Errorf("cannot install Stunner-Gateway-Operator: %w", err)
		}
	} else {
		log.Info("NOT installing Stunner-Gateway chart", "cli flag install-stunner-gateway-chart", o.shouldInstallCharts.ShouldInstallStunnerGatewayChart)
	}

	if o.shouldInstallCharts.ShouldInstallEnvoyGatewayChart {
		log.Info("start installing Envoy-Gateway")
		if err := external.EnvoyGatewayOperatorChart.InstallChart(o.ctx, o.logger); err != nil {
			return fmt.Errorf("cannot install Envoy-Gateway operator: %w", err)
		}
	} else {
		log.Info("NOT installing Envoy-Gateway chart", "cli flag install-envoy-gateway-chart", o.shouldInstallCharts.ShouldInstallEnvoyGatewayChart)
	}

	if o.shouldInstallCharts.ShouldInstallCertManagerChart {
		log.Info("start installing Cert-Manager")
		if err := external.CertManagerChart.InstallChart(o.ctx, o.logger); err != nil {
			return fmt.Errorf("cannot install Cert-Manager: %w", err)
		}
	} else {
		log.Info("NOT installing Cert-Manager chart", "cli flag install-cert-manager-chart", o.shouldInstallCharts.ShouldInstallCertManagerChart)
	}

	log.Info("starting LiveKitOperator controller")
	if err := controllers.RegisterLiveKitMeshController(o.mgr, o.renderCh, o.logger); err != nil {
		return fmt.Errorf("cannot register LiveKitOperator controller: %w", err)
	}

	o.eventLoop(ctx)

	return nil
}

func (o *Operator) eventLoop(ctx context.Context) {

	go func() {
		defer close(o.operatorCh)
		for {
			select {
			case e := <-o.operatorCh:
				o.log.Info("e", "e", e)
				switch e.GetType() {
				case event.TypeUpdate:
					o.updaterCh <- e
					//TODO

				}
			case <-ctx.Done():
				return
			}
		}
	}()

}

func HandleCleanup(ctx context.Context) error {
	//FIXME error handling https://trstringer.com/concurrent-error-handling-go/
	log := logr.Logger{}.WithName("Shutdown")
	<-ctx.Done()
	log.Info("Shutting down")
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		if external.EnvoyGatewayOperatorChart.IsInstalled() {
			err := external.EnvoyGatewayOperatorChart.UninstallChart()
			if err != nil {
				log.Error(err, "error while deleting envoy gateway operator chart on shutdown")
			}
			external.EnvoyGatewayOperatorChart.SetInstalled(false)
		}
		wg.Done()
	}()
	go func() {
		if external.CertManagerChart.IsInstalled() {
			err := external.CertManagerChart.UninstallChart()
			if err != nil {
				log.Error(err, "error while deleting cert manager chart on shutdown")
			}
			external.CertManagerChart.SetInstalled(false)
		}
		wg.Done()
	}()
	go func() {
		if external.StunnerGatewayOperatorChart.IsInstalled() {
			err := external.StunnerGatewayOperatorChart.UninstallChart()
			if err != nil {
				log.Error(err, "error while deleting stunner gateway operator chart on shutdown")
			}
			external.StunnerGatewayOperatorChart.SetInstalled(false)
		}
		wg.Done()
	}()
	wg.Wait()
	return nil
}

func (o *Operator) GetOperatorChannel() chan event.Event {
	return o.operatorCh
}
