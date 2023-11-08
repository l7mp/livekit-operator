package updater

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/l7mp/livekit-operator/internal/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Config struct {
	Manager manager.Manager
	Logger  logr.Logger
}

type Updater struct {
	ctx       context.Context
	manager   manager.Manager
	updaterCh chan event.Event
	log       logr.Logger
}

func NewUpdater(cfg Config) *Updater {
	return &Updater{
		manager:   cfg.Manager,
		updaterCh: make(chan event.Event, 10),
		log:       cfg.Logger.WithName("updater"),
	}
}

func (u *Updater) Start(ctx context.Context) error {
	u.ctx = ctx

	go func() {
		defer close(u.updaterCh)
		for {
			select {
			case e := <-u.updaterCh:
				if e.GetType() != event.TypeUpdate {
					u.log.Info("renderer thread received unknown event",
						"event", e.GetType().String())
					continue
				}
				u.log.Info("Update event received")
				ev := e.(*event.Update)
				err := u.processUpdate(ev)
				if err != nil {
					u.log.Error(err, "could not process update event", "event", ev.String())
				}

			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (u *Updater) GetUpdaterChannel() chan event.Event {
	return u.updaterCh
}

func (u *Updater) processUpdate(e *event.Update) error {
	gen := e.Generation
	u.log.Info("processing update event", "generation", gen, "update", e.String())

	uq := e.UpsertQueue

	for _, svc := range uq.Services.GetAll() {
		if op, err := u.upsertService(svc, gen); err != nil {
			u.log.Error(err, "cannot update service", "operation", op)
		}
	}

	for _, dp := range uq.Deployments.GetAll() {
		if op, err := u.upsertDeployment(dp, gen); err != nil {
			u.log.Error(err, "cannot update service", "operation", op)
		}
	}
	return nil
}
