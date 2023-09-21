package operator

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/l7mp/livekit-operator/internal/controllers"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Config struct {
	ControllerName string
	Manager        manager.Manager
	Logger         logr.Logger
}

type Operator struct {
	ctx         context.Context
	mgr         manager.Manager
	log, logger logr.Logger
}

func NewOperator(config Config) *Operator {
	return &Operator{
		mgr:    config.Manager,
		logger: config.Logger,
	}
}

func (o *Operator) Start(ctx context.Context) error {
	log := o.logger.WithName("operator")
	o.log = log
	o.ctx = ctx
	log.Info("Starting operator")

	log.Info("starting LiveKitOperator controller")
	if err := controllers.RegisterLiveKitMeshController(o.mgr, o.logger); err != nil {
		return fmt.Errorf("cannot register LiveKitOperator controller: %w", err)
	}

	return nil
}
