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

type Config struct {
	ControllerName      string
	Manager             manager.Manager
	RenderCh            chan event.Event
	UpdaterCh           chan event.Event
	ShouldInstallCharts bool
	Logger              logr.Logger
}

type Operator struct {
	ctx                             context.Context
	mgr                             manager.Manager
	renderCh, operatorCh, updaterCh chan event.Event
	shouldInstallCharts             bool
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

	if o.shouldInstallCharts {
		log.Info("start installing Envoy-Gateway")
		if err := external.EnvoyGatewayOperatorChart.InstallChart(o.ctx, o.logger); err != nil {
			return fmt.Errorf("cannot install Envoy-Gateway operator: %w", err)
		}

		log.Info("start installing Cert-Manager")
		if err := external.CertManagerChart.InstallChart(o.ctx, o.logger); err != nil {
			return fmt.Errorf("cannot install Cert-Manager: %w", err)
		}

		log.Info("start installing Stunner-Gateway-Operator")
		if err := external.StunnerGatewayOperatorChart.InstallChart(o.ctx, o.logger); err != nil {
			return fmt.Errorf("cannot install Stunner-Gateway-Operator: %w", err)
		}

		//log.Info("start installing ExternalDNS")
		//if err := external.ExternalDNSChart.InstallChart(o.ctx, o.logger); err != nil {
		//	return fmt.Errorf("cannot install ExternalDNS: %w", err)
		//}
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
	log := logr.Logger{}.WithName("Shutdown")
	<-ctx.Done()
	log.Info("Shutting down")
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		err := external.EnvoyGatewayOperatorChart.UninstallChart()
		if err != nil {
			log.Info("Error while deleting chart on shutdown")
		}
		wg.Done()
	}()
	go func() {
		err := external.CertManagerChart.UninstallChart()
		if err != nil {
			log.Info("Error while deleting chart on shutdown")
		}
		wg.Done()
	}()
	go func() {
		err := external.StunnerGatewayOperatorChart.UninstallChart()
		if err != nil {
			log.Info("Error while deleting chart on shutdown")
		}
		wg.Done()
	}()
	wg.Wait()
	return nil
}

func (o *Operator) GetOperatorChannel() chan event.Event {
	return o.operatorCh
}
