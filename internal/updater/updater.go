package updater

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/l7mp/livekit-operator/internal/event"
	//corev1 "k8s.io/api/core/v1"
	//ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

	for _, cm := range uq.ConfigMaps.GetAll() {
		if op, err := u.upsertConfigMap(cm, gen); err != nil {
			u.log.Error(err, "cannot update configmap", "operation", op)
		}
	}

	for _, svc := range uq.Services.GetAll() {
		if op, err := u.upsertService(svc, gen); err != nil {
			u.log.Error(err, "cannot update service", "operation", op)
		}
	}

	for _, dp := range uq.Deployments.GetAll() {
		if op, err := u.upsertDeployment(dp, gen); err != nil {
			u.log.Error(err, "cannot update deployment", "operation", op)
		}
	}

	for _, s := range uq.Secrets.GetAll() {
		if op, err := u.upsertSecret(s, gen); err != nil {
			u.log.Error(err, "cannot update secret", "operation", op)
		}
	}

	for _, i := range uq.Issuer.GetAll() {
		if op, err := u.upsertIssuer(i, gen); err != nil {
			u.log.Error(err, "cannot update issuer", "operation", op)
		}
	}

	// TODO this might fail if the svc is not ready yet, there is need for a separate goroutine and watch for the svc
	for _, ss := range uq.StatefulSets.GetAll() {
		if op, err := u.upsertStatefulSet(ss, gen); err != nil {
			u.log.Error(err, "cannot update statefulset", "operation", op)
		}
	}

	for _, gwc := range uq.GatewayClasses.GetAll() {
		if op, err := u.upsertGatewayClass(gwc, gen); err != nil {
			u.log.Error(err, "cannot update gatewayclass", "operation", op)
		}
	}

	for _, gwc := range uq.GatewayConfigs.GetAll() {
		if op, err := u.upsertGatewayConfigs(gwc, gen); err != nil {
			u.log.Error(err, "cannot update gatewayconfigs", "operation", op)
		}
	}

	for _, gw := range uq.Gateways.GetAll() {
		if op, err := u.upsertGateway(gw, gen); err != nil {
			u.log.Error(err, "cannot update gateway", "operation", op)
		}
	}

	for _, udpr := range uq.UDPRoutes.GetAll() {
		if op, err := u.upsertUDPRoute(udpr, gen); err != nil {
			u.log.Error(err, "cannot update udproute", "operation", op)
		}
	}

	for _, httpr := range uq.HTTPRoutes.GetAll() {
		if op, err := u.upsertHTTPRoute(httpr, gen); err != nil {
			u.log.Error(err, "cannot update httproute", "operation", op)
		}
	}

	for _, sa := range uq.ServiceAccounts.GetAll() {
		if op, err := u.upsertServiceAccount(sa, gen); err != nil {
			u.log.Error(err, "cannot update serviceaccount", "operation", op)
		}
	}

	for _, r := range uq.ClusterRoles.GetAll() {
		if op, err := u.upsertClusterRole(r, gen); err != nil {
			u.log.Error(err, "cannot update clusterrole", "operation", op)
		}
	}

	for _, rb := range uq.ClusterRoleBindings.GetAll() {
		if op, err := u.upsertClusterRoleBinding(rb, gen); err != nil {
			u.log.Error(err, "cannot update clusterrolebinding", "operation", op)
		}
	}

	return nil
}
